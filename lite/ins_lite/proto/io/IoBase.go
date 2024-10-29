package io

import (
	"encoding/hex"
	"strings"
)

type BufferReader interface {
	ReadByte() byte
	ReadBytes(l uint32) []byte
	ReadInt() int32
	ReadLong() int64
	ReadShort() int16
	ReadString() string
	ReadFloat32() float32
	ReadFloat64() float64
	ReadVarInt32() int32
	ReadVarUInt32() uint32
	ReadVarUInt64(isBit63 bool) uint64
	ReadRemain() []byte
	EOF() bool
	PeekRemain() []byte
	Offset() int
}

type BufferWriter interface {
	WriteByte(d byte)
	WriteBytes(d []byte)
	WriteInt(d int32)
	WriteFloat32(d float32)
	WriteFloat64(d float64)
	WriteLong(d int64)
	WriteShort(d int16)
	WriteString(d string)
	WriteVarInt32(l int32)
	WriteVarUInt32(l uint32)
	WriteVarInt64(v int64, isSigned bool)
}

type IoInterface interface {
	BufferReader
	BufferWriter
}

func DecodeHexData(data []byte) []byte {
	s := string(data)
	s = strings.ReplaceAll(s, " ", "")
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, "\r", "")
	s = strings.ReplaceAll(s, "\t", "")
	decodeString, _ := hex.DecodeString(s)
	return decodeString
}
