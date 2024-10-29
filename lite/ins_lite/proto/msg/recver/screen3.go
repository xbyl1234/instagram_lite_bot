package recver

import (
	"CentralizedControl/common/log"
	"CentralizedControl/ins_lite/proto/io"
	"CentralizedControl/ins_lite/proto/types"
)

type SubScreenSub3Item1 struct {
	Unknow1 types.VarInt32
	Unknow2 types.VarInt32
}

type SubScreenSub3Item struct {
	Unknow1 types.ListValue[SubScreenSub3Item1, types.VarUInt32]
	Unknow2 types.ListValue[types.VarInt32, types.VarUInt32]
}

type SubScreenSub3CustomItem struct {
	BitByteFlags types.BitByteFlags `ins_init:"BitByteFlags(4)"`
	Unknow1      int64
	Unknow2      types.VarInt32
	Unknow3      types.VarInt32
}

type SubScreenSub3Custom struct {
	Count    uint32
	FlagByte []byte
	Flags1   []bool
	Flags2   []bool
	Data     []*SubScreenSub3CustomItem
}

func (this *SubScreenSub3Custom) Write(to io.BufferWriter) {

}

func (this *SubScreenSub3Custom) Read(from io.BufferReader) {
	this.Count = from.ReadVarUInt32()
	if this.Count == 0 {
		return
	}
	this.Count -= 1
	this.Flags1 = make([]bool, this.Count)
	this.Flags2 = make([]bool, this.Count)
	var r byte
	for i := 0; i < int(this.Count); i++ {
		if i*2%8 == 0 {
			r = from.ReadByte()
		}
		this.Flags1[i] = types.BitByteFlagsJudge(r, i*2)
		this.Flags2[i] = types.BitByteFlagsJudge(r, i*2+1)
	}
	this.Data = make([]*SubScreenSub3CustomItem, 0)
	for i := 0; i < int(this.Count); i++ {
		if this.Flags1[i] {
			item := &SubScreenSub3CustomItem{}
			types.ReadMsg(from, item)
			this.Data = append(this.Data, item)
		}
	}
}

type SubScreen3Body struct {
	BitByteFlags             types.BitByteFlags                                        `ins_init:"BitByteFlags(138)"`
	Unknow0                  types.VarInt32                                            `ins:"get_flag(0)"`
	Unknow1                  types.VarInt32                                            `ins:"get_flag(1)"`
	Unknow2                  byte                                                      `ins:"get_flag(2)"`
	Unknow4                  types.VarInt32                                            `ins:"get_flag(4)"`
	Unknow7                  types.VarInt32                                            `ins:"get_flag(7)"`
	Unknow8                  types.VarInt32                                            `ins:"get_flag(8)"`
	Unknow9                  types.VarInt32                                            `ins:"get_flag(9)"`
	Unknow10                 types.VarInt32                                            `ins:"get_flag(10)"`
	Unknow11                 types.ListPadding[types.VarInt32, types.VarUInt32]        `ins:"get_flag(11)"`
	Unknow12                 types.ListPadding[types.VarUInt32String, types.VarUInt32] `ins:"get_flag(12)"`
	Unknow13                 types.VarUInt32String                                     `ins:"get_flag(13)"`
	Unknow14                 types.VarUInt32String                                     `ins:"get_flag(14)"`
	Unknow15                 int64                                                     `ins:"get_flag(15)"`
	Unknow16                 SubScreenSub3Item                                         `ins:"get_flag(16)"`
	Unknow18                 SubScreenSub3Item                                         `ins:"get_flag(18)"`
	Unknow20                 types.VarInt32                                            `ins:"get_flag(20)"`
	Unknow21                 types.VarInt32                                            `ins:"get_flag(21)"`
	Unknow25                 types.VarInt32                                            `ins:"get_flag(25)"`
	Unknow26                 types.VarUInt32String                                     `ins:"get_flag(26)"`
	Unknow27                 types.VarInt32                                            `ins:"get_flag(27)"`
	SubmitDataType           byte                                                      `ins:"get_flag(28)"`
	Unknow29                 types.VarInt32                                            `ins:"get_flag(29)"`
	Unknow30                 types.VarInt32                                            `ins:"get_flag(30)"`
	Unknow31                 byte                                                      `ins:"get_flag(31)"`
	Unknow32                 types.VarInt32                                            `ins:"get_flag(32)"`
	Unknow33                 byte                                                      `ins:"get_flag(33)"`
	Unknow34                 types.VarInt32                                            `ins:"get_flag(34)"`
	Unknow35                 types.VarInt32                                            `ins:"get_flag(35)"`
	Unknow36                 types.VarUInt32String                                     `ins:"get_flag(36)"`
	Unknow37                 types.VarUInt32String                                     `ins:"get_flag(37)"`
	Unknow38                 byte                                                      `ins:"get_flag(38)"`
	Unknow39                 types.VarInt32                                            `ins:"get_flag(39)"`
	Unknow40                 byte                                                      `ins:"get_flag(40)"`
	Unknow41                 types.VarInt32                                            `ins:"get_flag(41)"`
	Unknow42                 types.VarInt32                                            `ins:"get_flag(42)"`
	DefaultText              types.VarUInt32String                                     `ins:"get_flag(43)"`
	Unknow44                 types.VarUInt32String                                     `ins:"get_flag(44)"`
	Unknow45                 types.VarInt32                                            `ins:"get_flag(45)"`
	Unknow46                 types.VarInt32                                            `ins:"get_flag(46)"`
	Unknow47                 types.VarInt32                                            `ins:"get_flag(47)"`
	Unknow48                 types.VarInt32                                            `ins:"get_flag(48)"`
	Unknow49                 types.VarUInt32String                                     `ins:"get_flag(49)"`
	Unknow50                 byte                                                      `ins:"get_flag(50)"`
	Unknow55                 types.VarInt32                                            `ins:"get_flag(55)"`
	Unknow56                 byte                                                      `ins:"get_flag(56)"`
	Unknow57                 byte                                                      `ins:"get_flag(57)"`
	Unknow58                 byte                                                      `ins:"get_flag(58)"`
	OnLostFocusScreenCmdCode types.VarInt32                                            `ins:"get_flag(59)"`
	OnFocusedScreenCmdCode2  types.VarInt32                                            `ins:"get_flag(60)"`
	Unknow61                 types.VarInt32                                            `ins:"get_flag(61)"`
	Unknow62                 types.VarInt32                                            `ins:"get_flag(62)"`
	Unknow63                 types.VarInt32                                            `ins:"get_flag(63)"`
	OnFocusedScreenCmdCode1  types.VarInt32                                            `ins:"get_flag(64)"`
	Unknow65                 types.VarInt32                                            `ins:"get_flag(65)"`
	Unknow66                 byte                                                      `ins:"get_flag(66)"`
	Unknow67                 byte                                                      `ins:"get_flag(67)"`
	Unknow68                 types.VarInt32                                            `ins:"get_flag(68)"`
	Unknow70                 byte                                                      `ins:"get_flag(70)"`
	Unknow73                 types.ListPadding[types.VarInt32, types.VarUInt32]        `ins:"get_flag(73)"`
	Unknow74                 types.ListPadding[types.VarInt32, types.VarUInt32]        `ins:"get_flag(74)"`
	Unknow75                 types.ListPadding[types.VarInt32, types.VarUInt32]        `ins:"get_flag(75)"`
	Unknow76                 types.ListPadding[types.VarInt32, types.VarUInt32]        `ins:"get_flag(76)"`
	Unknow77                 byte                                                      `ins:"get_flag(77)"`
	Unknow78                 types.VarInt32                                            `ins:"get_flag(78)"`
	Unknow79                 types.ListPadding[types.VarInt32, types.VarUInt32]        `ins:"get_flag(79)"`
	Unknow80                 types.ListPadding[types.VarInt32, types.VarUInt32]        `ins:"get_flag(80)"`
	Unknow81                 types.ListPadding[types.VarInt32, types.VarUInt32]        `ins:"get_flag(81)"`
	AllowInputRule           types.ListPadding[types.VarUInt32String, types.VarInt32]  `ins:"get_flag(82)"`
	Unknow84                 types.VarInt32                                            `ins:"get_flag(84)"`
	Unknow85                 SubScreenSub3Custom                                       `ins:"get_flag(85)"`
	Unknow117                types.VarUInt32String                                     `ins:"get_flag(117)"`
	Unknow118                types.VarInt32                                            `ins:"get_flag(118)"`
	Unknow120                types.VarUInt32String                                     `ins:"get_flag(120)"`
	Unknow127                types.VarUInt32String                                     `ins:"get_flag(127)"`
	Unknow129                types.VarInt32                                            `ins:"get_flag(129)"`
	MInputSubmitFlag         types.VarInt32                                            `ins:"get_flag(131)"`
	Unknow134                types.VarInt32                                            `ins:"get_flag(134)"`
	Unknow136                byte                                                      `ins:"get_flag(136)"`
}

// text input view

type SubScreen3 struct {
	ScreenImpl     `ins:"false" json:"-"`
	SubScreen1     `json:"sub_screen_1"`
	SubScreen3Body `json:"sub_screen_3_body"`
}

func (this *SubScreen3) ReadChange(from io.BufferReader, changeType int, parentScreen *SubScreen) {
	log.Debug("ReadChange SubScreen3 type: %d, id: %d, offset: %d", changeType, parentScreen.GetScreenId(), from.Offset())
	this.SubScreen1.ReadChange(from, changeType, parentScreen)
}

func (this *SubScreen3) UpdateScreen(from io.BufferReader) {
	if this.ScreenImpl != nil {
		this.ScreenImpl.UpdateScreen(from)
		return
	}
	this.Read(from)
}

func (this *SubScreen3) Write(to io.BufferWriter) {

}

func (this *SubScreen3) Read(from io.BufferReader) {
	types.ReadMsg(from, &this.SubScreen1)
	types.ReadMsg(from, &this.SubScreen3Body)
	this.InitImpl()
}

func (this *SubScreen3) InitImpl() {
	this.SubScreen1.ScreenImpl = this
}

func (this *SubScreen3) GetSubScreen() *SubScreenArray {
	return nil
}

func (this *SubScreen3) GetDisplayActionScreenCmdCode() int32 {
	return this.DisplayActionScreenCmdCode.Value
}

func (this *SubScreen3) GetYNFlag() bool {
	return this.BitByteFlags.GetFlags(122)
}

func (this *SubScreen3) IsEncrypt() bool {
	return this.BitByteFlags.GetFlags(24)
}
