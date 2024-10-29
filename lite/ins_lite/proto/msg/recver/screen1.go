package recver

import (
	"CentralizedControl/common/log"
	"CentralizedControl/ins_lite/proto/io"
	"CentralizedControl/ins_lite/proto/types"
)

//标签窗口

type SubScreen1Body struct {
	BitByteFlags               types.BitByteFlags                                        `ins_init:"BitByteFlags(23)"`
	Unknow1                    types.VarUInt32String                                     `ins:"get_flag(1)"`
	Unknow2                    types.VarUInt32String                                     `ins:"get_flag(2)"`
	Unknow9                    types.VarUInt32String                                     `ins:"get_flag(9)"`
	Unknow10                   float64                                                   `ins:"get_flag(10)"`
	DisplayActionScreenCmdCode types.VarInt32                                            `ins:"get_flag(11)"`
	Unknow12                   byte                                                      `ins:"get_flag(12)"`
	Unknow13                   types.VarInt32                                            `ins:"get_flag(13)"`
	ShowText                   types.VarUInt32String                                     `ins:"get_flag(14)"`
	Unknow15                   types.VarInt32                                            `ins:"get_flag(15)"`
	Unknow16                   byte                                                      `ins:"get_flag(16)"`
	Unknow17                   types.VarInt32                                            `ins:"get_flag(17)"`
	Unknow18                   types.VarInt32                                            `ins:"get_flag(18)"`
	Unknow19                   types.ListPadding[types.VarInt32, types.VarUInt32]        `ins:"get_flag(19)"`
	Unknow20                   byte                                                      `ins:"get_flag(20)"`
	Unknow21                   byte                                                      `ins:"get_flag(21)"`
	Unknow22                   types.ListPadding[types.VarUInt32String, types.VarUInt32] `ins:"get_flag(22)"`
}

type SubScreen1 struct {
	ScreenImpl     `ins:"false" json:"-"`
	ScreenBase     `json:"screen_base"`
	SubScreen1Body `json:"sub_screen_1_body"`
}

func (this *SubScreen1) ReadChange(from io.BufferReader, changeType int, parentScreen *SubScreen) {
	log.Debug("ReadChange SubScreen1 type: %d, id: %d, offset: %d", changeType, parentScreen.GetScreenId(), from.Offset())
	this.ScreenBase.ReadChange(from, changeType, parentScreen)
	if changeType == 6 {
		types.ReadMsg(from, &this.SubScreen1Body)
	}
}

func (this *SubScreen1) UpdateScreen(from io.BufferReader) {
	if this.ScreenImpl != nil {
		this.ScreenImpl.UpdateScreen(from)
		return
	}
	this.Read(from)
}

func (this *SubScreen1) Write(to io.BufferWriter) {

}

func (this *SubScreen1) Read(from io.BufferReader) {
	types.ReadMsg(from, &this.ScreenBase)
	types.ReadMsg(from, &this.SubScreen1Body)
	this.InitImpl()
}

func (this *SubScreen1) InitImpl() {
	this.ScreenBase.ScreenImpl = this
}

func (this *SubScreen1) GetSubScreen() *SubScreenArray {
	return nil
}

func (this *SubScreen1) GetDisplayActionScreenCmdCode() int32 {
	return this.DisplayActionScreenCmdCode.Value
}
