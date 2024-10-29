package grpc_lite

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"

	"encoding/binary"
	"golang.org/x/net/http2"
)

var sevenbits = [...]byte{
	0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
	0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
	0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f,
	0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3a, 0x3b, 0x3c, 0x3d, 0x3e, 0x3f,
	0x40, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4a, 0x4b, 0x4c, 0x4d, 0x4e, 0x4f,
	0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0x59, 0x5a, 0x5b, 0x5c, 0x5d, 0x5e, 0x5f,
	0x60, 0x61, 0x62, 0x63, 0x64, 0x65, 0x66, 0x67, 0x68, 0x69, 0x6a, 0x6b, 0x6c, 0x6d, 0x6e, 0x6f,
	0x70, 0x71, 0x72, 0x73, 0x74, 0x75, 0x76, 0x77, 0x78, 0x79, 0x7a, 0x7b, 0x7c, 0x7d, 0x7e, 0x7f,
}

// AppendUleb128 appends v to b using unsigned LEB128 encoding.
func AppendUleb128(b []byte, v uint64) []byte {
	// If it's less than or equal to 7-bit
	if v < 0x80 {
		return append(b, sevenbits[v])
	}

	for {
		c := uint8(v & 0x7f)
		v >>= 7

		if v != 0 {
			c |= 0x80
		}

		b = append(b, c)

		if c&0x80 == 0 {
			break
		}
	}

	return b
}

// AppendSleb128 appends v to b using signed LEB128 encoding.
func AppendSleb128(b []byte, v int64) []byte {
	// If it's less than or equal to 7-bit
	if v >= 0 && v <= 0x3f {
		return append(b, sevenbits[v])
	} else if v < 0 && v >= ^0x3f {
		return append(b, sevenbits[0x80+v])
	}

	for {
		c := uint8(v & 0x7f)
		s := uint8(v & 0x40)
		v >>= 7

		if (v != -1 || s == 0) && (v != 0 || s != 0) {
			c |= 0x80
		}

		b = append(b, c)

		if c&0x80 == 0 {
			break
		}
	}

	return b
}

// DecodeUleb128 decodes b to u with unsigned LEB128 encoding and returns the
// number of bytes read. On error (bad encoded b), n will be 0 and therefore u
// must not be trusted.
func DecodeUleb128(b []byte) (u uint64, n uint8) {
	l := uint8(len(b) & 0xff)
	// The longest LEB128 encoded sequence is 10 byte long (9 0xff's and 1 0x7f)
	// so make sure we won't overflow.
	if l > 10 {
		l = 10
	}

	var i uint8
	for i = 0; i < l; i++ {
		u |= uint64(b[i]&0x7f) << (7 * i)
		if b[i]&0x80 == 0 {
			n = uint8(i + 1)
			return
		}
	}

	return
}

// DecodeSleb128 decodes b to s with signed LEB128 encoding and returns the
// number of bytes read. On error (bad encoded b), n will be 0 and therefore s
// must not be trusted.
func DecodeSleb128(b []byte) (s int64, n uint8) {
	l := uint8(len(b) & 0xff)
	if l > 10 {
		l = 10
	}

	var i uint8
	for i = 0; i < l; i++ {
		s |= int64(b[i]&0x7f) << (7 * i)
		if b[i]&0x80 == 0 {
			// If it's signed
			if b[i]&0x40 != 0 {
				s |= ^0 << (7 * (i + 1))
			}
			n = uint8(i + 1)
			return
		}
	}

	return
}

type GunConn struct {
	reader io.Reader
	writer io.Writer
	closer io.Closer
	local  net.Addr
	remote net.Addr
	// mu protect done
	mu   sync.Mutex
	done chan struct{}

	toRead []byte
	readAt int
}

type Client struct {
	client  *http.Client
	url     *url.URL
	headers http.Header
}

type Config struct {
	RemoteAddr  string
	ServerName  string
	ServiceName string
	Cleartext   bool
	tlsConfig   *tls.Config
}

func NewGunClient(config *Config) *Client {
	var dialFunc func(network, addr string, cfg *tls.Config) (net.Conn, error) = nil
	if config.Cleartext {
		dialFunc = func(network, addr string, cfg *tls.Config) (net.Conn, error) {
			return net.Dial(network, addr)
		}
	} else {
		dialFunc = func(network, addr string, cfg *tls.Config) (net.Conn, error) {
			pconn, err := net.Dial(network, addr)
			if err != nil {
				return nil, err
			}

			cn := tls.Client(pconn, cfg)
			if err := cn.Handshake(); err != nil {
				return nil, err
			}
			state := cn.ConnectionState()
			if p := state.NegotiatedProtocol; p != http2.NextProtoTLS {
				return nil, errors.New("http2: unexpected ALPN protocol " + p + "; want q" + http2.NextProtoTLS)
			}
			return cn, nil
		}
	}

	if config.tlsConfig == nil && config.ServerName != "" {
		config.tlsConfig = new(tls.Config)
		config.tlsConfig.ServerName = config.ServerName
		config.tlsConfig.NextProtos = []string{"h2"}
	}

	client := &http.Client{
		Transport: &http2.Transport{
			DialTLS:            dialFunc,
			TLSClientConfig:    config.tlsConfig,
			AllowHTTP:          false,
			DisableCompression: true,
			ReadIdleTimeout:    0,
			PingTimeout:        0,
		},
	}

	var serviceName = "GunService"
	if config.ServiceName != "" {
		serviceName = config.ServiceName
	}

	return &Client{
		client: client,
		url: &url.URL{
			Scheme: "https",
			Host:   config.RemoteAddr,
			Path:   fmt.Sprintf("/%s/Tun", serviceName),
		},
		headers: http.Header{
			"content-type": []string{"application/grpc"},
			"user-agent":   []string{"grpc-go/1.36.0"},
			"te":           []string{"trailers"},
		},
	}
}

type ChainedClosable []io.Closer

// Close implements io.Closer.Close().
func (cc ChainedClosable) Close() error {
	for _, c := range cc {
		_ = c.Close()
	}
	return nil
}

func (cli *Client) DialConn() (net.Conn, error) {
	reader, writer := io.Pipe()
	request := &http.Request{
		Method:     http.MethodPost,
		Body:       reader,
		URL:        cli.url,
		Proto:      "HTTP/2",
		ProtoMajor: 2,
		ProtoMinor: 0,
		Header:     cli.headers,
	}
	anotherReader, anotherWriter := io.Pipe()
	go func() {
		defer anotherWriter.Close()
		response, err := cli.client.Do(request)
		if err != nil {
			return
		}
		_, _ = io.Copy(anotherWriter, response.Body)
	}()

	return newGunConn(anotherReader, writer, ChainedClosable{reader, writer, anotherReader}, nil, nil), nil
}

var (
	ErrInvalidLength = errors.New("invalid length")
)

func newGunConn(reader io.Reader, writer io.Writer, closer io.Closer, local net.Addr, remote net.Addr) *GunConn {
	if local == nil {
		local = &net.TCPAddr{
			IP:   []byte{0, 0, 0, 0},
			Port: 0,
		}
	}
	if remote == nil {
		remote = &net.TCPAddr{
			IP:   []byte{0, 0, 0, 0},
			Port: 0,
		}
	}
	return &GunConn{
		reader: reader,
		writer: writer,
		closer: closer,
		local:  local,
		remote: remote,
		done:   make(chan struct{}),
	}
}

func (g *GunConn) isClosed() bool {
	select {
	case <-g.done:
		return true
	default:
		return false
	}
}

func (g *GunConn) Read(b []byte) (n int, err error) {
	if g.toRead != nil {
		n = copy(b, g.toRead[g.readAt:])
		g.readAt += n
		if g.readAt >= len(g.toRead) {
			g.toRead = nil
		}
		return n, nil
	}
	buf := make([]byte, 5)
	n, err = io.ReadFull(g.reader, buf)
	if err != nil {
		return 0, err
	}
	//log.Printf("GRPC Header: %x", buf[:n])
	grpcPayloadLen := binary.BigEndian.Uint32(buf[1:])
	//log.Printf("GRPC Payload Length: %d", grpcPayloadLen)

	buf = make([]byte, grpcPayloadLen)
	n, err = io.ReadFull(g.reader, buf)
	if err != nil {
		return 0, io.ErrUnexpectedEOF
	}
	protobufPayloadLen, protobufLengthLen := DecodeUleb128(buf[1:])
	//log.Printf("Protobuf Payload Length: %d, Length Len: %d", protobufPayloadLen, protobufLengthLen)
	if protobufLengthLen == 0 {
		return 0, ErrInvalidLength
	}
	if grpcPayloadLen != uint32(protobufPayloadLen)+uint32(protobufLengthLen)+1 {
		return 0, ErrInvalidLength
	}
	n = copy(b, buf[1+protobufLengthLen:])
	if n < int(protobufPayloadLen) {
		g.toRead = buf
		g.readAt = 1 + int(protobufLengthLen) + n
	}
	return n, nil
}

func (g *GunConn) Write(b []byte) (n int, err error) {
	if g.isClosed() {
		return 0, io.ErrClosedPipe
	}
	protobufHeader := AppendUleb128([]byte{0x0A}, uint64(len(b)))
	grpcHeader := make([]byte, 5)
	grpcPayloadLen := uint32(len(protobufHeader) + len(b))
	binary.BigEndian.PutUint32(grpcHeader[1:5], grpcPayloadLen)
	_, err = io.Copy(g.writer, io.MultiReader(bytes.NewReader(grpcHeader), bytes.NewReader(protobufHeader), bytes.NewReader(b)))
	if f, ok := g.writer.(http.Flusher); ok {
		f.Flush()
	}
	return len(b), err
}

func (g *GunConn) Close() error {
	g.mu.Lock()
	defer g.mu.Unlock()
	select {
	case <-g.done:
		return nil
	default:
		close(g.done)
		return g.closer.Close()
	}
}

func (g *GunConn) LocalAddr() net.Addr {
	return g.local
}

func (g *GunConn) RemoteAddr() net.Addr {
	return g.remote
}

func (g *GunConn) SetDeadline(t time.Time) error {
	return nil
}

func (g *GunConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (g *GunConn) SetWriteDeadline(t time.Time) error {
	return nil
}
