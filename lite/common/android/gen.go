package android

import (
	"CentralizedControl/common/google"
	"CentralizedControl/common/utils"
	"fmt"
	"strings"
)

func ChoiceFirstName(randTool *utils.RandTool) string {
	return ChoiceUsername(randTool)
}

func ChoiceSecondName(randTool *utils.RandTool) string {
	return ChoiceUsername(randTool)
}

func ChoiceUsername(randTool *utils.RandTool) string {
	return utils.ChoseOne2(randTool, Resource.username)
}

func ChoiceIco(randTool *utils.RandTool) string {
	return utils.ChoseOne2(randTool, Resource.ico)
}

func GenGoogleEmail(randTool *utils.RandTool) string {
	return randTool.GenString(utils.CharSet_abc, 8) + "@gmail.com"
}

func genAndroidId(randTool *utils.RandTool) string {
	return randTool.GenString(utils.CharSet_16_Num, 16)
}

func randomMeid(randTool *utils.RandTool) string {
	meid := "35503104"
	for i := 0; i < 6; i++ {
		meid += fmt.Sprintf("%d", randTool.GenNumber(0, 9))
	}
	sum := 0
	for i := 0; i < len(meid); i++ {
		c := int(meid[i] - '0')
		if (len(meid)-i-1)%2 == 0 {
			c *= 2
			c = c%10 + c/10
		}
		sum += c
	}
	check := (100 - sum) % 10
	meid += fmt.Sprintf("%d", check)
	return meid
}

type GenDeviceConfig func(device *Device)

func DeviceConfigHasGms(hasGms bool, hasSafetyNet bool, hasGoogleAccount bool, isInstallFromVending bool) GenDeviceConfig {
	return func(device *Device) {
		if hasGms {
			device.Google.HasGoogleGms = true
			device.Google.IsInstallFromVending = isInstallFromVending
			device.Google.AdvertisingId = device.RandTool.GenUUID()
			device.Google.LimitAdTracking = false
		}
		device.Google.HasSafetyNet = hasSafetyNet
		device.Google.HasGoogleAccount = hasGoogleAccount
		device.GmsClient = google.GmsManger.GetClient("")
		if hasGoogleAccount {
			if device.Google.LoginAccount == "" {
				var err error
				device.Google.LoginAccount, err = device.GmsClient.GetAccount()
				if err != nil {
					panic(err)
				}
			}
		}
	}
}

func DeviceConfigGenPerm(allowPerms []string) GenDeviceConfig {
	return func(device *Device) {
		device.Permission = &Permission{
			Perm: map[string]*PermissionInfo{},
		}
		for _, item := range AllPermNames {
			device.Permission.Perm[item] = &PermissionInfo{
				CurStatic:    PermStatusUnReq,
				ShouldStatic: PermStatusDenied,
			}
		}
		for _, item := range allowPerms {
			device.Permission.Perm[item].ShouldStatic = PermStatusGranted
		}
	}
}

func DeviceConfigGenPhoneBook() GenDeviceConfig {
	return func(device *Device) {
		count := device.RandTool.GenNumber(5, 20)
		device.PhoneBook = &PhoneBook{
			Item: make([]PhoneBookItem, count),
		}
		for i := 0; i < count; i++ {
			device.PhoneBook.Item[i] = createPhoneBookItem(device.RandTool, device.Env.Country)
		}
	}
}

var MobileNetTyps = []int{
	NetworkSubTypeNr,    //5g
	NetworkSubTypeLte,   //4g
	NetworkSubTypeHspap, //3g
}

func DeviceConfigGenNetwork(useWifi bool) GenDeviceConfig {
	return func(device *Device) {
		if useWifi {
			device.Network.NetworkType = NetworkTypeTypeWifi
			device.Network.NetworkSubType = NetworkSubTypeUnknown
		} else {
			device.Network.NetworkType = NetworkTypeTypeMobile
			device.Network.NetworkSubType = utils.ChoseOne2(device.RandTool, MobileNetTyps)
		}
	}
}

func (this *Device) InitDevice(country string, funcs ...GenDeviceConfig) {
	cty := Resource.deviceResource.Country[country]
	if cty == nil {
		//panic("not find country " + country)
		country = "us"
		cty = Resource.deviceResource.Country[country]
	}
	sim := Resource.deviceResource.Sim[country]
	if sim == nil {
		//panic("not find sim " + country)
		sim = Resource.deviceResource.Sim["us"]
	}
	this.Env.Country = cty.Country
	this.Env.Language = utils.ChoseOne2(this.RandTool, cty.Language)
	this.Env.Timezone = cty.Timezone
	this.Env.NotificationsEnabled = true

	sim1Info := utils.ChoseOne2(this.RandTool, sim)
	mccmnc := utils.ChoseOne2(this.RandTool, sim1Info.MccMnc)
	sim1 := Sim{
		Apn:                      utils.ChoseOne2(this.RandTool, mccmnc.Apn),
		Mcc:                      mccmnc.Mcc,
		Mnc:                      mccmnc.Mnc,
		CarrierId:                sim1Info.CarrierId,
		AreaCode:                 sim1Info.AreaCode,
		GsmSimOperatorAlpha:      sim1Info.SimOperatorNameEn,
		GsmSimOperatorIsoCountry: sim1Info.SimCountryIso,
		GsmSimOperatorNumeric:    mccmnc.Mcc + mccmnc.Mnc,
		GsmOperatorAlpha:         sim1Info.SimOperatorNameCn,
		GsmOperatorIsoCountry:    sim1Info.SimCountryIso,
		GsmOperatorNumeric:       mccmnc.Mcc + mccmnc.Mnc,
	}
	this.Phone.Sim = []Sim{sim1}
	//this.Phone.Esn = this.randTool.GenString() 8 bytes
	this.Phone.Meid = randomMeid(this.RandTool)

	this.Google.LoginAccount = ""
	this.Battery.Level = this.RandTool.GenNumber(int(float32(this.Battery.Scale)*0.2), this.Battery.Scale)
	this.Battery.Status = BatteryStatusDischarging
	for _, item := range funcs {
		item(this)
	}

	if this.Packages.Pkg == nil {
		this.Packages.Pkg = make(map[string]*PkgItem)
	}
	this.Packages.Set(PkgNameVending, &PkgItem{
		Name:        PkgNameVending,
		VersionStr:  Resource.AppConfig.VendingVersionStr,
		VersionCode: Resource.AppConfig.VendingVersionCode,
	})
	this.Packages.Set(PkgNameGms, &PkgItem{
		Name:        PkgNameGms,
		VersionStr:  Resource.AppConfig.GmsVersionStr,
		VersionCode: Resource.AppConfig.GmsVersionCode,
	})
	//this.Packages.Set(PkgNameGsf, &PkgItem{
	//	Name:        PkgNameGsf,
	//	VersionStr:  Resource.AppConfig.GsfVersionStr,
	//	VersionCode: Resource.AppConfig.GsfVersionCode,
	//})
	installSource := ""
	if this.Google.IsInstallFromVending {
		installSource = PkgNameVending
	}
	this.Packages.Set(PkgNameInstagramLite, &PkgItem{
		Name:             PkgNameInstagramLite,
		VersionStr:       Resource.AppConfig.InsLiteVersionStr,
		VersionCode:      Resource.AppConfig.InsLiteVersionCode,
		FirstInstallTime: 0,
		LastUpgradeTime:  0,
		TargetSdkVersion: 0,
		InstallSource:    installSource,
		Uid:              this.RandTool.GenNumber(10100, 10300),
		AndroidId:        genAndroidId(this.RandTool),
	})
	this.Google.DroidguardVersionStr = Resource.AppConfig.DroidguardVersionStr
	this.Google.DroidguardVersionCode = Resource.AppConfig.DroidguardVersionCode

	this.Props["ro.serialno"] = strings.ToUpper(this.RandTool.GenString(utils.CharSet_16_Num, 16))
	this.UsedProps.init(this.Props)

	this.Wifi.Ssid = this.RandTool.GenString(utils.CharSet_All, 10)
	this.Wifi.Bssid = this.RandTool.GenString(utils.CharSet_16_Num, 12)
	this.Wifi.LocalMac = this.RandTool.GenString(utils.CharSet_16_Num, 12)
}

func GetAndroidDevice(name string, id int64) *Device {
	var device Device
	randTool := utils.CreateRandTool(id)
	if name == "" {
		device = utils.ChoseOne2(randTool, Resource.androidDevice)
	} else {
		find := false
		for idx := range Resource.androidDevice {
			if Resource.androidDevice[idx].DeviceName == name {
				device = Resource.androidDevice[idx]
				find = true
				break
			}
		}
		if !find {
			loadOneDevices(name, &device)
		}
	}
	device.DeviceId = id
	device.RandTool = randTool
	return &device
}
