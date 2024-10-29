package facebook

import "CentralizedControl/common"

const (
	AccountStateNotCreate = 0
	AccountStateNotLogin  = 1
	AccountStateLogin     = 2
	AccountStateBan       = 3
)

type Cookies struct {
	Id int64 `json:"id" gorm:"primaryKey"`
	//account state
	AccountState int
	//account info
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Passwd    string `json:"passwd"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	TwoFAKey  string `json:"two_fa_key"`
	QrCodeUri string `json:"qr_code_uri"`

	//cookies
	Authorization  string `json:"authorization"`
	AccessToken    string `json:"access_token"`
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

	//device info
	AndroidId           string `json:"android_id"`
	UserAgent           string `json:"user_agent"` //runtime gen
	Locale              string `json:"locale"`     //runtime gen
	Density             string `json:"density"`
	Width               string `json:"width"`
	Height              string `json:"height"`
	Mcc                 string `json:"mcc"`
	Mnc                 string `json:"mnc"`
	Country             string `json:"country"`
	Language            string `json:"language"`
	SimOperatorName     string `json:"sim_operator_name"`
	BuildId             string `json:"build_id"`
	ProductModel        string `json:"product_model"`
	ProductManufacturer string `json:"product_manufacturer"`
	ProductBrand        string `json:"product_brand"`
	AndroidSdkInt       string `json:"android_sdk_int"`

	//app info
	AppVersion    string `json:"app_version"`
	AppVersionInt string `json:"app_version_int"`
	AppProduct    string `json:"app_product"`

	//flags
	Had2FA string `json:"had_2fa" gorm:"column:had_2fa"`

	//other
	ProxyId      string `json:"proxy_id"`
	ProxyCountry string `json:"proxy_country"`
	ProxyIp      string `json:"proxy_ip"`

	//time
	CreateTime common.CustomTime `json:"create_time" gorm:"autoCreateTime"`
	UpdatedAt  common.CustomTime `json:"updated_at" gorm:"autoUpdateTime"`
}

type TempCookies struct {
	WaterfallId string `json:"waterfall_id"`
}
