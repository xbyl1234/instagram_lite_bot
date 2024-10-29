package recver

import (
	"CentralizedControl/common/log"
	"CentralizedControl/ins_lite/proto/io"
	"CentralizedControl/ins_lite/proto/types"
)

type ScreenDecodeBodyItem8Item struct {
	Unknow0 types.VarInt32
	Unknow1 types.VarInt32
}

type ScreenDecodeBodyItem8 struct {
	Unknow0 types.ListValue[ScreenDecodeBodyItem8Item, types.VarUInt32]
	Unknow1 types.ListValue[types.VarInt32, types.VarUInt32]
}
type ScreenDecodeBodyItem10Item struct {
	Unknow0 types.VarInt32
	Unknow1 types.ListPadding[types.VarInt32, types.VarUInt32]
}

type ScreenDecodeBodyItem10 struct {
	Unknow0 types.ListValue[ScreenDecodeBodyItem10Item, types.VarUInt32]
	Unknow1 types.ListValue[types.VarInt32, types.VarUInt32]
}

type ScreenDecodeBodyItem73ItemItem struct {
	Unknow0 types.VarInt32
	Unknow1 types.ListPadding[types.VarInt32, types.VarUInt32]
}

type ScreenDecodeBodyItem73Item struct {
	Unknow0 types.ListValue[ScreenDecodeBodyItem73ItemItem, types.VarUInt32]
	Unknow1 types.ListValue[types.VarInt32, types.VarUInt32]
}

type ScreenDecodeBodyItem73 struct {
	BitByteFlags types.BitByteFlags         `ins_init:"BitByteFlags(2)"`
	Unknow0      ScreenDecodeBodyItem73Item `ins:"get_flag(0)"`
}

func (this *ScreenDecodeBodyItem73) IsNull() bool {
	return this.BitByteFlags.Flags == nil
}

func (this *ScreenDecodeBodyItem73) IsValueEmpty() bool {
	return this.BitByteFlags.GetFlags(0)
}

type ScreenDecodeBodyItem5 struct {
	BitByteFlags types.BitByteFlags `ins_init:"BitByteFlags(6)"`
	Unknow1      types.VarInt32     `ins:"get_flag(1)"`
	Unknow2      types.VarInt32     `ins:"get_flag(2)"`
	Unknow3      types.VarInt32     `ins:"get_flag(3)"`
	Unknow4      byte               `ins:"get_flag(4)"`
	Unknow5      types.VarInt32     `ins:"get_flag(5)"`
}

type ScreenDecodeBody struct {
	BitByteFlags           types.BitByteFlags                                 `ins_init:"BitByteFlags(77)"`
	Unknow3                types.VarInt32                                     `ins:"get_flag(3)"`
	Unknow4                types.VarInt32                                     `ins:"get_flag(4)"`
	Unknow5                ScreenDecodeBodyItem5                              `ins:"get_flag(5)"`
	Unknow8                ScreenDecodeBodyItem8                              `ins:"get_flag(8)"`
	Unknow10               ScreenDecodeBodyItem10                             `ins:"get_flag(10)"`
	Unknow12               byte                                               `ins:"get_flag(12)"`
	Unknow13               types.VarInt32                                     `ins:"get_flag(13)"`
	Unknow15               int64                                              `ins:"get_flag(15)"`
	ScreenName             types.VarUInt32String                              `ins:"get_flag(16)"`
	Unknow20               types.VarInt32                                     `ins:"get_flag(20)"`
	ScreenUrl              types.VarUInt32String                              `ins:"get_flag(21)"`
	Unknow22               types.VarUInt32String                              `ins:"get_flag(22)"`
	RunScreenCode          types.VarInt32                                     `ins:"get_flag(25)"`
	Unknow26               types.VarInt32                                     `ins:"get_flag(26)"`
	Unknow27               types.ListPadding[byte, types.VarUInt32]           `ins:"get_flag(27)"`
	Unknow28               types.ListPadding[types.VarInt32, types.VarUInt32] `ins:"get_flag(28)"`
	Unknow29               types.ListPadding[types.VarInt32, types.VarUInt32] `ins:"get_flag(29)"`
	Unknow30               types.ListPadding[types.VarInt32, types.VarUInt32] `ins:"get_flag(30)"`
	Unknow31               types.ListPadding[types.VarInt32, types.VarUInt32] `ins:"get_flag(31)"`
	Unknow32               types.ListPadding[types.VarInt32, types.VarUInt32] `ins:"get_flag(32)"`
	Unknow34               types.VarInt32                                     `ins:"get_flag(34)"`
	Unknow35               types.VarInt32                                     `ins:"get_flag(35)"`
	Unknow36               types.VarUInt32String                              `ins:"get_flag(36)"`
	Unknow37               types.VarInt32                                     `ins:"get_flag(37)"`
	Unknow38               types.VarInt32                                     `ins:"get_flag(38)"`
	Unknow39               types.VarInt32                                     `ins:"get_flag(39)"`
	Unknow41               types.VarInt32                                     `ins:"get_flag(41)"`
	Unknow42               types.VarInt32                                     `ins:"get_flag(42)"`
	Unknow43               types.VarInt32                                     `ins:"get_flag(43)"`
	Unknow49               types.VarInt32                                     `ins:"get_flag(49)"`
	Unknow52               byte                                               `ins:"get_flag(52)"`
	Unknow53               byte                                               `ins:"get_flag(53)"`
	Unknow55               byte                                               `ins:"get_flag(55)"`
	Unknow56               byte                                               `ins:"get_flag(56)"`
	Unknow57               byte                                               `ins:"get_flag(57)"`
	Unknow64               byte                                               `ins:"get_flag(64)"`
	Unknow65               types.VarInt32                                     `ins:"get_flag(65)"`
	Unknow66               byte                                               `ins:"get_flag(66)"`
	Unknow67               types.VarInt32                                     `ins:"get_flag(67)"`
	Unknow70               byte                                               `ins:"get_flag(70)"`
	ScreenDecodeBodyItem73 ScreenDecodeBodyItem73                             `ins:"get_flag(73)"`
	Unknow75               types.VarInt32                                     `ins:"get_flag(75)"`
}

type ScreenDecode struct {
	ScreenImpl       `ins:"false" json:"-"`
	SubScreen2       `json:"sub_screen2"`
	ScreenDecodeBody `json:"screen_decode_body"`
	ScreenCmdArray   `json:"screen_cmd_array"`
}

func (this *ScreenDecode) Write(to io.BufferWriter) {

}

func (this *ScreenDecode) Read(from io.BufferReader) {
	types.ReadMsg(from, &this.SubScreen2)
	types.ReadMsg(from, &this.ScreenDecodeBody)
	types.ReadMsg(from, &this.ScreenCmdArray)
}

func (this *ScreenDecode) ReadChange(from io.BufferReader, _ int, _ *SubScreen) {
	var changeType byte
	var flag2 int
	var flag3 types.VarInt32
	changeType = from.ReadByte()
	log.Debug("ReadChange ScreenDecode type: %d, id: %d, offset: %d", changeType, this.SubScreen2.ScreenBase.ScreenId, from.Offset())
	if (changeType & 0x80) != 0 {
		types.ReadMsg(from, &flag3)
		if flag3.Value&1 == 0 {
			flag2 = 0
		} else {
			flag2 = 1
		}
		changeType = changeType & 0x7f
	} else {
		flag2 = 0
	}
	this.SubScreen2.ReadChange(from, int(changeType), WrapSubScreen(&this.SubScreen2))
	if changeType == 3 {
		if flag2 == 0 || from.ReadByte() != 0 {
			types.ReadMsg(from, &this.ScreenDecodeBody)
		}
		types.ReadMsg(from, &this.ScreenCmdArray)
	}
}

func (this *ScreenDecode) UpdateScreen(from io.BufferReader) {

}

func (this *ScreenDecode) GetSubScreen() *SubScreenArray {
	return this.SubScreenArray.Copy()
}

func (this *ScreenDecode) GetAllSubScreenCpy() *SubScreenArray {
	return GetAllSubScreenByCpy(CreateSubScreen(2, &this.SubScreen2))
}

func (this *ScreenDecode) GetAllSubScreenSource() *SubScreenArray {
	return GetAllSubScreenSource(CreateSubScreen(2, &this.SubScreen2))
}
