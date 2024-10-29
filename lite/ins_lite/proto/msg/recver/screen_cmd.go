package recver

import (
	"CentralizedControl/ins_lite/proto/io"
	"CentralizedControl/ins_lite/proto/types"
)

type ScreenCmd struct {
	ResourceId int32
	ToScreenId int32
	Flags0     byte
	Flags1     byte
	Unknow1    types.ListValue[uint16, byte]
	Unknow2    int16
	Unknow3    int16
	Unknow4    string
}

func (this *ScreenCmd) HasExternData() bool {
	return this.Flags0 != 0
}

func (this *ScreenCmd) Write(to io.BufferWriter) {

}

func (this *ScreenCmd) Read(from io.BufferReader) {
	this.ResourceId = from.ReadInt()
	this.ToScreenId = from.ReadInt()
	if !from.EOF() {
		this.Flags0 = from.ReadByte()
		if this.Flags0&0x80 != 0 {
			this.Flags1 = from.ReadByte()
		}
		if this.Flags0&16 != 0 {
			types.ReadMsg(from, &this.Unknow1)
		}
		if this.Flags0&64 != 0 {
			this.Unknow2 = from.ReadShort()
		}
		if this.Flags0&0x80 != 0 && this.Flags1&1 != 0 {
			this.Unknow3 = from.ReadShort()
		}
		if this.Flags1&4 != 0 {
			this.Unknow4 = from.ReadString()
		}
	}
}

type ScreenSpeak struct {
	SpeakText           string
	AppLanguage         types.ListValue[string, byte]
	IfOpenTtsView       bool
	ScreenCmdIdxOnStart int16
	ScreenCmdIdxOnStop  int16
	ScreenCmdIdxOnDone  int16
	ScreenCmdIdxOnError int16
	IfUseGoogleTts      bool
}

func (this *ScreenSpeak) Write(to io.BufferWriter) {

}

func (this *ScreenSpeak) Read(from io.BufferReader) {
	types.ReadMsg(from, &this.SpeakText)
	types.ReadMsg(from, &this.AppLanguage)
	types.ReadMsg(from, &this.IfOpenTtsView)
	types.ReadMsg(from, &this.ScreenCmdIdxOnStart)
	types.ReadMsg(from, &this.ScreenCmdIdxOnStop)
	types.ReadMsg(from, &this.ScreenCmdIdxOnDone)
	types.ReadMsg(from, &this.ScreenCmdIdxOnError)
	if !from.EOF() {
		types.ReadMsg(from, &this.IfUseGoogleTts)
	}
}
