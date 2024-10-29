package recver

import (
	"CentralizedControl/ins_lite/proto/io"
	"CentralizedControl/ins_lite/proto/types"
)

type SubScreen7Body struct {
	BitByteFlags types.BitByteFlags    `ins_init:"BitByteFlags(27)"`
	SubScreen    SubScreen             `ins:"get_flag(6)"`
	Unknow1      types.VarUInt32String `ins:"get_flag(7)"`
	Unknow2      types.VarUInt32String `ins:"get_flag(8)"`
	Unknow3      types.VarInt32        `ins:"get_flag(9)"`
	Unknow4      int64                 `ins:"get_flag(10)"`
	Unknow5      types.VarInt32        `ins:"get_flag(11)"`
	Unknow6      types.VarInt32        `ins:"get_flag(12)"`
	Unknow7      types.VarUInt32String `ins:"get_flag(13)"`
	Unknow8      types.VarUInt32String `ins:"get_flag(14)"`
	Unknow9      types.VarUInt32String `ins:"get_flag(15)"`
	Unknow10     types.VarInt32        `ins:"get_flag(16)"`
	Unknow11     types.VarInt32        `ins:"get_flag(21)"`
	Unknow12     types.VarInt32        `ins:"get_flag(22)"`
	Unknow13     int64                 `ins:"get_flag(23)"`
	Unknow14     types.VarInt32        `ins:"get_flag(24)"`
	Unknow15     types.VarInt32        `ins:"get_flag(26)"`
}

type SubScreen7 struct {
	ScreenImpl     `ins:"false" json:"-"`
	ScreenBase     `json:"screen_base"`
	SubScreen7Body `json:"sub_screen7_body"`
}

func (this *SubScreen7) ReadChange(from io.BufferReader, changeType int, parentScreen *SubScreen) {
	panic("SubScreen7 ReadChange")
}

func (this *SubScreen7) UpdateScreen(from io.BufferReader) {
	panic("SubScreen7 UpdateScreen")
}

func (this *SubScreen7) Write(to io.BufferWriter) {

}

func (this *SubScreen7) Read(from io.BufferReader) {
	types.ReadMsg(from, &this.ScreenBase)
	types.ReadMsg(from, &this.SubScreen7Body)
	this.InitImpl()
}

func (this *SubScreen7) InitImpl() {
	this.ScreenBase.ScreenImpl = this
}

func (this *SubScreen7) GetSubScreen() *SubScreenArray {
	return nil
}

func (this *SubScreen7) GetDisplayActionScreenCmdCode() int32 {
	return 0
}
