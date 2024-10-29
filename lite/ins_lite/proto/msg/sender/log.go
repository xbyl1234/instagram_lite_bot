package sender

type LiteLog struct {
	Name      string  `json:"name"`
	Time      float64 `json:"time"`
	SessionId string  `json:"session_id"`
	LogId     string  `json:"log_id"`
	Module    string  `json:"module"`
	Extra     any     `json:"extra,omitempty"`
}

var LogDownloadModuleList = []string{
	"blokscamera",
	"boost",
	"fizz",
	"inappbrowser",
	"libunwindstack",
	"mediacompositionplayer",
	"mediastreaming",
	"mns",
	"mnshttp",
	"msys",
	"msysinfra",
	"profilo",
	"pytorch",
	"rtc",
	"s_blokscamera_boost",
	"s_blokscamera_msysinfra",
	"s_blokscamera_rtc",
	"s_boost_fizz_mediastreaming",
	"s_boost_mediastreaming",
	"s_boost_profilo",
	"s_fizz_mediastreaming",
	"s_fizz_msys",
	"s_inappbrowser_mediastreaming",
	"s_inappbrowser_mediastreaming_rtc",
	"s_libunwindstack_profilo",
	"s_mediastreaming_msys_profilo_rtc",
	"s_mediastreaming_msys_rtc",
	"s_mediastreaming_msysinfra",
	"s_mediastreaming_rtc",
	"s_mnshttp_msys",
	"s_mnshttp_msysinfra",
	"s_msys_rtc",
	"s_pytorch_rtc",
	"shared_fizz_ms_profilo",
	"uiqr",
}

type LogDownloadInfo struct {
	Name                    string  `json:"name"`
	DownloadState           *string `json:"downloadState,omitempty"`
	Packaging               *string `json:"packaging,omitempty"`
	LastForegroundUse       *string `json:"lastForegroundUse,omitempty"`
	LastEntryTime           *string `json:"lastEntryTime,omitempty"`
	FirstEntryTime          *string `json:"firstEntryTime,omitempty"`
	FirstRequestWasPrefetch *string `json:"firstRequestWasPrefetch,omitempty"`
	IsSplitInstalled        *bool   `json:"isSplitInstalled,omitempty"`
	IsDownloaded            *bool   `json:"isDownloaded,omitempty"`
	DownloadAttempted       *bool   `json:"downloadAttempted,omitempty"`
	ForegroundRequested     *bool   `json:"foregroundRequested,omitempty"`
	//	isSplitInstalled        bool
	InitialInstallRequestTimestamp *int64 `json:"initialInstallRequestTimestamp,omitempty"`
	InstallLatency                 *int64 `json:"installLatency,omitempty"`
	InitialLoadTimestamp           *int64 `json:"initialLoadTimestamp,omitempty"`
	LastLoadTimestamp              *int64 `json:"lastLoadTimestamp,omitempty"`
	InitialUninstallRequestTime    *int64 `json:"initialUninstallRequestTime,omitempty"`
	LastUninstallCompleteTime      *int64 `json:"lastUninstallCompleteTime,omitempty"`
	InitialPrefetchTime            *int64 `json:"initialPrefetchTime,omitempty"`
	LastPrefetchTime               *int64 `json:"lastPrefetchTime,omitempty"`
	InitialDownloadTime            *int64 `json:"initialDownloadTime,omitempty"`
}

type LogModuleInfo struct {
	Downloader string             `json:"downloader"`
	Modules    []*LogDownloadInfo `json:"modules"`
}

func (this *LogModuleInfo) GetModule(name string) *LogDownloadInfo {
	for idx := range this.Modules {
		if this.Modules[idx].Name == name {
			return this.Modules[idx]
		}
	}
	return nil
}

type LogSimInfo struct {
	Index             int         `json:"index"` //getActiveSubscriptionInfoForSimSlotIndex
	State             string      `json:"state"`
	Carrier           string      `json:"carrier"`
	SimCarrierName    interface{} `json:"sim_carrier_name"`
	SimDisplayName    interface{} `json:"sim_display_name"`
	CarrierCountryIso string      `json:"carrier_country_iso"`
	PhoneType         string      `json:"phone_type"`
	NetworkType       string      `json:"network_type"`
	CountryIso        string      `json:"country_iso"`
	Operator          interface{} `json:"operator"`
	SimOperatorName   interface{} `json:"sim_operator_name"`
	PhoneNumber       interface{} `json:"phone_number"`
	SerialNumber      interface{} `json:"serial_number"`
	SubscriberId      interface{} `json:"subscriber_id"`
	DeviceLocale      string      `json:"device_locale"`
	NetworkInfo       interface{} `json:"network_info"` //getActiveNetworkInfo(); networkInfo0.getExtraInfo
}

type PkgInstallStatus struct {
	InstallationStatus string `json:"installation_status"`
	PackageName        string `json:"package_name"`
	Version            int    `json:"version"`
}

type InsPkgStatus struct {
	InstagramAndroidInstallationStatus bool `json:"instagram_android_installation_status"`
	InstagramLiteInstallationStatus    bool `json:"instagram_lite_installation_status"`
}

type LogDeviceInfo struct {
	Pk                               string              `json:"pk"`
	Carrier                          string              `json:"carrier"`                                        //telephonyManager0.getNetworkOperatorName()
	CarrierCountryIso                string              `json:"carrier_country_iso"`                            //telephonyManager0.getNetworkCountryIso()
	NetworkType                      string              `json:"network_type"`                                   //telephonyManager0.getNetworkType() 1.GPRS EDGE UMTS CDMA EVDO_0 EVDO_A 1xRTT HSDPA HSUPA HSPA IDEN EVDO_B LTE EHRPD 15.HSPAP UNKNOWN
	PhoneType                        string              `json:"phone_type"`                                     //telephonyManager0.getPhoneType NONE GSM CDMA UNKNOWN SIP
	SimCountryIso                    string              `json:"sim_country_iso"`                                //telephonyManager0.getSimCountryIso()
	SimOperator                      string              `json:"sim_operator,omitempty"`                         //telephonyManager0.getSimState() == 5 getSimOperatorName()
	Locale                           string              `json:"locale"`                                         //Locale
	DeviceType                       string              `json:"device_type"`                                    //Build.MODEL
	Brand                            string              `json:"brand"`                                          //Build.BRAND
	Manufacturer                     string              `json:"manufacturer"`                                   //Build.MANUFACTURER
	OsType                           string              `json:"os_type"`                                        //"Android"
	OsVer                            string              `json:"os_ver"`                                         //Build.VERSION.RELEASE
	CpuAbi                           string              `json:"cpu_abi"`                                        //Build.CPU_ABI
	CpuAbi2                          string              `json:"cpu_abi2"`                                       //Build.CPU_ABI2
	CpuAbiList                       string              `json:"cpu_abilist"`                                    //Build.SUPPORTED_ABIS
	CpuBoardPlatform                 string              `json:"cpu_board_platform"`                             //ro.board.platform
	CpuBuildBoard                    string              `json:"cpu_build_board"`                                //Build.BOARD
	CpuMediatekPlatform              string              `json:"cpu_mediatek_platform"`                          //ro.mediatek.platform
	CpuChipName                      string              `json:"cpu_chip_name"`                                  //ro.chipname
	TouchFeature                     bool                `json:"touch_feature"`                                  //hasSystemFeature("android.hardware.touchscreen")
	RamLowFeature                    bool                `json:"ram_low_feature"`                                //hasSystemFeature("android.hardware.ram.low")
	TouchConfig                      int                 `json:"touch_config"`                                   //context0.getResources().getConfiguration().touchscreen
	KeyboardConfig                   int                 `json:"keyboard_config"`                                //context0.getResources().getConfiguration().keyboard
	UnreliableCoreCount              int                 `json:"unreliable_core_count"`                          //Math.max(Runtime.getRuntime().availableProcessors(), 1);
	ReliableCoreCount                int                 `json:"reliable_core_count"`                            // v = 0BD.A0V("/sys/devices/system/cpu/").listFiles(new 08l()).length;
	CpuMaxFreq                       int                 `json:"cpu_max_freq"`                                   // max /sys/devices/system/cpu/cpu%d/cpufreq/cpuinfo_max_freq
	LowPowerCpuMaxFreq               int                 `json:"low_power_cpu_max_freq"`                         // min /sys/devices/system/cpu/cpu%d/cpufreq/cpuinfo_max_freq
	TotalMem                         int64               `json:"total_mem"`                                      // getMemoryInfo(activityManager$MemoryInfo0);.totalMem;
	YearClass                        int                 `json:"year_class"`                                     //1Tj.A00
	CgroupsSupported                 bool                `json:"cgroups_supported"`                              // new File("/dev/cpuctl/tasks").exists();
	FirstInstallTime                 string              `json:"first_install_time"`                             //
	LastUpgradeTime                  string              `json:"last_upgrade_time"`                              //
	InstallLocation                  string              `json:"install_location"`                               //internal_storage sdcard external_storage
	Density                          float64             `json:"density"`                                        //getResources().getDisplayMetrics().density
	ScreenWidth                      int                 `json:"screen_width"`                                   //getSystemService("window")).getDefaultDisplay().getSize(point0).x
	ScreenHeight                     int                 `json:"screen_height"`                                  //getSystemService("window")).getDefaultDisplay().getSize(point0).y
	FrontCamera                      bool                `json:"front_camera"`                                   //packageManager0.hasSystemFeature("android.hardware.camera.front")
	RearCamera                       bool                `json:"rear_camera"`                                    //packageManager0.hasSystemFeature("android.hardware.camera")
	AllowsNonMarketInstalls          string              `json:"allows_non_market_installs"`                     // Settings.Secure.getString(context0.getContentResolver(), "install_non_market_apps"));
	AndroidId                        string              `json:"android_id"`                                     //
	IsBackgroundRestricted           *bool               `json:"is_background_restricted,omitempty"`             // sdk int >= 28:  isBackgroundRestricted()
	AppStandbyBucket                 *string             `json:"app_standby_bucket,omitempty"`                   // sdk int >= 28: usageStatsManager0.getAppStandbyBucket(); 10: "STANDBY_BUCKET_ACTIVE" 20:  "STANDBY_BUCKET_WORKING_SET" 30: "STANDBY_BUCKET_FREQUENT" 40: "STANDBY_BUCKET_RARE" 45: "STANDBY_BUCKET_RESTRICTED" "STANDBY_BUCKET_UNKNOWN_" + v2
	IsBatteryOptimized               bool                `json:"is_battery_optimized"`                           //isIgnoringBatteryOptimizations(s6) ^ 1
	IsPowerSaveMode                  bool                `json:"is_power_save_mode"`                             //isPowerSaveMode()
	OpenglVersion                    int                 `json:"opengl_version"`                                 //getDeviceConfigurationInfo().reqGlEsVersion
	GooglePlayServicesInstallation   string              `json:"google_play_services_installation"`              //0 SERVICE_DISABLED 1 SERVICE_ENABLED 2 SERVICE_INVALID SERVICE_MISSING
	GooglePlayServicesVersion        int                 `json:"google_play_services_version"`                   //
	GoogleAccounts                   int                 `json:"google_accounts"`                                //AccountManager.get(context0).getAccountsByType("com.google").length
	Installer                        string              `json:"installer"`                                      //getPackageManager().getInstallerPackageName(context0.getPackageName());
	OriginalInstaller                string              `json:"original_installer"`                             //getPackageManager().getInstallerPackageName(context0.getPackageName())
	AmazonAppStoreInstallationStatus *PkgInstallStatus   `json:"amazon_app_store_installation_status,omitempty"` //  "com.amazon.venezia", "com.amazon.mShop.android"
	DeviceTotalSpace                 string              `json:"device_total_space"`                             //data getBlockCount()) * ((long)statFs0.getBlockSize
	ExternalCacheSize                string              `json:"external_cache_size"`                            // ontext1.getExternalCacheDir();
	DeviceFreeSpace                  string              `json:"device_free_space"`                              //Environment.getDataDirectory().getPath()getAvailableBlocks()) * ((long)statFs0.getBlockSize(
	CacheSize                        string              `json:"cache_size"`                                     //dir all size ontext1.getCacheDir().getCanonicalFile(
	AppDataSize                      string              `json:"app_data_size"`                                  //context1.getFilesDir().getParentFile   -  cache_size
	ExternalAppDataSize              string              `json:"external_app_data_size"`                         //getExternalFilesDir(Environment.DIRECTORY_DOWNLOADS).getParentFile()  -  external_cache_size
	SdFreeSpace                      string              `json:"sd_free_space"`                                  //getExternalStorageDirectory().getPath() statFs1.getAvailableBlocks()) * ((long)statFs1.getBlockSize())
	SdTotalSpace                     string              `json:"sd_total_space"`                                 //tatFs1.getBlockSize()) * ((long)statFs1.getBlockCount());
	GooglePlayStore                  []*PkgInstallStatus `json:"google_play_store"`                              //"com.android.vending", "com.google.market", "com.google.android.finsky"
	GsfInstallationStatus            *PkgInstallStatus   `json:"gsf_installation_status"`                        //"com.google.android.gsf"
	InstagramAppsInstallationStatus  *InsPkgStatus       `json:"instagram_apps_installation_status"`             //  PropertiesStore_0C7.A04(0xE34, false):{com.instagram.android  com.instagram.lite  }
	SimInfo                          []*LogSimInfo       `json:"sim_info"`                                       //
	ModuleInfo                       *LogModuleInfo      `json:"moduleInfo"`                                     //
	AdvertiserId                     *string             `json:"advertiser_id,omitempty"`
	AllowAdsTracking                 bool                `json:"allow_ads_tracking"`                    //
	TargetSdk                        int                 `json:"target_sdk"`                            //context0.getApplicationInfo().targetSdkVersion
	VersionCode                      int                 `json:"version_code"`                          //
	LauncherName                     string              `json:"launcher_name"`                         // Intent intent0 = new Intent("android.intent.action.MAIN"); intent0.addCategory("android.intent.category.HOME"); ResolveInfo resolveInfo0 = context0.getPackageManager().resolveActivity(intent0, 0x10000);
	PermissionGroups                 int                 `json:"permissionGroups"`                      // 0 android.permission.READ_CALENDAR 1 android.permission.CAMERA 2 android.permission.READ_CONTACTS 3 android.permission.ACCESS_FINE_LOCATION 4 android.permission.RECORD_AUDIO 5 android.permission.READ_PHONE_STATE 6 android.permission.BODY_SENSORS 7 android.permission.READ_SMS 8 android.permission.READ_EXTERNAL_STORAGE
	AccessibilityEnabled             bool                `json:"accessibility_enabled"`                 //isEnabled()
	TouchExplorationEnabled          bool                `json:"touch_exploration_enabled"`             //isTouchExplorationEnabled()
	AndroidStrongboxAvailable        *bool               `json:"android_strongbox_available,omitempty"` //sdk>=28 : hasSystemFeature("android.hardware.strongbox_keystore"
	AndroidTeeAvailable              *bool               `json:"android_tee_available,omitempty"`       //sdk>=31  hasSystemFeature("android.hardware.hardware_keystore")
}
