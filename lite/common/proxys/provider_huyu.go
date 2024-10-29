package proxys

import (
	"CentralizedControl/common/log"
	"CentralizedControl/common/utils"
	"fmt"
)

type HuyuProvider struct {
	Url string
}

func CreateHuyuProvider() *HuyuProvider {
	return &HuyuProvider{
		Url: "",
	}
}

func (this *HuyuProvider) GetProxy(region string, asn string) Proxy {
	proxy, err := CreateSocks5Proxy(fmt.Sprintf("socks5://a6a9a0f0331d1f9598335-zone-custom-region-us-session-%s-sessTime-1:3915d88603588c853bc87191b676b89be607d7da@p2.mangoproxy.com:2333", utils.GenString(utils.CharSet_abc, 6)))
	if err != nil {
		log.Error("rola CreateSocks5Proxy error: %v", err)
		return nil
	}
	return proxy
}
