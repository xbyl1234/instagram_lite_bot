package http_helper

import (
	"crypto/tls"
	"net/http"
	"net/http/cookiejar"
	"time"
)

type ProxyType int
type HttpConfigFun func(c *http.Client)

func GetTransport(c *http.Client) *http.Transport {
	if c.Transport == nil {
		c.Transport = &http.Transport{}
	}
	return c.Transport.(*http.Transport)
}

func EnableHttp2() HttpConfigFun {
	return func(c *http.Client) {
		tr := GetTransport(c)
		tr.ForceAttemptHTTP2 = true
		//err := http2.ConfigureTransport(tr)
		//if err != nil {
		//	fmt.Printf("EnableHttp2 %v \n", err)
		//}
	}
}

func DisableHttpSslPinng() HttpConfigFun {
	return func(c *http.Client) {
		tr := GetTransport(c)
		if tr.TLSClientConfig == nil {
			tr.TLSClientConfig = &tls.Config{}
		}
		tr.TLSClientConfig.InsecureSkipVerify = true
	}
}

func DisableRedirect() HttpConfigFun {
	return func(c *http.Client) {
		c.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}
}

func FacebookTls() HttpConfigFun {
	return func(c *http.Client) {
		tr := GetTransport(c)
		if tr.TLSClientConfig == nil {
			tr.TLSClientConfig = &tls.Config{}
		}
		//tr.TLSClientConfig.NextProtos = []string{"h2", "h2-fb", "http/1.1"}
		tr.TLSClientConfig.MinVersion = tls.VersionTLS13
		tr.TLSClientConfig.MaxVersion = tls.VersionTLS13
		//tr.TLSClientConfig.CipherSuites = []uint16{
		//	tls.TLS_AES_128_GCM_SHA256,
		//	tls.TLS_AES_256_GCM_SHA384,
		//	tls.TLS_CHACHA20_POLY1305_SHA256,
		//}
	}
}

func NeedJar() HttpConfigFun {
	return func(c *http.Client) {
		jar, _ := cookiejar.New(nil)
		c.Jar = jar
	}
}

func HttpTimeout(sec int) HttpConfigFun {
	return func(c *http.Client) {
		c.Timeout = time.Duration(sec) * time.Second
	}
}

func CreateGoHttpClient(httpConfigs ...HttpConfigFun) *http.Client {
	tr := &http.Transport{}

	httpClient := &http.Client{
		Transport: tr,
	}

	for _, config := range httpConfigs {
		config(httpClient)
	}
	return httpClient
}
