package types

import "CentralizedControl/ins_lite/proto/io"

type VarUInt32String struct {
	Value string
}

func (this *VarUInt32String) Write(to io.BufferWriter) {
	l := uint32(len(this.Value))
	to.WriteVarUInt32(l + 1)
	if l > 0 {
		to.WriteBytes([]byte(this.Value))
	}
}

func (this *VarUInt32String) Read(from io.BufferReader) {
	l := from.ReadVarUInt32()
	if l > 0 {
		if l == 65451 {
			println("")
		}
		this.Value = string(from.ReadBytes(l - 1))
	}
}
