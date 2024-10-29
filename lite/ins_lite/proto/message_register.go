package proto

import (
	"CentralizedControl/ins_lite/proto/msg/common"
	"CentralizedControl/ins_lite/proto/msg/recver"
	"CentralizedControl/ins_lite/proto/msg/sender"
	"fmt"
	"reflect"
)

const (
	//inner musk
	MuskSubCmd = 0x0000000100000000

	//common
	MsgCodeTextMsg = 330 //Bloks recv?

	//send
	MsgCodeAppInitMsg                      = 1
	MsgCodeInitiativeActionMsg             = 2
	MsgCodeUnknow3                         = 3 //X.msg_deal_0K2.A0c sub cmd code 1
	MsgCodeClientEventLog                  = 4
	MsgCodeSendSystemPropertiesMsg         = 6
	MsgCodeBrowserAction                   = 7
	MsgCodeLikeSendImageRecv               = 8
	MsgCodeReConnect                       = 39
	MsgCodeSendLastSomeEventInterval       = 56 //when recv 55  now - start times?
	MsgCodePreSendPhoneBook                = 59
	MsgCodeSendPhonebook                   = 60
	MsgCodeSendPing                        = 64 //when recv 65
	MsgCodePassivityActionMsg              = 83 //when recv 89?  ScreenReceived
	MsgCodeSendInstallReferrer             = 85 // 273
	MsgCodeFirebaseInstanceId              = 86 //event type 37
	MsgCodeNetworkTypeChangeReporter       = 96
	MsgCodeDurationTracking                = 111
	MsgCodeConnBandwidthQuality            = 126
	MsgCodeActivityResumeOrStop            = 127
	MsgCodeSendNetworkInfo                 = 152
	MsgCodeDownloadableResources           = 167
	MsgCodeSendPermResult                  = 168
	MsgCodeDownLoadModule                  = 170
	MsgCodeSendPkgInfo                     = 175
	MsgCodeAppModuleDownload               = 195
	MsgCodeSendPhoneIdMsg                  = 204
	MsgCodeFetchImage                      = 232
	MsgCodeClientStoredServerAccessibleKey = 253 //when recv sub 214 send
	MsgCodeFetchBloks                      = 257
	MsgCodeUploadMediaInfo                 = 260
	MsgCodeActivityResumed                 = 294
	MsgCodeSafetyNetData                   = 304
	MsgCodeLoggedUserIdChange              = 307
	MsgCodeSendMcLiteConfigsAndParamsList  = 312
	MsgCodeAdvertiserId                    = 321 // when recv 329 and other
	MsgCodeDownloadVideo                   = 325
	MsgCodeSendGoogleOAuthToken            = 352 //google email?
	MsgCodeImageQuality                    = 356 //image input stream?

	//recv
	MsgCodeScreenReceived                 = 11
	MsgCodeGetSystemPropertiesMsg         = 13
	MsgCodePropStoreConfig                = 15
	MsgCodeRectangularBackgroundConfig16  = 16
	MsgCodeUnknow19                       = 19
	MsgCodePass40                         = 40
	MsgCodeRecvScreenDiff                 = 42
	MsgCodeRecvSubCmd                     = 45
	MsgCodeRectangularBackgroundConfig50  = 50
	MsgCodeSessionMessage                 = 53
	MsgCodePropStore54                    = 54
	MsgCodeGetLastSomeEventInterval       = 55
	MsgCodeGetPing                        = 65
	MsgCodeConfig89                       = 89
	MsgCodeUpdateApp                      = 105
	MsgCodeHoneyAnalyticsSession          = 115
	MsgCodeAboutRecvQueueAction135        = 135
	MsgCodeRecvImage                      = 158
	MsgCodeTransactionID                  = 160
	MsgCodeDownload                       = 172
	MsgCodeRecvUserIco                    = 182
	MsgCodeCreateNotification             = 187
	MsgCodePersistentTranslatedStringsMap = 193 //语言映射
	//MsgCodeGetClientStoredServerAccessibleKey = 214
	MsgCodePropStore223                                  = 223
	MsgCodeHandleMessageOxygenAcceptTosOnSuccessfulLogin = 235
	MsgCodeRecvBloks                                     = 258
	MsgCode                                              = 272
	MsgCodeBitFlag275                                    = 275
	MsgCodeUnknow296                                     = 296
	MsgCodeRequireMcLiteConfigsAndParamsList             = 311

	//RecvSubCmd
	MsgCodeSubCmdGetPhoneBook                 = 39 | MuskSubCmd
	MsgCodeSubCmdCheckCallingOrSelfPermission = 116 | MuskSubCmd
	MsgCodeSubCmdCheckPermAndApplyForPerm     = 117 | MuskSubCmd
	MsgCodeSubCmdGetPkgInfo                   = 120 | MuskSubCmd
	MsgCodeSubCmdRequestAndroidPermissions    = 129 | MuskSubCmd
	MsgCodeSubCmdGetPhoneIdMsg                = 147 | MuskSubCmd
	MsgCodeSubCmdSetLoggerSettings            = 155 | MuskSubCmd
	MsgCodeSubCmdCleanCookieManager           = 201 | MuskSubCmd
	MsgCodeSubCmdRemoveStorageKey             = 213 | MuskSubCmd
	MsgCodeSubCmdGetStorageHeaders            = 214 | MuskSubCmd
	MsgCodeSubCmdRequireSafetyData            = 311 | MuskSubCmd
	MsgCodeSubCmdGetGoogleOAuthToken          = 347 | MuskSubCmd
)

type MessageType interface {
	sender.AppInitMsg |
		sender.SendPhoneIdMsg |
		sender.ActivityResumed |
		sender.SendSystemPropertiesMsg |
		sender.SendStorageHeaders |
		sender.ActionMsg |
		sender.PassivityActionMsg |
		sender.InitiativeActionMsg |
		sender.BrowserAction |
		sender.GoogleOAuthToken |
		sender.SendPkgInfo |
		sender.ClientEventLog |
		sender.NetworkInfo |
		sender.SendPhonebook |
		sender.FetchImage |
		sender.LoggedUserIdChange |
		sender.PermResult |
		sender.ConnBandwidthQuality |
		sender.NetworkTypeChangeReporter |
		sender.SafetyNetData |
		sender.ReConnect |
		sender.PreSendPhoneBook |
		sender.SendInstallReferrer |
		sender.AdvertiserId |
		common.TextMsg |
		recver.RequireSafetyData |
		recver.SessionMessage |
		recver.GetPhoneBook |
		recver.GetSystemPropertiesMsg |
		recver.GetStorageHeaders |
		recver.ScreenBase |
		//recver.ScreenNotFirst |
		recver.ScreenReceived |
		recver.ScreenDecode |
		recver.GetPkgInfo
}

func init() {
	SendCode2InfoRegister = make(map[uint64]*CodeRegisterInfo)
	SendName2InfoRegister = make(map[string]*CodeRegisterInfo)
	RecvCode2InfoRegister = make(map[uint64]*CodeRegisterInfo)
	RecvName2InfoRegister = make(map[string]*CodeRegisterInfo)

	RegisterMessageCode(true, MsgCodeAppInitMsg, "AppInitMsg", reflect.TypeOf((*sender.AppInitMsg)(nil)).Elem())
	RegisterMessageCode(true, MsgCodeInitiativeActionMsg, "InitiativeActionMsg", reflect.TypeOf((*sender.InitiativeActionMsg)(nil)).Elem())
	RegisterMessageCode(true, MsgCodeClientEventLog, "ClientEventLog", reflect.TypeOf((*sender.ClientEventLog)(nil)).Elem())
	RegisterMessageCode(true, MsgCodeSendSystemPropertiesMsg, "SendSystemPropertiesMsg", reflect.TypeOf((*sender.SendSystemPropertiesMsg)(nil)).Elem())
	RegisterMessageCode(true, MsgCodeBrowserAction, "BrowserAction", reflect.TypeOf((*sender.BrowserAction)(nil)).Elem())
	RegisterMessageCode(true, MsgCodeLikeSendImageRecv, "LikeSendImageRecv", reflect.TypeOf((*sender.LikeSendImageRecv)(nil)).Elem())
	RegisterMessageCode(true, MsgCodeSendLastSomeEventInterval, "SendLastSomeEventInterval", reflect.TypeOf((*sender.SendLastSomeEventInterval)(nil)).Elem())
	RegisterMessageCode(true, MsgCodeSendPing, "SendPing", reflect.TypeOf((*common.NoData)(nil)).Elem())
	RegisterMessageCode(true, MsgCodePassivityActionMsg, "PassivityActionMsg", reflect.TypeOf((*sender.PassivityActionMsg)(nil)).Elem())
	RegisterMessageCode(true, MsgCodeFirebaseInstanceId, "FirebaseInstanceId", reflect.TypeOf((*sender.FirebaseInstanceId)(nil)).Elem())
	RegisterMessageCode(true, MsgCodeNetworkTypeChangeReporter, "NetworkTypeChangeReporter", reflect.TypeOf((*sender.NetworkTypeChangeReporter)(nil)).Elem())
	RegisterMessageCode(true, MsgCodeConnBandwidthQuality, "ConnBandwidthQuality", reflect.TypeOf((*sender.ConnBandwidthQuality)(nil)).Elem())
	RegisterMessageCode(true, MsgCodeActivityResumeOrStop, "ActivityResumeOrStop", reflect.TypeOf((*sender.ActivityResumeOrStop)(nil)).Elem())
	RegisterMessageCode(true, MsgCodeDownloadableResources, "DownloadableResources", reflect.TypeOf((*sender.DownloadableResources)(nil)).Elem())
	RegisterMessageCode(true, MsgCodeDownLoadModule, "DownLoadModule", reflect.TypeOf((*sender.DownLoadModule)(nil)).Elem())
	RegisterMessageCode(true, MsgCodeSendPkgInfo, "SendPkgInfo", reflect.TypeOf((*sender.SendPkgInfo)(nil)).Elem())
	RegisterMessageCode(true, MsgCodeAppModuleDownload, "AppModuleDownload", reflect.TypeOf((*sender.AppModuleDownload)(nil)).Elem())
	RegisterMessageCode(true, MsgCodeSendPhoneIdMsg, "SendPhoneIdMsg", reflect.TypeOf((*sender.SendPhoneIdMsg)(nil)).Elem())
	RegisterMessageCode(true, MsgCodeFetchImage, "FetchImage", reflect.TypeOf((*sender.FetchImage)(nil)).Elem())
	RegisterMessageCode(true, MsgCodeClientStoredServerAccessibleKey, "ClientStoredServerAccessibleKey", reflect.TypeOf((*sender.SendStorageHeaders)(nil)).Elem())
	RegisterMessageCode(true, MsgCodeFetchBloks, "FetchBloks", reflect.TypeOf((*sender.FetchBloks)(nil)).Elem())
	RegisterMessageCode(true, MsgCodeActivityResumed, "ActivityResumed", reflect.TypeOf((*sender.ActivityResumed)(nil)).Elem())
	RegisterMessageCode(true, MsgCodeSafetyNetData, "SafetyNetData", reflect.TypeOf((*sender.SafetyNetData)(nil)).Elem())
	RegisterMessageCode(true, MsgCodeLoggedUserIdChange, "LoggedUserIdChange", reflect.TypeOf((*sender.LoggedUserIdChange)(nil)).Elem())
	RegisterMessageCode(true, MsgCodeAdvertiserId, "MsgCodeAdvertiserId", reflect.TypeOf((*sender.AdvertiserId)(nil)).Elem())
	RegisterMessageCode(true, MsgCodeDownloadVideo, "DownloadVideo", reflect.TypeOf((*sender.DownloadVideo)(nil)).Elem())
	RegisterMessageCode(true, MsgCodeSendGoogleOAuthToken, "SendGoogleOAuthToken", reflect.TypeOf((*sender.GoogleOAuthToken)(nil)).Elem())
	RegisterMessageCode(true, MsgCodeImageQuality, "ImageQuality", reflect.TypeOf((*sender.ImageQuality)(nil)).Elem())
	RegisterMessageCode(true, MsgCodeTextMsg, "SendTextMsg", reflect.TypeOf((*common.TextMsg)(nil)).Elem())
	RegisterMessageCode(true, MsgCodeSendNetworkInfo, "SendNetworkInfo", reflect.TypeOf((*sender.NetworkInfo)(nil)).Elem())
	RegisterMessageCode(true, MsgCodeSendPhonebook, "SendPhonebook", reflect.TypeOf((*sender.SendPhonebook)(nil)).Elem())
	RegisterMessageCode(true, MsgCodePreSendPhoneBook, "PreSendPhoneBook", reflect.TypeOf((*sender.PreSendPhoneBook)(nil)).Elem())
	RegisterMessageCode(true, MsgCodeSendPermResult, "PermResult", reflect.TypeOf((*sender.PermResult)(nil)).Elem())
	RegisterMessageCode(true, MsgCodeUploadMediaInfo, "UploadMediaInfo", reflect.TypeOf((*sender.UploadMediaInfo)(nil)).Elem())
	RegisterMessageCode(true, MsgCodeReConnect, "ReConnect", reflect.TypeOf((*sender.ReConnect)(nil)).Elem())
	RegisterMessageCode(true, MsgCodeSendInstallReferrer, "InstallReferrer", reflect.TypeOf((*sender.SendInstallReferrer)(nil)).Elem())
	RegisterMessageCode(true, MsgCodeDurationTracking, "DurationTracking", reflect.TypeOf((*sender.DurationTracking)(nil)).Elem())

	RegisterMessageCode(false, MsgCodeSessionMessage, "SessionMessage", reflect.TypeOf((*recver.SessionMessage)(nil)).Elem())
	RegisterMessageCode(false, MsgCodeGetSystemPropertiesMsg, "GetSystemPropertiesMsg", reflect.TypeOf((*recver.GetSystemPropertiesMsg)(nil)).Elem())
	RegisterMessageCode(false, MsgCodePass40, "Pass40", reflect.TypeOf((*common.NoData)(nil)).Elem())
	RegisterMessageCode(false, MsgCodeGetLastSomeEventInterval, "GetLastSomeEventInterval", reflect.TypeOf((*common.NoData)(nil)).Elem())
	RegisterMessageCode(false, MsgCodeGetPing, "GetPing", reflect.TypeOf((*common.NoData)(nil)).Elem())
	RegisterMessageCode(false, MsgCodeConfig89, "Config89", reflect.TypeOf((*recver.Config89)(nil)).Elem())
	RegisterMessageCode(false, MsgCodeUpdateApp, "UpdateApp", reflect.TypeOf((*recver.UpdateApp)(nil)).Elem())
	RegisterMessageCode(false, MsgCodeAboutRecvQueueAction135, "Unknow135", reflect.TypeOf((*common.NoData)(nil)).Elem())
	//RegisterMessageCode(false, MsgCode_GetClientStoredServerAccessibleKey, "GetClientStoredServerAccessibleKey", reflect.TypeOf((*recver.GetClientStoredServerAccessibleKey)(nil)).Elem())
	RegisterMessageCode(false, MsgCodePropStore223, "PropStore223", reflect.TypeOf((*recver.PropStore223)(nil)).Elem())
	RegisterMessageCode(false, MsgCodeUnknow19, "Unknow19", reflect.TypeOf((*recver.Unknow19)(nil)).Elem())
	RegisterMessageCode(false, MsgCodeRecvScreenDiff, "ReceScreenDiff", reflect.TypeOf((*recver.ScreenDiff)(nil)).Elem())
	RegisterMessageCode(false, MsgCodeHoneyAnalyticsSession, "PasswordEncryption", reflect.TypeOf((*common.NoData)(nil)).Elem())
	RegisterMessageCode(false, MsgCodeTransactionID, "TransactionID", reflect.TypeOf((*recver.TransactionID)(nil)).Elem())
	RegisterMessageCode(false, MsgCodeRecvUserIco, "RecvUserIco", reflect.TypeOf((*recver.RecvUserIco)(nil)).Elem())
	RegisterMessageCode(false, MsgCodeCreateNotification, "Notification", reflect.TypeOf((*recver.Notification)(nil)).Elem())
	RegisterMessageCode(false, MsgCodeHandleMessageOxygenAcceptTosOnSuccessfulLogin, "HandleMessageOxygenAcceptTosOnSuccessfulLogin", reflect.TypeOf((*common.NoData)(nil)).Elem())
	RegisterMessageCode(false, MsgCodeRecvBloks, "RecvBloks", reflect.TypeOf((*recver.RecvBloks)(nil)).Elem())
	RegisterMessageCode(false, MsgCodeUnknow296, "Unknow296", reflect.TypeOf((*common.NoData)(nil)).Elem())
	RegisterMessageCode(false, MsgCodeTextMsg, "RecvTextMsg", reflect.TypeOf((*common.TextMsg)(nil)).Elem())
	RegisterMessageCode(false, MsgCodeRecvSubCmd, "RecvSubCmd", reflect.TypeOf((*recver.RecvSubCmd)(nil)).Elem())
	RegisterMessageCode(false, MsgCodeScreenReceived, "ScreenReceived", reflect.TypeOf((*recver.ScreenReceivedHeader)(nil)).Elem())
	RegisterMessageCode(false, MsgCodeRecvImage, "RecvImage", reflect.TypeOf((*recver.RecvImage)(nil)).Elem())
	RegisterMessageCode(false, MsgCodePropStore54, "PropStore54", reflect.TypeOf((*recver.PropStore54)(nil)).Elem())
	RegisterMessageCode(false, MsgCodePersistentTranslatedStringsMap, "PersistentTranslatedStringsMap", reflect.TypeOf((*recver.PersistentTranslatedStringsMap)(nil)).Elem())
	//RegisterMessageCode(false, MsgCode_PropStoreConfig, "", reflect.TypeOf((*recver.WindowConfig)(nil)).Elem())
	//RegisterMessageCode(false, MsgCode_RectangularBackgroundConfig16, "", reflect.TypeOf((*recver.RectangularBackgroundConfig16)(nil)).Elem())
	//RegisterMessageCode(false, MsgCode_RectangularBackgroundConfig50, "", reflect.TypeOf((*recver.RectangularBackgroundConfig50)(nil)).Elem())
	//RegisterMessageCode(false, MsgCode_Config54, "", reflect.TypeOf((*recver.Config54)(nil)).Elem())
	//RegisterMessageCode(false, MsgCode_Download, "", reflect.TypeOf((*recver.Download)(nil)).Elem())
	//RegisterMessageCode(false, MsgCode_StringMapConfig, "", reflect.TypeOf((*recver.StringMapConfig)(nil)).Elem())
	//RegisterMessageCode(false, MsgCode_Unknow275, "", reflect.TypeOf((*recver.Unknow275)(nil)).Elem())

	//sub cmd
	RegisterMessageCode(false, MsgCodeSubCmdRemoveStorageKey, "RemoveStorageKey", reflect.TypeOf((*recver.RemoveStorageKey)(nil)).Elem())
	RegisterMessageCode(false, MsgCodeSubCmdGetStorageHeaders, "GetStorageHeaders", reflect.TypeOf((*recver.GetStorageHeaders)(nil)).Elem())
	RegisterMessageCode(false, MsgCodeSubCmdGetPkgInfo, "GetPkgInfo", reflect.TypeOf((*recver.GetPkgInfo)(nil)).Elem())
	RegisterMessageCode(false, MsgCodeSubCmdRequestAndroidPermissions, "RequestPerm", reflect.TypeOf((*recver.RequestPerm)(nil)).Elem())
	RegisterMessageCode(false, MsgCodeSubCmdSetLoggerSettings, "LoggerSettings", reflect.TypeOf((*recver.LoggerSettings)(nil)).Elem())
	RegisterMessageCode(false, MsgCodeSubCmdRequireSafetyData, "RequireSafetyData", reflect.TypeOf((*recver.RequireSafetyData)(nil)).Elem())
	RegisterMessageCode(false, MsgCodeSubCmdGetPhoneBook, "RequireSafetyData", reflect.TypeOf((*recver.GetPhoneBook)(nil)).Elem())
	RegisterMessageCode(false, MsgCodeSubCmdCheckPermAndApplyForPerm, "CheckPermAndApplyForPerm", reflect.TypeOf((*recver.CheckPermAndApplyForPerm)(nil)).Elem())
	RegisterMessageCode(false, MsgCodeSubCmdGetGoogleOAuthToken, "GoogleOAuthToken", reflect.TypeOf((*recver.GoogleOAuthToken)(nil)).Elem())
}

type CodeRegisterInfo struct {
	IsSend bool
	Code   uint64
	Desc   string
	Type   reflect.Type
}

var SendCode2InfoRegister map[uint64]*CodeRegisterInfo
var SendName2InfoRegister map[string]*CodeRegisterInfo

var RecvCode2InfoRegister map[uint64]*CodeRegisterInfo
var RecvName2InfoRegister map[string]*CodeRegisterInfo

func RegisterMessageCode(isSend bool, code uint64, desc string, _type reflect.Type) {
	info := &CodeRegisterInfo{
		IsSend: isSend,
		Code:   code,
		Desc:   desc,
		Type:   _type,
	}
	if isSend {
		SendCode2InfoRegister[code] = info
		SendName2InfoRegister[_type.Name()] = info
	} else {
		RecvCode2InfoRegister[code] = info
		RecvName2InfoRegister[_type.Name()] = info
	}
}

func GetMessageName(isSend bool, code uint64) string {
	info := GetMessageInfo(isSend, code)
	if info != nil {
		return info.Desc
	} else {
		return fmt.Sprintf("%d", code)
	}
}

func GetMessageInfo(isSend bool, code uint64) *CodeRegisterInfo {
	var reg map[uint64]*CodeRegisterInfo
	if isSend {
		reg = SendCode2InfoRegister
	} else {
		reg = RecvCode2InfoRegister
	}
	return reg[code]
}

func GetMessageCode(isSend bool, typeName string) uint64 {
	var reg map[string]*CodeRegisterInfo
	if isSend {
		reg = SendName2InfoRegister
	} else {
		reg = RecvName2InfoRegister
	}
	return reg[typeName].Code
}

func CreateMsgByCode(isSend bool, code uint64) *reflect.Value {
	info := GetMessageInfo(isSend, code)
	if info == nil {
		return nil
	}
	r := reflect.New(info.Type)
	return &r
}
