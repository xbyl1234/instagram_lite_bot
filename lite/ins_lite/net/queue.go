package net

import (
	"container/list"
	"sync"
)

type ChanQueue[DataType any] struct {
	readChan chan DataType
	msgList  list.List
	lock     sync.Mutex
}

func (this *ChanQueue[DataType]) PushFront(data DataType) {
	this.lock.Lock()
	defer this.lock.Unlock()
	select {
	case this.readChan <- data:
	default:
		this.msgList.PushFront(data)
	}
}

func (this *ChanQueue[DataType]) Push(data DataType) {
	this.lock.Lock()
	defer this.lock.Unlock()
	select {
	case this.readChan <- data:
	default:
		this.msgList.PushBack(data)
	}
}

func (this *ChanQueue[DataType]) Get() chan DataType {
	this.lock.Lock()
	defer this.lock.Unlock()
	if this.msgList.Len() > 0 {
		e := this.msgList.Front()
		this.msgList.Remove(e)
		select {
		case this.readChan <- e.Value.(DataType):
		default:
		}
	}
	return this.readChan
}

func (this *ChanQueue[DataType]) Close() {
	close(this.readChan)
}

func CreateChanQueue[DataType any]() *ChanQueue[DataType] {
	return &ChanQueue[DataType]{
		readChan: make(chan DataType),
		msgList:  list.List{},
	}
}
