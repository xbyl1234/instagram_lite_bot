package types

import (
	"CentralizedControl/ins_lite/proto/io"
	"reflect"
)

type DescribeValue[Desc IntLike] struct {
	// byte 1 int32 2 string 3
	describe Desc
	Value    interface{}
}

func (this *DescribeValue[Desc]) Write(to io.BufferWriter) {
	dv := int64(0)
	switch this.Value.(type) {
	case byte:
		dv = 1
	case int32:
		dv = 2
	case string:
		dv = 3
	}

	d := reflect.ValueOf(&this.describe)
	dd := d.Elem()
	setIntValue(d, dd, dv)
	GetWriteFunc(dd.Type())(to, dd)

	switch this.Value.(type) {
	case byte:
		to.WriteByte(this.Value.(byte))
	case int32:
		to.WriteInt(this.Value.(int32))
	case string:
		to.WriteString(this.Value.(string))
	}
}

func (this *DescribeValue[Desc]) Read(from io.BufferReader) {
	d := reflect.ValueOf(&this.describe)
	dd := d.Elem()
	GetReadFunc(d.Type())(from, dd)

	dv := getIntValue(d, dd)

	switch dv {
	case 1:
		this.Value = from.ReadByte()
	case 2:
		this.Value = from.ReadInt()
	case 3:
		this.Value = from.ReadString()
	}
}
