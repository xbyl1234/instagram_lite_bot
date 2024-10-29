package obfs

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/v2rayA/shadowsocksR/ssr"
	"math/rand"
	"strings"
)

var (
	requestPath = []string{
		"", "",
		"login.php?redir=", "",
		"register.php?code=", "",
		"?keyword=", "",
		"search?src=typd&q=", "&lang=en",
		"s?ie=utf-8&f=8&rsv_bp=1&rsv_idx=1&ch=&bar=&wd=", "&rn=",
		"post.php?id=", "&goto=view.php",
	}
	requestUserAgent = []string{
		"Mozilla/5.0 (Windows NT 6.3; WOW64; rv:40.0) Gecko/20100101 Firefox/40.0",
		"Mozilla/5.0 (Windows NT 6.3; WOW64; rv:40.0) Gecko/20100101 Firefox/44.0",
		"Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/535.11 (KHTML, like Gecko) Ubuntu/11.10 Chromium/27.0.1453.93 Chrome/27.0.1453.93 Safari/537.36",
		"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:35.0) Gecko/20100101 Firefox/35.0",
		"Mozilla/5.0 (compatible; WOW64; MSIE 10.0; Windows NT 6.2)",
		"Mozilla/5.0 (Windows; U; Windows NT 6.1; en-US) AppleWebKit/533.20.25 (KHTML, like Gecko) Version/5.0.4 Safari/533.20.27",
		"Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 6.3; Trident/7.0; .NET4.0E; .NET4.0C)",
		"Mozilla/5.0 (Windows NT 6.3; Trident/7.0; rv:11.0) like Gecko",
		"Mozilla/5.0 (Linux; Android 4.4; Nexus 5 Build/BuildID) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/30.0.0.0 Mobile Safari/537.36",
		"Mozilla/5.0 (iPad; CPU OS 5_0 like Mac OS X) AppleWebKit/534.46 (KHTML, like Gecko) Version/5.1 Mobile/9A334 Safari/7534.48.3",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 5_0 like Mac OS X) AppleWebKit/534.46 (KHTML, like Gecko) Version/5.1 Mobile/9A334 Safari/7534.48.3",
	}
)

// HttpSimple http_simple obfs encapsulate
type httpSimplePost struct {
	ssr.ServerInfo
	rawTransSent     bool
	rawTransReceived bool
	userAgentIndex   int
	methodGet        bool // true for get, false for post
}

func init() {
	register("http_simple", newHttpSimple)
}

// newHttpSimple create a http_simple object
func newHttpSimple() IObfs {

	t := &httpSimplePost{
		rawTransSent:     false,
		rawTransReceived: false,
		userAgentIndex:   rand.Intn(len(requestUserAgent)),
		methodGet:        true,
	}
	return t
}

func (t *httpSimplePost) SetServerInfo(s *ssr.ServerInfo) {
	t.ServerInfo = *s
}

func (t *httpSimplePost) GetServerInfo() (s *ssr.ServerInfo) {
	return &t.ServerInfo
}

func (t *httpSimplePost) SetData(data interface{}) {

}

func (t *httpSimplePost) GetData() interface{} {
	return nil
}

func (t *httpSimplePost) boundary() (ret string) {

	set := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	for i := 0; i < 32; i++ {
		ret = fmt.Sprintf("%s%c", ret, set[rand.Intn(len(set))])
	}
	return
}

func (t *httpSimplePost) data2URLEncode(data []byte) (ret string) {
	for i := 0; i < len(data); i++ {
		ret = fmt.Sprintf("%s%%%s", ret, hex.EncodeToString([]byte{data[i]}))
	}
	return
}

func (t *httpSimplePost) Encode(data []byte) (encodedData []byte, err error) {
	if t.rawTransSent {
		return data, nil
	}

	dataLength := len(data)
	var headData []byte
	if headSize := t.IVLen + t.HeadLen; dataLength-headSize > 64 {
		headData = make([]byte, headSize+rand.Intn(64))
	} else {
		headData = make([]byte, dataLength)
	}
	copy(headData, data[0:len(headData)])
	requestPathIndex := rand.Intn(len(requestPath)/2) * 2
	host := t.Host
	var customHead string

	if len(t.Param) > 0 {
		customHeads := strings.Split(t.Param, "#")
		if len(customHeads) > 2 {
			customHeads = customHeads[0:2]
		}
		param := t.Param
		if len(customHeads) > 1 {
			customHead = customHeads[1]
			param = customHeads[0]
		}
		hosts := strings.Split(param, ",")
		if len(hosts) > 0 {
			host = strings.TrimSpace(hosts[rand.Intn(len(hosts))])
		}
	}
	method := "GET /"
	if !t.methodGet {
		method = "POST /"
	}
	httpBuf := fmt.Sprintf("%s%s%s%s HTTP/1.1\r\nHost: %s:%d\r\n",
		method,
		requestPath[requestPathIndex],
		t.data2URLEncode(headData),
		requestPath[requestPathIndex+1],
		host,
		t.Port)
	if len(customHead) > 0 {
		httpBuf = httpBuf + strings.Replace(customHead, "\\n", "\r\n", -1) + "\r\n\r\n"
	} else {
		var contentType string
		if !t.methodGet {
			contentType = "Content-Type: multipart/form-data; boundary=" + t.boundary() + "\r\n"
		}
		httpBuf = httpBuf +
			"User-Agent: " + requestUserAgent[t.userAgentIndex] + "\r\n" +
			"Accept: text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8\r\n" +
			"Accept-Language: en-US,en;q=0.8\r\n" +
			"Accept-Encoding: gzip, deflate\r\n" +
			contentType +
			"DNT: 1\r\n" +
			"Connection: keep-alive\r\n" +
			"\r\n"
	}

	if len(headData) < dataLength {
		encodedData = make([]byte, len(httpBuf)+(dataLength-len(headData)))
		copy(encodedData, []byte(httpBuf))
		copy(encodedData[len(httpBuf):], data[len(headData):])
	} else {
		encodedData = []byte(httpBuf)
	}
	t.rawTransSent = true

	return
}

func (t *httpSimplePost) Decode(data []byte) (decodedData []byte, needSendBack bool, err error) {
	if t.rawTransReceived {
		return data, false, nil
	}

	pos := bytes.Index(data, []byte("\r\n\r\n"))
	if pos > 0 {
		decodedData = make([]byte, len(data)-pos-4)
		copy(decodedData, data[pos+4:])
		t.rawTransReceived = true
	}
	return decodedData, false, nil
}

func (t *httpSimplePost) GetOverhead() int {
	return 0
}
