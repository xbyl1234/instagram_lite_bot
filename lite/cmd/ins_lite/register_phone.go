package main

import (
	"CentralizedControl/common/android"
	"CentralizedControl/common/log"
	"CentralizedControl/common/phone"
	"CentralizedControl/common/proxys"
	"CentralizedControl/common/utils"
	"CentralizedControl/ins_lite"
	"CentralizedControl/ins_lite/proto"
	"CentralizedControl/ins_lite/proto/msg/sender"
	"fmt"
	"time"
)

func registerPhone(phoneNumber *phone.Number, proxy proxys.Proxy, country string) (regResult RegisterResult) {
	defer func() {
		r := recover()
		if r != nil {
			regResult.Error = fmt.Sprintf("recover error: %v", r)
			regResult.IsSuccess = false
		}
	}()
	device := android.GetAndroidDevice("", time.Now().UnixNano())
	//device := android.GetAndroidDevice("")
	device.InitDevice(country,
		android.DeviceConfigGenPerm([]string{android.GetAccounts, android.ReadContacts}),
		android.DeviceConfigHasGms(true, true, false, true),
		android.DeviceConfigGenPhoneBook(),
		android.DeviceConfigGenNetwork(true))
	regResult.DeviceName = device.DeviceName
	ck := ins_lite.CreateNewCookies(device)
	if ck == nil {
		panic("ck is null")
	}
	client := ins_lite.CreateNewInsLiteClient(ck, proxy)
	if client == nil {
		panic("ins lite client is null")
	}
	event := client.GetWaitEvent(proto.MsgCodeAboutRecvQueueAction135)
	client.SendInitAppMsg()
	event.Wait()
	ActivityResumed := &proto.Message[sender.ActivityResumed]{}
	ActivityResumed.Body.GenDeviceTimeId = client.Cookies.GenDeviceTimeId
	client.SendMsg(ActivityResumed)
	client.ReporterNetworkTypeChange()
	client.SendLoggedUserIdChange()

	client.MustWaitScreenRecvFinish(ins_lite.ScreenNameInstagramLoginEntrypoint)
	log.Info("ScreenNameInstagramLoginEntrypoint recv finish!")
	client.SleepForRecvNewView()

	client.ClickButton(ins_lite.WindowIdLandingSignupButton)
	client.SendConnBandwidthQuality()
	client.SleepForClickBtn()

	client.MustWaitScreenRecvFinish(ins_lite.ScreenNameIgCarbonRegistration)
	log.Info("ScreenNameIgCarbonRegistration recv finish!")
	client.SleepForRecvNewView()

	updateEvent := client.GetScreenUpdateFinishEvent(ins_lite.ScreenNameIgCarbonRegistration)
	emailInput := client.GetCurrentScreen().GetScreenByWindowId(ins_lite.WindowIdEmailInput)
	if emailInput != nil {
		client.ClickButton(ins_lite.WindowIdPhoneTab)
		client.SleepForClickBtn()
		updateEvent.MustWait30Second()
	}
	//phone input lost focus
	//email input get focus
	regResult.PhoneNumber = phoneNumber.Number
	regResult.AreaCode = phoneNumber.AreaCode
	client.PutSubmitString("", ins_lite.WindowIdPhoneInput, phoneNumber.Number)
	client.SleepForInputText()

	client.ClickButton(ins_lite.WindowIdPhoneInput)
	client.SleepForClickBtn()

	client.MustWaitScreenRecvFinish(ins_lite.ConfirmPhoneNumberEntrypoint)
	log.Info("ConfirmPhoneNumberEntrypoint recv finish!")
	regResult.HadSendCode = true
	client.SleepForRecvNewView()

	code, err := phoneNumber.SyncGetCode()
	if err != nil {
		panic(fmt.Sprintf("get code error: %v", err))
	}

	regResult.HadGetCode = true
	client.PutSubmitString("", ins_lite.WindowIdConfirmContactPointInput, code)
	client.SleepForInputText()

	client.ClickButton(ins_lite.WindowIdNextButton)
	client.SleepForClickBtn()

	client.MustWaitScreenRecvFinish(ins_lite.ScreenNameRegistrationNameAndPasswordEntrypoint)
	log.Info("ScreenNameRegistrationNameAndPasswordEntrypoint recv finish!")
	client.SleepForRecvNewView()

	username := android.ChoiceUsername(device.RandTool)
	passwd := utils.GenString(utils.CharSet_All, 12)
	regResult.Username = username
	regResult.Passwd = passwd
	client.PutSubmitString("", ins_lite.WindowIdFullNameInput, username)
	client.SleepForClickBtn()

	client.PutSubmitString("", ins_lite.WindowIdPasswordInput, passwd)
	client.SleepForClickBtn()

	client.PutSubmitString("", ins_lite.WindowIdSavePasswd, "y")
	client.SleepForClickBtn()

	client.ClickButton(ins_lite.WindowIdContinueWithSyncButton)
	client.SleepForClickBtn()

	client.MustWaitScreenRecvFinish(ins_lite.ScreenNameRegistrationBirthdayEntrypoint)
	log.Info("ScreenNameRegistrationBirthdayEntrypoint recv finish!")
	client.SleepForRecvNewView()

	year := utils.GenNumber(1997, 2001)
	month := utils.GenNumber(1, 12)
	day := utils.GenNumber(1, 27)
	yearBtnId := fmt.Sprintf(ins_lite.WindowIdYear, year)
	monthBtnId := fmt.Sprintf(ins_lite.WindowIdMonth, month)
	dayBtnId := fmt.Sprintf(ins_lite.WindowIdDay, day)

	updateEvent = client.GetScreenUpdateFinishEvent(ins_lite.ScreenNameRegistrationBirthdayEntrypoint)
	client.ClickButton(yearBtnId)
	client.SleepForClickBtn()
	updateEvent.MustWait30Second()

	updateEvent = client.GetScreenUpdateFinishEvent(ins_lite.ScreenNameRegistrationBirthdayEntrypoint)
	client.ClickButton(monthBtnId)
	client.SleepForClickBtn()
	updateEvent.MustWait30Second()

	updateEvent = client.GetScreenUpdateFinishEvent(ins_lite.ScreenNameRegistrationBirthdayEntrypoint)
	client.ClickButton(dayBtnId)
	client.SleepForClickBtn()
	updateEvent.MustWait30Second()

	client.ClickButton(ins_lite.WindowIdDataPickerNextButton)

	client.MustWaitScreenRecvFinish(ins_lite.ScreenNameRegistrationWelcomeToInstagramEntrypoint)
	log.Info("ScreenNameRegistrationWelcomeToInstagramEntrypoint recv finish!")
	client.SleepForRecvNewView()

	client.ClickButton(ins_lite.WindowIdNextButton)
	client.SleepForClickBtn()

	recvSuccess, recvScreen := client.WaitMultiScreenRecvFinish(ins_lite.ScreenNameConnectToFacebookEntrypoint,
		ins_lite.ScreenNameAddProfilePhotoEntrypoint, ins_lite.ScreenNameBloksShell, ins_lite.ScreenNameCarbonCheckpointHandler)
	if recvSuccess && (recvScreen == ins_lite.ScreenNameConnectToFacebookEntrypoint || recvScreen == ins_lite.ScreenNameAddProfilePhotoEntrypoint) {
		log.Info("reg finish, email: %s, username: %s, passwd: %s", phoneNumber.Number, username, passwd)
		regResult.IsSuccess = true
	} else {
		if recvScreen == "" {
			recvScreen = "timeout"
		}
		log.Error("reg failed ban: %s, email: %s, username: %s, passwd: %s", recvScreen, phoneNumber.Number, username, passwd)
		regResult.IsSuccess = false
		regResult.Error = recvScreen
	}
	return regResult
}
