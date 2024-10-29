package types

import (
	"CentralizedControl/common/goexpr"
	"CentralizedControl/common/log"
	"CentralizedControl/ins_lite/proto/io"
	"reflect"
)

func WriteValueAndGetInt(to io.BufferWriter, v reflect.Value) int64 {
	dv := v.Elem()
	GetWriteFunc(dv.Type())(to, dv)
	return getIntValue(v, dv)
}

func ReadValueAndGetInt(from io.BufferReader, v reflect.Value) int64 {
	dv := v.Elem()
	GetReadFunc(dv.Type())(from, dv)
	return getIntValue(v, dv)
}

var engine *goexpr.Engine
var IntAbleType reflect.Type
var CustomWriterType reflect.Type

type ExprContext struct {
	CurStructObj any
	Reader       *io.BufferReader
}

func init() {
	IntAbleType = reflect.TypeOf((*IntAble)(nil)).Elem()
	CustomWriterType = reflect.TypeOf((*CustomWriter)(nil)).Elem()
	engine = goexpr.NewEngine(func(path string, inputArgs goexpr.InputArgs) interface{} {
		s := inputArgs.(reflect.Value)
		v := s.FieldByName(path)
		t := v.Type()
		switch t.Kind() {
		case reflect.Struct:
			if reflect.PointerTo(t).Implements(IntAbleType) {
				return v.Addr().Interface().(IntAble).Get()
			}
			return v.Interface()
		default:
			return v.Interface()
		}
	})

	engine.AddFunc("get_flag", func(inputArgs goexpr.InputArgs, v ...interface{}) interface{} {
		s := inputArgs.(reflect.Value)
		agrsV := s.FieldByName("BitByteFlags")
		m := agrsV.Addr().Interface().(*BitByteFlags)
		return m.GetFlags(int(goexpr.Int64Val(v[0])))
	})

	engine.AddFunc("BitByteFlags", func(inputArgs goexpr.InputArgs, v ...interface{}) interface{} {
		return goexpr.Int64Val(v[0])
	})
}

func StructWrite(to io.BufferWriter, v reflect.Value) {
	if reflect.PointerTo(v.Type()).Implements(CustomWriterType) {
		v.Addr().Interface().(CustomWriter).Write(to)
	} else {
		fieldNum := v.NumField()
		t := v.Type()
		for i := 0; i < fieldNum; i++ {
			fieldType := t.Field(i)
			if !fieldType.IsExported() {
				log.Debug("write idx: %d, name: %s, not export pass", i, fieldType.Name)
				continue
			}
			tagValue, ok := fieldType.Tag.Lookup("ins")
			if ok {
				if !engine.Execute(tagValue, v).(bool) {
					log.Debug("write idx: %d, name: %s, tags %s is false", i, fieldType.Name, tagValue)
					continue
				}
			}
			field := v.Field(i)
			log.Debug("write idx: %d, read %s", i, fieldType.Name)
			GetWriteFunc(field.Type())(to, field)
		}
	}
}

func GetWriteFunc(t reflect.Type) func(to io.BufferWriter, v reflect.Value) {
	switch t.Kind() {
	case reflect.Bool:
		return func(to io.BufferWriter, v reflect.Value) {
			if v.Bool() {
				to.WriteByte(1)
			} else {
				to.WriteByte(0)
			}
		}
	case reflect.Uint8: //byte
		return func(to io.BufferWriter, v reflect.Value) {
			to.WriteByte(byte(v.Uint()))
		}
	case reflect.Int16:
		return func(to io.BufferWriter, v reflect.Value) {
			to.WriteShort(int16(v.Int()))
		}
	case reflect.Int32:
		return func(to io.BufferWriter, v reflect.Value) {
			to.WriteInt(int32(v.Int()))
		}
	case reflect.Int64:
		return func(to io.BufferWriter, v reflect.Value) {
			to.WriteLong(int64(v.Int()))
		}
	case reflect.Uint16:
		return func(to io.BufferWriter, v reflect.Value) {
			to.WriteShort(int16(v.Uint()))
		}
	case reflect.Uint32:
		return func(to io.BufferWriter, v reflect.Value) {
			to.WriteInt(int32(v.Uint()))
		}
	case reflect.Uint64:
		return func(to io.BufferWriter, v reflect.Value) {
			to.WriteLong(int64(v.Uint()))
		}
	case reflect.String:
		return func(to io.BufferWriter, v reflect.Value) {
			to.WriteString(v.String())
		}
	case reflect.Float32:
		return func(to io.BufferWriter, v reflect.Value) {
			to.WriteFloat32(float32(v.Float()))
		}
	case reflect.Float64:
		return func(to io.BufferWriter, v reflect.Value) {
			to.WriteFloat64(v.Float())
		}
	case reflect.Pointer:
		return func(to io.BufferWriter, v reflect.Value) {
			GetWriteFunc(v.Elem().Type())(to, v.Elem())
		}
	case reflect.Struct:
		return StructWrite
	default:
		return func(to io.BufferWriter, v reflect.Value) {
			log.Warn("do nothing")
		}
	}
}

func StructRead(from io.BufferReader, v reflect.Value) {
	log.Debug("read struct name: %s, offset: %d", v.Type().Name(), from.Offset())
	customWriterType := reflect.TypeOf((*CustomWriter)(nil)).Elem()
	if reflect.PointerTo(v.Type()).Implements(customWriterType) {
		v.Addr().Interface().(CustomWriter).Read(from)
	} else {
		fieldNum := v.NumField()
		t := v.Type()
		for i := 0; i < fieldNum; i++ {
			fieldValue := v.Field(i)
			fieldType := t.Field(i)
			if !fieldType.IsExported() {
				log.Debug("read idx: %d, name: %s, not export pass", i, fieldType.Name)
				continue
			}
			tagValue, ok := fieldType.Tag.Lookup("ins")
			if ok {
				if !engine.Execute(tagValue, v).(bool) {
					log.Debug("read idx: %d, name: %s, tags %s is false", i, fieldType.Name, tagValue)
					continue
				}
			}
			tagInitValue, ok := fieldType.Tag.Lookup("ins_init")
			if ok {
				count := engine.Execute(tagInitValue, v).(int64)
				m := fieldValue.FieldByName("FlagCount")
				m.SetInt(count)
			}
			log.Debug("read idx: %d, read %s", i, fieldType.Name)
			GetReadFunc(fieldValue.Type())(from, fieldValue)
		}
	}
}

func GetReadFunc(t reflect.Type) func(from io.BufferReader, v reflect.Value) {
	switch t.Kind() {
	case reflect.Bool:
		return func(from io.BufferReader, v reflect.Value) {
			v.SetBool(from.ReadByte() != 0)
		}
	case reflect.Uint8: //byte
		return func(from io.BufferReader, v reflect.Value) {
			v.SetUint(uint64(from.ReadByte()))
		}
	case reflect.Int16:
		return func(from io.BufferReader, v reflect.Value) {
			v.SetInt(int64(from.ReadShort()))
		}
	case reflect.Int32:
		return func(from io.BufferReader, v reflect.Value) {
			v.SetInt(int64(from.ReadInt()))
		}
	case reflect.Int64:
		return func(from io.BufferReader, v reflect.Value) {
			v.SetInt(from.ReadLong())
		}
	case reflect.Uint16:
		return func(from io.BufferReader, v reflect.Value) {
			v.SetUint(uint64(from.ReadShort()))
		}
	case reflect.Uint32:
		return func(from io.BufferReader, v reflect.Value) {
			v.SetUint(uint64(from.ReadInt()))
		}
	case reflect.Uint64:
		return func(from io.BufferReader, v reflect.Value) {
			v.SetUint(uint64(from.ReadLong()))
		}
	case reflect.String:
		return func(from io.BufferReader, v reflect.Value) {
			v.SetString(from.ReadString())
		}
	case reflect.Float32:
		return func(from io.BufferReader, v reflect.Value) {
			v.SetFloat(float64(from.ReadFloat32()))
		}
	case reflect.Float64:
		return func(from io.BufferReader, v reflect.Value) {
			v.SetFloat(from.ReadFloat64())
		}
	case reflect.Pointer:
		return func(from io.BufferReader, v reflect.Value) {
			GetReadFunc(v.Elem().Type())(from, v.Elem())
		}
	case reflect.Struct:
		return StructRead
	default:
		return func(from io.BufferReader, v reflect.Value) {
			log.Warn("do nothing: " + v.Type().Name())
		}
	}
}
