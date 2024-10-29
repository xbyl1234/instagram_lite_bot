package ins_lite

import (
	"CentralizedControl/common/log"
	"CentralizedControl/ins_lite/proto/msg/recver"
)

func (this *InsLiteClient) ClickButton(windowsId string) error {
	recvScreen := this.GetCurrentScreen()
	log.Info("cur screen: %s, click button: %s", recvScreen.GetScreenName(), windowsId)
	btnScreen := recvScreen.GetScreenByWindowId(windowsId)
	if btnScreen == nil {
		panic("not find window")
	}
	idx := btnScreen.GetScreenClickCmdCodeIdx()
	if idx == 0 {
		panic("screen code idx is 0")
	}
	return this.callScreenEvent(recvScreen, uint16(idx), "Btn_"+windowsId)
}

func (this *InsLiteClient) CallScreenEvent(screen *recver.ScreenReceived, code uint16, eventName string) error {
	return this.callScreenEvent(screen, code, eventName)
}
