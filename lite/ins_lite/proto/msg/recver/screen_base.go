package recver

import (
	"CentralizedControl/common/log"
	"CentralizedControl/ins_lite/proto/io"
	"CentralizedControl/ins_lite/proto/types"
)

type ScreenBaseItem9 struct {
	BitByteFlags types.BitByteFlags `ins_init:"BitByteFlags(10)"`
	Unknow0      types.VarInt32     `ins:"get_flag(0)"`
	Unknow1      types.VarInt32     `ins:"get_flag(1)"`
	Unknow2      types.VarInt32     `ins:"get_flag(2)"`
	Unknow3      types.VarInt32     `ins:"get_flag(3)"`
	Unknow4      types.VarInt32     `ins:"get_flag(4)"`
	Unknow5      types.VarInt32     `ins:"get_flag(5)"`
	Unknow6      types.VarInt32     `ins:"get_flag(6)"`
	Unknow7      types.VarInt32     `ins:"get_flag(7)"`
	Unknow8      types.VarInt32     `ins:"get_flag(8)"`
	Unknow9      types.VarInt32     `ins:"get_flag(9)"`
}

type ScreenBaseItem14Item struct {
	Unknow1 types.VarInt32
	Unknow2 types.VarInt32
}

type ScreenBaseItem14 struct {
	Unknow1 types.ListValue[ScreenBaseItem14Item, types.VarUInt32]
	Unknow2 types.ListValue[types.VarInt32, types.VarUInt32]
}

type ScreenBaseItem21 struct {
	BitByteFlags types.BitByteFlags `ins_init:"BitByteFlags(47)"`
	Unknow0      byte               `ins:"get_flag(0)"`
	Unknow1      float64            `ins:"get_flag(1)"`
	Unknow2      float64            `ins:"get_flag(2)"`
	Unknow3      float64            `ins:"get_flag(3)"`
	Unknow4      float64            `ins:"get_flag(4)"`
	Unknow5      float64            `ins:"get_flag(5)"`
	Unknow6      float64            `ins:"get_flag(6)"`
	Unknow7      float64            `ins:"get_flag(7)"`
	Unknow8      float64            `ins:"get_flag(8)"`
	Unknow9      float64            `ins:"get_flag(9)"`
	Unknow10     float64            `ins:"get_flag(10)"`
	Unknow11     float64            `ins:"get_flag(11)"`
	Unknow12     float64            `ins:"get_flag(12)"`
	Unknow13     float64            `ins:"get_flag(13)"`
	Unknow14     float64            `ins:"get_flag(14)"`
	Unknow15     float64            `ins:"get_flag(15)"`
	Unknow16     float64            `ins:"get_flag(16)"`
	Unknow17     float64            `ins:"get_flag(17)"`
	Unknow18     float64            `ins:"get_flag(18)"`
	Unknow19     float64            `ins:"get_flag(19)"`
	Unknow20     float64            `ins:"get_flag(20)"`
	Unknow21     float64            `ins:"get_flag(21)"`
	Unknow22     float64            `ins:"get_flag(22)"`
	Unknow23     float64            `ins:"get_flag(23)"`
	Unknow24     float64            `ins:"get_flag(24)"`
	Unknow25     float64            `ins:"get_flag(25)"`
	Unknow26     float64            `ins:"get_flag(26)"`
	Unknow27     float64            `ins:"get_flag(27)"`
	Unknow28     float64            `ins:"get_flag(28)"`
	Unknow29     float64            `ins:"get_flag(29)"`
	Unknow30     float64            `ins:"get_flag(30)"`
	Unknow31     float64            `ins:"get_flag(31)"`
	Unknow32     float64            `ins:"get_flag(32)"`
	Unknow33     float64            `ins:"get_flag(33)"`
	Unknow34     float64            `ins:"get_flag(34)"`
	Unknow35     float64            `ins:"get_flag(35)"`
	Unknow36     float64            `ins:"get_flag(36)"`
	Unknow37     float64            `ins:"get_flag(37)"`
	Unknow38     float64            `ins:"get_flag(38)"`
	Unknow39     float64            `ins:"get_flag(39)"`
	Unknow40     float64            `ins:"get_flag(40)"`
	Unknow41     byte               `ins:"get_flag(41)"`
	Unknow42     float64            `ins:"get_flag(42)"`
	Unknow43     float64            `ins:"get_flag(43)"`
	Unknow44     float64            `ins:"get_flag(44)"`
	Unknow45     float64            `ins:"get_flag(45)"`
	Unknow46     types.VarInt32     `ins:"get_flag(46)"`
}

type ScreenBaseItem105 struct {
	BitByteFlags types.BitByteFlags `ins_init:"BitByteFlags(4)"`
	Unknow0      types.VarInt32     `ins:"get_flag(0)"`
	Unknow1      types.VarInt32     `ins:"get_flag(1)"`
	Unknow2      types.VarInt32     `ins:"get_flag(2)"`
	Unknow3      types.VarInt32     `ins:"get_flag(3)"`
}

type ScreenBaseItem150 struct {
	BitByteFlags types.BitByteFlags `ins_init:"BitByteFlags(23)"`
	Unknow2      float32            `ins:"get_flag(2)"`
	Unknow3      types.VarInt32     `ins:"get_flag(3)"`
	Unknow4      types.VarInt32     `ins:"get_flag(4)"`
	Unknow5      types.VarInt32     `ins:"get_flag(5)"`
	Unknow6      types.VarInt32     `ins:"get_flag(6)"`
	Unknow7      float32            `ins:"get_flag(7)"`
	Unknow8      float32            `ins:"get_flag(8)"`
	Unknow9      float32            `ins:"get_flag(9)"`
	Unknow10     float32            `ins:"get_flag(10)"`
	Unknow11     float32            `ins:"get_flag(11)"`
	Unknow12     float32            `ins:"get_flag(12)"`
	Unknow13     float32            `ins:"get_flag(13)"`
	Unknow14     float32            `ins:"get_flag(14)"`
	Unknow15     float32            `ins:"get_flag(15)"`
	Unknow16     float32            `ins:"get_flag(16)"`
	Unknow17     types.VarInt32     `ins:"get_flag(17)"`
	Unknow18     types.VarInt32     `ins:"get_flag(18)"`
	Unknow19     float32            `ins:"get_flag(19)"`
	Unknow20     float32            `ins:"get_flag(20)"`
	Unknow21     float32            `ins:"get_flag(21)"`
	Unknow22     float32            `ins:"get_flag(22)"`
}

// X.0kz
type ScreenBase struct {
	ScreenImpl            `ins:"false" json:"-"`
	BitByteFlags          types.BitByteFlags    `ins_init:"BitByteFlags(153)"`
	Unknow3               byte                  `ins:"get_flag(3)"`
	Unknow4               types.VarInt32        `ins:"get_flag(4)"`
	ScreenId              types.VarInt32        `ins:"get_flag(5)"`
	Unknow6               types.VarInt32        `ins:"get_flag(6)"`
	Unknow7               types.VarInt32        `ins:"get_flag(7)"`
	Unknow9               ScreenBaseItem9       `ins:"get_flag(9)"`
	Unknow11              byte                  `ins:"get_flag(11)"`
	Unknow12              types.VarInt32        `ins:"get_flag(12)"`
	Unknow13              types.VarInt32        `ins:"get_flag(13)"`
	Unknow14              ScreenBaseItem14      `ins:"get_flag(14)"`
	Unknow16              byte                  `ins:"get_flag(16)"`
	Unknow17              int64                 `ins:"get_flag(17)"`
	Unknow18              types.VarUInt32String `ins:"get_flag(18)"`
	Unknow20              byte                  `ins:"get_flag(20)"`
	Unknow21              ScreenBaseItem21      `ins:"get_flag(21)"`
	Unknow24              types.VarUInt32String `ins:"get_flag(24)"`
	Unknow26              types.VarInt32        `ins:"get_flag(26)"`
	Unknow28              types.VarInt32        `ins:"get_flag(28)"`
	Unknow35              float32               `ins:"get_flag(35)"`
	Unknow36              float32               `ins:"get_flag(36)"`
	Unknow38              float64               `ins:"get_flag(38)"`
	Unknow42              byte                  `ins:"get_flag(42)"`
	WindowType            types.VarUInt32String `ins:"get_flag(62)"`
	TitleEng              types.VarUInt32String `ins:"get_flag(63)"`
	TitleCh               types.VarUInt32String `ins:"get_flag(64)"`
	Unknow65              types.VarInt32        `ins:"get_flag(65)"`
	Unknow66              types.VarInt32        `ins:"get_flag(66)"`
	WindowId              types.VarUInt32String `ins:"get_flag(68)"`
	Unknow69              byte                  `ins:"get_flag(69)"`
	ClickRunScreenCmdCode types.VarInt32        `ins:"get_flag(70)"`
	Unknow71              types.VarInt32        `ins:"get_flag(71)"`
	Unknow72              types.VarInt32        `ins:"get_flag(72)"`
	Unknow73              types.VarInt32        `ins:"get_flag(73)"`
	Unknow74              types.VarInt32        `ins:"get_flag(74)"`
	Unknow75              types.VarInt32        `ins:"get_flag(75)"`
	Unknow76              types.VarInt32        `ins:"get_flag(76)"`
	Unknow77              types.VarInt32        `ins:"get_flag(77)"`
	Unknow78              types.VarInt32        `ins:"get_flag(78)"`
	Unknow79              types.VarInt32        `ins:"get_flag(79)"`
	Unknow80              types.VarInt32        `ins:"get_flag(80)"`
	Unknow81              byte                  `ins:"get_flag(81)"`
	Unknow82              types.VarInt32        `ins:"get_flag(82)"`
	Unknow84              types.VarInt32        `ins:"get_flag(84)"`
	Unknow85              types.VarInt32        `ins:"get_flag(85)"`
	Unknow86              types.VarInt32        `ins:"get_flag(86)"`
	Unknow87              types.VarInt32        `ins:"get_flag(87)"`
	Unknow88              byte                  `ins:"get_flag(88)"`
	LikeActionResourceId  types.VarInt32        `ins:"get_flag(89)"`
	Unknow90              int64                 `ins:"get_flag(90)"`
	Unknow91              types.VarUInt32String `ins:"get_flag(91)"`
	Unknow92              byte                  `ins:"get_flag(92)"`
	Unknow93              byte                  `ins:"get_flag(93)"`
	Unknow94              byte                  `ins:"get_flag(94)"`
	ImageId               uint64                `ins:"get_flag(95)"`
	Unknow96              byte                  `ins:"get_flag(96)"`
	Unknow97              types.VarInt32        `ins:"get_flag(97)"`
	Unknow98              byte                  `ins:"get_flag(98)"`
	Unknow99              byte                  `ins:"get_flag(99)"`
	Unknow100             byte                  `ins:"get_flag(100)"`
	Unknow101             byte                  `ins:"get_flag(101)"`
	Unknow102             types.VarInt32        `ins:"get_flag(102)"`
	Unknow103             types.VarInt32        `ins:"get_flag(103)"`
	Unknow104             types.VarInt32        `ins:"get_flag(104)"`
	Unknow105             ScreenBaseItem105     `ins:"get_flag(105)"`
	Unknow107             byte                  `ins:"get_flag(107)"`
	Unknow108             types.VarInt32        `ins:"get_flag(108)"`
	Unknow109             byte                  `ins:"get_flag(109)"`
	Unknow110             types.VarInt32        `ins:"get_flag(110)"`
	Unknow111             byte                  `ins:"get_flag(111)"`
	Unknow112             types.VarInt32        `ins:"get_flag(112)"`
	Unknow113             byte                  `ins:"get_flag(113)"`
	Unknow114             types.VarInt32        `ins:"get_flag(114)"`
	Unknow115             byte                  `ins:"get_flag(115)"`
	Unknow116             types.VarInt32        `ins:"get_flag(116)"`
	Unknow117             byte                  `ins:"get_flag(117)"`
	Unknow118             byte                  `ins:"get_flag(118)"`
	Height                types.VarInt32        `ins:"get_flag(119)"`
	Width                 types.VarInt32        `ins:"get_flag(120)"`
	Unknow121             types.VarInt32        `ins:"get_flag(121)"`
	Unknow122             types.VarInt32        `ins:"get_flag(122)"`
	Unknow123             byte                  `ins:"get_flag(123)"`
	Unknow125             types.VarUInt32String `ins:"get_flag(125)"`
	SomeScreenCmdIdx      types.VarInt32        `ins:"get_flag(126)"`
	Unknow127             types.VarInt32        `ins:"get_flag(127)"`
	Unknow128             types.VarInt32        `ins:"get_flag(128)"`
	Unknow129             types.VarInt32        `ins:"get_flag(129)"`
	Unknow130             types.VarInt32        `ins:"get_flag(130)"`
	Unknow131             types.VarInt32        `ins:"get_flag(131)"`
	Unknow132             types.VarInt32        `ins:"get_flag(132)"`
	Unknow133             byte                  `ins:"get_flag(133)"`
	Unknow135             types.VarInt32        `ins:"get_flag(135)"`
	Unknow136             byte                  `ins:"get_flag(136)"`
	Unknow137             types.VarInt32        `ins:"get_flag(137)"`
	Unknow140             types.VarInt32        `ins:"get_flag(140)"`
	Unknow141             types.VarInt32        `ins:"get_flag(141)"`
	Unknow142             types.VarInt32        `ins:"get_flag(142)"`
	Unknow143             types.VarInt32        `ins:"get_flag(143)"`
	Unknow144             types.VarInt32        `ins:"get_flag(144)"`
	Unknow145             types.VarInt32        `ins:"get_flag(145)"`
	Unknow147             types.VarUInt32String `ins:"get_flag(147)"`
	Unknow148             types.VarUInt32String `ins:"get_flag(148)"`
	Unknow149             types.VarUInt32String `ins:"get_flag(149)"`
	Unknow150             ScreenBaseItem150     `ins:"get_flag(150)"`
}

func (this *ScreenBase) ReadChange(from io.BufferReader, changeType int, parentScreen *SubScreen) {
	log.Debug("ReadChange ScreenBase type: %d, id: %d, offset: %d", changeType, parentScreen.GetScreenId(), from.Offset())
	switch changeType {
	case 1:
		panic("change type error")
	case 5:
		from.ReadShort()
	case 6:
		types.ReadMsg(from, this)
	default:
		this.ScreenImpl.UpdateScreen(from)
	}
}

func (this *ScreenBase) UpdateScreen(from io.BufferReader) {
	if this.ScreenImpl != nil {
		this.ScreenImpl.UpdateScreen(from)
	}
}

func (this *ScreenBase) InitImpl() {

}

func (this *ScreenBase) GetSubScreen() *SubScreenArray {
	return nil
}

func (this *ScreenBase) GetDisplayActionScreenCmdCode() int32 {
	return 0
}

func (this *ScreenBase) GetIsLikeResIdChildFlag() bool {
	return this.BitByteFlags.GetFlags(48)
}
