package proto

import (
	"CentralizedControl/ins_lite/proto/io"
	"CentralizedControl/ins_lite/proto/types"
	"encoding/json"
	"os"
	"reflect"
	"time"
)

type MessageInterface interface {
	ReadFrom(data []byte)
	WriteTo() []byte
	GetCode(isSend bool) uint64
	SetCode(code uint64)
	GetTime() int64
	GetSenderIdx() int32
	GetRecverIdx() int32
	SetSenderIdx(int32)
	SetRecverIdx(int32)
	GetBody() any
}

type Message[Body MessageType] struct {
	Code      uint64
	Body      Body
	SenderIdx int32
	RecverIdx int32
	Time      time.Time
}

func (this *Message[Body]) GetBody() any {
	return &this.Body
}

func (this *Message[Body]) CreateFromBytes(body []byte) {
	reader := io.CreateReaderBuffer(body)
	types.ReadMsg(reader, &this.Body)
}

func (this *Message[Body]) GetCode(isSend bool) uint64 {
	if this.Code == 0 {
		this.Code = GetMessageCode(isSend, reflect.TypeOf(this.Body).Name())
	}
	return this.Code
}

func (this *Message[Body]) SetCode(code uint64) {
	this.Code = code
}

func (this *Message[Body]) ReadFromHexFile(path string) error {
	file, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	this.ReadFrom(io.DecodeHexData(file))
	return nil
}

func (this *Message[Body]) ReadFromFile(path string) error {
	file, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	this.ReadFrom(file)
	return nil
}

func (this *Message[Body]) ReadFromJsonFile(path string) error {
	file, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return this.ReadFromJson(file)
}

func (this *Message[Body]) ReadFromJson(jsonByte []byte) error {
	return json.Unmarshal(jsonByte, &this.Body)
}

func (this *Message[Body]) ReadFrom(data []byte) {
	reader := io.CreateReaderBuffer(data)
	types.ReadMsg(reader, &this.Body)
}

func (this *Message[Body]) WriteTo() []byte {
	writer := io.CreateWriter(512)
	types.WriteMsg(writer, &this.Body)
	return writer.GetBytes()
}

func (this *Message[Body]) GetSenderIdx() int32 {
	return this.SenderIdx
}

func (this *Message[Body]) GetRecverIdx() int32 {
	return this.RecverIdx
}

func (this *Message[Body]) SetSenderIdx(i int32) {
	this.SenderIdx = i
}

func (this *Message[Body]) SetRecverIdx(i int32) {
	this.RecverIdx = i
}

func (this *Message[Body]) GetTime() int64 {
	return this.Time.UnixMilli()
}

type MessageC[Body MessageType] struct {
	Message[Body]
	ClientId int64
	Magic    int32
}
