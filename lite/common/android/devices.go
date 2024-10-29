package android

import (
	"CentralizedControl/common/google"
	"CentralizedControl/common/utils"
	"strings"
)

type PkgItem struct {
	Name             string `json:"name"`
	VersionStr       string `json:"version_str"`
	VersionCode      int    `json:"version_code"`
	FirstInstallTime int64  `json:"first_install_time"`
	LastUpgradeTime  int64  `json:"last_upgrade_time"`
	TargetSdkVersion int    `json:"target_sdk_version"`
	InstallSource    string `json:"install_source"`
	Uid              int    `json:"uid"`
	AndroidId        string `json:"android_id"`
}

type Packages struct {
	Pkg             map[string]*PkgItem `json:"pkg"`
	LauncherPkgName string              `json:"launcher_pkg_name"`
	HomePkgName     string              `json:"home_pkg_name"`
	TTSPkgName      string              `json:"tts_pkg_name"`
}

func (this *Packages) GetSelf() *PkgItem {
	return this.Get(PkgNameInstagramLite)
}

func (this *Packages) Set(name string, pkg *PkgItem) {
	this.Pkg[name] = pkg
}

func (this *Packages) Get(name string) *PkgItem {
	if this.Pkg == nil {
		return nil
	}
	return this.Pkg[name]
}

func (this *Packages) HasTTS() bool {
	return this.TTSPkgName != ""
}

type Local struct {
	Country  string `json:"country"`
	Language string `json:"language"`
}

type Env struct {
	AssetsLocales        []string `json:"assets_locales"`
	AvailableLocales     []Local  `json:"available_locales"`
	Country              string   `json:"country"`
	Language             string   `json:"language"`
	Timezone             string   `json:"timezone"`
	UserAgent            string   `json:"user_agent"`
	HttpAgent            string   `json:"http_agent"`
	NotificationsEnabled bool     `json:"notifications_enabled"`
}

func (this *Env) GetLanguageLocale() string {
	return this.Language + "_" + strings.ToUpper(this.Country)
}

type Sim struct {
	SlotIndex                int    `json:"slot_index"`
	Apn                      string `json:"apn"`
	Mcc                      string `json:"mcc"`
	Mnc                      string `json:"mnc"`
	CarrierId                int    `json:"carrier_id"`
	AreaCode                 string `json:"area_code"`
	GsmSimOperatorAlpha      string `json:"gsm_sim_operator_alpha"`
	GsmSimOperatorIsoCountry string `json:"gsm_sim_operator_iso_country"`
	GsmSimOperatorNumeric    string `json:"gsm_sim_operator_numeric"`
	GsmOperatorAlpha         string `json:"gsm_operator_alpha"`
	GsmOperatorIsoCountry    string `json:"gsm_operator_iso_country"`
	GsmOperatorNumeric       string `json:"gsm_operator_numeric"`
}

type Phone struct {
	Sim       []Sim  `json:"sim"`
	PhoneType int    `json:"phone_type"` //getPhoneType
	Esn       string `json:"esn"`
	Meid      string `json:"meid"`
}

type Feature struct {
	Flags          int    `json:"flags"`
	Name           string `json:"name"`
	ReqGlEsVersion int    `json:"reqGlEsVersion"`
	Version        int    `json:"version"`
}

type Configuration struct {
	ColorMode            int     `json:"colorMode"`
	DensityDpi           int     `json:"densityDpi"`
	FontScale            float64 `json:"fontScale"`
	FontWeightAdjustment int     `json:"fontWeightAdjustment"`
	HardKeyboardHidden   int     `json:"hardKeyboardHidden"`
	Keyboard             int     `json:"keyboard"`
	KeyboardHidden       int     `json:"keyboardHidden"`
	Navigation           int     `json:"navigation"`
	NavigationHidden     int     `json:"navigationHidden"`
	Orientation          int     `json:"orientation"`
	ScreenHeightDp       int     `json:"screenHeightDp"`
	ScreenLayout         int     `json:"screenLayout"`
	ScreenWidthDp        int     `json:"screenWidthDp"`
	//SemMobileKeyboardCovered int     `json:"semMobileKeyboardCovered"`
	//SmallestScreenWidthDp    int     `json:"smallestScreenWidthDp"`
	Touchscreen int `json:"touchscreen"`
	UiMode      int `json:"uiMode"`
}

type ActivityMemory struct {
	TotalMem  int64 `json:"totalMem"`
	AvailMem  int   `json:"availMem"`
	Threshold int   `json:"threshold"`
	LowMemory bool  `json:"lowMemory"`
}

type CpuInfo struct {
	MinFrequency int `json:"min_frequency"`
	MaxFrequency int `json:"max_frequency"`
}

type ConfigurationInfo struct {
	ReqTouchScreen   int `json:"reqTouchScreen"`
	ReqKeyboardType  int `json:"reqKeyboardType"`
	ReqNavigation    int `json:"reqNavigation"`
	ReqInputFeatures int `json:"reqInputFeatures"`
	ReqGlEsVersion   int `json:"reqGlEsVersion"`
}

type WindowSize struct {
	DensityDpi    int     `json:"density_dpi"`
	RealWidth     int     `json:"real_width"`
	RealHeight    int     `json:"real_height"`
	AppWidth      int     `json:"app_width"`
	AppHeight     int     `json:"app_height"`
	StatusBar     int     `json:"status_bar"`
	NavigationBar int     `json:"navigation_bar"`
	RefreshRate   float64 `json:"refresh_rate"`
	Scale         float64 `json:"scale"`
	XDpi          float64 `json:"xdpi"`
	YDpi          float64 `json:"ydpi"`
}

type Google struct {
	HasGoogleGms          bool   `json:"has_google_gms"`
	HasGoogleAccount      bool   `json:"has_google_account"`
	HasSafetyNet          bool   `json:"has_safety_net"`
	IsInstallFromVending  bool   `json:"is_install_from_vending"`
	LoginAccount          string `json:"login_account"`
	AdvertisingId         string `json:"advertising_id"`
	LimitAdTracking       bool   `json:"limit_ad_tracking"`
	DroidguardVersionStr  string `json:"droidguard_version_str"`
	DroidguardVersionCode int    `json:"droidguard_version_code"`
}

type Props struct {
	AbiList             string `json:"ro.product.cpu.abilist"`
	AbiList32           string `json:"ro.product.cpu.abilist32"`
	AbiList64           string `json:"ro.product.cpu.abilist64"`
	ProductBoard        string `json:"ro.product.board"`
	ProductModel        string `json:"ro.product.model"`
	ProductBrand        string `json:"ro.product.brand"`
	ProductManufacturer string `json:"ro.product.manufacturer"`
	BuildVersionRelease string `json:"ro.build.version.release"`
	BuildProduct        string `json:"ro.build.product"`
	BuildId             string `json:"ro.build.id"`
	Hardware            string `json:"ro.hardware"`
	ProductDevice       string `json:"ro.product.device"`
	BootHardware        string `json:"ro.boot.hardware"`
	MediatekPlatform    string `json:"ro.mediatek.platform"`
	BoardPlatform       string `json:"ro.board.platform"`
	ChipName            string `json:"ro.chipname"`
	Serial              string `json:"ro.serialno"`
}

func (this *Props) init(p map[string]string) {
	this.AbiList = p["ro.product.cpu.abilist"]
	this.AbiList32 = p["ro.product.cpu.abilist32"]
	this.AbiList64 = p["ro.product.cpu.abilist64"]
	this.ProductBoard = p["ro.product.board"]
	this.ProductModel = p["ro.product.model"]
	this.ProductBrand = p["ro.product.brand"]
	this.ProductManufacturer = p["ro.product.manufacturer"]
	this.BuildVersionRelease = p["ro.build.version.release"]
	this.BuildProduct = p["ro.build.product"]
	this.BuildId = p["ro.build.id"]
	this.Hardware = p["ro.hardware"]
	this.ProductDevice = p["ro.product.device"]
	this.BootHardware = p["ro.boot.hardware"]
	this.MediatekPlatform = p["ro.mediatek.platform"]
	this.BoardPlatform = p["ro.board.platform"]
	this.ChipName = p["ro.chipname"]
	this.Serial = p["ro.serialno"]
}

func (this *Props) GetBuildCpuAbi() string {
	sp := strings.Split(this.AbiList64, ",")
	return sp[0]
}

func (this *Props) GetBuildCpuAbi2() string {
	sp := strings.Split(this.AbiList64, ",")
	if len(sp) > 1 {
		return sp[1]
	}
	return ""
}

type Network struct {
	NetworkType            int  `json:"network_type"`
	NetworkSubType         int  `json:"network_sub_type"`
	IsActiveNetworkMetered bool `json:"is_active_network_metered"`
}

func (this *Network) IsWifi() bool {
	return this.NetworkType == NetworkTypeTypeWifi
}

type RuntimeMemory struct {
	MaxMemory   int `json:"max_memory"`
	TotalMemory int `json:"total_memory"`
	FreeMemory  int `json:"free_memory"`
}

type Memory struct {
	MemoryClass    int            `json:"memory_class"`
	ActivityMemory ActivityMemory `json:"activity_memory"`
	RuntimeMemory  RuntimeMemory  `json:"runtime_memory"`
}

type Battery struct {
	Level  int `json:"level"`
	Scale  int `json:"scale"`
	Status int `json:"status"`
}

type System struct {
	CGroupsSupported               bool `json:"c_groups_supported"` // exists /dev/cpuctl/tasks
	BackgroundRestricted           bool `json:"background_restricted"`
	AppStandbyBucket               int  `json:"app_standby_bucket"`
	IsIgnoringBatteryOptimizations bool `json:"is_ignoring_battery_optimizations"`
	IsPowerSaveMode                bool `json:"is_power_save_mode"`
}

type Setting struct {
	Kv map[string]string
}

func (this *Setting) GetSetting(name string, _default string) string {
	v, ok := this.Kv[name]
	if ok {
		return v
	}
	return _default
}

type DiskInfo struct {
	TotalSize int `json:"total_size"`
	FreeSize  int `json:"free_size"`
}

type DirInfo struct {
	TotalSize int `json:"total_size"`
}

type Disk struct {
	Data               DiskInfo `json:"data"`
	Sdcard             DiskInfo `json:"sdcard"`
	CacheDir           DirInfo  `json:"cache_dir"`
	ExternalCacheDir   DirInfo  `json:"external_cache_dir"`
	AppDataDir         DirInfo  `json:"app_data_dir"`
	ExternalAppDataDir DirInfo  `json:"external_app_data_dir"`
}

type OpenGl struct {
	EglExtensions []string `json:"egl_extensions"`
}

type Wifi struct {
	Ssid     string `json:"ssid"`
	Bssid    string `json:"bssid"`
	LocalMac string `json:"local_mac"`
}

type Device struct {
	GmsClient         *google.GmsClient `json:"-"`
	RandTool          *utils.RandTool   `json:"-"`
	DeviceId          int64             `json:"device_id"`
	DeviceName        string            `json:"device_name"`
	DeviceType        int               `json:"device_type"`
	Props             map[string]string `json:"props"`
	UsedProps         Props             `json:"-"`
	Features          []Feature         `json:"features"`
	SharedLibrary     []string          `json:"shared_library"`
	Configuration     Configuration     `json:"configuration"`
	ConfigurationInfo ConfigurationInfo `json:"configuration_info"`
	Cpus              []CpuInfo         `json:"cpus"`
	Memory            Memory            `json:"memory"`
	Google            Google            `json:"google"`
	Packages          Packages          `json:"pkg_infos"`
	Env               Env               `json:"env"`
	Phone             Phone             `json:"phone"`
	Wifi              Wifi              `json:"wifi"`
	WindowSize        WindowSize        `json:"window_size"`
	SdkInt            int               `json:"sdk_int"`
	Permission        *Permission       `json:"permission"`
	PhoneBook         *PhoneBook        `json:"phone_book"`
	Network           Network           `json:"network"`
	Battery           Battery           `json:"battery"`
	System            System            `json:"system"`
	Settings          Setting           `json:"settings"`
	OpenGl            OpenGl            `json:"opengl"`
	Disk              Disk              `json:"disk"`
}

func (this *Device) HasFeature(feature string) bool {
	for i := range this.Features {
		if this.Features[i].Name == feature {
			return true
		}
	}
	return false
}

func (this *Device) CpuCount() int {
	return len(this.Cpus)
}

func (this *Device) MaxCpuFreq() int {
	_max := 0
	for i := range this.Cpus {
		if this.Cpus[i].MaxFrequency > _max {
			_max = this.Cpus[i].MaxFrequency
		}
	}
	return _max
}

func (this *Device) MinCpuFreq() int {
	_min := -1
	for i := range this.Cpus {
		if this.Cpus[i].MinFrequency < _min {
			_min = this.Cpus[i].MaxFrequency
		}
	}
	return _min
}

var sdkInt2AndroidVersionMap = map[int]int{
	33: 13,
	32: 12,
	31: 12,
	30: 11,
	29: 10,
}

func SdkInt2AndroidVersion(sdkInt int) int {
	return sdkInt2AndroidVersionMap[sdkInt]
}
