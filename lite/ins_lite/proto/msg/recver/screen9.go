package recver

import (
	"CentralizedControl/ins_lite/proto/io"
	"CentralizedControl/ins_lite/proto/types"
)

type SubScreen9SubScreen struct {
	Unknow    types.VarInt32
	Flag      byte
	SubScreen SubScreen
	Unknow1   types.ListValue[types.VarInt32, types.VarUInt32]
}

func (this *SubScreen9SubScreen) Write(to io.BufferWriter) {

}

func (this *SubScreen9SubScreen) Read(from io.BufferReader) {
	types.ReadMsg(from, &this.Unknow)
	types.ReadMsg(from, &this.Flag)
	if types.BitByteFlagsJudge(this.Flag, 1) {
		types.ReadMsg(from, &this.SubScreen)
	}
}

type SubScreen9Body struct {
	BitByteFlags        types.BitByteFlags                                    `ins_init:"BitByteFlags(4)"`
	SubScreen9SubScreen types.ListValue[SubScreen9SubScreen, types.VarUInt32] `ins:"get_flag(0)"`
	SubScreen           SubScreen                                             `ins:"get_flag(2)"`
	Unknow              types.VarInt32                                        `ins:"get_flag(3)"`
}

type SubScreen9 struct {
	ScreenImpl     `ins:"false" json:"-"`
	ScreenBase     `json:"screen_base"`
	SubScreen9Body `json:"sub_screen9_body"`
}

func (this *SubScreen9) ReadChange(from io.BufferReader, changeType int, parentScreen *SubScreen) {
	panic("SubScreen9 ReadChange")
}

func (this *SubScreen9) UpdateScreen(from io.BufferReader) {
	panic("SubScreen9 UpdateScreen")
}

func (this *SubScreen9) Write(to io.BufferWriter) {

}

func (this *SubScreen9) Read(from io.BufferReader) {
	types.ReadMsg(from, &this.ScreenBase)
	types.ReadMsg(from, &this.SubScreen9Body)
	this.InitImpl()
}

func (this *SubScreen9) InitImpl() {
	this.ScreenBase.ScreenImpl = this
}

func (this *SubScreen9) GetSubScreen() *SubScreenArray {
	return nil
}

func (this *SubScreen9) GetDisplayActionScreenCmdCode() int32 {
	return 0
}
