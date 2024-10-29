package ins_lite

import (
	"CentralizedControl/common"
	"container/list"
	"fmt"
	"sync"
	"time"
)

type EventTypes interface {
	string | uint64
}

type Event[EventType EventTypes] struct {
	EventName EventType
	Obj       any
	event     *common.Event
}

func WaitForTime[EventType EventTypes](event []*Event[EventType], timeout time.Duration) (bool, int) {
	es := make([]*common.Event, len(event))
	for idx := range event {
		es[idx] = event[idx].event
	}
	return common.SelectEvent(es, timeout)
}

func (this *Event[EventType]) Wait() {
	this.event.Wait()
}

func (this *Event[EventType]) Wait30Second() bool {
	return this.WaitForTime(time.Second * 30)
}

func (this *Event[EventType]) MustWait30Second() {
	if !this.Wait30Second() {
		panic(fmt.Sprintf("Event wait after 30s, event: %v", this.EventName))
	}
}

func (this *Event[EventType]) WaitForTime(timeout time.Duration) bool {
	return this.event.WaitForTime(timeout)
}

func (this *Event[EventType]) Target() {
	this.event.Signal()
}

type EventList[EventType EventTypes] struct {
	waitEvent     *list.List
	waitEventLock sync.Mutex
}

func CreateWaitEvent[EventType EventTypes]() *EventList[EventType] {
	return &EventList[EventType]{
		waitEvent:     list.New(),
		waitEventLock: sync.Mutex{},
	}
}

func (this *EventList[EventType]) TargetEvent(name EventType) {
	this.waitEventLock.Lock()
	for e := this.waitEvent.Front(); e != nil; {
		n := e.Next()
		if e.Value.(*Event[EventType]).EventName == name {
			e.Value.(*Event[EventType]).Target()
			this.waitEvent.Remove(e)
		}
		e = n
	}
	this.waitEventLock.Unlock()
}

func (this *EventList[EventType]) ReleaseEvent(event *Event[EventType]) {
	this.waitEventLock.Lock()
	for e := this.waitEvent.Front(); e != nil; {
		n := e.Next()
		if e.Value.(*Event[EventType]) == event {
			this.waitEvent.Remove(e)
		}
		e = n
	}
	this.waitEventLock.Unlock()
}

func (this *EventList[EventType]) GetEvent(name EventType) *Event[EventType] {
	e := &Event[EventType]{
		EventName: name,
		event:     common.CreateEventWait(false),
	}
	this.waitEventLock.Lock()
	this.waitEvent.PushBack(e)
	this.waitEventLock.Unlock()
	return e
}
