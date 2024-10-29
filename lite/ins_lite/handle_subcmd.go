package ins_lite

import (
	"CentralizedControl/common/android"
	"CentralizedControl/common/log"
	"CentralizedControl/common/utils"
	"CentralizedControl/ins_lite/net"
	"CentralizedControl/ins_lite/proto"
	"CentralizedControl/ins_lite/proto/io"
	"CentralizedControl/ins_lite/proto/msg/recver"
	"CentralizedControl/ins_lite/proto/msg/sender"
	"strings"
)

func (this *InsLiteClient) registerSubHandle() {
	this.RegisterEvent(proto.MsgCodeSubCmdGetStorageHeaders, this.subCmdRecvGetStorageHeaders)
	this.RegisterEvent(proto.MsgCodeSubCmdGetPhoneIdMsg, this.subCmdRecvGetPhoneIdMsg)
	this.RegisterEvent(proto.MsgCodeSubCmdGetGoogleOAuthToken, this.subCmdRecvGetGoogleOAuthToken)
	this.RegisterEvent(proto.MsgCodeSubCmdRequestAndroidPermissions, this.subCmdRecvRequestAndroidPermissions)
	this.RegisterEvent(proto.MsgCodeSubCmdGetPkgInfo, this.subCmdRecvGetPkgInfo)
	this.RegisterEvent(proto.MsgCodeSubCmdCleanCookieManager, this.subCmdRecvDoNothing)
	this.RegisterEvent(proto.MsgCodeSubCmdRemoveStorageKey, this.subCmdRecvRemoveStorageKey)
	this.RegisterEvent(proto.MsgCodeSubCmdSetLoggerSettings, this.subCmdRecvSetLoggerSettings)
	this.RegisterEvent(proto.MsgCodeSubCmdRequireSafetyData, this.subCmdRecvRequireSafetyData)
	this.RegisterEvent(proto.MsgCodeSubCmdGetPhoneBook, this.subCmdRecvGetPhoneBook)
	this.RegisterEvent(proto.MsgCodeSubCmdCheckPermAndApplyForPerm, this.subCmdRecvCheckPermAndApplyForPerm)
}

func (this *InsLiteClient) callSubCmd(code uint64, data []byte) error {
	this.PostSubEvent(&net.MsgWrap{Reader: io.CreateReaderBuffer(data)}, code|proto.MuskSubCmd)
	return nil
}

func (this *InsLiteClient) recvSubCmd(msgEvent *net.MsgWrap) error {
	body := msgEvent.Body.Interface().(*recver.RecvSubCmd)
	this.PostSubEvent(msgEvent, uint64(body.SubCode.Value)|proto.MuskSubCmd)
	return nil
}

func (this *InsLiteClient) subCmdRecvDoNothing(msgEvent *net.MsgWrap) error {
	return nil
}

func (this *InsLiteClient) subCmdRecvGetPhoneIdMsg(msgEvent *net.MsgWrap) error {
	SendPhoneIdMsg := &proto.Message[sender.SendPhoneIdMsg]{}
	SendPhoneIdMsg.Body.PhoneId = this.Cookies.PhoneId
	this.SendMsg(SendPhoneIdMsg)
	return nil
}

func (this *InsLiteClient) subCmdRecvGetStorageHeaders(msgEvent *net.MsgWrap) error {
	body := msgEvent.Body.Interface().(*recver.GetStorageHeaders)
	s := &proto.Message[sender.SendStorageHeaders]{}
	s.Body.IntKey = body.IntKey
	s.Body.StringKey = body.Key
	if body.Key == "zero_storage_headers_v2" || strings.HasPrefix(body.Key, "client_stored_server_accessible_key_") {
		s.Body.Value = this.getSharedPreferencesString(body.Key, "")
	} else {
		s.Body.Value = ""
	}
	this.SendMsg(s)
	return nil
}

func (this *InsLiteClient) subCmdRecvGetGoogleOAuthToken(msgEvent *net.MsgWrap) error {
	recv := msgEvent.Body.Interface().(*recver.GoogleOAuthToken)
	msg := &proto.Message[sender.GoogleOAuthToken]{}
	body := &msg.Body
	if !this.Cookies.Permission.IsAllow(android.ReadContacts) {
		body.Status = 0
		body.FailedError = "READ_CONTACTS_PERMISSION_NOT_AVAILABLE"
		this.SendMsg(msg)
		return nil
	}
	if !this.Cookies.Permission.IsAllow(android.GetAccounts) {
		body.Status = 0
		body.FailedError = "GET_ACCOUNTS_PERMISSION_NOT_AVAILABLE"
		this.SendMsg(msg)
		return nil
	}

	if this.Cookies.Google.LoginAccount == "" {
		body.Status = 0
		body.FailedError = "NO_ACCOUNT_IN_DEVICE"
	} else {
		token := this.Cookies.GmsClient.GetAccountToken(this.Cookies.Google.LoginAccount, recv.Scope, android.PkgNameInstagramLite)
		if token != nil {
			body.Status = 1
			item := sender.GoogleOAuthTokenItem{}
			item.Account = this.Cookies.Google.LoginAccount
			item.Cookies = token.AuthToken
			body.Token.Put(item)
		} else {
			body.Status = 0
			body.FailedError = "EXCEPTION"
			log.Error("get oauth token error!")
		}
	}
	this.SendMsg(msg)
	return nil
}

func (this *InsLiteClient) subCmdRecvGetPkgInfo(msgEvent *net.MsgWrap) error {
	body := msgEvent.Body.Interface().(*recver.GetPkgInfo)
	msg := &proto.Message[sender.SendPkgInfo]{}
	for _, pkgName := range body.Pkg.Value {
		var value byte
		if this.Cookies.Packages.Get(pkgName) == nil {
			value = 0
		} else {
			value = 1
		}
		msg.Body.PkgInfo.Put(pkgName, value)
	}
	this.SendMsg(msg)
	return nil
}

func (this *InsLiteClient) subCmdRecvRemoveStorageKey(msgEvent *net.MsgWrap) error {
	body := msgEvent.Body.Interface().(*recver.RemoveStorageKey)
	_ = body
	return nil
}

func (this *InsLiteClient) subCmdRecvRequestAndroidPermissions(msgEvent *net.MsgWrap) error {
	body := msgEvent.Body.Interface().(*recver.RequestPerm)
	permReq := make([]byte, 0)
	for _, perm := range body.PermType {
		if !this.checkHasPerm(perm) {
			permReq = append(permReq, perm)
		}
	}
	if len(permReq) == 0 {
		return this.callSubCmd(uint64(body.OnGrantedIdx.Value), body.OnGrantedData.Value)
	}
	result := this.requestPermissions(permReq)
	granted := make([]byte, 0)
	for idx := range result {
		if result[idx] == android.PermStatusGranted {
			granted = append(granted, permReq[idx])
		}
	}
	if len(granted) > 0 {
		return this.callSubCmd(uint64(body.OnGrantedIdx.Value), body.OnGrantedData.Value)
	} else {
		return this.callSubCmd(uint64(body.OnDeniedIdx.Value), body.OnDeniedData.Value)
	}
}

func (this *InsLiteClient) subCmdRecvSetLoggerSettings(msgEvent *net.MsgWrap) error {
	body := msgEvent.Body.Interface().(*recver.LoggerSettings)
	_ = body
	return nil
}

func (this *InsLiteClient) subCmdRecvRequireSafetyData(msgEvent *net.MsgWrap) error {
	msgBody := msgEvent.Body.Interface().(*recver.RequireSafetyData)
	safetyNet := &proto.Message[sender.SafetyNetData]{}
	if this.Cookies.Google.HasSafetyNet {
		attest, err := this.Cookies.GmsClient.SafetyNetAttest(msgBody.ApiKey, utils.Base64Encode(msgBody.Nonce.Value))
		if err == nil {
			safetyNet.Body.Const0 = 0
			safetyNet.Body.Const1 = 1
			safetyNet.Body.Unknow = attest
		} else {
			log.Error("SafetyNetAttest error: %v", err)
			safetyNet.Body.Const0 = 0
			safetyNet.Body.Const1 = 0
			safetyNet.Body.Unknow = "NETWORK_ERROR"
		}
	} else {
		safetyNet.Body.Const0 = 0
		safetyNet.Body.Const1 = 0
		safetyNet.Body.Unknow = "NETWORK_ERROR"
	}
	//UNKNOWN
	//SUCCESS
	//SERVICE_MISSING
	//SERVICE_VERSION_UPDATE_REQUIRED
	//SERVICE_DISABLED
	//SIGN_IN_REQUIRED
	//INVALID_ACCOUNT
	//RESOLUTION_REQUIRED
	//NETWORK_ERROR
	//INTERNAL_ERROR
	//SERVICE_INVALID
	//DEVELOPER_ERROR
	//LICENSE_CHECK_FAILED
	//CANCELED
	//TIMEOUT
	//INTERRUPTED
	//API_UNAVAILABLE
	//SIGN_IN_FAILED
	//SERVICE_UPDATING
	//SERVICE_MISSING_PERMISSION
	//RESTRICTED_PROFILE
	//API_VERSION_UPDATE_REQUIRED
	//RESOLUTION_ACTIVITY_NOT_FOUND
	//API_DISABLED
	//API_DISABLED_FOR_CONNECTION
	//UNFINISHED
	//DRIVE_EXTERNAL_STORAGE_REQUIRED
	this.SendMsg(safetyNet)
	return nil
}

func (this *InsLiteClient) subCmdRecvGetPhoneBook(msgEvent *net.MsgWrap) error {
	getPhoneBook := msgEvent.Body.Interface().(*recver.GetPhoneBook)
	preSendPhoneBook := &proto.Message[sender.PreSendPhoneBook]{}
	sendPhoneBook := &proto.Message[sender.SendPhonebook]{}

	body1 := preSendPhoneBook.Body
	body1.Idx = getPhoneBook.Idx
	body1.Const0 = 0
	body1.Unknow4 = ""

	body := &sendPhoneBook.Body
	body.Idx = getPhoneBook.Idx
	body.Const0 = 0
	body.Phonebook = "Phonebook"
	body.Count = int16(len(this.Cookies.PhoneBook.Item))
	for _, item := range this.Cookies.PhoneBook.Item {
		sendItem := sender.PhonebookItem{}
		sendItem.PutFirstName120(item.FirstName)
		sendItem.PutSecondName119(item.SecondName)
		sendItem.PutFullName105(item.FirstName + " " + item.SecondName)
		for _, email := range item.Email {
			sendItem.PutEmail103(email)
		}
		for _, Phone := range item.Phone {
			sendItem.PutPhone115(Phone)
		}
		body.Item.Put(sendItem)
	}

	body1.SendPhoneBookSize = int32(len(sendPhoneBook.WriteTo()) - 2)
	this.SendMsg(preSendPhoneBook)
	this.SendMsg(sendPhoneBook)
	return nil
}

func (this *InsLiteClient) subCmdRecvCheckPermAndApplyForPerm(msgEvent *net.MsgWrap) error {
	body := msgEvent.Body.Interface().(*recver.CheckPermAndApplyForPerm)
	if this.checkHasPerm(body.PermType) {
		if this.curScreenName != "" {
			this.callScreenEvent(this.GetCurrentScreen(), uint16(body.OnGrantedScreenCmdIdx), "Tag_RecvCheckPermAndApplyForPerm")
		}
		if body.HasPkgIdx {
			this.sendPermResult(int32(body.PkgIdx), false, true, false)
		}
		if body.ReqPermType == 0 {
			return nil
		}
	}

	var result []int
	switch body.ReqPermType {
	case 0:
		result = this.requestPermissions([]byte{body.ReqPermType})
	case 1:
		result = this.requestPermissions([]byte{body.ReqPermType})
	}

	if body.ReqPermType == 0 {
		if result[0] == android.PermStatusGranted {
			this.callScreenEvent(this.GetCurrentScreen(), uint16(body.OnGrantedScreenCmdIdx), "Tag_RecvCheckPermAndApplyForPerm")
		} else {
			this.callScreenEvent(this.GetCurrentScreen(), uint16(body.OnDeniedScreenCmdIdx), "Tag_RecvCheckPermAndApplyForPerm")
		}
	}
	if body.HasPkgIdx {
		if result[0] == android.PermStatusGranted {
			this.sendPermResult(int32(body.PkgIdx), false, true, false)
		} else {
			this.sendPermResult(int32(body.PkgIdx), false, false, this.shouldNotShowPermissionDialog(body.ReqPermType))
		}
	}
	return nil
}
