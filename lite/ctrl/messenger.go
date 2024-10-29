package ctrl

import (
	"CentralizedControl/common/log"
	"CentralizedControl/common/proxys"
	"CentralizedControl/messenger"
)

type MessengerManager struct {
	proxys      *proxys.ProxyManage
	triggerWork chan *messenger.Cookies
}

func CreateMessengerManager(proxys *proxys.ProxyManage) *MessengerManager {
	msgMgr := &MessengerManager{
		proxys:      proxys,
		triggerWork: make(chan *messenger.Cookies, 10),
	}
	go msgMgr.Works()
	return msgMgr
}

func (this *MessengerManager) OnCreateAccount(cookies *messenger.Cookies) {
	select {
	case this.triggerWork <- cookies:
	default:
	}
}

func (this *MessengerManager) work(cookies *messenger.Cookies) {
	msg := messenger.CreateMessenger(cookies)
	p := this.proxys.GetProxyFromId(cookies.VpnId)
	if p == nil {
		log.Warn("not find proxy id: %s", cookies.VpnId)
		p = this.proxys.GetSubscribeProxy()
	}
	//msg.SetProxy(p)
	twoFA, err := msg.Set2FA()
	if err != nil {
		log.Error("messenger account %s set 2fa error: %v", cookies.Email)
		return
	}
	cookies.TwoFAKey = twoFA.KeyText
	cookies.QrCodeUri = twoFA.QrCodeUri
	cookies.Had2FA = "1"
	UpdateMessenger2FA(cookies)
}

func (this *MessengerManager) Works() {
	for cookies := range this.triggerWork {
		this.work(cookies)
	}
}
