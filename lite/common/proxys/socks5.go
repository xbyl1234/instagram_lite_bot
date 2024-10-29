package proxys

import (
	"CentralizedControl/common"
	"fmt"
	"golang.org/x/net/proxy"
)

type Socks5Proxy struct {
	*ProxyInfo
	*ProxyUrlInfo
}

func CreateSocks5Proxy(url string) (*Socks5Proxy, error) {
	config := parseFromUrl(url)
	if config == nil {
		return nil, common.NerError("parse error")
	}
	config.Protocol = 2
	info := &Socks5Proxy{
		ProxyUrlInfo: config,
	}
	info.ProxyInfo = createProxyInfo("socks5", "", config, info)
	return info, nil
}

func (this *Socks5Proxy) GetDialer() Dialer {
	var auth = &proxy.Auth{}
	if this.NeedAuth {
		auth.User = this.Username
		auth.Password = this.Passwd
	} else {
		auth = nil
	}
	dialer, _ := proxy.SOCKS5("tcp", this.Ip+":"+fmt.Sprintf("%d", this.Port), auth, proxy.Direct)
	return dialer
}
