package proxys

import (
	"CentralizedControl/common/fastjson"
	"CentralizedControl/common/http_helper"
	"CentralizedControl/common/log"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"golang.org/x/net/http2"
	"net"
	"net/http"
	"time"
)

type Dialer interface {
	Dial(network, address string) (net.Conn, error)
}

type IpTestInfo struct {
	Ip      string
	Country string
	Isp     string
	Org     string
	As      string
	Mobile  string
	Proxy   string
}

type Proxy interface {
	GetDialer() Dialer
	GetProxy() func(*http.Client)
	GetName() string
	GetId() string
	GetType() string
	GetTypeInt() int
	Test() *IpTestInfo
	ToJson() string
	GetSsr() *SsrProxy
	GetData() any
}

type ProxyInfo struct {
	Type      string
	Name      string
	IpInfo    *IpTestInfo
	data      any
	LiveTime  time.Duration
	StartTime time.Time
	Proxy
}

const (
	ProxySocks = 0
	ProxyHttp  = 1
	ProxySs    = 2
	ProxySsr   = 3
	ProxyVmess = 4
)

var ProxyTypeMap = map[string]int{
	"socks": ProxySocks,
	"http":  ProxyHttp,
	"ss":    ProxySs,
	"ssr":   ProxySsr,
	"vmess": ProxyVmess,
}

func (this *ProxyInfo) GetTypeInt() int {
	return ProxyTypeMap[this.GetType()]
}

func (this *ProxyInfo) GetData() any {
	return this.data
}

func (this *ProxyInfo) GetId() string {
	marshal, err := json.Marshal(this.data)
	if err != nil {
	}
	return fmt.Sprintf("%x", md5.Sum(marshal))
}

func (this *ProxyInfo) GetSsr() *SsrProxy {
	return this.data.(*SsrProxy)
}

func (this *ProxyInfo) GetName() string {
	return this.Name
}

func (this *ProxyInfo) GetProxy() func(*http.Client) {
	return func(c *http.Client) {
		tr := http_helper.GetTransport(c)
		tr.Dial = this.GetDialer().Dial
		err := http2.ConfigureTransport(tr)
		if err != nil {
			fmt.Printf("new client %v \n", err)
		}
	}
}

// {"ip":"61.141.64.234","country":"CN","asn":{"asnum":4134,"org_name":"Chinanet"},"geo":{"city":"Huizhou","region":"GD","region_name":"Guangdong","postal_code":"","latitude":23.1072,"longitude":114.3997,"tz":"Asia/Shanghai","lum_city":"huizhou","lum_region":"gd"}}

func (this *ProxyInfo) Test() *IpTestInfo {
	defer func() {
		if r := recover(); r != nil {
			log.Error("ProxyInfo.Test recover error: %v", r)
			this.IpInfo = nil
		}
	}()

	c := http_helper.CreateGoHttpClient(this.GetProxy(),
		http_helper.HttpTimeout(10),
		http_helper.EnableHttp2(),
		http_helper.DisableRedirect())

	resp, err := http_helper.HttpDo(c, &http_helper.RequestOpt{
		ReqUrl: "http://lumtest.com/myip.json",
	})
	if err != nil {
		log.Error("test proxy %s error: %v", this.Name, err)
		return nil
	}
	parse, err := fastjson.Parse(resp)
	if err != nil {
		log.Error("test proxy %s error: %v", this.Name, err)
		return nil
	}
	this.IpInfo = &IpTestInfo{
		Ip:      parse.GetString("ip"),
		Country: parse.GetString("country"),
		Isp:     "",
		Org:     parse.Get("asn").GetString("org_name"),
		As:      "",
		Mobile:  "",
		Proxy:   "",
	}
	return this.IpInfo
}

func (this *ProxyInfo) GetType() string {
	return this.Type
}

func (this *ProxyInfo) ToJson() string {
	marshal, err := json.Marshal(this.data)
	if err != nil {
		return ""
	}
	return string(marshal)
}

func createProxyInfo(_type string, name string, data any, impl any) *ProxyInfo {
	return &ProxyInfo{
		Name:      name,
		data:      data,
		Type:      _type,
		StartTime: time.Now(),
		Proxy:     impl.(Proxy),
	}
}
