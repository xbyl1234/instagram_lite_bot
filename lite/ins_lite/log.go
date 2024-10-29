package ins_lite

import (
	"CentralizedControl/common/android"
	"CentralizedControl/ins_lite/proto"
	"CentralizedControl/ins_lite/proto/msg/sender"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

func getNetworkType(networkType int) string {
	switch networkType {
	case 1:
		return "GPRS"
	case 2:
		return "EDGE"
	case 3:
		return "UMTS"
	case 4:
		return "CDMA"
	case 5:
		return "EVDO_0"
	case 6:
		return "EVDO_A"
	case 7:
		return "1xRTT"
	case 8:
		return "HSDPA"
	case 9:
		return "HSUPA"
	case 10:
		return "HSPA"
	case 11:
		return "IDEN"
	case 12:
		return "EVDO_B"
	case 13:
		return "LTE"
	case 14:
		return "EHRPD"
	case 15:
		return "HSPAP"
	default:
		return "UNKNOWN"
	}
}

func (this *InsLiteClient) SendJsonLog(logName string, module string, extra any) {
	liteLog := &sender.LiteLog{}
	liteLog.SessionId = this.Cookies.GenDeviceTimeId
	if liteLog.SessionId == "" {
		liteLog.SessionId = "client_event"
	}
	liteLog.Time = float64(time.Now().UnixMilli())
	liteLog.Name = logName
	liteLog.LogId = ""
	liteLog.Module = module
	liteLog.Extra = extra

	s := &proto.Message[sender.ClientEventLog]{}
	body := &s.Body
	marshal, _ := json.Marshal(liteLog)
	body.Log = string(marshal)
	body.Type1 = 1
	body.Type2 = 255
	this.SendMsg(s)
}

func getPhoneType(phoneType int) string {
	switch phoneType {
	case 0:
		return "NONE"
	case 1:
		return "GSM"
	case 2:
		return "CDMA"
	case 3:
		return "SIP"
	default:
		return "UNKNOWN"
	}
}

func getTimeFormatStr(v int64, timezone string) string {
	location, err := time.LoadLocation(timezone)
	if err != nil {
		panic("Invalid timezone")
	}
	date := time.Unix(v, 0).In(location)
	return date.Format("2006-01-02T15:04:05.999") + date.Format("-07:00")
}

func getAppStandbyBucket(s int) string {
	switch s {
	case 10:
		return "STANDBY_BUCKET_ACTIVE"
	case 20:
		return "STANDBY_BUCKET_WORKING_SET"
	case 30:
		return "STANDBY_BUCKET_FREQUENT"
	case 40:
		return "STANDBY_BUCKET_RARE"
	case 45:
		return "STANDBY_BUCKET_RESTRICTED"
	default:
		return fmt.Sprintf("STANDBY_BUCKET_UNKNOWN_%d", s)
	}
}
func getSimState(state int) string {
	switch 5 {
	case 1:
		return "ABSENT"
	case 2:
		return "PIN_REQUIRED"
	case 3:
		return "PUK_REQUIRED"
	case 4:
		return "NETWORK_LOCKED"
	case 5:
		return "READY"
	default:
		return "UNKNOWN"
	}
}

func (this *InsLiteClient) makeSimInfo(sim android.Sim) *sender.LogSimInfo {
	result := &sender.LogSimInfo{
		Index:             sim.SlotIndex,
		State:             getSimState(5),
		Carrier:           sim.GsmOperatorAlpha,
		SimCarrierName:    nil,
		SimDisplayName:    nil,
		CarrierCountryIso: sim.GsmOperatorIsoCountry,
		PhoneType:         getPhoneType(this.Cookies.Phone.PhoneType),
		NetworkType:       getNetworkType(this.Cookies.Network.NetworkType),
		CountryIso:        sim.GsmSimOperatorIsoCountry,
		Operator:          sim.GsmSimOperatorNumeric,
		SimOperatorName:   sim.GsmSimOperatorAlpha,
		DeviceLocale:      this.Cookies.Env.GetLanguageLocale(),
	}
	if this.Cookies.Network.IsWifi() {
		result.NetworkInfo = ""
	} else {

	}
	if this.Cookies.Permission.IsAllow(android.ReadPhoneState) {
		result.PhoneNumber = ""
		result.SerialNumber = "SecurityException"
		result.SubscriberId = "SecurityException"
	}
	return result
}

func (this *InsLiteClient) getClassYears() int {
	var classYear int
	var totalMem = this.Cookies.Memory.ActivityMemory.TotalMem
	var cpuCount = this.Cookies.CpuCount()
	var cpuMaxFreq = this.Cookies.MaxCpuFreq()
	if totalMem <= 0x30000000 {
		classYear = 2010
		if cpuCount <= 1 {
			classYear = 2009
		}
	} else {
		classYear = 2012
		if totalMem > 0x40000000 {
			if totalMem > 0x60000000 {
				if totalMem <= 0x80000000 {
					classYear = 2013
					return classYear
				}
				if totalMem <= 0xC0000000 {
					classYear = 2014
				} else {
					classYear = 2016
					if totalMem <= 0x140000000 {
						classYear = 2015
						return classYear
					}
				}
			} else if cpuMaxFreq >= 1800000 {
				classYear = 2013
				return classYear
			}
		} else if cpuMaxFreq < 1300000 {
			classYear = 2011
		}
	}
	return classYear
}

func (this *InsLiteClient) createModuleDownloadInfo() *sender.LogModuleInfo {
	downloadInfo := &sender.LogModuleInfo{}
	downloadInfo.Downloader = "Facebook"
	downloadInfo.Modules = make([]*sender.LogDownloadInfo, len(sender.LogDownloadModuleList))
	for idx := range sender.LogDownloadModuleList {
		downloadInfo.Modules[idx] = &sender.LogDownloadInfo{
			Name: sender.LogDownloadModuleList[idx],
		}
	}

	setValue := func(item *sender.LogDownloadInfo) {
		item.FirstRequestWasPrefetch = new(string)
		item.InitialInstallRequestTimestamp = new(int64)
		item.InitialPrefetchTime = new(int64)
		item.LastPrefetchTime = new(int64)
		_time := time.Now().UnixMilli()
		*item.FirstRequestWasPrefetch = "true"
		*item.InitialInstallRequestTimestamp = _time
		*item.InitialPrefetchTime = _time
		*item.LastPrefetchTime = _time
	}
	setValue(downloadInfo.GetModule("fizz"))
	setValue(downloadInfo.GetModule("s_fizz_mediastreaming"))
	setValue(downloadInfo.GetModule("s_fizz_msys"))
	setValue(downloadInfo.GetModule("shared_fizz_ms_profilo"))
	return downloadInfo
}

func (this *InsLiteClient) createLogDeviceInfo() *sender.LogDeviceInfo {
	deviceInfo := &sender.LogDeviceInfo{
		Pk:                      this.Cookies.Pk,
		Carrier:                 this.Cookies.Phone.Sim[0].GsmOperatorAlpha,
		CarrierCountryIso:       this.Cookies.Phone.Sim[0].GsmOperatorIsoCountry,
		NetworkType:             getNetworkType(this.Cookies.Network.NetworkType),
		PhoneType:               getPhoneType(this.Cookies.Phone.PhoneType),
		SimCountryIso:           this.Cookies.Phone.Sim[0].GsmSimOperatorIsoCountry,
		SimOperator:             this.Cookies.Phone.Sim[0].GsmSimOperatorAlpha,
		Locale:                  this.Cookies.Env.GetLanguageLocale(),
		DeviceType:              this.Cookies.UsedProps.ProductModel,
		Brand:                   this.Cookies.UsedProps.ProductBrand,
		Manufacturer:            this.Cookies.UsedProps.ProductManufacturer,
		OsType:                  "Android",
		OsVer:                   this.Cookies.UsedProps.BuildVersionRelease,
		CpuAbi:                  this.Cookies.UsedProps.GetBuildCpuAbi(),
		CpuAbi2:                 this.Cookies.UsedProps.GetBuildCpuAbi2(),
		CpuAbiList:              this.Cookies.UsedProps.AbiList,
		CpuBoardPlatform:        this.Cookies.UsedProps.BoardPlatform,
		CpuBuildBoard:           this.Cookies.UsedProps.ProductBoard,
		CpuMediatekPlatform:     this.Cookies.UsedProps.MediatekPlatform,
		CpuChipName:             this.Cookies.UsedProps.ChipName,
		TouchFeature:            this.Cookies.HasFeature("android.hardware.touchscreen"),
		RamLowFeature:           this.Cookies.HasFeature("android.hardware.ram.low"),
		TouchConfig:             this.Cookies.Configuration.Touchscreen,
		KeyboardConfig:          this.Cookies.Configuration.Keyboard,
		UnreliableCoreCount:     this.Cookies.CpuCount(),
		ReliableCoreCount:       this.Cookies.CpuCount(),
		CpuMaxFreq:              this.Cookies.MaxCpuFreq(),
		LowPowerCpuMaxFreq:      this.Cookies.MinCpuFreq(),
		TotalMem:                this.Cookies.Memory.ActivityMemory.TotalMem,
		YearClass:               this.getClassYears(),
		CgroupsSupported:        this.Cookies.System.CGroupsSupported,
		FirstInstallTime:        getTimeFormatStr(this.Cookies.Packages.GetSelf().FirstInstallTime, "America/New_York"),
		LastUpgradeTime:         getTimeFormatStr(this.Cookies.Packages.GetSelf().LastUpgradeTime, "America/New_York"),
		InstallLocation:         "internal_storage",
		Density:                 this.Cookies.WindowSize.Scale,
		ScreenWidth:             this.Cookies.WindowSize.AppWidth,
		ScreenHeight:            this.Cookies.WindowSize.AppHeight,
		FrontCamera:             this.Cookies.HasFeature("android.hardware.camera.front"),
		RearCamera:              this.Cookies.HasFeature("android.hardware.camera"),
		AllowsNonMarketInstalls: this.Cookies.Settings.GetSetting("install_non_market_apps", "0"),
		AndroidId:               this.Cookies.Packages.GetSelf().AndroidId,
		IsBatteryOptimized:      !this.Cookies.System.IsIgnoringBatteryOptimizations,
		IsPowerSaveMode:         this.Cookies.System.IsPowerSaveMode,
		OpenglVersion:           this.Cookies.ConfigurationInfo.ReqGlEsVersion,
		DeviceTotalSpace:        strconv.Itoa(this.Cookies.Disk.Data.TotalSize),
		ExternalCacheSize:       strconv.Itoa(this.Cookies.Disk.ExternalCacheDir.TotalSize),
		DeviceFreeSpace:         strconv.Itoa(this.Cookies.Disk.Data.FreeSize),
		CacheSize:               strconv.Itoa(this.Cookies.Disk.CacheDir.TotalSize),
		AppDataSize:             strconv.Itoa(this.Cookies.Disk.AppDataDir.TotalSize),
		ExternalAppDataSize:     strconv.Itoa(this.Cookies.Disk.ExternalAppDataDir.TotalSize - this.Cookies.Disk.ExternalCacheDir.TotalSize),
		SdFreeSpace:             strconv.Itoa(this.Cookies.Disk.Sdcard.FreeSize),
		SdTotalSpace:            strconv.Itoa(this.Cookies.Disk.Sdcard.TotalSize),
		AllowAdsTracking:        !this.Cookies.Google.LimitAdTracking,
		TargetSdk:               this.Cookies.Packages.GetSelf().TargetSdkVersion,
		VersionCode:             this.Cookies.Packages.GetSelf().VersionCode,
		LauncherName:            this.Cookies.Packages.LauncherPkgName,
		AccessibilityEnabled:    false,
		TouchExplorationEnabled: false,
		ModuleInfo:              this.createModuleDownloadInfo(),
	}

	if this.Cookies.SdkInt >= 28 {
		deviceInfo.IsBackgroundRestricted = &this.Cookies.System.BackgroundRestricted
		appStandbyBucket := getAppStandbyBucket(this.Cookies.System.AppStandbyBucket)
		deviceInfo.AppStandbyBucket = &appStandbyBucket
	}

	if this.Cookies.Google.HasGoogleGms {
		deviceInfo.GooglePlayServicesInstallation = "SERVICE_ENABLED"
		deviceInfo.GooglePlayServicesVersion = this.Cookies.Packages.Get("").VersionCode
		deviceInfo.GoogleAccounts = 1
		deviceInfo.AdvertiserId = &this.Cookies.Google.AdvertisingId

		pkgInstalls := make([]*sender.PkgInstallStatus, 0)
		pkgInstalls = append(pkgInstalls, &sender.PkgInstallStatus{
			InstallationStatus: "SERVICE_ENABLED",
			PackageName:        android.PkgNameVending,
			Version:            this.Cookies.Packages.Get(android.PkgNameVending).VersionCode,
		})
		pkg := this.Cookies.Packages.Get(android.PkgNameGoogleMarket)
		if pkg != nil {
			pkgInstalls = append(pkgInstalls, &sender.PkgInstallStatus{
				InstallationStatus: "SERVICE_ENABLED",
				PackageName:        android.PkgNameGoogleMarket,
				Version:            this.Cookies.Packages.Get(android.PkgNameGoogleMarket).VersionCode,
			})
		}
		pkg = this.Cookies.Packages.Get(android.PkgNameFinsky)
		if pkg != nil {
			pkgInstalls = append(pkgInstalls, &sender.PkgInstallStatus{
				InstallationStatus: "SERVICE_ENABLED",
				PackageName:        android.PkgNameFinsky,
				Version:            this.Cookies.Packages.Get(android.PkgNameFinsky).VersionCode,
			})
		}
		deviceInfo.GooglePlayStore = pkgInstalls

		pkg = this.Cookies.Packages.Get(android.PkgNameGsf)
		if pkg != nil {
			deviceInfo.GsfInstallationStatus = &sender.PkgInstallStatus{
				InstallationStatus: "SERVICE_ENABLED",
				PackageName:        android.PkgNameGsf,
				Version:            this.Cookies.Packages.Get(android.PkgNameGsf).VersionCode,
			}
		}
	} else {
		//"SERVICE_DISABLED""SERVICE_INVALID"
		deviceInfo.GooglePlayServicesInstallation = "SERVICE_MISSING"
		deviceInfo.GooglePlayServicesVersion = -1
		deviceInfo.GoogleAccounts = 0
		deviceInfo.AdvertiserId = nil
	}

	if this.Cookies.Packages.GetSelf().InstallSource != "" {
		deviceInfo.Installer = this.Cookies.Packages.GetSelf().InstallSource
		deviceInfo.OriginalInstaller = this.Cookies.Packages.GetSelf().InstallSource
	} else {
		deviceInfo.Installer = "UNKNOWN"
		deviceInfo.OriginalInstaller = "UNKNOWN"
	}

	if this.Cookies.SdkInt >= 28 {
		androidStrongboxAvailable := this.Cookies.HasFeature("android.hardware.strongbox_keystore")
		deviceInfo.AndroidStrongboxAvailable = &androidStrongboxAvailable
		if this.Cookies.SdkInt >= 31 {
			androidTeeAvailable := this.Cookies.HasFeature("android.hardware.hardware_keystore")
			deviceInfo.AndroidTeeAvailable = &androidTeeAvailable
		}
	}

	permissionGroups := 0
	perms := []string{
		android.ReadCalendar,
		android.Camera,
		android.ReadContacts,
		android.AccessFineLocation,
		android.RecordAudio,
		android.ReadPhoneState,
		android.BodySensors,
		android.ReadSms,
		android.ReadExternalStorage,
	}
	for idx := range perms {
		if this.Cookies.Permission.IsAllow(perms[idx]) {
			permissionGroups |= 1 << idx
		}
	}
	deviceInfo.PermissionGroups = permissionGroups

	if this.Cookies.PropStore54.GetBool(0xE34, false) {
		insPkgStatus := &sender.InsPkgStatus{
			InstagramLiteInstallationStatus: true,
		}
		pkg := this.Cookies.Packages.Get(android.PkgNameInstagram)
		if pkg != nil {
			insPkgStatus.InstagramAndroidInstallationStatus = true
		}
		deviceInfo.InstagramAppsInstallationStatus = insPkgStatus
	}
	simInfo := make([]*sender.LogSimInfo, 0)
	simInfo = append(simInfo, this.makeSimInfo(this.Cookies.Phone.Sim[0]))
	deviceInfo.SimInfo = simInfo
	//ModuleInfo:              sender.LogModuleInfo{},
	//AmazonAppStoreInstallationStatus
	return deviceInfo
}
