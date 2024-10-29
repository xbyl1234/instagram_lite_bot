package recver

import (
	"CentralizedControl/common/log"
	"CentralizedControl/ins_lite/proto/io"
	"CentralizedControl/ins_lite/proto/types"
)

type UnknowScreen0PyItem18 struct {
	BitByteFlags types.BitByteFlags `ins_init:"BitByteFlags(5)"`
	Unknow0      byte               `ins:"get_flag(0)"`
	Unknow1      byte               `ins:"get_flag(1)"`
	Unknow2      byte               `ins:"get_flag(2)"`
	Unknow3      byte               `ins:"get_flag(3)"`
	Unknow4      byte               `ins:"get_flag(4)"`
}

// 0QN-A03: X.0Py
type SubScreen2Body struct {
	BitByteFlags types.BitByteFlags    `ins_init:"BitByteFlags(62)"`
	Unknow1      types.VarUInt32String `ins:"get_flag(1)"`
	Unknow3      types.VarInt32        `ins:"get_flag(3)"`
	Unknow4      types.VarInt32        `ins:"get_flag(4)"`
	Unknow5      types.VarInt32        `ins:"get_flag(5)"`
	Unknow6      types.VarInt32        `ins:"get_flag(6)"`
	Unknow7      types.VarInt32        `ins:"get_flag(7)"`
	Unknow8      types.VarInt32        `ins:"get_flag(8)"`
	Unknow11     uint16                `ins:"get_flag(11)"`
	Unknow12     uint16                `ins:"get_flag(12)"`
	Unknow15     SubScreen             `ins:"get_flag(15)"`
	Unknow18     UnknowScreen0PyItem18 `ins:"get_flag(18)"`
	Unknow37     types.VarInt32        `ins:"get_flag(37)"`
	Unknow38     types.VarInt32        `ins:"get_flag(38)"`
	Unknow39     types.VarInt32        `ins:"get_flag(39)"`
	Unknow40     byte                  `ins:"get_flag(40)"`
	Unknow41     byte                  `ins:"get_flag(41)"`
	Unknow42     types.VarInt32        `ins:"get_flag(42)"`
	Unknow43     types.VarInt32        `ins:"get_flag(43)"`
	Unknow44     byte                  `ins:"get_flag(44)"`
	Unknow45     types.VarInt32        `ins:"get_flag(45)"`
	Unknow46     types.VarInt32        `ins:"get_flag(46)"`
	Unknow47     types.VarInt32        `ins:"get_flag(47)"`
	Unknow48     types.VarInt32        `ins:"get_flag(48)"`
	Unknow49     byte                  `ins:"get_flag(49)"`
	Unknow50     uint16                `ins:"get_flag(50)"`
	Unknow52     types.VarInt32        `ins:"get_flag(52)"`
	Unknow53     types.VarInt32        `ins:"get_flag(53)"`
	Unknow56     types.VarInt32        `ins:"get_flag(56)"`
	Unknow57     types.VarInt32        `ins:"get_flag(57)"`
	Unknow58     types.VarInt32        `ins:"get_flag(58)"`
	Unknow60     types.VarInt32        `ins:"get_flag(60)"`
	Unknow61     float32               `ins:"get_flag(61)"`
}

type SubScreen2 struct {
	ScreenImpl     `ins:"false" json:"-"`
	ScreenBase     `json:"screen_base"`
	SubScreen2Body `json:"sub_screen_2_body"`
	SubScreenArray `json:"sub_screen_array"` //窗口里面的控件
}

var loop = 0

func (this *SubScreen2) ReadChange(from io.BufferReader, changeType int, parentScreen *SubScreen) {
	loop++
	log.Debug("ReadChange SubScreen2 type: %d, id: %d, offset: %d", changeType, parentScreen.GetScreenId(), from.Offset())
	if changeType != 1 {
		switch changeType {
		case 2, 3:
			types.ReadMsg(from, &this.ScreenBase)
			types.ReadMsg(from, &this.SubScreen2Body)
		case 5:
			from.ReadShort()
		case 6:
			types.ReadMsg(from, &this.ScreenBase)
			types.ReadMsg(from, &this.SubScreen2Body)
		}
		count := from.ReadShort()
		type insertPare struct {
			idx    int
			screen *SubScreen
		}
		deleteIdx := make([]int, 0)
		insertScreen := make([]insertPare, 0)
		subScreen := parentScreen.GetSubScreenRaw()
		for i := int16(0); i < count; i++ {
			deep := from.ReadShort()
			subChangeType := from.ReadByte()
			log.Debug("read ChangeInfo  parent: %d, changeType: %d, deep: %d, idx: %d, sub count: %d, offset: %d",
				parentScreen.Type, subChangeType, deep, deep, parentScreen.GetSubScreenRaw().Count(), from.Offset())

			switch subChangeType {
			case 1:
				deleteIdx = append(deleteIdx, int(deep))
				log.Debug("delete loop %d, index: %d", loop, deep)
			case 2:
				screenType := from.ReadByte()
				newScreen := CreateNewSubScreen(int(screenType))
				ReadSubScreenChange(from, newScreen, 2)
				insertScreen = append(insertScreen, insertPare{
					idx:    int(deep),
					screen: newScreen,
				})
				log.Debug("insert loop %d, index: %d, screen: %d, %s", loop, deep, screenType, Tree2String(DumpViewTree(newScreen), true))
			case 3, 4, 5, 6:
				screen := subScreen.Get(int(deep))
				ReadSubScreenChange(from, screen, int(subChangeType))
				log.Debug("update loop %d, index: %d", loop, deep)
			}
		}
		for _, item := range deleteIdx {
			subScreen.SubScreen[item].Value = nil
		}
		for idx := len(insertScreen) - 1; idx >= 0; idx-- {
			subScreen.Insert(insertScreen[idx].idx, insertScreen[idx].screen)
		}
		subScreen.Delete()
	}
	loop--
}

func (this *SubScreen2) UpdateScreen(from io.BufferReader) {
	if this.ScreenImpl != nil {
		this.ScreenImpl.UpdateScreen(from)
		return
	}
	this.Read(from)
}

func (this *SubScreen2) Write(to io.BufferWriter) {

}

func (this *SubScreen2) Read(from io.BufferReader) {
	types.ReadMsg(from, &this.ScreenBase)
	types.ReadMsg(from, &this.SubScreen2Body)
	types.ReadMsg(from, &this.SubScreenArray)
	this.InitImpl()
}

func (this *SubScreen2) InitImpl() {
	this.ScreenBase.ScreenImpl = this
}

func (this *SubScreen2) GetSubScreen() *SubScreenArray {
	return this.SubScreenArray.Copy()
}

func (this *SubScreen2) GetDisplayActionScreenCmdCode() int32 {
	return 0
}
