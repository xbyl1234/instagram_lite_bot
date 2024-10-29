package sender

import (
	"CentralizedControl/ins_lite/proto/io"
	"CentralizedControl/ins_lite/proto/types"
)

type MatePackageInfo struct {
	Unknow1 string
	Unknow2 string
	Unknow3 byte
}

type SettingItem struct {
	Type  types.VarUInt32
	Value types.AnyType
}

func (this *SettingItem) Write(to io.BufferWriter) {
	types.WriteMsg(to, &this.Type)
	types.WriteMsg(to, &this.Value)
}

func (this *SettingItem) Read(from io.BufferReader) {
	types.ReadMsg(from, &this.Type)
	this.Value.Type = GetSettingNameEnumTypeByValue1(int(this.Type.Value))
	types.ReadMsg(from, &this.Value)
}

type Setting struct {
	FirstLunch bool
	Settings   types.ListValue[SettingItem, types.VarUInt32]
}

func (this *Setting) Write(to io.BufferWriter) {
	if this.FirstLunch {
		to.WriteByte(0)
	} else {
		types.WriteMsg(to, &this.Settings)
	}
}

func (this *Setting) Read(from io.BufferReader) {
	remain := from.PeekRemain()
	if remain[0] == 0 {
		this.FirstLunch = true
		from.ReadByte()
	} else {
		types.ReadMsg(from, &this.Settings)
	}
}

type AppInitMsg struct {
	Width  int16 // getSystemService("display")).getDisplay(0)
	Height int16 // getSystemService("display")).getDisplay(0) - status_bar_height
	// 				sdk >=30     Insets insets0 = windowManager0.getCurrentWindowMetrics().getWindowInsets().getInsetsIgnoringVisibility(WindowInsets.Type.systemBars());
	//             				 point0.y = windowMetrics0.getBounds().height() - insets0.top - insets0.bottom;
	MemMaxSubTotal               int64 //A0I runtime0.maxMemory() - runtime0.totalMemory()
	MemoryClass                  int64 //A0L getSystemService("activity")).getMemoryClass() * 0x400 * 0x400
	PlatformRequestInfo          byte  //platform_request_info prop store 5
	Random                       int32 //random0.nextInt(0x7FFFFFFE) + 1
	NowTime                      int64 //System.currentTimeMillis()
	Flag                         types.VarUInt32
	Constant24                   byte
	Version                      string
	AboutLoginFlag               byte
	AboutLogin1                  string                         `ins:"AboutLoginFlag != 0"`
	AboutLogin2                  string                         `ins:"AboutLoginFlag != 0"`
	AboutLogin3                  string                         `ins:"AboutLoginFlag != 0"`
	HasMsg74StrArrayData         int16                          // 0,1
	Msg74StrArrayData            types.ListValue[string, int16] `ins:"HasMsg74StrArrayData != 0"`
	Flag2                        int64
	LoginFailureCounter          int16
	IgLiteUrl                    string
	StartUpUseTime               int32 //毫秒
	Local                        string
	AgentAndIsInstallFromGoogle  string //System.getProperty("http.agent") is_install_from_google modular=2
	Model                        string //Build.MODEL
	DensityDpi                   int16
	Constant0                    byte
	FontMap                      types.MapValue[string, int16, int16] // Arrays.sort(arr_byte)
	FontMapCrc32                 int32
	PrefKeyFontExperimentHash    int32
	PrefKeyNewFontExperimentHash int32
	PrefKeyNewFontExperimentData types.MapValue[int32, types.MapValue[int32, types.DescribeValue[byte], byte], byte] //暂时为kon
	Network                      byte
	TimeZone                     string
	Radom116                     int16
	AdvertisingId                string //gaid
	MatePackageInfo1             types.ListValue[MatePackageInfo, int16]
	FbMeta                       string
	Constant1                    byte
	FacebookOrcaVersionName      string
	FacebookMliteVersionName     string
	FacebookKatanaVersionName    string
	DataAvailableSize            uint16 //getDataDirectory
	SdcardAvailableSize          uint16 //getExternalStorageDirectory
	PropertiesstoreV02Crc        int32
	Battery                      float32
	Charge                       byte
	CpuAbi                       string
	Constant0x20199628           int32
	Settings                     Setting
	BuildSdkInt                  types.VarInt32                        //Build.VERSION.SDK_INT
	BuildVersionRelease          string                                //Build_VERSION_RELEASE
	ClientEncryptionSecret       types.ListValue[byte, types.VarInt32] // [2 random-16 0]
	NotificationsEnabled         byte                                  //getSystemService("notification")).areNotificationsEnabled()
	LimitAdTrackingFlag          types.VarInt32
	InstanceCount                int32
	Constant0_                   int16
	InAppBrowserModuleExist      byte
	AndroidId                    string
	GenDeviceTimeId              string //long v1 = ((long)new SecureRandom().nextInt()) & 0xFFFFFFFFL | ((long)(((int)v))) << 0x20;
	SimCountryIso                string //getSimCountryIso();getNetworkCountryIso();
	CallFrom                     byte
	BloksVersionID               string
	E2eTest                      string
	TrackingState                TrackingState
	McQueryHashBin               types.ListValue[byte, types.VarInt32]
	NullBytes                    types.ListValue[byte, types.VarInt32] //mc_query_hash_bin相关
	MatePackageInfo2             types.ListValue[MatePackageInfo, int16]
	GoogleEmailAccount           string
	SimState                     int32
	Uid                          int32
	GoogleAccountInfo            string
}
