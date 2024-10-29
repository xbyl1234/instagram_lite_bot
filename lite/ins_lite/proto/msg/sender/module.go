package sender

import "CentralizedControl/ins_lite/proto/types"

type SendStorageHeaders struct {
	IntKey    int32
	StringKey string
	Value     string
}

type SendSubmitData struct {
	Type byte
	Data string
}

type SendActionUnUsedData struct {
	Unknow1 int64
	Unknow2 int16
	Unknow3 int16
	Unknow4 byte
}

type TrackingState struct {
	Flag        types.VarUInt32
	MarkerId    int32          `ins:"(Flag & 4) == 4"`
	InstanceKey types.VarInt32 `ins:"(Flag & 8) == 8"`
}

func CreateTrackingState(hasMarkerId bool, markerId int32, hasInstanceKey bool, instanceKey int32, isMarkerOn bool, isTracing bool) *TrackingState {
	result := &TrackingState{}
	var flag int32 = 0
	if isMarkerOn {
		flag |= 1
	}
	if isTracing {
		flag |= 2
	}
	//markerId
	if hasMarkerId {
		flag |= 4
		result.MarkerId = markerId
	}
	//instanceKey
	if hasInstanceKey {
		flag |= 8
		result.InstanceKey.Set(int64(instanceKey))
	}
	result.Flag.Set(int64(flag))
	return result
}

func (this *TrackingState) IsDisableTracking() bool {
	return this.Flag.Value == 0
}

type ActionMsg struct {
	FromScreenId              int32
	ToScreenId                int32
	ResourceId                int32
	LikeActionResourceId      int16
	RespMsgData               types.ListValue[SendSubmitData, int16]
	Flag1                     types.VarUInt32
	UnusedData                types.ListValue[SendActionUnUsedData, int16] `ins:"(Flag1 & 1) == 1"`
	LikeZeroOrScreenCmdResult int32                                        `ins:"(Flag1 & 8) == 8"`
	InstanceKey               int32
	Const                     int32
	Time                      int64
	TrackingState             TrackingState `ins:"Flag1 & 64 != 0"`
}

type PassivityActionMsg struct {
	ActionMsg
}

type InitiativeActionMsg struct {
	ActionMsg
}

type BrowserAboutItem1 struct {
	Key   int32
	Value types.ListValue[int32, types.VarUInt32]
}

type BrowserAction struct {
	ScreenId                   int32
	Time                       int64
	Const0                     int32
	Flags                      types.VarUInt32
	Unknow5                    types.VarUInt32 `ins:"(Flags & 4) != 0"`
	Navigation                 byte            `ins:"(Flags & 128) != 0"`
	ConstFlag                  types.VarUInt32
	HasScreenDecodeBodyItem73  byte                                                `ins:"(Flags & 1024) != 0"`
	ScreenDecodeBodyItem73Data types.ListValue[BrowserAboutItem1, types.VarUInt32] `ins:"(Flags & 1024) != 0 && HasScreenDecodeBodyItem73==1"`
	Const0_2                   int32                                               `ins:"(Flags & 1024) != 0 && HasScreenDecodeBodyItem73==1"`
	IsBackground               byte
	GenDeviceTimeId            string `ins:"(Flags & 4096) != 0"`
	SomeConfig                 string `ins:"(Flags & 8192) != 0"`
	BloksScreenName            string `ins:"(Flags & 16384) != 0"`
}

type FetchImage struct {
	Const1   byte
	ImageId  uint64
	Part     uint32
	Const1_2 byte
	LikeType int16
}

type DownloadableResources struct {
	Resources types.MapValue[string, byte, int16]
}

type DownloadVideoItem struct {
	VideoId          string
	Url              string
	DownloadProgress int64
}

type DownloadVideo struct {
	Video types.ListValue[DownloadVideoItem, int32]
}

type ClientEventLog struct {
	Type1 int16
	Type2 int16
	Log   string
	Index types.VarUInt32
}

type AppModuleDownloadItem struct {
	Key      string
	Padding1 int32
	Padding2 int16
	ConstStr string
	Padding3 byte
	Id       int32
}

type AppModuleDownload struct {
	Unknow1 types.VarUInt32
	Unknow2 byte
	Unknow3 byte
	Value   types.ListValue[AppModuleDownloadItem, int16]
}

type SendPkgInfo struct {
	PkgInfo types.MapValue[string, byte, int32]
}

type SendLastSomeEventInterval struct {
	Interval int64
}

const (
	ChangeType_STABLE                = 1
	ChangeType_UNSTABLE              = 2
	ChangeType_STABLE_AFTER_UNSTABLE = 3
)
const (
	ActiveNetworkMetered_Unknow = 1
	ActiveNetworkMetered_True   = 2
	ActiveNetworkMetered_False  = 3
)

type NetworkTypeChangeReporter struct {
	ChangeType             byte
	NowTime                int64
	NetworkType            string `ins:"ChangeType != 2"`
	NetworkSubType         string `ins:"ChangeType != 2"`
	TimeInterval           int64  `ins:"ChangeType != 2 && ChangeType == 3"`
	Unkonw2                int16  `ins:"ChangeType != 2 && ChangeType == 3"` //count?
	IsActiveNetworkMetered byte   `ins:"ChangeType != 2"`
}

var (
	QualityType_VERY_POOR = 1
	QualityType_POOR      = 2
	QualityType_MODERATE  = 3
	QualityType_GOOD      = 4
	QualityType_EXCELLENT = 5
	QualityType_UNKNOWN   = 0
	QualityType           = []int{
		QualityType_VERY_POOR,
		QualityType_POOR,
		QualityType_MODERATE,
		QualityType_GOOD,
		QualityType_EXCELLENT,
		QualityType_UNKNOWN,
	}
)

type ConnBandwidthQuality struct {
	QualityType byte
	Bandwidth   int32
}

type LoggedUserIdChange struct {
	Unknow types.ListValue[string, int32]
}

const (
	GoogleOAuthToken_OK                                     = 0
	GoogleOAuthToken_GET_ACCOUNTS_PERMISSION_NOT_AVAILABLE  = 1
	GoogleOAuthToken_READ_CONTACTS_PERMISSION_NOT_AVAILABLE = 2
	GoogleOAuthToken_GET_ACCOUNT_MANAGER_FAILED             = 3
	GoogleOAuthToken_NO_ACCOUNT_IN_DEVICE                   = 4
	GoogleOAuthToken_EXCEPTION                              = 5
)

type GoogleOAuthTokenItem struct {
	Account string
	Status  byte
	Cookies string
}

type GoogleOAuthToken struct {
	Status      byte
	FailedError string                                       `ins:"Status == 0"`
	Token       types.ListValue[GoogleOAuthTokenItem, int16] `ins:"Status != 0"`
}

type LikeSendImageRecv struct {
	Unknow1 types.ListValue[int32, int16]
	Unknow2 byte
}

type FirebaseInstanceId struct {
	FcmRegistrationId    string
	FbLitePhoneIdStoreId string
	FMC                  string
	Key                  string
	KeyId                string
	Algorithm            string
	ConstNull            string
}

type ActivityResumeOrStop struct {
	Unknow1 byte
	Unknow2 byte
	Unknow3 string
	Unknow4 int16
	NowTime int64
}

type FetchBloks struct {
	Unknow1 types.VarUInt32
	Unknow2 types.VarUInt32
	Unknow3 string
	Unknow4 types.MapValue[string, string, types.VarUInt32]
	Unknow5 types.VarUInt32
	Unknow6 types.VarUInt32
	Unknow7 int32 `ins:"(Unknow6 & 4) == 4"`
	Unknow8 int32 `ins:"(Unknow6 & 8) == 8"`
}

type SafetyNetData struct {
	Const0 byte
	Const1 byte
	Unknow string
}

type AdvertiserId struct {
	Id    string
	Flags byte
}

type TextMsg struct {
	Unknow1 string
	Unknow2 string
}

const (
	DEGRADED  = 0
	POOR      = 1
	MODERATE  = 2
	GOOD      = 3
	EXCELLENT = 4
	UNKNOWN   = 5
)

type ImageQuality struct {
	Unknow1 string
	Unknow2 int32
}

type UnknowItem1 struct {
	Unknow1 byte
	Unknow2 string
}

type UnknowItem2 struct {
	Unknow1 int64
	Unknow2 int16
	Unknow3 int16
	Unknow4 byte
}

type DownLoadModule struct {
	Unknow1 int32
	Unknow2 types.ListValue[string, int32]
	Unknow3 types.VarUInt32
}

type SendMcLiteConfigsAndParamsList struct {
	FileData types.ListValue[byte, types.VarUInt32]
}

type NetworkInfo struct {
	Const1_0           byte
	OldScreenId1       int64
	TimeInterval1      int32
	Flag               byte
	Const2_5           int32
	Const3_0           int32
	QualityType        int32
	Const4_1           int32
	DecodeBodyDataSize int32
	Const5_2           int32
	WifiIsConnected    int32 //true:3, false:1
	Const6_3           int32
	PhoneType          int32
	Const7_4           int32
	NetworkSubType     int32
	TimeInterval2      int64
	NewScreenId        int64
	OldScreenId2       int32
}

type MsgCodeUnknow3 struct {
	Unknow1 int32
	Unknow2 int64
	Unknow3 int32
	Flag    types.VarUInt32
	Unknow4 int32 `ins:"(Flag & 1) == 1"`
	Unknow5 byte  `ins:"(Flag & 8) != 0"`
	Unknow6 TrackingState
}

type PreSendPhoneBook struct {
	Idx               int16 //Unknow1
	SendPhoneBookSize int32 //remain size?
	Const0            int64
	Unknow4           string `ins:"Unknow4 != ''"`
}

type PhonebookLine struct {
	Type   byte
	Const0 byte
	Value  string
}

type PhonebookItem struct {
	Count         byte
	FullName105   PhonebookLine
	SecondName119 PhonebookLine
	FirstName120  PhonebookLine
	Email103      []PhonebookLine
	Phone115      []PhonebookLine
}

func (this *PhonebookItem) PutFullName105(v string) {
	this.FullName105 = PhonebookLine{
		Type:   105,
		Const0: 0,
		Value:  v,
	}
}

func (this *PhonebookItem) PutSecondName119(v string) {
	this.SecondName119 = PhonebookLine{
		Type:   119,
		Const0: 0,
		Value:  v,
	}
}

func (this *PhonebookItem) PutFirstName120(v string) {
	this.FirstName120 = PhonebookLine{
		Type:   120,
		Const0: 0,
		Value:  v,
	}
}

func (this *PhonebookItem) PutEmail103(v string) {
	this.Email103 = append(this.Email103, PhonebookLine{
		Type:   103,
		Const0: 0,
		Value:  v,
	})
}

func (this *PhonebookItem) PutPhone115(v string) {
	this.Phone115 = append(this.Phone115, PhonebookLine{
		Type:   115,
		Const0: 0,
		Value:  v,
	})
}

type SendPhonebook struct {
	Idx       int16 //Unknow1
	Const0    byte
	Phonebook string
	Count     int16
	Item      types.ListValue[PhonebookItem, int16]
}

type PermResult struct {
	Idx                             int32
	Unknow1                         byte
	IsAllow                         byte
	DontShowRequestPermissionDialog byte
}

const (
	MediaTypePhoto   = 0
	MediaTypeVideo   = 1
	MediaTypeGif     = 2
	MediaTypeSticker = 3
	MediaTypeUnknown = 4
	MediaTypeFile    = 5
	MediaTypeXma     = 6
	MediaTypeAudio   = 11
)

type UploadMediaInfo struct {
	Unknow1   string
	MediaType byte
	Unknow2   int32
}

type ActivityResumed struct {
	GenDeviceTimeId string
}

type SendPhoneIdMsg struct {
	PhoneId string
}

type PropertiesItem struct {
	Name     string
	HasValue byte
	Value    string `ins:"HasValue != 0"`
}

type SendSystemPropertiesMsg struct {
	Properties types.ListValue[PropertiesItem, int16]
}

type ReConnect struct {
}

type SendInstallReferrer struct {
	InstallReferrer string
}

type DurationTracking struct {
	trackingId int32
	Duration   int32
	ImageState int32
	Unknow1    byte
	Unknow2    int64
}
