package common

import (
	"reflect"
	"sync"
	"time"
)

type Event struct {
	event        chan byte
	isExpires    bool
	autoReusable bool
	lock         sync.Mutex
}

func SelectEvent(event []*Event, timeout time.Duration) (bool, int) {
	for _, e := range event {
		e.lock.Lock()
	}
	cases := make([]reflect.SelectCase, len(event)+1)
	for i, ch := range event {
		if ch.isExpires {
			for _, e := range event {
				e.lock.Lock()
			}
			return true, i
		}
		cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch.event)}
	}
	for _, e := range event {
		e.lock.Unlock()
	}
	cases[len(cases)-1] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(time.After(timeout))}
	chosen, _, _ := reflect.Select(cases)
	return chosen != len(cases)-1, chosen
}

func (this *Event) WaitForTime(timeout time.Duration) bool {
	var e chan byte
	this.lock.Lock()
	if this.isExpires {
		this.lock.Unlock()
		return true
	}
	e = this.event
	this.lock.Unlock()
	select {
	case <-e:
		return true
	case <-time.After(timeout):
		return false
	}
}

func (this *Event) Wait() {
	var e chan byte
	this.lock.Lock()
	if this.isExpires {
		this.lock.Unlock()
		return
	}
	e = this.event
	this.lock.Unlock()
	select {
	case <-e:
	}
}

func (this *Event) Signal() {
	this.lock.Lock()
	defer this.lock.Unlock()
	if this.isExpires {
		panic("event is expires!")
	}
	close(this.event)
	if this.autoReusable {
		this.event = make(chan byte)
	} else {
		this.isExpires = true
	}
}

func (this *Event) ReSet() {
	this.lock.Lock()
	defer this.lock.Unlock()
	if this.isExpires {
		this.event = make(chan byte)
		this.isExpires = false
	} else {
		panic("event had not targeted!")
	}
}

func CreateEventWait(autoReusable bool) *Event {
	return &Event{
		event:        make(chan byte),
		autoReusable: autoReusable,
		isExpires:    false,
	}
}
