package proxys

import (
	"CentralizedControl/common/http_helper"
	"CentralizedControl/common/log"
	"fmt"
	"net/http"
)

var RotatingResi = 0
var PhoneNet = 1
var Datacenter = 2

type RolaProxy struct {
	Type            int
	RotatingResiUrl string
	PhoneNetUrl     string
	DatacenterUrl   string
	client          *http.Client
}

func CreateRolaPool(ipType int) *RolaProxy {
	return &RolaProxy{
		Type:            ipType,
		RotatingResiUrl: "http://list.rola.info:8088/user_get_ip_list?token=kU9b2G8cJyhoVBXE1684123961109&qty=1&country=%s&state=&city=&time=10&format=json&protocol=socks5&filter=1",
		PhoneNetUrl:     "http://list.rola.info:8088/user_get_ip_list?token=kU9b2G8cJyhoVBXE1684123961109&type=4g&qty=1&country=%s&time=10&format=json&protocol=socks5&filter=1",
		DatacenterUrl:   "http://list.rola.info:8088/user_get_ip_list?token=kU9b2G8cJyhoVBXE1684123961109&type=datacenter&qty=1&time=10&country=%s&format=json&protocol=socks5&filter=1",
		client:          &http.Client{},
	}
}

type RolaResp struct {
	Code int      `json:"code"`
	Data []string `json:"data"`
	Msg  string   `json:"msg"`
}

func (this *RolaProxy) GetProxy(region string, asn string) Proxy {
	var url string
	switch this.Type {
	case RotatingResi:
		url = this.RotatingResiUrl
	case PhoneNet:
		url = this.PhoneNetUrl
	case Datacenter:
		url = this.DatacenterUrl
	}
	url = fmt.Sprintf(url, region)
	var resp RolaResp
	err := http_helper.HttpDoJson(this.client, &http_helper.RequestOpt{
		IsPost: false,
		ReqUrl: url,
	}, &resp)
	if err != nil {
		log.Error("rola get proxy error: %v", err)
		return nil
	}
	if len(resp.Data) != 1 {
		log.Error("rola get proxy error: data is null")
		return nil
	}
	proxy, err := CreateSocks5Proxy("socks://" + resp.Data[0])
	if err != nil {
		log.Error("rola CreateSocks5Proxy error: %v", err)
		return nil
	}
	return proxy
}
