package ins_lite

import (
	"CentralizedControl/common/log"
	"CentralizedControl/ins_lite/proto/io"
	"CentralizedControl/ins_lite/proto/msg/recver"
	"sync"
)

type ScreenEventData struct {
	Screen  *recver.ScreenReceived
	Source  string
	CmdIdx  uint16
	CmdCode uint64
	CmdData []byte
	Reader  *io.Reader
}

type ScreenEventCallBack = func(event *ScreenEventData) error

type ScreenEvent struct {
	eventMap           map[uint64]ScreenEventCallBack
	eventLock          sync.Mutex
	defaultMsgDealFunc ScreenEventCallBack
}

func CreateScreenEvent(defaultDealFunc ScreenEventCallBack) *ScreenEvent {
	return &ScreenEvent{
		eventMap:           map[uint64]ScreenEventCallBack{},
		eventLock:          sync.Mutex{},
		defaultMsgDealFunc: defaultDealFunc,
	}
}

func (this *ScreenEvent) RegisterEvent(code uint64, callback ScreenEventCallBack) {
	this.eventLock.Lock()
	defer this.eventLock.Unlock()
	this.eventMap[code] = callback
}

func (this *ScreenEvent) callEvent(msg *ScreenEventData) error {
	var err error
	callback, ok := this.eventMap[msg.CmdCode]
	if ok {
		err = callback(msg)
	}
	if !ok {
		err = this.defaultMsgDealFunc(msg)
	}
	if err != nil {
		log.Error("deal recv msgCode: %d, error: %v", msg.CmdCode, err)
	}
	return err
}
