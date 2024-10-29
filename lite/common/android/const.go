package android

const (
	BatteryStatusCharging    = 2 //充电中
	BatteryStatusDischarging = 3 //放电中
	BatteryStatusNotCharging = 4 //未充电
	BatteryStatusFull        = 5 //已充满
	BatteryStatusUnknown     = 1 //状态未知
)

const (
	SimStateUnknown  = 0
	SimStateAbsent   = 1
	SimStateReady    = 5
	SimStateNotReady = 6
)

const (
	DeviceTypeOther   = -1
	DeviceTypeXiaoMi  = 0
	DeviceTypeHuaWei  = 1
	DeviceTypeHonor   = 2
	DeviceTypeOppo    = 3
	DeviceTypeVivo    = 4
	DeviceTypeRealMe  = 5
	DeviceTypeSamsung = 6
	DeviceTypeRedmi   = 7
	DeviceTypeSony    = 9
	DeviceTypeAsus    = 10
)

const (
	PhoneTypeNone = 0
	PhoneTypeGsm  = 1
	PhoneTypeCdma = 2
	PhoneTypeSip  = 3
	PhoneTypeIms  = 5
)

const (
	NetworkTypeTypeMobile     = 0
	NetworkTypeTypeWifi       = 1
	NetworkTypeMobileMms      = 2
	NetworkTypeMobileSupl     = 3
	NetworkTypeMobileDun      = 4
	NetworkTypeMobileHipri    = 5
	NetworkTypeWiMax          = 6
	NetworkTypeBluetooth      = 7
	NetworkTypeDummy          = 8
	NetworkTypeEthernet       = 9
	NetworkNetworkTypeTypeVpn = 17
)

const (
	NetworkSubTypeUnknown = 0
	NetworkSubTypeGprs    = 1
	NetworkSubTypeEdge    = 2
	NetworkSubTypeUmts    = 3
	NetworkSubTypeCdma    = 4
	NetworkSubTypeEvdo0   = 5
	NetworkSubTypeEvdoA   = 6
	NetworkSubType1xRtt   = 7
	NetworkSubTypeHsdpa   = 8
	NetworkSubTypeHsupa   = 9
	NetworkSubTypeHspa    = 10
	NetworkSubTypeIden    = 11
	NetworkSubTypeEvdoB   = 12
	NetworkSubTypeLte     = 13
	NetworkSubTypeEhrpd   = 14
	NetworkSubTypeHspap   = 15
	NetworkSubTypeGsm     = 16
	NetworkSubTypeTdScdma = 17
	NetworkSubTypeIwlan   = 18
	NetworkSubTypeLteCa   = 19
	NetworkSubTypeNr      = 20
)

var networkSubType2Name = map[int]string{
	NetworkSubTypeGprs:    "GPRS",
	NetworkSubTypeEdge:    "EDGE",
	NetworkSubTypeUmts:    "UMTS",
	NetworkSubTypeHsdpa:   "HSDPA",
	NetworkSubTypeHsupa:   "HSUPA",
	NetworkSubTypeHspa:    "HSPA",
	NetworkSubTypeCdma:    "CDMA",
	NetworkSubTypeEvdo0:   "CDMA - EvDo rev. 0",
	NetworkSubTypeEvdoA:   "CDMA - EvDo rev. A",
	NetworkSubTypeEvdoB:   "CDMA - EvDo rev. B",
	NetworkSubType1xRtt:   "CDMA - 1xRTT",
	NetworkSubTypeLte:     "LTE",
	NetworkSubTypeEhrpd:   "CDMA - eHRPD",
	NetworkSubTypeIden:    "iDEN",
	NetworkSubTypeHspap:   "HSPA+",
	NetworkSubTypeGsm:     "GSM",
	NetworkSubTypeTdScdma: "TD_SCDMA",
	NetworkSubTypeIwlan:   "IWLAN",
	NetworkSubTypeLteCa:   "LTE_CA",
	NetworkSubTypeNr:      "NR",
}

func NetworkSubType2Name(t int) string {
	return networkSubType2Name[t]
}

const (
	PkgNameVending       = "com.android.vending"
	PkgNameGms           = "com.google.android.gms"
	PkgNameGsf           = "com.google.android.gsf"
	PkgNameGoogleMarket  = "com.google.market"
	PkgNameFinsky        = "com.google.android.finsky"
	PkgNameAmazonVenezia = "com.amazon.venezia"
	PkgNameAmazonShop    = "com.amazon.mShop.android"
	PkgNameInstagram     = "com.instagram.android"
	PkgNameInstagramLite = "com.instagram.lite"
)
