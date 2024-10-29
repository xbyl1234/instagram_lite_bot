package android

const (
	ReadContacts                = "android.permission.READ_CONTACTS"
	WriteContacts               = "android.permission.WRITE_CONTACTS"
	GetAccounts                 = "android.permission.GET_ACCOUNTS"
	ReadCallLog                 = "android.permission.READ_CALL_LOG"
	ReadPhoneState              = "android.permission.READ_PHONE_STATE"
	CallPhone                   = "android.permission.CALL_PHONE"
	ReadCalendar                = "android.permission.READ_CALENDAR"
	WriteCalendar               = "android.permission.WRITE_CALENDAR"
	Camera                      = "android.permission.CAMERA"
	AccessCoarseLocation        = "android.permission.ACCESS_COARSE_LOCATION"
	AccessFineLocation          = "android.permission.ACCESS_FINE_LOCATION"
	WriteExternalStorage        = "android.permission.WRITE_EXTERNAL_STORAGE"
	RecordAudio                 = "android.permission.RECORD_AUDIO"
	ReadSms                     = "android.permission.READ_SMS"
	ReadExternalStorage         = "android.permission.READ_EXTERNAL_STORAGE" //sdk < 33
	ReadMediaImages             = "android.permission.READ_MEDIA_IMAGES"     //sdk > 33
	ReadMediaVideo              = "android.permission.READ_MEDIA_VIDEO"
	PostNotifications           = "android.permission.POST_NOTIFICATIONS"
	AnswerPhoneCalls            = "android.permission.ANSWER_PHONE_CALLS"
	ReadMediaVisualUserSelected = "android.permission.READ_MEDIA_VISUAL_USER_SELECTED"
	BodySensors                 = "android.permission.BODY_SENSORS"
)

var AllPermNames = []string{
	ReadContacts,
	WriteContacts,
	GetAccounts,
	ReadCallLog,
	ReadPhoneState,
	CallPhone,
	ReadCalendar,
	WriteCalendar,
	Camera,
	AccessCoarseLocation,
	AccessFineLocation,
	WriteExternalStorage,
	RecordAudio,
	ReadSms,
	ReadExternalStorage,
	ReadMediaImages,
	ReadMediaVideo,
	PostNotifications,
	AnswerPhoneCalls,
	ReadMediaVisualUserSelected,
	BodySensors,
}

const (
	PermStatusUnReq   = -1
	PermStatusDenied  = 0
	PermStatusGranted = 1
)

type PermissionInfo struct {
	CurStatic    int
	ShouldStatic int
}

type Permission struct {
	Perm map[string]*PermissionInfo
}

func (this *Permission) IsAllow(name string) bool {
	return this.Perm[name].CurStatic == PermStatusGranted
}

func (this *Permission) Require(name string) int {
	this.Perm[name].CurStatic = this.Perm[name].ShouldStatic
	return this.Perm[name].CurStatic
}

func CreateAndroidPermission(allowPerms []string) *Permission {
	p := &Permission{
		Perm: map[string]*PermissionInfo{},
	}
	for _, item := range AllPermNames {
		p.Perm[item] = &PermissionInfo{
			CurStatic:    PermStatusUnReq,
			ShouldStatic: PermStatusDenied,
		}
	}
	for _, item := range allowPerms {
		p.Perm[item].ShouldStatic = PermStatusGranted
	}
	return p
}
