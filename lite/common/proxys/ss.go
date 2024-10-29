package proxys

import (
	"CentralizedControl/common/log"
	"github.com/mzz2017/gg/dialer/shadowsocks"
)

type SsProxy struct {
	*ProxyInfo
	*shadowsocks.Shadowsocks
}
type SsProxyInfo struct {
	ExperimentReducedIvHeadEntropy bool   `json:"experimentReducedIvHeadEntropy"`
	Method                         string `json:"method"`
	Password                       string `json:"password"`
	Plugin                         string `json:"plugin"`
	ExtraType                      int    `json:"extraType"`
	Group                          string `json:"group"`
	Name                           string `json:"name"`
	ProfileId                      string `json:"profileId"`
	ServerAddress                  string `json:"serverAddress"`
	ServerPort                     int    `json:"serverPort"`
}

func Conv2SsProxyInfo(ss *shadowsocks.Shadowsocks) *SsProxyInfo {
	return &SsProxyInfo{
		ExperimentReducedIvHeadEntropy: false,
		Method:                         ss.Cipher,
		Password:                       ss.Password,
		Plugin:                         "",
		ExtraType:                      0,
		Group:                          "",
		Name:                           ss.Name,
		ProfileId:                      "",
		ServerAddress:                  ss.Server,
		ServerPort:                     ss.Port,
	}
}

func CreateSsClient(link string) (client *SsProxy, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Error("CreateSsClient recover error: %v", err)
			client = nil
			err = r.(error)
		}
	}()
	client = &SsProxy{}
	client.Shadowsocks, err = shadowsocks.ParseSSURL(link)
	if err != nil {
		return nil, err
	}
	client.ProxyInfo = createProxyInfo("ss", client.Shadowsocks.Name, Conv2SsProxyInfo(client.Shadowsocks), client)
	return client, nil
}

func (this *SsProxy) GetDialer() Dialer {
	dialer, err := this.Shadowsocks.Dialer()
	if err != nil {
		log.Warn("SsProxy GetDialer error: %v", err)
	}
	return dialer
}
