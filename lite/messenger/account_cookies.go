package messenger

import (
	"CentralizedControl/common"
	"CentralizedControl/common/android"
)

type Cookies struct {
	Id int64 `json:"id" gorm:"primaryKey"`
	//account info
	android.Account
	android.Device
	android.IpInfo
	android.VpnInfo

	TwoFAKey  string `json:"two_fa_key"`
	QrCodeUri string `json:"qr_code_uri"`
	//cookies
	Authorization  string `json:"authorization"`
	AccountType    string `json:"account_type"` // int
	AccountId      string `json:"account_id"`   // uid int64
	MachineId      string `json:"machine_id"`
	PhoneId        string `json:"phone_id"`
	DeviceId       string `json:"device_id"`
	FamilyDeviceId string `json:"family_device_id"`
	SessionCookies string `json:"session_cookies"`
	Secret         string `json:"secret"`
	Claim          string `json:"claim"`
	SessionKey     string `json:"session_key"`

	//flags
	Had2FA string `json:"had_2fa" gorm:"column:had_2fa"`

	//time
	CreateTime common.CustomTime `json:"create_time" gorm:"autoCreateTime"`
	UpdatedAt  common.CustomTime `json:"updated_at" gorm:"autoUpdateTime"`
}

type SessionCookies struct {
	Name             string `json:"name"`
	Value            string `json:"value"`
	Expires          string `json:"expires"`
	ExpiresTimestamp int    `json:"expires_timestamp"`
	Domain           string `json:"domain"`
	Path             string `json:"path"`
	Secure           bool   `json:"secure"`
	Httponly         *bool  `json:"httponly"`
	Samesite         string `json:"samesite"`
}

//func ConvDeviceFile2Cookies(ck string) *Cookies {
//	var msgCk Cookies
//	_ = json.Unmarshal([]byte(ck), &msgCk)
//	msgCk.Locale = strings.ToLower(msgCk.Language) + "_" + strings.ToUpper(msgCk.Country)
//	msgCk.AccountType = "0"
//	msgCk.SdkInt = "9"
//	msgCk.AppVersion = "417.0.0.12.64"
//	msgCk.AppVersionCode = "493307655"
//	if msgCk.DeviceId == "" {
//		msgCk.DeviceId = msgCk.PhoneId
//	}
//	if msgCk.FamilyDeviceId == "" {
//		msgCk.FamilyDeviceId = msgCk.PhoneId
//	}
//
//	var sessionCookies []SessionCookies
//	_ = json.Unmarshal([]byte(msgCk.SessionCookies), &sessionCookies)
//	for _, item := range sessionCookies {
//		if item.Name == "datr" {
//			msgCk.MachineId = item.Value
//		}
//	}
//
//	//"Dalvik/2.1.0 (Linux; U; Android 9; SM-N9700 Build/PQ3B.190801.06161913) [FBAN/Orca-Android;FBAV/391.2.0.20.404;FBPN/com.facebook.orca;FBLC/zh_CN;FBBV/437533963;FBCR/CHINA MOBILE;FBMF/samsung;FBBD/samsung;FBDV/SM-N9700;FBSV/9;FBCA/x86:armeabi-v7a;FBDM/{density=3.0,width=1080,height=1920};FB_FW/1;]"
//	// Dalvik/2.1.0 (Linux; U; Android 9; RMX1931  Build/%s) [FBAN/Orca-Android;FBAV/417.0.0.12.64; FBPN/com.facebook.orca;FBLC/zh_CN;FBBV/493307655;FBCR/CHINA MOBILE;FBMF/realme; FBBD/realme; FBDV/RMX1931; FBSV/9;FBCA/x86:armeabi-v7a;FBDM/{density=3.0,width=1080,height=1920};FB_FW/1;]
//	// Dalvik/2.1.0 (Linux; U; Android 9; ro.product.model Build/PQ3B.190801.06161913) [FBAN/Orca-Android;FBAV/417.0.0.12.64;FBPN/com.facebook.orca;FBLC/zh_CN;FBBV/493307655;FBCR/CHINA MOBILE;FBMF/ro.product.manufacturer;FBBD/ro.product.brand;FBDV/ro.product.model;FBSV/9;FBCA/x86:armeabi-v7a;FBDM/{density=3.0,width=1080,height=1920};FB_FW/1;]
//
//	useragent := fmt.Sprintf("Dalvik/2.1.0 (Linux; U; Android %s; %s Build/%s) "+
//		"[FBAN/Orca-Android;FBAV/%s;FBPN/com.facebook.orca;FBLC/%s;FBBV/%s;FBCR/%s;"+
//		"FBMF/%s;FBBD/%s;FBDV/%s;FBSV/9;FBCA/x86:armeabi-v7a;FBDM/{density=%s,width=%s,height=%s};"+
//		"FB_FW/1;]",
//		msgCk.SdkInt,
//		msgCk.ProductModel,        //ro.product.model
//		msgCk.BuildId,             //ro.build.id
//		msgCk.AppVersion,          //417.0.0.12.64
//		msgCk.Locale,              //zh_CN
//		msgCk.AppVersionCode,       //493307655
//		msgCk.GsmOperatorAlpha,    //gsm.operator.alpha
//		msgCk.ProductManufacturer, //ro.product.manufacturer
//		msgCk.ProductBrand,        //ro.product.brand
//		msgCk.ProductModel,        //ro.product.model
//		msgCk.Scale,
//		msgCk.Width,
//		msgCk.Height,
//	)
//	msgCk.UserAgent = useragent
//	return &msgCk
//}

type TempCookies struct {
}

//func CreateTempCookies() *TempCookies {
//
//}
