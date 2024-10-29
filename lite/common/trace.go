package common

import (
	"CentralizedControl/common/log"
	"bytes"
	"fmt"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type TraceType string
type TraceSetType string

type TraceInfo struct {
	success int32
	failed  int32
}

type TraceSetInfo struct {
	count int32
}

type Trace struct {
	trace    map[TraceType]*TraceInfo
	traceSet map[TraceSetType]*TraceSetInfo
	file     *os.File
	dy       int
	lock     sync.Mutex
}

func (this *Trace) AddSuccess(key TraceType) {
	info := this.trace[key]
	atomic.AddInt32(&info.success, 1)
}

func (this *Trace) AddFailed(key TraceType) {
	info := this.trace[key]
	atomic.AddInt32(&info.failed, 1)
}

func (this *Trace) SetCount(key TraceSetType, count int32) {
	this.lock.Lock()
	this.traceSet[key].count = count
	this.lock.Unlock()
}

func (this *Trace) AddCount(key TraceSetType) {
	info := this.traceSet[key]
	atomic.AddInt32(&info.count, 1)
}

func (this *Trace) Log() {
	LogTrace(this)
}

func CreateTrace(logName string, traceKeys []TraceType, setKeys []TraceSetType, dy int) *Trace {
	location, _ := time.LoadLocation("Asia/Shanghai")
	timeName := time.Now().In(location).Format("2006-01-02 15:04:05")
	timeName = strings.ReplaceAll(timeName, ":", "_")
	file, err := os.OpenFile(fmt.Sprintf("./trace/%s_%s.txt", logName, timeName), os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		panic(err)
	}

	t := &Trace{
		trace:    make(map[TraceType]*TraceInfo),
		traceSet: make(map[TraceSetType]*TraceSetInfo),
		file:     file,
		dy:       dy,
	}

	for _, item := range traceKeys {
		t.trace[item] = new(TraceInfo)
	}

	for _, item := range setKeys {
		t.traceSet[item] = new(TraceSetInfo)
	}

	go LogTrace(t)
	return t
}

func StartLogTrace(t *Trace) {
	for true {
		LogTrace(t)
		time.Sleep(time.Duration(t.dy) * time.Second)
	}
}

func LogTrace(t *Trace) {
	buff := bytes.NewBufferString("")

	for k, v := range t.trace {
		log.Debug(fmt.Sprintf("trace:    %s success: %d, failed: %d", k, v.success, v.failed))
		buff.WriteString(fmt.Sprintf("%s \t success: %d, failed: %d\n", k, v.success, v.failed))
	}
	for k, v := range t.traceSet {
		log.Debug(fmt.Sprintf("trace:    %s : %d", k, v.count))
		buff.WriteString(fmt.Sprintf("%s \t : %d\n", k, v.count))
	}

	t.file.WriteAt(buff.Bytes(), 0)
	t.file.Sync()
}
