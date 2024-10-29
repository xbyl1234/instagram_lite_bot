package http

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"sync"
)

type Conn struct {
	net.Conn

	proxy *HttpProxy
	addr  string

	chShakeFinished chan struct{}
	muShake         sync.Mutex
	reqBuf          io.ReadWriter
}

func NewConn(c net.Conn, proxy *HttpProxy, addr string) *Conn {
	return &Conn{
		Conn:            c,
		proxy:           proxy,
		addr:            addr,
		chShakeFinished: make(chan struct{}),
		reqBuf:          nil,
	}
}

func (c *Conn) Write(b []byte) (n int, err error) {
	c.muShake.Lock()
	select {
	case <-c.chShakeFinished:
		c.muShake.Unlock()
		return c.Conn.Write(b)
	default:
		// Handshake
		defer c.muShake.Unlock()
		_, firstLine, _ := bufio.ScanLines(b, true)
		isHttpReq := regexp.MustCompile(`^\S+ \S+ HTTP/[\d.]+$`).Match(firstLine)

		var req *http.Request
		if isHttpReq {
			// HTTP Request

			if c.reqBuf == nil {
				c.reqBuf = bytes.NewBuffer(b)
			} else {
				c.reqBuf.Write(b)
			}
			req, err = http.ReadRequest(bufio.NewReader(c.reqBuf))
			if err != nil {
				if errors.Is(err, io.ErrUnexpectedEOF) {
					// Request more data.
					return len(b), nil
				}
				// Error
				return 0, err
			}
			// Clear the buf
			c.reqBuf = nil

			req.URL.Scheme = "http"
			req.URL.Host = c.addr
		} else {
			// Arbitrary TCP

			// HACK. http.ReadRequest also does this.
			reqURL, err := url.Parse("http://" + c.addr)
			if err != nil {
				return 0, err
			}
			reqURL.Scheme = ""

			req, err = http.NewRequest("CONNECT", reqURL.String(), nil)
			if err != nil {
				return 0, err
			}
		}
		req.Close = false
		if c.proxy.HaveAuth {
			req.Header.Set("Proxy-Authorization", base64.StdEncoding.EncodeToString([]byte(c.proxy.Username+":"+c.proxy.Password)))
		}
		// https://www.rfc-editor.org/rfc/rfc7230#appendix-A.1.2
		// As a result, clients are encouraged not to send the Proxy-Connection header field in any requests.
		if len(req.Header.Values("Proxy-Connection")) > 0 {
			req.Header.Del("Proxy-Connection")
		}

		err = req.WriteProxy(c.Conn)
		if err != nil {
			return 0, err
		}

		if isHttpReq {
			// Allow Read here to void race.
			close(c.chShakeFinished)
			return len(b), nil
		} else {
			// We should read tcp connection here, and we will be guaranteed higher priority by chShakeFinished.
			resp, err := http.ReadResponse(bufio.NewReader(c.Conn), req)
			if err != nil {
				if resp != nil {
					resp.Body.Close()
				}
				return 0, err
			}
			resp.Body.Close()
			// Allow Read here to avoid race.
			close(c.chShakeFinished)
			if resp.StatusCode != 200 {
				err = fmt.Errorf("connect server using proxy error, StatusCode [%d]", resp.StatusCode)
				return 0, err
			}
			return c.Conn.Write(b)
		}
	}
}

func (c *Conn) Read(b []byte) (n int, err error) {
	<-c.chShakeFinished
	return c.Conn.Read(b)
}
