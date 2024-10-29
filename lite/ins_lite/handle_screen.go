package ins_lite

import (
	"CentralizedControl/common/log"
	"CentralizedControl/ins_lite/net"
	"CentralizedControl/ins_lite/proto/io"
	"CentralizedControl/ins_lite/proto/msg/recver"
	"CentralizedControl/ins_lite/proto/types"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
)

const (
	ScreenMsgCodeSendAction2  = 2
	ScreenMsgCodeRedirect     = 15
	ScreenMsgCodeSendAction40 = 40
	ScreenMsgCodeSendAction54 = 54
	ScreenMsgCodeSendAction76 = 76
	ScreenMsgCodeSpeakText    = 269

	//screen cmd code to sub cmd
	ScreenMsgCodeCheckCallingOrSelfPermission = 116
)

func (this *InsLiteClient) registerScreenHandle() {
	this.screenEvent.RegisterEvent(ScreenMsgCodeSendAction2, this.handleScreenMsgCodeSendAction)
	this.screenEvent.RegisterEvent(ScreenMsgCodeSendAction40, this.handleScreenMsgCodeSendAction)
	this.screenEvent.RegisterEvent(ScreenMsgCodeSendAction54, this.handleScreenMsgCodeSendAction)
	this.screenEvent.RegisterEvent(ScreenMsgCodeSendAction76, this.handleScreenMsgCodeSendAction)
	this.screenEvent.RegisterEvent(ScreenMsgCodeRedirect, this.handleScreenMsgCodeRedirect)
	this.screenEvent.RegisterEvent(ScreenMsgCodeSpeakText, this.handleSpeakText)
	this.screenEvent.RegisterEvent(ScreenMsgCodeCheckCallingOrSelfPermission, this.handleCheckCallingOrSelfPermission)
}

func (this *InsLiteClient) DefaultScreenDealFunc(event *ScreenEventData) error {
	log.Warn("undisposed screen cmd code: %d", event.CmdCode)
	return nil
}

func (this *InsLiteClient) OnScreenRecvFinish(screen *recver.ScreenReceived, allSubScreen *recver.SubScreenArray) error {
	if screen.GetScreenName() == ScreenNameIgCarbonRegistration {
		//phoneInput := screen.GetScreenByWindowId(WindowIdPhoneInput)
		//if phoneInput == nil {
		//	panic("not find phoneInput")
		//}
		//err := this.callScreenEvent(screen, uint16(phoneInput.ToSubScreen3().SubScreen3Body.OnFocusedScreenCmdCode1.Value))
		//if err != nil {
		//	return err
		//}
		//err = this.callScreenEvent(screen, uint16(phoneInput.ToSubScreen3().SubScreen3Body.OnFocusedScreenCmdCode2.Value))
		//if err != nil {
		//	return err
		//}

		this.sendNetworkInfo(screen.GetScreenId(), this.getScreenIdByName(this.curScreenName), screen.DecodeBodyDataSize)
	}
	return nil
}

func (this *InsLiteClient) recvScreenDiff(msgEvent *net.MsgWrap) error {
	body := msgEvent.Body.Interface().(*recver.ScreenDiff)
	screen := this.getScreenById(body.ScreenId)
	if screen == nil {
		panic(fmt.Sprintf("not find screen %d", body.ScreenId))
	}
	screen.DecodeBody.ReadChange(msgEvent.Reader, 0, nil)
	all := screen.GetAllSubScreen()
	this.updateSubmitInfo(screen, all)
	//for test
	marshal, _ := json.Marshal(screen)
	log.Info("update screen: %s", marshal)
	this.targetScreenUpdateFinish(screen.GetScreenName())
	return nil
}

func (this *InsLiteClient) recvScreenReceived(msgEvent *net.MsgWrap) error {
	screenHeader := msgEvent.Body.Interface().(*recver.ScreenReceivedHeader)
	var screen *recver.ScreenReceived
	log.Info("screen %d recv part %d", screenHeader.ScreenId, screenHeader.PartNumber)
	if screenHeader.PartNumber.Value == 0 {
		screen = &recver.ScreenReceived{}
		screen.Header = screenHeader
		screen.DecodeBodyDataSize = len(msgEvent.Reader.PeekRemain())
		types.ReadMsg(msgEvent.Reader, &screen.DecodeBody)
		this.addScreen(screen)
	} else {
		screen = this.getScreenById(screenHeader.ScreenId)
		screen.ReadNextPart(msgEvent.Reader)
	}
	if !screenHeader.IsFinish() {
		log.Warn("screen %d not recv finish", screenHeader.ScreenId)
		return nil
	}
	marshal, _ := json.Marshal(screen)
	log.Info("ScreenBody: %s", string(marshal))
	if screen.GetScreenName() == ScreenNameIgError {
		panic("recv igliteerror screen")
	}
	log.Info("recv screen id: %d, name: %s", screen.GetScreenId(), screen.GetScreenName())
	all := screen.GetAllSubScreen()
	this.updateSubmitInfo(screen, all)
	//执行窗口初始命令
	if screen.GetRunScreenCode() != 0 {
		err := this.callScreenEvent(screen, uint16(screen.GetRunScreenCode()), "Tag_RunScreenCode")
		if err != nil {
			panic(err)
		}
	}

	//展示窗口消息
	if screen.GetScreenId() != LoadingScreenId {
		log.Info("send browser action")
		this.sendBrowserAction(screen, &screen.Header.NavigationData,
			true, "", 0, "")
	}

	for idx := 0; idx < all.Count(); idx++ {
		item := all.Get(idx)
		code := item.GetDisplayActionScreenCmdCode()
		if code != 0 {
			err := this.callScreenEvent(screen, uint16(code), "Tag_RunDisplayAction")
			if err != nil {
				panic(err)
			}
		}
	}

	err := this.OnScreenRecvFinish(screen, all)
	if err != nil {
		return err
	}

	this.setCurScreen(screen)
	this.targetScreenRecvFinish(screen.GetScreenName())
	return nil
}

func (this *InsLiteClient) callScreenEvent(screen *recver.ScreenReceived, cmdIdx uint16, source string) error {
	if cmdIdx == 0 {
		return nil
	}
	code, data := screen.GetScreenCmdByIdx(int32(cmdIdx))
	log.Info("run screen event: %s, screen id: %d, name: %s, cmd idx: %d, cmd: %d, %s",
		source, screen.GetScreenId(), screen.GetScreenName(), cmdIdx, code, hex.EncodeToString(data))
	reader := io.CreateReaderBuffer(data)
	e := &ScreenEventData{
		Screen:  screen,
		CmdIdx:  cmdIdx,
		CmdCode: uint64(code),
		CmdData: data,
		Reader:  reader,
		Source:  source,
	}
	return this.screenEvent.callEvent(e)
}

func (this *InsLiteClient) handleScreenMsgCodeSendAction(event *ScreenEventData) error {
	screenCmd := &recver.ScreenCmd{}
	types.ReadMsg(event.Reader, screenCmd)
	var isCode83 bool
	if (!event.Screen.DecodeBody.ScreenDecodeBody.BitByteFlags.GetFlags(17) ||
		!event.Screen.DecodeBody.ScreenDecodeBody.BitByteFlags.GetFlags(18)) &&
		(screenCmd.Flags0&4 == 0) {
		isCode83 = false
	} else {
		isCode83 = true
	}
	var instanceKey int32
	if !isCode83 || screenCmd.Unknow1.Size != 0 || screenCmd.Flags0&32 != 0 {
		instanceKey = this.getTrackingInstanceKey()
	} else {
		instanceKey = -1
	}
	likeZeroOrScreenCmdResult := int32(screenCmd.Flags0 & 8)

	var likeActionResourceId int32 = 0
	if event.Source != "" {
		if strings.Index(event.Source, "Btn_") == 0 {
			windowName := strings.ReplaceAll(event.Source, "Btn_", "")
			subScreen := event.Screen.GetScreenByWindowId(windowName)
			if subScreen == nil {
				panic("window name should not null: " + windowName)
			}
			for subScreen != nil {
				if subScreen.Type == 2 && subScreen.GetBaseScreen().GetIsLikeResIdChildFlag() && subScreen.GetBaseScreen().LikeActionResourceId.Value == 30001 {
					likeActionResourceId = 30001
					break
				}
				subScreen = subScreen.Parent
			}
		} else {
			subScreens := event.Screen.GetAllSubScreen()
			for idx := 0; idx < subScreens.Count(); idx++ {
				subScreen := subScreens.Get(idx)
				if subScreen.Type == 2 && subScreen.GetBaseScreen().GetIsLikeResIdChildFlag() && subScreen.GetBaseScreen().LikeActionResourceId.Value == 30001 {
					likeActionResourceId = 30001
					break
				}
				if subScreen.Type == 2 && subScreen.GetBaseScreen().GetIsLikeResIdChildFlag() {
					idx = 0
					subScreens = subScreen.GetSubScreenCpy()
				}
			}
		}
		log.Info("get %s likeActionResourceId:  %v", event.Source, likeActionResourceId)
	} else {
		log.Error("source is null")
	}

	act := &SubmitParams{
		event:                     event,
		cmdData:                   screenCmd,
		isCode83:                  isCode83,
		screenCmdHasExternData:    false, //有问题
		likeNotIsCode83:           !isCode83,
		likeZeroOrScreenCmdResult: likeZeroOrScreenCmdResult,
		instanceKeyParam:          instanceKey,
		likeActionResourceId:      likeActionResourceId, // 有问题
		submitData:                this.getSubmitInfo(event.Screen.GetScreenName()),
		unusedData:                nil,
	}
	this.submitWindowData(act)
	return nil
}

func (this *InsLiteClient) handleScreenMsgCodeRedirect(event *ScreenEventData) error {
	for len(event.Reader.PeekRemain()) >= 2 {
		idx := event.Reader.ReadUShort()
		err := this.callScreenEvent(event.Screen, idx, event.Source)
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *InsLiteClient) handleSpeakText(event *ScreenEventData) error {
	screenCmd := &recver.ScreenSpeak{}
	types.ReadMsg(event.Reader, screenCmd)
	var err error
	if this.Cookies.Packages.HasTTS() {
		err = this.callScreenEvent(event.Screen, uint16(screenCmd.ScreenCmdIdxOnStart), "Tag_SpeakText")
		if err != nil {
			return err
		}
		err = this.callScreenEvent(event.Screen, uint16(screenCmd.ScreenCmdIdxOnDone), "Tag_SpeakText")
		if err != nil {
			return err
		}
		err = this.callScreenEvent(event.Screen, uint16(screenCmd.ScreenCmdIdxOnStop), "Tag_SpeakText")
		if err != nil {
			return err
		}
	} else {
		err = this.callScreenEvent(event.Screen, uint16(screenCmd.ScreenCmdIdxOnError), "Tag_SpeakText")
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *InsLiteClient) handleCheckCallingOrSelfPermission(event *ScreenEventData) error {
	screenCmd := &recver.CheckCallingOrSelfPermission{}
	types.ReadMsg(event.Reader, screenCmd)
	if this.checkHasPerm(screenCmd.PermType) {
		return this.callScreenEvent(event.Screen, uint16(screenCmd.SubCmdCodeIdx), "Tag_CheckCallingOrSelfPermission")
	}
	if len(event.Reader.PeekRemain()) >= 2 {
		subCmdCodeIdx := event.Reader.ReadUShort()
		return this.callScreenEvent(event.Screen, subCmdCodeIdx, "Tag_CheckCallingOrSelfPermission")
	}
	this.sendPermResult(screenCmd.PkgIdx, true, false, this.shouldNotShowPermissionDialog(screenCmd.PermType))
	return nil
}
