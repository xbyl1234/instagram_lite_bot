package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"testing"
	"time"

	tls "github.com/refraction-networking/utls"
	"golang.org/x/net/http2"
)

var (
	dialTimeout   = time.Duration(15) * time.Second
	sessionTicket = []uint8(`Here goes phony session ticket: phony enough to get into ASCII range
Ticket could be of any length, but for camouflage purposes it's better to use uniformly random contents
and common length. See https://tlsfingerprint.io/session-tickets`)
)

// var requestHostname = "tls.peet.ws" // speaks http2 and TLS 1.3
var requestHostname = "baidu.com" // speaks http2 and TLS 1.3
// var requestAddr = "205.185.123.167:443"
var requestAddr = "39.156.66.10:443"

func HttpGetCustom(hostname string, addr string) (*http.Response, error) {
	config := tls.Config{ServerName: hostname}
	dialConn, err := net.DialTimeout("tcp", addr, dialTimeout)
	if err != nil {
		return nil, fmt.Errorf("net.DialTimeout error: %+v", err)
	}
	uTlsConn := tls.UClient(dialConn, &config, tls.HelloCustom)
	defer uTlsConn.Close()

	//err = uTlsConn.ApplyPreset(
	//	GenTlsConfig("771,4865-4866-4867-49195-49196-52393-49199-49200-52392-49161-49162-49171-49172-156-157-47-53,0-23-65281-10-11-35-5-13-51-45-43-21,29-23-24,0",
	//		&InsLiteExtensionConfig))

	if err != nil {
		return nil, fmt.Errorf("uTlsConn.Handshake() error: %+v", err)
	}

	err = uTlsConn.Handshake()
	if err != nil {
		return nil, fmt.Errorf("uTlsConn.Handshake() error: %+v", err)
	}

	return httpGetOverConn(uTlsConn, uTlsConn.ConnectionState().NegotiatedProtocol)
}

var roller *tls.Roller

// this example creates a new roller for each function call,
// however it is advised to reuse the Roller
func HttpGetGoogleWithRoller() (*http.Response, error) {
	var err error
	if roller == nil {
		roller, err = tls.NewRoller()
		if err != nil {
			return nil, err
		}
	}

	// As of 2018-07-24 this tries to connect with Chrome, fails due to ChannelID extension
	// being selected by Google, but not supported by utls, and seamlessly moves on to either
	// Firefox or iOS fingerprints, which work.
	c, err := roller.Dial("tcp4", requestHostname+":443", requestHostname)
	if err != nil {
		return nil, err
	}

	return httpGetOverConn(c, c.ConnectionState().NegotiatedProtocol)
}

func TestTls(t *testing.T) {
	var response *http.Response
	var err error
	//
	//response, err = HttpGetDefault(requestHostname, requestAddr)
	//if err != nil {
	//	fmt.Printf("#> HttpGetDefault failed: %+v\n", err)
	//} else {
	//	fmt.Printf("#> HttpGetDefault response: %+s\n", dumpResponseNoBody(response))
	//}

	//response, err = HttpGetByHelloID(requestHostname, requestAddr, tls.HelloChrome_62)
	//if err != nil {
	//	fmt.Printf("#> HttpGetByHelloID(HelloChrome_62) failed: %+v\n", err)
	//} else {
	//	fmt.Printf("#> HttpGetByHelloID(HelloChrome_62) response: %+s\n", dumpResponseNoBody(response))
	//}
	//
	//response, err = HttpGetConsistentRandomized(requestHostname, requestAddr)
	//if err != nil {
	//	fmt.Printf("#> HttpGetConsistentRandomized() failed: %+v\n", err)
	//} else {
	//	fmt.Printf("#> HttpGetConsistentRandomized() response: %+s\n", dumpResponseNoBody(response))
	//}
	//
	//response, err = HttpGetExplicitRandom(requestHostname, requestAddr)
	//if err != nil {
	//	fmt.Printf("#> HttpGetExplicitRandom failed: %+v\n", err)
	//} else {
	//	fmt.Printf("#> HttpGetExplicitRandom response: %+s\n", dumpResponseNoBody(response))
	//}
	//
	//response, err = HttpGetTicket(requestHostname, requestAddr)
	//if err != nil {
	//	fmt.Printf("#> HttpGetTicket failed: %+v\n", err)
	//} else {
	//	fmt.Printf("#> HttpGetTicket response: %+s\n", dumpResponseNoBody(response))
	//}
	//
	//response, err = HttpGetTicketHelloID(requestHostname, requestAddr, tls.HelloFirefox_56)
	//if err != nil {
	//	fmt.Printf("#> HttpGetTicketHelloID(HelloFirefox_56) failed: %+v\n", err)
	//} else {
	//	fmt.Printf("#> HttpGetTicketHelloID(HelloFirefox_56) response: %+s\n", dumpResponseNoBody(response))
	//}
	//
	response, err = HttpGetCustom(requestHostname, requestAddr)
	if err != nil {
		fmt.Printf("#> HttpGetCustom() failed: %+v\n", err)
	} else {
		fmt.Printf("#> HttpGetCustom() response: %+s\n", dumpResponseNoBody(response))
	}
	////
	//for i := 0; i < 5; i++ {
	//	response, err = HttpGetGoogleWithRoller()
	//	if err != nil {
	//		fmt.Printf("#> HttpGetGoogleWithRoller() #%v failed: %+v\n", i, err)
	//	} else {
	//		fmt.Printf("#> HttpGetGoogleWithRoller() #%v response: %+s\n",
	//			i, dumpResponseNoBody(response))
	//	}
	//}
	//
	//forgeConn()

	return
}

func httpGetOverConn(conn net.Conn, alpn string) (*http.Response, error) {
	//u, _ := url.Parse("https://tls.peet.ws/api/tls")
	u, _ := url.Parse("https://tls.peet.ws/")
	req := &http.Request{
		Method: "GET",
		URL:    u,
		Header: make(http.Header),
		Host:   requestHostname,
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36")

	switch alpn {
	case "h2":
		req.Proto = "HTTP/2.0"
		req.ProtoMajor = 2
		req.ProtoMinor = 0

		tr := http2.Transport{}
		cConn, err := tr.NewClientConn(conn)
		if err != nil {
			return nil, err
		}
		return cConn.RoundTrip(req)
	case "http/1.1", "":
		req.Proto = "HTTP/1.1"
		req.ProtoMajor = 1
		req.ProtoMinor = 1

		err := req.Write(conn)
		if err != nil {
			return nil, err
		}
		return http.ReadResponse(bufio.NewReader(conn), req)
	default:
		return nil, fmt.Errorf("unsupported ALPN: %v", alpn)
	}
}

func dumpResponseNoBody(response *http.Response) string {
	resp, err := httputil.DumpResponse(response, true)
	if err != nil {
		return fmt.Sprintf("failed to dump response: %v", err)
	}
	return string(resp)
}
