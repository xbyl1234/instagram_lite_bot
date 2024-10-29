package recver

import (
	"CentralizedControl/common/log"
	"CentralizedControl/ins_lite/proto/io"
	"CentralizedControl/ins_lite/proto/types"
)

type SubScreen13Body struct {
	BitByteFlags types.BitByteFlags `ins_init:"BitByteFlags(3)"`
	Unknow1      types.VarInt32     `ins:"get_flag(1)"`
}

type SubScreen13 struct {
	ScreenImpl      `ins:"false" json:"-"`
	SubScreen2      `json:"sub_screen_2"`
	SubScreen13Body `json:"sub_screen_13_body"`
}

func (this *SubScreen13) ReadChange(from io.BufferReader, changeType int, parentScreen *SubScreen) {
	log.Debug("ReadChange SubScreen13 type: %d, id: %d, offset: %d", changeType, parentScreen.GetScreenId(), from.Offset())
	this.SubScreen2.ReadChange(from, changeType, parentScreen)
	types.ReadMsg(from, &this.SubScreen13Body)
}

func (this *SubScreen13) UpdateScreen(from io.BufferReader) {
	if this.ScreenImpl != nil {
		this.ScreenImpl.UpdateScreen(from)
		return
	}
	this.Read(from)
}

func (this *SubScreen13) Write(to io.BufferWriter) {

}

func (this *SubScreen13) Read(from io.BufferReader) {
	types.ReadMsg(from, &this.SubScreen2)
	types.ReadMsg(from, &this.SubScreen13Body)
	this.InitImpl()
}

func (this *SubScreen13) InitImpl() {
	this.SubScreen2.ScreenImpl = this
}

func (this *SubScreen13) GetSubScreen() *SubScreenArray {
	return this.SubScreenArray.Copy()
}

func (this *SubScreen13) GetDisplayActionScreenCmdCode() int32 {
	return 0
}
