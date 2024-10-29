package proxys

import (
	"CentralizedControl/common"
	"CentralizedControl/common/http_helper"
	"CentralizedControl/common/log"
	"CentralizedControl/common/utils"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type ProxyConfig struct {
	Subscribe []string `json:"subscribe"`
}

type ZennoProxy struct {
	minPort    int
	maxPort    int
	curPort    int
	cachePort  map[string]int
	httpClient *http.Client
}

type ProxyManage struct {
	config ProxyConfig
	lock   sync.Mutex
	//vpn
	subscribeProvider []*SubscribeProvider
	subscribeIdMap    map[string]Proxy //ip
	subscribesId      []string
	subscribesCurIdx  int
	waitUpdateFinish  *common.Event

	//luminati
	luminati *LuminatiProvider
	//zenno
	*ZennoProxy
	mango *MangoProxy
}

func (this *ProxyManage) LoadSubscribeConfig() error {
	err := utils.LoadJsonFile("./proxy.json", &this.config)
	if err != nil {
		return err
	}
	for _, link := range this.config.Subscribe {
		this.subscribeProvider = append(this.subscribeProvider, CreateSubscribeProvider(link))
	}
	return err
}

func (this *ProxyManage) LoadLuminatiConfig() error {
	this.luminati = CreateLuminatiProvider("ips-zone5.txt", "zone5_ips.csv")
	return nil
}

func (this *ProxyManage) GetProxyFromId(id string) Proxy {
	return this.subscribeIdMap[id]
}

func (this *ProxyManage) GetSubscribeProxy() Proxy {
	this.lock.Lock()
	defer this.lock.Unlock()
	if this.subscribeIdMap == nil {
		return nil
	}
	_proxy := this.subscribeIdMap[this.subscribesId[this.subscribesCurIdx]]
	this.subscribesCurIdx++
	if this.subscribesCurIdx >= len(this.subscribesId) {
		this.subscribesCurIdx = 0
	}
	return _proxy
}

func (this *ProxyManage) GetOxylabsProxy(country string) Proxy {
	proxy, err := CreateHttpProxy(fmt.Sprintf("http://customer-xxbyl-sessid-%s:sda32134csdA@dc.us-pr.oxylabs.io:30000",
		//proxy, err := CreateHttpProxy(fmt.Sprintf("http://customer-xxbyl-sessid-%s:sda32134csdA@dc.hk-pr.oxylabs.io:16000",
		utils.GenString(utils.CharSet_abc, 10)))
	if err != nil {
		return nil
	}
	return proxy
}

type RefreshResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (this *ProxyManage) GetZenooProxy(clientId string, country string) Proxy {
	var port int
	this.lock.Lock()
	var ok bool
	if port, ok = this.cachePort[clientId]; !ok {
		if this.curPort == 0 || this.curPort >= this.maxPort {
			this.curPort = this.minPort
		} else {
			this.curPort += 1
		}
		port = this.curPort
		this.cachePort[clientId] = port
	}
	this.lock.Unlock()
	proxy, err := CreateSocks5Proxy(fmt.Sprintf("socks://asd:123123@47.242.117.34:%d", port))
	if err != nil {
		return nil
	}
	go func() {
		var resp RefreshResponse
		err = http_helper.HttpDoJson(this.httpClient, &http_helper.RequestOpt{
			IsPost: false,
			ReqUrl: fmt.Sprintf("http://47.242.117.34:12569/sv5-server/refresh_proxy_chain/248/%d", port),
		}, &resp)
		if err != nil {
			return
		}
		if !resp.Success {
			log.Error("refresh_proxy_chain %d error: %s", port, resp.Message)
		}
	}()
	return proxy
}

func (this *ProxyManage) GetLuminatiProxy(country string) Proxy {
	return this.luminati.GetProxy(country)
}

func (this *ProxyManage) GetMangoProxy(region string, asn string) Proxy {
	return this.mango.GetProxy(region, asn)
}

func (this *ProxyManage) UpdateSubscribe() {
	log.Info("start update proxy")
	testThread := 20
	testChan := make(chan Proxy, testThread)

	waitUpdate := sync.WaitGroup{}
	waitUpdate.Add(len(this.subscribeProvider))
	waitTest := sync.WaitGroup{}
	waitTest.Add(testThread)

	for _, provider := range this.subscribeProvider {
		go func(provider *SubscribeProvider) {
			defer func() {
				waitUpdate.Done()
			}()
			update, err := provider.Update()
			if err != nil && (update == nil || len(update) <= 0) {
				log.Error("disable subscribe %s", provider.GetLink())
				return
			}
			log.Info("subscribe %s count: %d", provider.GetLink(), len(update))
			for _, _proxy := range update {
				testChan <- _proxy
			}
		}(provider)
	}

	lock := sync.Mutex{}
	//dupId := make(map[string]Proxy)
	dupIp := make(map[string]Proxy)

	checkAndSet := func(key string, value Proxy, _map map[string]Proxy) bool {
		lock.Lock()
		defer lock.Unlock()
		if _map[key] != nil {
			return false
		}
		_map[key] = value
		return true
	}

	for i := 0; i < testThread; i++ {
		go func() {
			defer func() {
				waitTest.Done()
			}()
			for _proxy := range testChan {
				//if !checkAndSet(_proxy.GetId(), _proxy, dupId) {
				//	log.Warn("dup proxy id, name: %s", _proxy.GetName())
				//	continue
				//}
				testInfo := _proxy.Test()
				if testInfo == nil {
					continue
				}
				if !checkAndSet(testInfo.Ip, _proxy, dupIp) {
					log.Warn("dup proxy %s ip: %s", _proxy.GetName(), testInfo.Ip)
					continue
				}
			}
		}()
	}

	waitUpdate.Wait()
	close(testChan)
	waitTest.Wait()

	var ids []string
	var idMap = make(map[string]Proxy)
	for _, v := range dupIp {
		id := v.GetId()
		idMap[id] = v
		ids = append(ids, id)
	}

	log.Info("update proxy finish! all count: %d", len(ids))
	this.lock.Lock()
	this.subscribeIdMap = idMap
	this.subscribesId = ids
	this.subscribesCurIdx = 0
	this.lock.Unlock()
}

func (this *ProxyManage) WaitUpdateSubscribe() {
	this.waitUpdateFinish.Wait()
}

func (this *ProxyManage) UpdateSubscribeTimed() {
	this.UpdateSubscribe()
	this.waitUpdateFinish.Signal()
	for true {
		select {
		case <-time.After(time.Hour * 2):
			this.UpdateSubscribe()
		}
	}
}

func CreateProxysManage() *ProxyManage {
	log.Info("proxy server start...")
	mag := &ProxyManage{
		waitUpdateFinish: common.CreateEventWait(false),
		ZennoProxy: &ZennoProxy{
			minPort:    24000,
			maxPort:    24030,
			cachePort:  map[string]int{},
			httpClient: http_helper.CreateGoHttpClient(),
		},
		mango: CreateMangoProxy(),
	}
	err := mag.LoadSubscribeConfig()
	if err != nil {
		panic(err)
	}
	err = mag.LoadLuminatiConfig()
	if err != nil {
		panic(err)
	}
	go mag.UpdateSubscribeTimed()
	return mag
}
