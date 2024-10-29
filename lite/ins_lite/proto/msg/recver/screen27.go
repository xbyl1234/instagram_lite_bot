package recver

import (
	"CentralizedControl/ins_lite/proto/io"
	"CentralizedControl/ins_lite/proto/types"
)

type SubScreen27Body struct {
	BitByteFlags types.BitByteFlags                                        `ins_init:"BitByteFlags(33)"`
	Unknow0      int64                                                     `ins:"get_flag(0)"`
	Unknow1      types.VarInt32                                            `ins:"get_flag(7)"`
	Unknow2      types.VarInt32                                            `ins:"get_flag(8)"`
	Unknow3      types.VarInt32                                            `ins:"get_flag(9)"`
	Unknow4      types.VarInt32                                            `ins:"get_flag(10)"`
	Unknow5      types.VarInt32                                            `ins:"get_flag(11)"`
	Unknow6      types.VarInt32                                            `ins:"get_flag(12)"`
	Unknow7      types.VarInt32                                            `ins:"get_flag(13)"`
	Unknow8      byte                                                      `ins:"get_flag(14)"`
	Unknow9      types.VarUInt32String                                     `ins:"get_flag(15)"`
	Unknow10     types.VarUInt32String                                     `ins:"get_flag(16)"`
	Unknow11     types.VarUInt32String                                     `ins:"get_flag(17)"`
	Unknow12     types.VarUInt32String                                     `ins:"get_flag(18)"`
	Unknow13     types.VarUInt32String                                     `ins:"get_flag(19)"`
	Unknow14     types.VarUInt32String                                     `ins:"get_flag(20)"`
	Unknow15     types.VarInt32                                            `ins:"get_flag(22)"`
	Unknow16     types.ListPadding[types.VarUInt32String, types.VarUInt32] `ins:"get_flag(23)"`
	Unknow17     types.VarInt32                                            `ins:"get_flag(26)"`
	Unknow18     types.VarInt32                                            `ins:"get_flag(27)"`
	Unknow19     types.VarInt32                                            `ins:"get_flag(28)"`
	Unknow20     types.VarInt32                                            `ins:"get_flag(29)"`
	Unknow21     types.VarInt32                                            `ins:"get_flag(30)"`
	Unknow22     types.ListPadding[types.VarInt32, types.VarUInt32]        `ins:"get_flag(31)"`
}

type SubScreen27 struct {
	ScreenImpl      `ins:"false" json:"-"`
	ScreenBase      `json:"screen_base"`
	SubScreen27Body `json:"sub_screen27_body"`
}

func (this *SubScreen27) ReadChange(from io.BufferReader, changeType int, parentScreen *SubScreen) {
	panic("SubScreen27 ReadChange")
}

func (this *SubScreen27) UpdateScreen(from io.BufferReader) {
	panic("SubScreen27 UpdateScreen")
}

func (this *SubScreen27) Write(to io.BufferWriter) {

}

func (this *SubScreen27) Read(from io.BufferReader) {
	types.ReadMsg(from, &this.ScreenBase)
	types.ReadMsg(from, &this.SubScreen27Body)
	this.InitImpl()
}

func (this *SubScreen27) InitImpl() {
	this.ScreenImpl = this
}

func (this *SubScreen27) GetSubScreen() *SubScreenArray {
	return nil
}

func (this *SubScreen27) GetDisplayActionScreenCmdCode() int32 {
	return 0
}
