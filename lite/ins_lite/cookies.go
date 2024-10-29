package ins_lite

import (
	"CentralizedControl/common"
	"CentralizedControl/common/android"
	"CentralizedControl/common/utils"
	"CentralizedControl/ins_lite/proto/msg/recver"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

type TempCookies struct {
	NetQualityType int
	InstanceKey    uint32
}

type AppInfo struct {
	McQueryHashBin []byte `json:"mc_query_hash_bin"`
	BloksVersionID string `json:"bloks_version_id"`
}

type Session struct {
	ClientEncryptionSecret   []byte `json:"client_encryption_secret"`
	SessionId                int64  `json:"session_id"`
	TransientToken           int32  `json:"transient_token"`
	StickinessTokenTimeStamp int64  `json:"stickiness_token_time_stamp"`
}

const (
	StateNoLogin  = 0
	StateLoggedIn = 1
)

type Cookies struct {
	Id               int64 `json:"id" gorm:"primaryKey"`
	*android.Account `json:"account"`
	*android.IpInfo  `json:"ip_info"`
	*android.VpnInfo `json:"vpn_info"`
	*android.Device  `json:"device"`
	*AppInfo         `json:"app_info"`

	Temp            *TempCookies `json:"temp"`
	GenDeviceTimeId string       `json:"gen_device_time_id"` //session_id
	DeviceId        string       `json:"device_id"`          //imie
	PhoneId         string       `json:"phone_id"`
	Pk              string       `json:"pk"`

	//state
	State       int                `json:"state"`
	Session     Session            `json:"session"`
	Properties  map[string]string  `json:"properties"`
	PropStore54 recver.PropStore54 `json:"prop_store_54"`
	Prop223     map[int]string     `json:"prop_223"`
	prop223Lock sync.Mutex
	//time
	LunchCount int               `json:"lunch_count"`
	CreateTime common.CustomTime `json:"create_time" gorm:"autoCreateTime"`
	UpdatedAt  common.CustomTime `json:"updated_at" gorm:"autoUpdateTime"`
}

func CreateTempCookies() *TempCookies {
	ck := TempCookies{}
	return &ck
}

func CreateNewCookies(device *android.Device) *Cookies {
	ck := Cookies{
		Id:      0,
		Account: nil,
		IpInfo:  nil,
		VpnInfo: nil,
		Device:  device,
		AppInfo: &AppInfo{
			McQueryHashBin: utils.DecodeHex(android.Resource.AppConfig.McQueryHashBin),
			BloksVersionID: android.Resource.AppConfig.BloksVersionId,
		},
		Temp: CreateTempCookies(),
		//GenDeviceTimeId: fmt.Sprintf("%v-fg", uint64(math_rand.Int())&0xFFFFFFFF|uint64(time.Now().UnixMilli())<<0x20),
		GenDeviceTimeId: fmt.Sprintf("%s-fg", utils.GenString(utils.CharSet_123, 16)),
		DeviceId:        utils.GenUUID(),
		PhoneId:         utils.GenUUID(),
		Properties:      map[string]string{},
		PropStore54: recver.PropStore54{
			Props: map[int]any{},
		},
		Prop223:    map[int]string{},
		LunchCount: 0,
		CreateTime: common.CustomTime(time.Now()),
	}
	return &ck
}

func LoadCookies(path string) *Cookies {
	file, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	var ck Cookies
	err = json.Unmarshal(file, &ck)
	if err != nil {
		panic(err)
	}
	return &ck
}

func SaveCookies(path string, ck *Cookies) {
	marshal, err := json.Marshal(ck)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(path, marshal, 0666)
	if err != nil {
		panic(err)
	}
}
