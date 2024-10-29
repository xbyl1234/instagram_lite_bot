package recver

import (
	"CentralizedControl/common/log"
	"CentralizedControl/ins_lite/proto/io"
	"CentralizedControl/ins_lite/proto/types"
	"sync"
)

type PropStore223 struct {
	Props types.MapValue[byte, string, byte]
}

type PersistentTranslatedItem struct {
	Key   string
	Value string
}

type PersistentTranslatedStringsMap struct {
	Props types.ListValue[PersistentTranslatedItem, int32]
	Flag  bool
}

func (this *PersistentTranslatedStringsMap) Write(to io.BufferWriter) {

}

func (this *PersistentTranslatedStringsMap) Read(from io.BufferReader) {
	types.ReadMsg(from, &this.Props)
	if !from.EOF() {
		this.Flag = from.ReadByte() != 0
	}
}

type PropStore54 struct {
	Props map[int]any `json:"props"`
	lock  sync.Mutex
}

func (this *PropStore54) Write(to io.BufferWriter) {

}

func (this *PropStore54) Update(news *PropStore54) {
	this.lock.Lock()
	defer this.lock.Unlock()
	for k, v := range news.Props {
		this.Props[k] = v
	}
}

func (this *PropStore54) GetBytes(key int, _default []byte) []byte {
	return this.GetProp(key, _default).([]byte)
}

func (this *PropStore54) GetBool(key int, _default bool) bool {
	v, ok := this.GetPropNoDefault(key)
	if !ok {
		return _default
	}
	return v != 0
}

func (this *PropStore54) GetInt(key int, _default int64) int64 {
	value := this.GetProp(key, _default)
	switch value.(type) {
	case byte:
		return int64(value.(byte))
	case int8:
		return int64(value.(int8))
	case uint16:
		return int64(value.(uint16))
	case uint32:
		return int64(value.(uint32))
	case uint64:
		return int64(value.(uint64))
	case int16:
		return int64(value.(int16))
	case int32:
		return int64(value.(int32))
	case int64:
		return int64(value.(int64))
	}
	panic("unknow type")
}

func (this *PropStore54) GetStr(key int, _default any) string {
	return this.GetProp(key, _default).(string)
}

func (this *PropStore54) GetProp(key int, _default any) any {
	this.lock.Lock()
	defer this.lock.Unlock()
	v, ok := this.Props[key]
	if !ok {
		v = _default
	}
	return v
}

func (this *PropStore54) GetPropNoDefault(key int) (any, bool) {
	this.lock.Lock()
	defer this.lock.Unlock()
	v, ok := this.Props[key]
	return v, ok
}

func (this *PropStore54) putProp(key int, value any) {
	log.Debug("put %d : %v", key, value)
	this.Props[key] = value
}

func (this *PropStore54) readSub(from io.BufferReader, format int, key int) bool {
	_type := from.ReadByte()
	switch _type {
	case 1:
		var value uint64
		if format == 0 {
			value = uint64(from.ReadLong())
		} else {
			value = from.ReadVarUInt64(true)
		}
		this.putProp(key, value)
	case 2:
		this.putProp(key, from.ReadString())
	case 3:
		var value uint32
		if format == 0 {
			value = uint32(from.ReadInt())
		} else {
			value = from.ReadVarUInt32()
		}
		this.putProp(key, value)
	case 4:
		l := from.ReadShort()
		value := make([]string, l)
		for idx := l - 1; idx >= 0; idx-- {
			value[idx] = from.ReadString()
		}
		this.putProp(key, value)
	case 5:
		this.putProp(key, nil)
	case 6:
		l := from.ReadShort()
		value := make([]uint32, l)
		for idx := l - 1; idx >= 0; idx-- {
			value[idx] = uint32(from.ReadInt())
		}
		this.putProp(key, value)
	case 7:
		l := from.ReadShort()
		value := make([]byte, l)
		for idx := l - 1; idx >= 0; idx-- {
			value[idx] = from.ReadByte()
		}
		this.putProp(key, value)
	case 8:
		this.putProp(key, from.ReadByte())
	case 9:
		l := from.ReadShort()
		value := make([]uint64, l)
		for idx := l - 1; idx >= 0; idx-- {
			value[idx] = uint64(from.ReadLong())
		}
		this.putProp(key, value)
	case 10:
		subProp := &PropStore54{
			Props: map[int]any{},
		}
		subProp.read(from, 1)
		this.putProp(key, subProp)
	case 11:
		this.putProp(key, from.ReadFloat64())
	case 12:
		this.putProp(key, from.ReadVarUInt64(false))
	}
	return true
}

func (this *PropStore54) read(from io.BufferReader, format int) bool {
	var count = 2147483647
	if format == 1 {
		count = int(from.ReadVarUInt32())
	}
	var result = true
	for !from.EOF() && count > 0 {
		var key = 0
		if format == 0 {
			key = int(from.ReadShort())
			if key == 462 {
				return false
			}
		} else {
			key = int(from.ReadVarUInt32()*0x1000 + from.ReadVarUInt32())
		}
		if !this.readSub(from, format, key) {
			result = false
		}
		count--
	}
	return result
}

func (this *PropStore54) Read(from io.BufferReader) {
	this.Props = map[int]any{}
	this.read(from, 0)
}
