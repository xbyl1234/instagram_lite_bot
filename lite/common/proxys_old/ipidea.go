package proxys_old

import (
	"CentralizedControl/common"
	"CentralizedControl/common/http_helper"
	"CentralizedControl/common/log"
	"fmt"
	"net/http"
	url2 "net/url"
	"sync"
	"time"
)

type IdeaPool struct {
	ProxyImpl
	url         string
	proxyType   http_helper.ProxyType
	lastReqTime time.Time
	lock        sync.Mutex
	proxyList   []*common.Proxy
	proxyMask   []bool
	client      *http.Client
	Country     string
}

func InitIdeaPool(url string) (ProxyImpl, error) {
	var pool = &IdeaPool{}
	purl, err := url2.Parse(url)
	if err != nil {
		return nil, err
	}

	for key, value := range purl.Query() {
		if key == "protocol" {
			if value[0] == "socks5" {
				pool.proxyType = common.ProxySocket
			} else {
				pool.proxyType = common.ProxyHttp
			}
		} else if key == "regions" {
			pool.Country = value[0]
		}
	}

	pool.url = url
	pool.client = &http.Client{}
	//common.DebugHttpClient(pool.client)
	return pool, nil
}

type IdeaResp struct {
	Code      int    `json:"code"`
	Success   bool   `json:"success"`
	Msg       string `json:"msg"`
	RequestIp string `json:"request_ip"`
	Data      []struct {
		Ip   string `json:"ip"`
		Port int    `json:"port"`
	} `json:"data"`
}

func (this *IdeaPool) RequestProxy() bool {
	resp := &IdeaResp{}
	for true {
		err := http_helper.HttpDoJson(this.client, &http_helper.RequestOpt{
			ReqUrl: this.url,
			IsPost: false,
		}, resp)
		if err != nil {
			log.Error("idea proxy request error: %v", err)
			time.Sleep(1 * time.Second)
			continue
		}
		if !resp.Success {
			log.Error("idea proxy request error: %v", resp.Msg)
			time.Sleep(3 * time.Second)
			continue
		}

		break
	}
	if len(resp.Data) == 0 {
		log.Error("idea request proxy list is null!")
		return false
	}
	this.proxyList = make([]*common.Proxy, len(resp.Data))
	this.proxyMask = make([]bool, len(resp.Data))

	for index := range resp.Data {
		dp := &resp.Data[index]
		this.proxyList[index] = &common.Proxy{
			ID:        "idea",
			Ip:        dp.Ip,
			Port:      fmt.Sprintf("%d", dp.Port),
			Username:  "",
			Passwd:    "",
			Rip:       dp.Ip,
			ProxyType: this.proxyType,
			NeedAuth:  false,
			Country:   this.Country,
		}
		this.proxyMask[index] = true
	}

	log.Info("idea request proxy list success!")
	return true
}

func (this *IdeaPool) get() *common.Proxy {
	this.lock.Lock()
	defer this.lock.Unlock()
	index := 0
	find := false
	for index = range this.proxyMask {
		if this.proxyMask[index] {
			find = true
			break
		}
	}

	if find {
		this.proxyMask[index] = false
		return this.proxyList[index]
	}
	if this.RequestProxy() {
		this.proxyMask[0] = false
		return this.proxyList[0]
	} else {
		return nil
	}
}

func (this *IdeaPool) GetNoRisk(busy bool, used bool) *common.Proxy {
	return this.get()
}

func (this *IdeaPool) Get() *common.Proxy {
	return this.get()
}
