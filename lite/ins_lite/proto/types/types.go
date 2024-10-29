package types

import (
	"CentralizedControl/ins_lite/proto/io"
	"reflect"
)

type CustomWriter interface {
	Write(to io.BufferWriter)
	Read(from io.BufferReader)
}

type IntAble interface {
	Set(v int64)
	Get() int64
}

type IntLike interface {
	byte | int8 | int16 | int32 | int64 | VarInt32 | VarUInt32
}

func setIntValue(v reflect.Value, dv reflect.Value, intValue int64) {
	switch dv.Type().Kind() {
	case reflect.Struct:
		v.Interface().(IntAble).Set(intValue)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		dv.SetUint(uint64(intValue))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		dv.SetInt(intValue)
	}
}

func getIntValue(v reflect.Value, dv reflect.Value) int64 {
	size := int64(0)
	switch dv.Type().Kind() {
	case reflect.Struct:
		size = v.Interface().(IntAble).Get()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		size = int64(dv.Uint())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		size = dv.Int()
	}
	return size
}

func SetIntValue(obj any, intValue int64) {
	d := reflect.ValueOf(obj)
	dd := d.Elem()
	setIntValue(d, dd, intValue)
}

func GetIntValue(obj any) int64 {
	d := reflect.ValueOf(obj)
	dd := d.Elem()
	return getIntValue(d, dd)
}
