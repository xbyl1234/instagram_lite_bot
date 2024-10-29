package instagram

import (
	"CentralizedControl/common"
	"encoding/json"
	"fmt"
	"strings"
)

type Status int

const (
	Logout    Status = 0
	Login     Status = 1
	ShadowBan Status = -1
	Ban       Status = -2
)

type Cookies struct {
	Id     int64  `json:"id" gorm:"primaryKey"`
	Status Status `json:"status"`
	//account info
	Passwd   string `json:"passwd"`
	Username string `json:"username"`
	Email    string `json:"email"`
	//cookies
	AccountId        string `json:"account_id"`
	DeviceId         string `json:"device_id"`
	FamilyId         string `json:"family_id"`
	Claim            string `json:"claim"`
	Rur              string `json:"rur"`
	Authorization    string `json:"authorization"`
	Mid              string `json:"mid"`
	DirectRegionHint string `json:"direct_region_hint"`
	IgUShbid         string `json:"ig_u_shbid" gorm:"column:ig_u_shbid"`
	IgUShbts         string `json:"ig_u_shbts" gorm:"column:ig_u_shbts"`
	//device info
	AndroidSdkInt string `json:"android_sdk_int"`
	AndroidId     string `json:"android_id"`
	UserAgent     string `json:"user_agent"` //runtime gen
	Scale         string `json:"scale"`
	Dpi           string `json:"dpi"`
	Width         string `json:"width"`
	Height        string `json:"height"`
	Mcc           string `json:"mcc"`
	Mnc           string `json:"mnc"`
	PhoneType     string `json:"phone_type"`
	//env
	Country  string `json:"country"`
	Language string `json:"language"`
	Timezone string `json:"timezone"`
	//sys prop
	BuildId             string `json:"build_id"`
	Hardware            string `json:"hardware"`
	BuildProduct        string `json:"build_product"`
	ProductModel        string `json:"product_model"`
	ProductManufacturer string `json:"product_manufacturer"`
	ProductBrand        string `json:"product_brand"`
	//app info
	AppVersion    string `json:"app_version"`
	AppVersionInt string `json:"app_version_int"`
	//other
	ProxyId      string `json:"proxy_id"`
	ProxyCountry string `json:"proxy_country"`
	ProxyIp      string `json:"proxy_ip"`
	//time
	CreateTime common.CustomTime `json:"create_time" gorm:"autoCreateTime"`
	UpdatedAt  common.CustomTime `json:"updated_at" gorm:"autoUpdateTime"`
}

type TempCookies struct {
	PigeonSessionId      string
	TimezoneOffset       string
	BandwidthTotaltimeMs int
}

func ConvDeviceFile2Cookies(ck string) *Cookies {
	var msgCk Cookies
	err := json.Unmarshal([]byte(ck), &msgCk)
	if err != nil {
		return nil
	}
	//									Instagram 275.0.0.27.98 Android (28/9; 320dpi; 720x1280; ro.product.manufacturer/ro.product.brand; ro.product.model; gracelte; qcom; zh_CN; 458229258)
	//									Instagram 275.0.0.27.98 Android (28/9; %d dpi; %dx%d; 	 %s; 									   %s; 			     %s; 	   %s;   %s; 458229258
	//									   Instagram 275.0.0.27.98 Android (28/9; 320dpi; 720x1280; vivo; V1824A; gracelte; qcom; zh_CN; 458229258)
	msgCk.UserAgent = fmt.Sprintf("Instagram 275.0.0.27.98 Android (28/9; %sdpi; %sx%s; %s; %s; %s; %s; %s; 458229258)",
		msgCk.Dpi,
		msgCk.Width,
		msgCk.Height,
		msgCk.ProductBrand,
		msgCk.ProductModel,
		msgCk.BuildProduct,
		msgCk.Hardware,
		strings.ToLower(msgCk.Language)+"_"+strings.ToUpper(msgCk.Country),
	)

	return &msgCk
}
