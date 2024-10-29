package types

import "CentralizedControl/ins_lite/proto/io"

type VarInt32 struct {
	Value int32
}

func (this *VarInt32) Write(to io.BufferWriter) {
	to.WriteVarInt32(this.Value)
}

func (this *VarInt32) Read(from io.BufferReader) {
	this.Value = from.ReadVarInt32()
}

func (this *VarInt32) Set(v int64) {
	this.Value = int32(v)
}

func (this *VarInt32) Get() int64 {
	return int64(this.Value)
}

func CreateVarInt32(v int) VarInt32 {
	return VarInt32{
		Value: int32(v),
	}
}
