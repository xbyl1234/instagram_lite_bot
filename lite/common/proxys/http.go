package proxys

import (
	"CentralizedControl/common"
	"CentralizedControl/common/proxys/http_proxy"
	"net/url"
	"strconv"
	"strings"
)

type ProxyUrlInfo struct {
	Ip       string `json:"serverAddress"`
	Port     int    `json:"serverPort"`
	Username string `json:"username"`
	Passwd   string `json:"password"`
	NeedAuth bool   `json:"need_auth"`
	Url      string `json:"url"`
	url      *url.URL
	//for Matsuri
	MatsuriType                             string        `json:"type"`
	MatsuriAllowInsecure                    bool          `json:"allowInsecure"`
	MatsuriAlpn                             string        `json:"alpn"`
	MatsuriCertificates                     string        `json:"certificates"`
	MatsuriEarlyDataHeaderName              string        `json:"earlyDataHeaderName"`
	MatsuriGrpcServiceName                  string        `json:"grpcServiceName"`
	MatsuriHeaderType                       string        `json:"headerType"`
	MatsuriHost                             string        `json:"host"`
	MatsuriMKcpSeed                         string        `json:"mKcpSeed"`
	MatsuriPacketEncoding                   int           `json:"packetEncoding"`
	MatsuriPath                             string        `json:"path"`
	MatsuriPinnedPeerCertificateChainSha256 string        `json:"pinnedPeerCertificateChainSha256"`
	MatsuriQuicKey                          string        `json:"quicKey"`
	MatsuriQuicSecurity                     string        `json:"quicSecurity"`
	MatsuriSecurity                         string        `json:"security"`
	MatsuriSni                              string        `json:"sni"`
	MatsuriUtlsFingerprint                  string        `json:"utlsFingerprint"`
	MatsuriUuid                             string        `json:"uuid"`
	MatsuriWsMaxEarlyData                   int           `json:"wsMaxEarlyData"`
	MatsuriWsUseBrowserForwarder            bool          `json:"wsUseBrowserForwarder"`
	MatsuriExtraType                        int           `json:"extraType"`
	MatsuriGroup                            string        `json:"group"`
	MatsuriName                             string        `json:"name"`
	MatsuriProfileId                        string        `json:"profileId"`
	MatsuriTags                             []interface{} `json:"tags"`
	//for socks5 =2
	Protocol int `json:"protocol"`
}

type HttpProxy struct {
	*ProxyInfo
	*ProxyUrlInfo
}

func CreateHttpProxy(url string) (*HttpProxy, error) {
	config := parseFromUrl(url)
	if config == nil {
		return nil, common.NerError("parse error")
	}
	info := &HttpProxy{
		ProxyUrlInfo: config,
	}
	info.ProxyInfo = createProxyInfo("http", "", config, info)
	return info, nil
}

func parseFromUrl(u string) *ProxyUrlInfo {
	info := &ProxyUrlInfo{
		Url:         u,
		MatsuriType: "tcp",
	}
	_url, err := url.Parse(u)
	if err != nil {
		return nil
	}
	info.url = _url
	info.Username = _url.User.Username()
	info.Passwd, info.NeedAuth = _url.User.Password()
	sp := strings.Split(_url.Host, ":")
	if len(sp) == 2 {
		info.Ip = sp[0]
		info.Port, _ = strconv.Atoi(sp[1])
	} else {
		return nil
	}
	return info
}

func (this *HttpProxy) GetDialer() Dialer {
	if this.NeedAuth {
		return http_proxy.New(this.url, http_proxy.WithProxyAuth(http_proxy.AuthBasic(this.Username, this.Passwd)))
	} else {
		return http_proxy.New(this.url)
	}
}

//func (this *Proxy) IsOutLiveTime() bool {
//	if this.LiveTime == -1 {
//		return false
//	}
//	if time.Duration(float64(time.Since(this.StartTime))*0.5) > this.LiveTime {
//		return true
//	}
//	return false
//}
