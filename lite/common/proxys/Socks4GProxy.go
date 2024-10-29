package proxys

import (
	"CentralizedControl/common/http_helper"
	"CentralizedControl/common/log"
	"net/http"
)

type Socks4GProxy struct {
	resetLink string
	client    *http.Client
}

func CreateSocks4GProxy() *Socks4GProxy {
	return &Socks4GProxy{
		resetLink: "https://api.ltesocks.io/v2/port/reset/81bea6e68b81db16c18146efa31477a42a44177a9f721b2d3a7d292009b1458c",
		client:    &http.Client{},
	}
}

func (this *Socks4GProxy) GetProxy(region string, asn string) Proxy {
	_, err := http_helper.HttpDo(this.client, &http_helper.RequestOpt{
		Params: map[string]string{},
		IsPost: false,
		ReqUrl: this.resetLink,
	})
	if err != nil {
		log.Warn("reset s5 error: %v", err)
		return nil
	}
	proxy, err := CreateSocks5Proxy("socks5://sasdasf323:safasfa3@ap1.socks.expert:22755")
	if err != nil {
		log.Warn("create s5 error: %v", err)
		return nil
	}
	return proxy
}
