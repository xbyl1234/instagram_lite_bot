package proxys_old

//
//type ProxyImpl interface {
//	Get() *common.Proxy
//	Remove(proxy *common.Proxy)
//	Dumps()
//}
//
//type ProxyConfigt struct {
//	Providers []struct {
//		ProviderName string           `json:"imap_provider"`
//		Url          string           `json:"url"`
//		Country      string           `json:"country"`
//		ProxyType    common.ProxyType `json:"proxy_type"`
//	} `json:"providers"`
//}
//
//type ProxyPoolt struct {
//	proxys map[string]ProxyImpl
//}
//
//func (this *ProxyPoolt) Get(country string) *common.Proxy {
//	if common.UseCharles {
//		return common.BurpHttpProxy
//	}
//
//	if country == "" {
//		for key := range this.proxys {
//			country = key
//			break
//		}
//	}
//
//	p := this.proxys[country]
//	if p == nil {
//		return nil
//	}
//	for true {
//		proxy := p.Get()
//		if !proxy.IsOutLiveTime() {
//			return proxy
//		}
//	}
//	return nil
//}
//
//var ProxyPool ProxyPoolt
//var ProxyConfig ProxyConfigt
//
//func InitProxyPool(configPath string) error {
//	err := common.LoadJsonFile(configPath, &ProxyConfig)
//	if err != nil || len(ProxyConfig.Providers) == 0 {
//		log.Error("load proxy config error: %v", err)
//		return err
//	}
//
//	ProxyPool.proxys = make(map[string]ProxyImpl)
//	for _, imap_provider := range ProxyConfig.Providers {
//		var _proxy ProxyImpl
//		var err error
//
//		switch imap_provider.ProviderName {
//		case "dove":
//			_proxy, err = InitDovePool(imap_provider.Url)
//			break
//		//case "luminati":
//		//	_proxy, err = InitLuminatiPool(imap_provider.Url)
//		//	break
//		case "idea":
//			_proxy, err = InitIdeaPool(imap_provider.Url)
//			break
//		case "rola":
//			_proxy, err = InitRolaPool(imap_provider.Url)
//			break
//		default:
//			return &common.MakeMoneyError{ErrStr: fmt.Sprintf("proxy config imap_provider error: %s",
//				imap_provider.ProviderName), ErrType: common.OtherError}
//		}
//		if err != nil {
//			return &common.MakeMoneyError{ErrStr: fmt.Sprintf("proxy config imap_provider error: %s",
//				imap_provider.ProviderName), ErrType: common.OtherError}
//		}
//
//		ProxyPool.proxys[imap_provider.Country] = _proxy
//	}
//	return err
//}
