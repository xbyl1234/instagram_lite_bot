package types

import (
	"CentralizedControl/ins_lite/proto/io"
)

type BitByteFlags struct {
	Flags     []byte `json:"-"`
	FlagCount int    `json:"-"`
}

func BitByteFlagsJudge(flag byte, idx int) bool {
	return flag&(1<<(idx%8)) != 0
}

func (this *BitByteFlags) MakeFlags(from io.BufferReader) {
	byteCount := (this.FlagCount / 8) + 1
	this.Flags = make([]byte, byteCount)
	for i := 0; i < byteCount; i++ {
		this.Flags[i] = from.ReadByte()
	}
	for i := 0; i < this.FlagCount; i++ {
		//log.Debug("judge: %x - %d, %v", this.Flags[i/8], i, BitByteFlagsJudge(this.Flags[i/8], i))
	}
}

func (this *BitByteFlags) GetFlags(idx int) bool {
	flag := this.Flags[idx/8]
	return BitByteFlagsJudge(flag, idx)
}

func (this *BitByteFlags) Write(to io.BufferWriter) {
	panic("not impl")
}

func (this *BitByteFlags) Read(from io.BufferReader) {
	this.MakeFlags(from)
}
