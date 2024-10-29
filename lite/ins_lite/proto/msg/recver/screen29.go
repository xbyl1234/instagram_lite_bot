package recver

import (
	"CentralizedControl/ins_lite/proto/io"
	"CentralizedControl/ins_lite/proto/types"
)

type SubScreen29 struct {
	ScreenImpl `ins:"false" json:"-"`
	ScreenBase `json:"screen_base"`
}

func (this *SubScreen29) ReadChange(from io.BufferReader, changeType int, parentScreen *SubScreen) {
	panic("SubScreen29 ReadChange")
}

func (this *SubScreen29) UpdateScreen(from io.BufferReader) {
	panic("SubScreen29 UpdateScreen")
}

func (this *SubScreen29) Write(to io.BufferWriter) {

}

func (this *SubScreen29) Read(from io.BufferReader) {
	types.ReadMsg(from, &this.ScreenBase)
	this.InitImpl()
}

func (this *SubScreen29) InitImpl() {
	this.ScreenImpl = this
}

func (this *SubScreen29) GetSubScreen() *SubScreenArray {
	return nil
}

func (this *SubScreen29) GetDisplayActionScreenCmdCode() int32 {
	return 0
}
