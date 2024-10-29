package types

import "CentralizedControl/ins_lite/proto/io"

type VarUInt32 struct {
	Value uint32
}

func (this *VarUInt32) Write(to io.BufferWriter) {
	to.WriteVarUInt32(this.Value)
}

func (this *VarUInt32) Read(from io.BufferReader) {
	this.Value = from.ReadVarUInt32()
}

func (this *VarUInt32) Set(v int64) {
	this.Value = uint32(v)
}

func (this *VarUInt32) Get() int64 {
	return int64(this.Value)
}

func CreateVarUInt32(i uint32) *VarUInt32 {
	return &VarUInt32{
		Value: i,
	}
}
