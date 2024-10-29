package temp

import (
	"CentralizedControl/common/email_server/base"
	"CentralizedControl/common/email_server/config"
	"CentralizedControl/common/http_helper"
	"CentralizedControl/common/log"
	"CentralizedControl/common/proxys"
	"errors"
	"net/http"
	"time"
)

type WebTempEmailData struct {
	base.EmailClientData
	Client http.Client
}

func (this *WebTempEmailData) SetProxy(p proxys.Proxy) {
	this.Proxy = p
	http_helper.EnableHttp2()(&this.Client)
	p.GetProxy()(&this.Client)
}

func (this *WebTempEmailData) Close() {
	this.Client.CloseIdleConnections()
}

func (this *WebTempEmailData) WaitForEmail(self base.EmailClient, from string, callback base.ExtractEmailCallback) (string, error) {
	this.Lock.Lock()
	defer this.Lock.Unlock()
	start := time.Now()
	for time.Since(start) < config.DefaultRetryConfig.RetryTimeoutDuration {
		var emails []base.RespEmail
		var err error
		emails, err = self.GetEmail(from)
		if emails != nil {
			for _, item := range emails {
				var code string
				code, err = callback(self, &item)
				if err == nil && code != "" {
					return code, nil
				}
			}
			return "", errors.New("no one accept")
		}
		log.Warn("wait for %s code...error: %v", this.EmailAddr, err)
		time.Sleep(config.DefaultRetryConfig.RetryDelayDuration)
	}
	return "", errors.New("require code timeout")
}

type ProviderWebTempEmail struct {
	_type string
	Proxy proxys.Proxy
}

func (this *ProviderWebTempEmail) GetType() string {
	return this._type
}

func (this *ProviderWebTempEmail) GetEmail() (base.EmailClient, error) {
	var result base.EmailClient
	switch this._type {
	case base.ProviderLinShiName:
		result = &LingshiMail{}
	case base.ProviderTempMailIo:
		result = &TempMailIo{}
	}
	if result == nil {
		return nil, errors.New("unknow provider")
	}
	result.SetProxy(this.Proxy)
	err := result.Init()
	return result, err
}

func CreateProviderWebTempEmail(name string, p proxys.Proxy) base.Provider {
	return &ProviderWebTempEmail{
		_type: name,
		Proxy: p,
	}
}
