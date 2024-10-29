package recver

import (
	"CentralizedControl/ins_lite/proto/io"
	"CentralizedControl/ins_lite/proto/types"
)

type SubScreen28BodyUnknow3 struct {
	BitByteFlags types.BitByteFlags    `ins_init:"BitByteFlags(5)"`
	Unknow0      types.VarUInt32String `ins:"get_flag(0)"`
	Unknow1      types.VarInt32        `ins:"get_flag(1)"`
	Unknow2      types.VarInt32        `ins:"get_flag(2)"`
	Unknow3      types.VarInt32        `ins:"get_flag(3)"`
	Unknow4      byte                  `ins:"get_flag(4)"`
}

type SubScreen28BodyUnknow2 struct {
	BitByteFlags types.BitByteFlags `ins_init:"BitByteFlags(26)"`
	Unknow1      types.VarInt32     `ins:"get_flag(0)"`
	Unknow2      types.VarInt32     `ins:"get_flag(1)"`
	Unknow3      types.VarInt32     `ins:"get_flag(2)"`
	Unknow4      types.VarInt32     `ins:"get_flag(3)"`
	Unknow5      types.VarInt32     `ins:"get_flag(4)"`
	Unknow6      types.VarInt32     `ins:"get_flag(5)"`
	Unknow7      types.VarInt32     `ins:"get_flag(6)"`
	Unknow8      types.VarInt32     `ins:"get_flag(7)"`
	Unknow9      types.VarInt32     `ins:"get_flag(8)"`
	Unknow10     types.VarInt32     `ins:"get_flag(9)"`
	Unknow11     types.VarInt32     `ins:"get_flag(10)"`
	Unknow12     types.VarInt32     `ins:"get_flag(11)"`
	Unknow13     types.VarInt32     `ins:"get_flag(12)"`
	Unknow14     types.VarInt32     `ins:"get_flag(13)"`
	Unknow15     types.VarInt32     `ins:"get_flag(14)"`
	Unknow16     types.VarInt32     `ins:"get_flag(16)"`
	Unknow17     types.VarInt32     `ins:"get_flag(17)"`
	Unknow18     types.VarInt32     `ins:"get_flag(18)"`
	Unknow19     types.VarInt32     `ins:"get_flag(19)"`
	Unknow20     types.VarInt32     `ins:"get_flag(20)"`
	Unknow21     types.VarInt32     `ins:"get_flag(22)"`
	Unknow22     types.VarInt32     `ins:"get_flag(23)"`
	Unknow23     types.VarInt32     `ins:"get_flag(24)"`
	Unknow24     types.VarInt32     `ins:"get_flag(25)"`
}

type SubScreen28BodyUnknow1Item struct {
	BitByteFlags types.BitByteFlags    `ins_init:"BitByteFlags(6)"`
	Unknow1      types.VarUInt32String `ins:"get_flag(0)"`
	Unknow2      types.VarInt32        `ins:"get_flag(1)"`
	Unknow3      types.VarUInt32String `ins:"get_flag(5)"`
}

type SubScreen28BodyUnknow1 struct {
	Count    uint32
	FlagByte []byte
	Flags1   []bool
	Flags2   []bool
	Data     []*SubScreen28BodyUnknow1Item
}

func (this *SubScreen28BodyUnknow1) Write(to io.BufferWriter) {

}

func (this *SubScreen28BodyUnknow1) Read(from io.BufferReader) {
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
	this.Data = make([]*SubScreen28BodyUnknow1Item, 0)
	for i := 0; i < int(this.Count); i++ {
		if this.Flags1[i] {
			item := &SubScreen28BodyUnknow1Item{}
			types.ReadMsg(from, item)
			this.Data = append(this.Data, item)
		}
	}
}

type SubScreen28Body struct {
	BitByteFlags types.BitByteFlags     `ins_init:"BitByteFlags(12)"`
	Unknow1      SubScreen28BodyUnknow1 `ins:"get_flag(0)"`
	Unknow2      SubScreen28BodyUnknow2 `ins:"get_flag(1)"`
	Unknow3      types.VarInt32         `ins:"get_flag(4)"`
	Unknow4      types.VarInt32         `ins:"get_flag(6)"`
	Unknow5      types.VarUInt32String  `ins:"get_flag(7)"`
	Unknow6      types.VarInt32         `ins:"get_flag(9)"`
	Unknow7      SubScreen28BodyUnknow3 `ins:"get_flag(10)"`
}

type SubScreen28 struct {
	ScreenImpl      `ins:"false" json:"-"`
	ScreenBase      `json:"screen_base"`
	SubScreen28Body `json:"sub_screen28_body"`
}

func (this *SubScreen28) ReadChange(from io.BufferReader, changeType int, parentScreen *SubScreen) {
	panic("SubScreen28 ReadChange")
}

func (this *SubScreen28) UpdateScreen(from io.BufferReader) {
	panic("SubScreen28 UpdateScreen")
}

func (this *SubScreen28) Write(to io.BufferWriter) {

}

func (this *SubScreen28) Read(from io.BufferReader) {
	types.ReadMsg(from, &this.ScreenBase)
	types.ReadMsg(from, &this.SubScreen28Body)
	this.InitImpl()
}

func (this *SubScreen28) InitImpl() {
	this.ScreenImpl = this
}

func (this *SubScreen28) GetSubScreen() *SubScreenArray {
	return nil
}

func (this *SubScreen28) GetDisplayActionScreenCmdCode() int32 {
	return 0
}
