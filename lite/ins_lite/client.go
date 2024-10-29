package ins_lite

import (
	"CentralizedControl/common"
	"CentralizedControl/common/log"
	"CentralizedControl/common/proxys"
	"CentralizedControl/ins_lite/net"
	"CentralizedControl/ins_lite/proto"
	"CentralizedControl/ins_lite/proto/msg/sender"
	"time"
)

type InsLiteClient struct {
	*net.InsLiteSocket
	*MsgRecvEvent
	*ScreenManager
	*common.RandomSleep
	Cookies         *Cookies
	EnableReConnect bool
}

func (this *InsLiteClient) initProp() {
	this.Cookies.Properties["com.android.imei"] = this.Cookies.DeviceId
	//this.Cookies.properties["modularity_properties"]=
	//this.Cookies.properties["modularity.state."] =
	//this.Cookies.properties["com.android.ffdpid.issynced"]
}

func (this *InsLiteClient) getSharedPreferencesInt(key string, def int) int {
	return this.getSharedPreferences(key, def).(int)
}

func (this *InsLiteClient) getSharedPreferencesString(key string, def string) string {
	return this.getSharedPreferences(key, def).(string)
}

func (this *InsLiteClient) setTransientToken(v int32) {
	this.Cookies.Session.TransientToken = v
	this.Cookies.Session.StickinessTokenTimeStamp = time.Now().Unix()
}

func (this *InsLiteClient) getTransientToken() int32 {
	if this.Cookies.PropStore54.GetBool(3143, true) {
		if this.Cookies.Packages.Get("com.facebook.liteqa") != nil {
			return 0x5FC1A0C
		}
	}
	transientToken := this.Cookies.Session.TransientToken
	timeOut := this.Cookies.PropStore54.GetInt(397, int64(1*time.Hour))
	if this.Cookies.Session.TransientToken != 0 &&
		time.Since(time.Unix(this.Cookies.Session.StickinessTokenTimeStamp, 0)).Milliseconds() < timeOut {
		return transientToken
	}

	v := this.Cookies.PropStore54.GetInt(258, -1)
	if v == -1 {
		if transientToken != 0 {
			transientToken = transientToken / 100000 * 100000
		}
		transientToken += 0x1869F
	} else {
		transientToken = int32(0x1869E + ((v & 0x7FFFFFFF) / 100000 * 100000))
	}
	this.setTransientToken(transientToken)
	return transientToken
}

func (this *InsLiteClient) getSharedPreferences(key string, def any) any {
	return def
}

func (this *InsLiteClient) OnSocketLost(isSend bool, err any) {
	log.Error("connect is_send: %v, lost: %v", isSend, err)
	if !this.EnableReConnect {
		return
	}
	for {
		err1 := this.ReConnect(this.getTransientToken())
		if err1 == nil {
			log.Info("ReConnect success")
			break
		}
		log.Error("ReConnect error: %v", err1)
		time.Sleep(1 * time.Second)
	}
}

func (this *InsLiteClient) OnSocketConnected(isReConnect bool) {
	log.Info("OnSocketConnected: %v", isReConnect)
	if isReConnect {
		s := &proto.Message[sender.ReConnect]{}
		this.SendMsgFront(s)
	}
}

func (this *InsLiteClient) SaveCookies(path string) {
	SaveCookies(path, this.Cookies)
}

func CreateNewInsLiteClient(ck *Cookies, p proxys.Proxy) *InsLiteClient {
	client := &InsLiteClient{
		EnableReConnect: false,
		Cookies:         ck,
		RandomSleep:     common.CreateRandomSleep(common.DefaultRandomConfig),
	}

	client.ScreenManager = CreateScreenManager(client)

	event := CreateMsgRecvEvent(client.DefaultMsgDealFunc)
	client.MsgRecvEvent = event

	s, err := net.CreateInsLiteSocketClient(client, p)
	if err != nil {
		return nil
	}
	client.InsLiteSocket = s

	client.registerHandle()
	client.initProp()
	err = s.Start(client.getTransientToken())
	if err != nil {
		return nil
	}
	return client
}
