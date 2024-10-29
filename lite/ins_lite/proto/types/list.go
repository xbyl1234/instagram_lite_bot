package types

import (
	"CentralizedControl/ins_lite/proto/io"
	"reflect"
)

type ListValue[Value any, Size IntLike] struct {
	Size  Size
	Value []Value
}

func (this *ListValue[Value, Size]) Put(value Value) {
	this.Value = append(this.Value, value)
	SetIntValue(&this.Size, int64(len(this.Value)))
}

func (this *ListValue[Value, Size]) Write(to io.BufferWriter) {
	size := WriteValueAndGetInt(to, reflect.ValueOf(&this.Size))
	//for i := int64(0); i < int64(len(this.Value)); i++ {
	for i := int64(0); i < size; i++ {
		v := reflect.ValueOf(&this.Value[i])
		GetWriteFunc(v.Type())(to, v)
	}
}

func (this *ListValue[Value, Size]) Read(from io.BufferReader) {
	size := ReadValueAndGetInt(from, reflect.ValueOf(&this.Size))
	this.Value = make([]Value, size)
	for i := 0; int64(i) < size; i++ {
		v := reflect.ValueOf(&this.Value[i])
		GetReadFunc(v.Type())(from, v)
	}
}

func CreateListValue[Value any, Size IntLike](l []Value) *ListValue[Value, Size] {
	ret := &ListValue[Value, Size]{
		Value: nil,
	}
	SetIntValue(&ret.Size, 0)
	if l == nil || len(l) == 0 {
		return ret
	}
	for _, item := range l {
		ret.Put(item)
	}
	return ret
}
