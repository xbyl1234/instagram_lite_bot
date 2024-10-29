package facebook

type Register struct {
	Fb *Facebook
}

//
//func createEmptyFacebook() *Facebook {
//	ck := &Cookies{
//		Id:             0,
//		AccountState:   AccountStateNotCreate,
//		Email:          "",
//		MachineId:      "",
//		PhoneId:        "",
//		SessionCookies: "",
//		Secret:         "",
//		Claim:          "",
//		SessionKey:     "",
//	}
//	ck.AndroidSdkInt = "9"
//	ck.AndroidId = utils.GenString(utils.CharSet_16_Num, 16)
//	ck.AppVersion = "417.0.0.33.65"
//	ck.AppVersionInt = "480086274"
//	ck.AppProduct = "350685531728"
//
//	ck.FamilyDeviceId = strings.ToLower(utils.GenUUID())
//	ck.DeviceId = ck.FamilyDeviceId
//	ck.AccessToken = ck.AppProduct + "|62f8ce9f74b12f84c123cc23437a4a32"
//	ck.Authorization = "OAuth " + ck.AccessToken
//
//	//ck.Username = android.Resource.ChoiceUsername()
//	//ck.FirstName = android.Resource.ChoiceUsername()
//	//ck.LastName = android.Resource.ChoiceUsername()
//
//	ck.Passwd = utils.GenString(utils.CharSet_abc, 3) +
//		utils.GenString(utils.CharSet_ABC, 3) +
//		utils.GenString(utils.CharSet_123, 3)
//
//	//device := android.Resource.GenDevice("hk")
//	//ck.Locale = device.Locale
//	//ck.Density = device.Density
//	//ck.Width = device.Width
//	//ck.Height = device.Height
//	//ck.Mcc = device.Mcc
//	//ck.Mnc = device.Mnc
//	//ck.Country = device.Country
//	//ck.Language = device.Language
//	//ck.SimOperatorName = device.SimOperatorName
//	//ck.BuildId = device.BuildId
//	//ck.ProductModel = device.ProductModel
//	//ck.ProductManufacturer = device.ProductManufacturer
//	//ck.ProductBrand = device.ProductBrand
//
//	ck.UserAgent = fmt.Sprintf("[FBAN/FB4A;"+
//		"FBAV/%s;"+
//		"FBBV/%s;"+
//		"FBDM/{density=%s,width=%s,height=%s};"+
//		"FBLC/%s;"+
//		"FBRV/0;"+
//		"FBCR/%s;"+
//		"FBMF/%s;"+
//		"FBBD/%s;"+
//		"FBPN/com.facebook.katana;"+
//		"FBDV/%s;"+
//		"FBSV/%s;"+
//		"FBOP/1;"+
//		"FBCA/x86:armeabi-v7a;]",
//		ck.AppVersion,
//		ck.AppVersionInt,
//		ck.Density,
//		ck.Width,
//		ck.Height,
//		ck.Locale,
//		ck.SimOperatorName,
//		ck.ProductManufacturer,
//		ck.ProductBrand,
//		ck.ProductModel,
//		ck.AndroidSdkInt,
//	)
//
//	return CreateFacebook(ck)
//}

//func CreateRegister() *Register {
//	return &Register{
//		Fb: createEmptyFacebook(),
//	}
//}

//
//func (this *Register) FamilyDeviceIDAppScopedDeviceIDSyncMutation() {
//	req := this.Fb.newApiRequest("FamilyDeviceIDAppScopedDeviceIDSyncMutation", "")
//	body := req.GetGraphqlBody()
//	body.AutoSetForm()
//	body.AutoCommonJson(body.GetRoot().Get("input"))
//	send, err := req.Send()
//	if err != nil {
//		return
//	}
//	_ = send
//}
//
//func (this *Register) ZeroHeadersPingParamsV2() {
//	req := this.Fb.newApiRequest("/zero_headers_ping_params_v2", "b9d6e6020acc6ced84dd22ea9dd164c8")
//	req.SetHeader("Authorization", "OAuth null")
//	body := req.GetFormBody()
//	body.AutoSetForm()
//	body.SetForm("logged_out_id", req.tempData.LoggingId)
//	send, err := req.Send()
//	if err != nil {
//		return
//	}
//	_ = send
//}
//
//func (this *Register) MobileConfigSessionLess() {
//	req := this.Fb.newApiRequest("/mobileconfigsessionless", "")
//	body := req.GetFormBody()
//	body.AutoSetForm()
//	send, err := req.Send()
//	if err != nil {
//		return
//	}
//	_ = send
//}
//
//func (this *Register) ProcessClientDataAndRedirect() {
//	req := this.Fb.newApiRequest("FbBloksActionRootQuery-com.bloks.www.bloks.caa.login.process_client_data_and_redirect", "")
//	body := req.GetGraphqlBody()
//	body.AutoCommonJson(body.GetSubParams())
//	sub := body.GetSubParams()
//	sub.SetString("qpl_join_id", strings.ToLower(utils.GenUUID()))
//	send, err := req.Send()
//	if err != nil {
//		return
//	}
//	_ = send
//}
//
//func (this *Register) PwdKeyFetch() {
//	req := this.Fb.newApiRequest("//pwd_key_fetch", "")
//	body := req.GetFormBody()
//	body.AutoSetForm()
//	send, err := req.Send()
//	if err != nil {
//		return
//	}
//	_ = send
//}

func (this *Register) MobileGatekeepers() {

}
