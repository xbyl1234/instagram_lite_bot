package types

import (
	"CentralizedControl/ins_lite/proto/io"
	"reflect"
)

type ListPadding[Value any, Size IntLike] struct {
	Size  Size
	Value []Value
}

func (this *ListPadding[Value, Size]) Put(value Value) {
	this.Value = append(this.Value, value)
	SetIntValue(&this.Size, int64(len(this.Value)))
}

func (this *ListPadding[Value, Size]) Write(to io.BufferWriter) {
	size := WriteValueAndGetInt(to, reflect.ValueOf(&this.Size))
	for i := int64(0); i < size; i++ {
		v := reflect.ValueOf(&this.Value[i])
		GetWriteFunc(v.Type())(to, v)
	}
}

func (this *ListPadding[Value, Size]) Read(from io.BufferReader) {
	size := ReadValueAndGetInt(from, reflect.ValueOf(&this.Size)) - 1
	if size == 0 {
		return
	}
	this.Value = make([]Value, size)
	for i := 0; int64(i) < size; i++ {
		v := reflect.ValueOf(&this.Value[i])
		GetReadFunc(v.Type())(from, v)
	}
}
