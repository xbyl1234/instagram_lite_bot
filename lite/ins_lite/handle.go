package ins_lite

import (
	"CentralizedControl/common/log"
	"CentralizedControl/common/utils"
	"CentralizedControl/ins_lite/net"
	"CentralizedControl/ins_lite/proto"
	"CentralizedControl/ins_lite/proto/msg/recver"
	"CentralizedControl/ins_lite/proto/msg/sender"
)

func (this *InsLiteClient) registerHandle() {
	this.RegisterEvent(proto.MsgCodeGetSystemPropertiesMsg, this.sendSystemPropertiesMsg)
	this.RegisterEvent(proto.MsgCodeSessionMessage, this.recvSessionMessage)
	this.RegisterEvent(proto.MsgCodeRecvSubCmd, this.recvSubCmd)
	this.RegisterEvent(proto.MsgCodeScreenReceived, this.recvScreenReceived)
	this.RegisterEvent(proto.MsgCodeRecvScreenDiff, this.recvScreenDiff)
	this.RegisterEvent(proto.MsgCodeRecvImage, this.recvImage)
	this.RegisterEvent(proto.MsgCodeTransactionID, this.recvTransactionID)
	this.RegisterEvent(proto.MsgCodeUpdateApp, this.recvUpdateApp)
	this.RegisterEvent(proto.MsgCodePropStore54, this.recvPropStore54)
	this.RegisterEvent(proto.MsgCodePropStore223, this.recvPropStore223)
	this.RegisterEvent(proto.MsgCodeHandleMessageOxygenAcceptTosOnSuccessfulLogin, this.onLoginSuccess)
	this.registerSubHandle()
	this.registerScreenHandle()
}

var IgnoreMsgCode = []int{
	proto.MsgCodePropStoreConfig,
	proto.MsgCodeRectangularBackgroundConfig16,
	proto.MsgCodePass40,
	proto.MsgCodeRectangularBackgroundConfig50,
	proto.MsgCodeConfig89,
	proto.MsgCodeAboutRecvQueueAction135,
	proto.MsgCodeBitFlag275,
}

func arrayHas(array []int, value int) bool {
	for _, item := range array {
		if item == value {
			return true
		}
	}
	return false
}

func (this *InsLiteClient) DefaultMsgDealFunc(msgEvent *net.MsgWrap) error {
	var level = log.LevelError
	if arrayHas(IgnoreMsgCode, int(msgEvent.Code)) {
		level = log.LevelWarn
	}
	if msgEvent.Parent != nil {
		log.Logs(level, "undisposed sub cmd code: %d", msgEvent.Code&^proto.MuskSubCmd)
	} else {
		log.Logs(level, "undisposed msg code: %d", msgEvent.Code)
	}
	return nil
}

func (this *InsLiteClient) recvSessionMessage(msgEvent *net.MsgWrap) error {
	body := msgEvent.Body.Interface().(*recver.SessionMessage)
	this.SetSession(body.SessionId, body.LzmaDictSize)
	this.Cookies.Session.ClientEncryptionSecret = body.ClientEncryptionSecret.Value
	this.Cookies.Session.SessionId = body.SessionId
	this.setTransientToken(body.TransientToken)
	return nil
}

func (this *InsLiteClient) sendSystemPropertiesMsg(msgEvent *net.MsgWrap) error {
	body := msgEvent.Body.Interface().(*recver.GetSystemPropertiesMsg)
	s := &proto.Message[sender.SendSystemPropertiesMsg]{}
	for _, item := range body.Properties.Value {
		v, ok := this.Cookies.Properties[item]
		if ok {
			s.Body.Properties.Put(sender.PropertiesItem{
				Name:     item,
				HasValue: 1,
				Value:    v,
			})
		} else {
			s.Body.Properties.Put(sender.PropertiesItem{
				Name:     item,
				HasValue: 0,
				Value:    "",
			})
		}
	}
	this.SendMsg(s)
	return nil
}

func (this *InsLiteClient) recvImage(msgEvent *net.MsgWrap) error {
	body := msgEvent.Body.Interface().(*recver.RecvImage)
	_ = body
	return nil
}

func (this *InsLiteClient) recvTransactionID(msgEvent *net.MsgWrap) error {
	body := msgEvent.Body.Interface().(*recver.TransactionID)
	_ = body
	return nil
}

func (this *InsLiteClient) recvUpdateApp(msgEvent *net.MsgWrap) error {
	return nil
}

func (this *InsLiteClient) SendConnBandwidthQuality() {
	s := &proto.Message[sender.ConnBandwidthQuality]{}
	s.Body.QualityType = byte(utils.ChoseOne(sender.QualityType))
	s.Body.Bandwidth = int32(utils.GenNumber(200, 2000))
	this.SendMsg(s)
}

func (this *InsLiteClient) recvPropStore54(msgEvent *net.MsgWrap) error {
	body := msgEvent.Body.Interface().(*recver.PropStore54)
	this.Cookies.PropStore54.Update(body)
	return nil
}

func (this *InsLiteClient) recvPropStore223(msgEvent *net.MsgWrap) error {
	body := msgEvent.Body.Interface().(*recver.PropStore223)
	this.Cookies.prop223Lock.Lock()
	defer this.Cookies.prop223Lock.Unlock()
	for idx := range body.Props.Kv {
		this.Cookies.Prop223[int(body.Props.Kv[idx].Key)] = body.Props.Kv[idx].Value
	}
	return nil
}

func (this *InsLiteClient) onLoginSuccess(msgEvent *net.MsgWrap) error {
	this.Cookies.State = StateLoggedIn
	this.SendAdvertiserId()
	this.SendAdvertiserId()
	this.SendAdvertiserId()
	this.SendAdvertiserId()
	return nil
}
