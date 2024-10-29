package ins_lite

import (
	"CentralizedControl/common/android"
)

var PermType2Name = map[int]string{
	0:  android.ReadContacts,
	1:  android.WriteContacts,
	2:  android.GetAccounts,
	3:  android.ReadCallLog,
	4:  android.ReadPhoneState,
	5:  android.CallPhone,
	6:  android.ReadCalendar,
	7:  android.WriteCalendar,
	8:  android.Camera,
	9:  android.AccessCoarseLocation,
	10: android.AccessFineLocation,
	11: android.WriteExternalStorage,
	12: android.RecordAudio,
	13: android.ReadSms,
	15: android.ReadExternalStorage,
	16: android.ReadMediaVideo,
	17: android.PostNotifications,
	18: android.AnswerPhoneCalls,
	19: android.ReadMediaVisualUserSelected,
}

func (this *InsLiteClient) shouldNotShowPermissionDialog(permType byte) bool {
	return false
}

func (this *InsLiteClient) checkHasPerm(permType byte) bool {
	return this.Cookies.Permission.IsAllow(PermType2Name[int(permType)])
}

func (this *InsLiteClient) requestPermissions(permTypes []byte) []int {
	result := make([]int, len(permTypes))
	for idx := range permTypes {
		result[idx] = this.Cookies.Permission.Require(PermType2Name[int(permTypes[idx])])
	}
	return result
}
