package ins_lite

import (
	"CentralizedControl/common/android"
	"CentralizedControl/common/utils"
	"CentralizedControl/ins_lite/config"
	"CentralizedControl/ins_lite/proto"
	"CentralizedControl/ins_lite/proto/msg/sender"
	"CentralizedControl/ins_lite/proto/types"
	"bytes"
	"fmt"
	"strings"
	"time"
)

func (this *InsLiteClient) getSimpleNetType() int {
	if this.Cookies.Network.NetworkType == android.NetworkTypeTypeWifi {
		return 5
	} else if this.Cookies.Network.NetworkType == android.NetworkTypeTypeMobile {
		switch this.Cookies.Network.NetworkSubType {
		case 1, 4, 7, 11, 2:
			return 2
		case 3, 6, 8, 9, 10, 12, 14, 15, 5:
			return 3
		case 13:
			return 4
		default:
		}
	}
	return 1
}

func (this *InsLiteClient) checkTtsAndInstallSource() int {
	var flag2 = 0
	if this.Cookies.Google.IsInstallFromVending {
		flag2 |= 0x2000000
	} else {
		flag2 |= 0x1000000
	}
	isTouchExplorationEnabled := false
	if isTouchExplorationEnabled {
		flag2 |= 0x4000000
	}
	hasTtsService := true
	if hasTtsService {
		flag2 |= 0x200000000
	}
	if this.Cookies.Packages.Get("com.google.android.tts") != nil {
		flag2 |= 0x800000000
	}
	return flag2
}

func (this *InsLiteClient) SendInitAppMsg() {
	ClientEncryptionSecret := bytes.NewBuffer([]byte{})
	if this.Cookies.Session.ClientEncryptionSecret != nil {
		ClientEncryptionSecret.Write(this.Cookies.Session.ClientEncryptionSecret)
	} else {
		ClientEncryptionSecret.WriteByte(2)
		ClientEncryptionSecret.WriteByte(0x10)
		ClientEncryptionSecret.Write(utils.GenBytes(16))
		ClientEncryptionSecret.WriteByte(0)
	}

	flag := 16 | 0x200 | 0x400
	likeLogin := false
	if likeLogin {
		flag |= 2
	}
	clientEncryptionSecretIsNotNull := true
	if clientEncryptionSecretIsNotNull {
		flag |= 8
	}
	isBackground := false
	if isBackground {
		flag |= 64
	}
	IsTracking := true
	if IsTracking {
		flag |= 0x100
	}
	if this.Cookies.SdkInt >= 29 && (this.Cookies.Configuration.UiMode&48) == 32 {
		flag |= 0x800
	}

	flag2 := 0x880000 | 0x100000 | 0x20000000 | 2
	if this.Cookies.PropStore54.GetBool(0xE96, true) {
		//version code > 12451000
		if this.Cookies.Google.HasGoogleGms {
			flag2 |= 0x40000000
		}
	}
	if this.Cookies.PropStore54.GetBool(0xF39, false) {
		flag2 |= this.checkTtsAndInstallSource()
	}

	AgentAndIsInstallFromGoogle := "SupportsFresco=1 "
	if this.Cookies.Google.IsInstallFromVending {
		AgentAndIsInstallFromGoogle += "modular=2 "
	} else {
		AgentAndIsInstallFromGoogle += "modular=3 "
	}
	AgentAndIsInstallFromGoogle += this.Cookies.Env.HttpAgent

	var networkFlag byte = 0
	switch this.getSimpleNetType() {
	case 0:
		networkFlag = 8
	case 1:
		networkFlag = 40
	case 2:
		networkFlag = 16
	case 3:
		networkFlag = 24
	case 4:
		networkFlag = 32
	case 5:
		networkFlag = 2
	}
	if this.Cookies.PropStore54.GetBool(0xF36, false) {
		if this.getSimpleNetType() != 0 {
			networkFlag = networkFlag | 4
		}
	} else {
		if this.Cookies.Network.IsActiveNetworkMetered {
			networkFlag = networkFlag | 4
		}
	}
	setting := sender.Setting{}
	if this.Cookies.LunchCount != 0 {
		setting.FirstLunch = false
	} else {
		setting.FirstLunch = true
	}
	var notificationsEnabled byte = 0
	if this.Cookies.Env.NotificationsEnabled {
		notificationsEnabled = 2
	} else {
		notificationsEnabled = 1
	}
	var trackingState *sender.TrackingState
	if this.Cookies.LunchCount == 0 {
		trackingState = sender.CreateTrackingState(true, 0x1450003, false, 0, false, false)
	} else {
		trackingState = sender.CreateTrackingState(true, 0x1450003, false, 0, true, false)
	}
	var charge byte
	switch this.Cookies.Battery.Status {
	case 2, 5:
		charge = 1
	default:
		charge = 0
	}
	abiList := strings.Split(this.Cookies.UsedProps.AbiList, ",")
	if len(abiList) > 2 {
		abiList = abiList[0:2]
	}
	cpuAbi := strings.Join(abiList, "|")

	AppInitMsg := &proto.MessageC[sender.AppInitMsg]{
		Message: proto.Message[sender.AppInitMsg]{
			Body: sender.AppInitMsg{
				Width:                        int16(this.Cookies.WindowSize.AppWidth),
				Height:                       int16(this.Cookies.WindowSize.AppHeight),
				MemMaxSubTotal:               int64(this.Cookies.Memory.RuntimeMemory.MaxMemory - this.Cookies.Memory.RuntimeMemory.TotalMemory),
				MemoryClass:                  int64(this.Cookies.Memory.MemoryClass * 1024 * 1024),
				PlatformRequestInfo:          0,
				Random:                       int32(utils.GenNumber(0, 0x7FFFFFFE)) + 1,
				NowTime:                      time.Now().UnixMilli(),
				Flag:                         *types.CreateVarUInt32(uint32(flag)),
				Constant24:                   24,
				Version:                      this.Cookies.Packages.GetSelf().VersionStr,
				AboutLoginFlag:               0,
				AboutLogin1:                  "",
				AboutLogin2:                  "",
				AboutLogin3:                  "",
				HasMsg74StrArrayData:         0,
				Msg74StrArrayData:            types.ListValue[string, int16]{},
				Flag2:                        int64(flag2),
				LoginFailureCounter:          int16(this.getSharedPreferencesInt("login_failure_counter", 0)),
				IgLiteUrl:                    config.InsHost,
				StartUpUseTime:               int32(utils.GenNumber(1000, 3000)),
				Local:                        this.Cookies.Env.GetLanguageLocale(),
				AgentAndIsInstallFromGoogle:  AgentAndIsInstallFromGoogle,
				Model:                        this.Cookies.UsedProps.ProductModel,
				DensityDpi:                   int16(this.Cookies.WindowSize.DensityDpi),
				Constant0:                    0,
				FontMap:                      types.MapValue[string, int16, int16]{},
				FontMapCrc32:                 0,
				PrefKeyFontExperimentHash:    int32(this.getSharedPreferencesInt("pref_key_font_experiment_hash", 0)),
				PrefKeyNewFontExperimentHash: int32(this.getSharedPreferencesInt("pref_key_new_font_experiment_hash", 0)),
				PrefKeyNewFontExperimentData: types.MapValue[int32, types.MapValue[int32, types.DescribeValue[byte], byte], byte]{},
				Network:                      networkFlag,
				TimeZone:                     this.Cookies.Env.Timezone,
				Radom116:                     int16(utils.GenNumber(0, 0xffff)),
				AdvertisingId:                "",
				MatePackageInfo1:             types.ListValue[sender.MatePackageInfo, int16]{},
				FbMeta:                       "",
				Constant1:                    1,
				FacebookOrcaVersionName:      "",
				FacebookMliteVersionName:     "",
				FacebookKatanaVersionName:    "",
				DataAvailableSize:            uint16(utils.GenNumber(20000, 50000)),
				SdcardAvailableSize:          uint16(utils.GenNumber(20000, 50000)),
				PropertiesstoreV02Crc:        0,
				Battery:                      float32(this.Cookies.Battery.Level) / float32(this.Cookies.Battery.Scale),
				Charge:                       charge,
				CpuAbi:                       cpuAbi,
				Constant0x20199628:           0x20199628,
				Settings:                     setting,
				BuildSdkInt:                  types.CreateVarInt32(this.Cookies.SdkInt),
				BuildVersionRelease:          fmt.Sprintf("%d", android.SdkInt2AndroidVersion(this.Cookies.SdkInt)),
				ClientEncryptionSecret:       *types.CreateListValue[byte, types.VarInt32](ClientEncryptionSecret.Bytes()),
				NotificationsEnabled:         notificationsEnabled,
				LimitAdTrackingFlag:          types.VarInt32{},
				InstanceCount:                0,
				Constant0_:                   0,
				InAppBrowserModuleExist:      0,
				AndroidId:                    this.Cookies.Packages.GetSelf().AndroidId,
				GenDeviceTimeId:              this.Cookies.GenDeviceTimeId,
				SimCountryIso:                this.Cookies.Phone.Sim[0].GsmSimOperatorIsoCountry,
				CallFrom:                     0,
				BloksVersionID:               this.Cookies.BloksVersionID,
				E2eTest:                      "",
				TrackingState:                *trackingState,
				McQueryHashBin:               *types.CreateListValue[byte, types.VarInt32](this.Cookies.McQueryHashBin),
				NullBytes:                    types.ListValue[byte, types.VarInt32]{},
				MatePackageInfo2:             types.ListValue[sender.MatePackageInfo, int16]{},
				GoogleEmailAccount:           "",
				SimState:                     android.SimStateReady,
				Uid:                          int32(this.Cookies.Packages.GetSelf().Uid),
				GoogleAccountInfo:            "",
			},
		},
	}

	//从谷歌商店下载的ins才可以访问谷歌服务
	if this.Cookies.Google.HasGoogleGms && this.Cookies.Google.IsInstallFromVending {
		AppInitMsg.Body.AdvertisingId = this.Cookies.Google.AdvertisingId
		if this.Cookies.Google.LimitAdTracking {
			AppInitMsg.Body.LimitAdTrackingFlag = types.CreateVarInt32(1)
		} else {
			AppInitMsg.Body.LimitAdTrackingFlag = types.CreateVarInt32(0)
		}
	} else {
		AppInitMsg.Body.AdvertisingId = ""
		AppInitMsg.Body.LimitAdTrackingFlag = types.CreateVarInt32(2)
	}
	if this.Cookies.Google.HasGoogleAccount && this.Cookies.Google.HasGoogleGms && this.Cookies.Permission.IsAllow(android.GetAccounts) {
		AppInitMsg.Body.GoogleEmailAccount = this.Cookies.Google.LoginAccount // a@a.c,a@a.c
	}

	AppInitMsg.ClientId = -1
	AppInitMsg.Magic = 0xcf3
	this.SendMsg(AppInitMsg)
}
