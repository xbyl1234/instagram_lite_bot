package recver

import (
	"CentralizedControl/common/log"
	"CentralizedControl/ins_lite/proto/io"
	"CentralizedControl/ins_lite/proto/types"
)

type SubScreen19Item6 struct {
	BitByteFlags types.BitByteFlags `ins_init:"BitByteFlags(5)"`
	Unknow0      types.VarInt32     `ins:"get_flag(0)"`
	Unknow3      float64            `ins:"get_flag(3)"`
	Unknow4      float64            `ins:"get_flag(4)"`
}

type SubScreen19Item9 struct {
	BitByteFlags types.BitByteFlags    `ins_init:"BitByteFlags(3)"`
	Unknow2      types.VarUInt32String `ins:"get_flag(2)"`
}

type SubScreen19Body struct {
	BitByteFlags types.BitByteFlags    `ins_init:"BitByteFlags(15)"`
	Unknow1      byte                  `ins:"get_flag(1)"`
	Unknow6      SubScreen19Item6      `ins:"get_flag(6)"`
	Unknow8      types.VarInt32        `ins:"get_flag(8)"`
	Unknow9      SubScreen19Item9      `ins:"get_flag(9)"`
	Unknow11     int64                 `ins:"get_flag(11)"`
	Unknow12     types.VarUInt32String `ins:"get_flag(12)"`
	Unknow13     types.VarInt32        `ins:"get_flag(13)"`
	Unknow14     byte                  `ins:"get_flag(14)"`
}

type SubScreen19 struct {
	ScreenImpl      `ins:"false" json:"-"`
	SubScreen1      `json:"sub_screen_1"`
	SubScreen19Body `json:"sub_screen_19_body"`
}

func (this *SubScreen19) ReadChange(from io.BufferReader, changeType int, parentScreen *SubScreen) {
	log.Debug("ReadChange SubScreen19 type: %d, id: %d, offset: %d", changeType, parentScreen.GetScreenId(), from.Offset())
	this.SubScreen1.ReadChange(from, changeType, parentScreen)
	if changeType == 6 {
		types.ReadMsg(from, &this.SubScreen19Body)
	}
}

func (this *SubScreen19) UpdateScreen(from io.BufferReader) {
	if this.ScreenImpl != nil {
		this.ScreenImpl.UpdateScreen(from)
		return
	}
	this.Read(from)
}

func (this *SubScreen19) Write(to io.BufferWriter) {

}

func (this *SubScreen19) Read(from io.BufferReader) {
	types.ReadMsg(from, &this.SubScreen1)
	types.ReadMsg(from, &this.SubScreen19Body)
	this.InitImpl()
}

func (this *SubScreen19) InitImpl() {
	this.SubScreen1.ScreenImpl = this
}

func (this *SubScreen19) GetSubScreen() *SubScreenArray {
	return nil
}

func (this *SubScreen19) GetDisplayActionScreenCmdCode() int32 {
	return this.DisplayActionScreenCmdCode.Value
}
