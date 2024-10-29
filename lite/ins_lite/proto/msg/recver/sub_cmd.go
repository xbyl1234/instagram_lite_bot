package recver

import (
	"CentralizedControl/ins_lite/proto/io"
	"CentralizedControl/ins_lite/proto/types"
)

type RecvSubCmd struct {
	SubCode types.VarUInt32
	Flags   byte
	Unknow2 int32 `ins:"(Flags & 1) != 0"`
}

type GetStorageHeaders struct {
	Key    string
	IntKey int32
}

type RequestPerm struct {
	OnGrantedIdx  types.VarUInt32
	OnGrantedData types.ListValue[byte, int32]
	OnDeniedIdx   types.VarUInt32
	OnDeniedData  types.ListValue[byte, int32]
	PermType      []byte
}

func (this *RequestPerm) Write(to io.BufferWriter) {

}

func (this *RequestPerm) Read(from io.BufferReader) {
	types.ReadMsg(from, &this.OnGrantedIdx)
	types.ReadMsg(from, &this.OnGrantedData)
	types.ReadMsg(from, &this.OnDeniedIdx)
	types.ReadMsg(from, &this.OnDeniedData)
	this.PermType = from.ReadRemain()
}

type RemoveStorageKey struct {
	Key string
}

type LoggerSettingItem1 struct {
	int32
	int16
}
type LoggerSettingItem2 struct {
	int32
	string
}
type LoggerSettingItem4 struct {
	types.ListValue[LoggerSettings, int32]
}
type LoggerSettingItem6 struct {
	Unknow1 int32
	Flag    int32
	Unknow2 int32 `ins:"Flag & 2 != 0"`
}

type LoggerSettingItem8 struct {
	Unknow1 int64
	Unknow2 byte
}
type LoggerSettingItem9 struct {
	Unknow1 int32
	Unknow2 types.ListValue[string, int32]
}

type LoggerSettingItem12 struct {
	Unknow1 int32
	Unknow2 int32
}
type LoggerSettingItem13 struct {
	Unknow1 int32
	Unknow2 string
	Unknow3 string
}
type LoggerSettingItem14 struct {
	Unknow1 int32
	Unknow2 string
	Unknow3 int64
}
type LoggerSettingItem15 struct {
	Unknow1 int32
	Unknow2 string
	Unknow3 bool
}
type LoggerSettingItem16 struct {
	Unknow1 int32
	Unknow2 string
	Unknow3 float64
}

type LoggerSettingItem struct {
	Unknow int32
}

type LoggerSettings struct {
	Flag    byte
	Unknow1 int                 `ins:"Flag == 7"`
	Unknow2 byte                `ins:"Flag == 7"`
	Item1   LoggerSettingItem1  `ins:"Flag == 1"`
	Item2   LoggerSettingItem2  `ins:"Flag == 2"`
	Item4   LoggerSettingItem4  `ins:"Flag == 4"`
	Item6   LoggerSettingItem6  `ins:"Flag == 6"`
	Item8   LoggerSettingItem8  `ins:"Flag == 8"`
	Item9   LoggerSettingItem9  `ins:"Flag == 9"`
	Item12  LoggerSettingItem12 `ins:"Flag == 3 || Flag == 12"`
	Item13  LoggerSettingItem13 `ins:"Flag == 13"`
	Item14  LoggerSettingItem14 `ins:"Flag == 14"`
	Item15  LoggerSettingItem15 `ins:"Flag == 15"`
	Item16  LoggerSettingItem16 `ins:"Flag == 16"`
	Item    LoggerSettingItem   `ins:"(Flag <1)|| (Flag >16)|| (Flag == 5)|| (Flag == 7)|| (Flag == 10)||(Flag == 11)"`
}

type RequireSafetyData struct {
	Nonce  types.ListValue[byte, types.VarUInt32]
	ApiKey string
}

type GetPhoneBook struct {
	UnUsed1 int32
	Idx     int16
	UnUsed2 int16
	Unknow2 byte
	Unknow3 int16
	UnUsed3 int16
	UnUsed4 types.ListValue[int32, int32]
}

type CheckCallingOrSelfPermission struct {
	PermType      byte
	SubCmdCodeIdx int16
	PkgIdx        int32
}

type CheckPermAndApplyForPerm struct {
	PermType              byte
	ReqPermType           byte  //0 requestPermissions, 1 startSettingActivity
	OnGrantedScreenCmdIdx int16 `ins:"ReqPermType==0||ReqPermType==1"`
	OnDeniedScreenCmdIdx  int16 `ins:"ReqPermType==0||ReqPermType==1"`
	PkgIdx                int16
	HasPkgIdx             bool
}

func (this *CheckPermAndApplyForPerm) Write(to io.BufferWriter) {

}

func (this *CheckPermAndApplyForPerm) Read(from io.BufferReader) {
	this.PermType = from.ReadByte()
	this.ReqPermType = from.ReadByte()
	if this.ReqPermType == 0 || this.ReqPermType == 1 {
		this.OnGrantedScreenCmdIdx = from.ReadShort()
		this.OnDeniedScreenCmdIdx = from.ReadShort()
	}
	if len(from.PeekRemain()) >= 2 {
		this.PkgIdx = from.ReadShort()
		this.HasPkgIdx = true
	} else {
		this.HasPkgIdx = false
	}
}

type GoogleOAuthToken struct {
	Scope string
}
