package types

import "CentralizedControl/ins_lite/proto/io"

const (
	Byte         = 8
	Bool         = 3
	VarInt       = 0
	Short        = 10
	VarLong      = 4
	Float        = 7
	String       = 2
	FixedInt     = 5
	FixedLong    = 9
	StringArray  = 6
	ByteArray    = 1
	RawByteArray = 11
	VarBitField  = 12
)

type AnyType struct {
	Len   int
	Type  int
	Value any
}

func (this *AnyType) Write(to io.BufferWriter) {
	switch this.Type {
	case Byte:
		to.WriteByte(this.Value.(byte))
	case Bool:
		if this.Value.(bool) {
			to.WriteByte(1)
		} else {
			to.WriteByte(0)
		}
	case VarInt:
		to.WriteVarUInt32(this.Value.(uint32))
	case Short:
		to.WriteShort(this.Value.(int16))
	case VarLong:
		to.WriteVarInt64(this.Value.(int64), false)
	case Float:
		to.WriteFloat32(this.Value.(float32))
	case String:
		to.WriteString(this.Value.(string))
	case FixedInt:
		to.WriteInt(this.Value.(int32))
	case FixedLong:
		to.WriteLong(this.Value.(int64))
	case StringArray:
		v := this.Value.([]string)
		to.WriteVarUInt32(uint32(len(v)))
		for idx := range v {
			to.WriteString(v[idx])
		}
	case ByteArray:
		v := this.Value.([]byte)
		to.WriteVarUInt32(uint32(len(v)))
		to.WriteBytes(v)
	case RawByteArray:
		to.WriteBytes(this.Value.([]byte))
	case VarBitField:
		v := this.Value.([]byte)
		var i = len(v) - 1
		for ; i > 0; i-- {
			if v[i] != 0 {
				break
			}
		}
		v = v[:i+1]
		to.WriteBytes(v)
	}
}

func (this *AnyType) Read(from io.BufferReader) {
	switch this.Type {
	case Byte:
		this.Value = from.ReadByte()
	case Bool:
		if from.ReadByte() != 0 {
			this.Value = true
		} else {
			this.Value = false
		}
	case VarInt:
		this.Value = from.ReadVarUInt32()
	case Short:
		this.Value = from.ReadShort()
	case VarLong:
		this.Value = from.ReadVarUInt64(false)
	case Float:
		this.Value = from.ReadFloat32()
	case String:
		this.Value = from.ReadString()
	case FixedInt:
		this.Value = from.ReadInt()
	case FixedLong:
		this.Value = from.ReadLong()
	case StringArray:
		l := from.ReadVarUInt32()
		v := make([]string, l)
		for i := 0; i < int(l); i++ {
			v[i] = from.ReadString()
		}
		this.Value = v
	case ByteArray:
		l := from.ReadVarUInt32()
		this.Value = from.ReadBytes(l)
	case RawByteArray:
		this.Value = from.ReadBytes(uint32(this.Len))
	case VarBitField:
		this.Value = from.ReadBytes(uint32(this.Len))
	}
}
