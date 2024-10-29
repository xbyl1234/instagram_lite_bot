package ins_lite

import (
	"CentralizedControl/common/log"
	"CentralizedControl/ins_lite/net"
	"CentralizedControl/ins_lite/proto"
	"CentralizedControl/ins_lite/proto/types"
	"encoding/hex"
	"encoding/json"
	"sync"
)

type EventCallBack = func(msg *net.MsgWrap) error

type MsgRecvEvent struct {
	eventMap           map[uint64]EventCallBack
	eventLock          sync.Mutex
	defaultMsgDealFunc EventCallBack
	eventList          *EventList[uint64]
}

func CreateMsgRecvEvent(defaultMsgDealFunc EventCallBack) *MsgRecvEvent {
	e := &MsgRecvEvent{
		defaultMsgDealFunc: defaultMsgDealFunc,
		eventMap:           map[uint64]EventCallBack{},
		eventLock:          sync.Mutex{},
		eventList:          CreateWaitEvent[uint64](),
	}
	return e
}

func (this *MsgRecvEvent) OnSocketConnect(msgEvent *net.MsgWrap) {
	return
}

func (this *MsgRecvEvent) logEvent(msgEvent *net.MsgWrap) {
	var logBody string
	var msgName string
	if msgEvent.Body != nil {
		marshal, _ := json.Marshal(msgEvent.Body.Interface())
		logBody = string(marshal)
		msgName = proto.GetMessageName(false, msgEvent.Code)
	} else {
		logBody = hex.EncodeToString(msgEvent.Reader.PeekRemain())
	}
	if msgEvent.Parent != nil {
		log.Info("recv sub msgCode: %d, msgName: %s, recvIdx: 0x%x, sendIdx: 0x%x, body: %s",
			msgEvent.Code&^proto.MuskSubCmd, msgName, msgEvent.RecvIdx, msgEvent.SendIdx, logBody)
	} else {
		log.Info("packet recv\t%s\t%s", msgName, logBody)
		//log.Info("recv msgCode: %d, msgName: %s, recvIdx: 0x%x, sendIdx: 0x%x, body: %s",
		//	msgEvent.Code, msgName, msgEvent.RecvIdx, msgEvent.SendIdx, logBody)
	}
}

func (this *MsgRecvEvent) parseMsgBody(msgEvent *net.MsgWrap) {
	msg := proto.CreateMsgByCode(false, msgEvent.Code)
	if msg != nil {
		types.ReadMsg(msgEvent.Reader, msg.Interface())
		msgEvent.Body = msg
	}
}

func (this *MsgRecvEvent) OnSocketRecv(msgEvent *net.MsgWrap) {
	//defer func() {
	//	if r := recover(); r != nil {
	//		log.Error("OnSocketRecv msg code: %d, error: %v, body: %s",
	//			msgEvent.Code, r, hex.EncodeToString(msgEvent.Reader.PeekRemain()))
	//	}
	//}()

	this.parseMsgBody(msgEvent)
	this.logEvent(msgEvent)
	this.callEvent(msgEvent)
	this.targetEvent(msgEvent)
}

func (this *MsgRecvEvent) PostSubEvent(parent *net.MsgWrap, subCode uint64) {
	//defer func() {
	//	if r := recover(); r != nil {
	//		log.Error("OnSocketRecv msg code: %d, sub code: %d, error: %v, body: %s",
	//			parent.Code, subCode&^proto.MuskSubCmd, r, hex.EncodeToString(parent.Reader.PeekRemain()))
	//	}
	//}()

	subMsg := &net.MsgWrap{
		Parent:  parent,
		Code:    subCode,
		Reader:  parent.Reader,
		RecvIdx: 0,
		SendIdx: 0,
	}
	this.parseMsgBody(subMsg)
	this.logEvent(subMsg)
	this.callEvent(subMsg)
	this.targetEvent(subMsg)
}

func (this *MsgRecvEvent) RegisterEvent(code uint64, callback EventCallBack) {
	this.eventLock.Lock()
	defer this.eventLock.Unlock()
	this.eventMap[code] = callback
}

func (this *MsgRecvEvent) GetWaitEvent(code uint64) *Event[uint64] {
	return this.eventList.GetEvent(code)
}

func (this *MsgRecvEvent) targetEvent(msg *net.MsgWrap) {
	this.eventList.TargetEvent(msg.Code)
}

func (this *MsgRecvEvent) WaitEvent(code uint64) {
	this.GetWaitEvent(code).Wait()
}

func (this *MsgRecvEvent) callEvent(msg *net.MsgWrap) {
	var err error
	callback, ok := this.eventMap[msg.Code]
	if ok {
		err = callback(msg)
	}
	if !ok {
		err = this.defaultMsgDealFunc(msg)
	}
	if err != nil {
		log.Error("deal recv msgCode: %d, error: %v", msg.Code, err)
	}
}
