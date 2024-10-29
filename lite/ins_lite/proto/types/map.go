package types

import (
	"CentralizedControl/ins_lite/proto/io"
	"reflect"
)

type KeyValue[Key any, Value any] struct {
	Key   Key
	Value Value
}

func (this *KeyValue[Key, Value]) Write(to io.BufferWriter) {
	k := reflect.ValueOf(&this.Key).Elem()
	v := reflect.ValueOf(&this.Value).Elem()
	GetWriteFunc(k.Type())(to, k)
	GetWriteFunc(v.Type())(to, v)
}

func (this *KeyValue[Key, Value]) Read(from io.BufferReader) {
	k := reflect.ValueOf(&this.Key).Elem()
	v := reflect.ValueOf(&this.Value).Elem()
	GetReadFunc(k.Type())(from, k)
	GetReadFunc(v.Type())(from, v)
}

type MapValue[Key any, Value any, Size any] struct {
	Size Size
	Kv   []KeyValue[Key, Value]
}

func (this *MapValue[Key, Value, Size]) Write(to io.BufferWriter) {
	size := WriteValueAndGetInt(to, reflect.ValueOf(&this.Size))
	for i := int64(0); i < size; i++ {
		this.Kv[i].Write(to)
	}
}

func (this *MapValue[Key, Value, Size]) Read(from io.BufferReader) {
	size := ReadValueAndGetInt(from, reflect.ValueOf(&this.Size))
	this.Kv = make([]KeyValue[Key, Value], size)
	for i := int64(0); i < size; i++ {
		this.Kv[i].Read(from)
	}
}

func (this *MapValue[Key, Value, Size]) Put(key Key, value Value) {
	kv := KeyValue[Key, Value]{
		Key:   key,
		Value: value,
	}
	this.Kv = append(this.Kv, kv)
	SetIntValue(&this.Size, int64(len(this.Kv)))
}
