package recver

import (
	"CentralizedControl/ins_lite/proto/io"
	"CentralizedControl/ins_lite/proto/types"
)

type ScreenDataArrayItem struct {
	Index2        int16
	Index         int16
	ScreenCmdCode types.VarUInt32
	ScreenCmdData types.ListValue[byte, int16]
}

type ScreenDataArrayItem2 struct {
	Index         int16
	ScreenCmdCode types.VarUInt32
	ScreenCmdData types.ListValue[byte, int16]
}

type ScreenCmdArray struct {
	Flags0     byte
	Flags1     uint16
	ScreenCmd1 types.ListValue[ScreenDataArrayItem, int16] `ins:"Flags1 > 0 && Flags0 != 0"`
	ScreenCmd2 []ScreenDataArrayItem2                      `ins:"Flags1 > 0 && Flags0 == 0"`
}

func (this *ScreenCmdArray) GetScreenCmdByIdx(idx int32) (code uint32, data []byte) {
	if this.ScreenCmd2 != nil {
		for i := 0; i < len(this.ScreenCmd2); i++ {
			item := this.ScreenCmd2[i]
			if item.Index == int16(idx) {
				return item.ScreenCmdCode.Value, item.ScreenCmdData.Value
			}
		}
	}
	if this.ScreenCmd1.Size != 0 {
		for i := 0; i < len(this.ScreenCmd1.Value); i++ {
			item := this.ScreenCmd1.Value[i]
			if item.Index == int16(idx) {
				return item.ScreenCmdCode.Value, item.ScreenCmdData.Value
			}
		}
	}
	return 0, nil
}

func (this *ScreenCmdArray) Write(to io.BufferWriter) {

}

func (this *ScreenCmdArray) Read(from io.BufferReader) {
	types.ReadMsg(from, &this.Flags0)
	types.ReadMsg(from, &this.Flags1)
	if this.Flags1 > 0 {
		if this.Flags0 != 0 {
			types.ReadMsg(from, &this.ScreenCmd1)
		} else {
			this.ScreenCmd2 = make([]ScreenDataArrayItem2, this.Flags1)
			for i := 0; i < int(this.Flags1); i++ {
				types.ReadMsg(from, &this.ScreenCmd2[i])
			}
		}
	}
}
