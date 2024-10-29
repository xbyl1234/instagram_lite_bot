package ins_lite

import (
	"CentralizedControl/common/log"
	"CentralizedControl/ins_lite/proto/msg/recver"
	"CentralizedControl/ins_lite/proto/types"
	"CentralizedControl/ins_lite/tools"
	"encoding/json"
)

func (this *InsLiteClient) TestAddScreen(data string) *recver.ScreenReceived {
	parse := tools.ParseLiteRecvStr(data)
	screenHeader := parse.Body.(*recver.ScreenReceivedHeader)
	var screen *recver.ScreenReceived
	log.Info("screen %d recv part %d", screenHeader.ScreenId, screenHeader.PartNumber)
	if screenHeader.PartNumber.Value == 0 {
		screen = &recver.ScreenReceived{}
		screen.Header = screenHeader
		types.ReadMsg(parse.Reader, &screen.DecodeBody)
		this.addScreen(screen)
	} else {
		screen = this.getScreenById(screenHeader.ScreenId)
		screen.ReadNextPart(parse.Reader)
	}
	if !screenHeader.IsFinish() {
		log.Warn("screen %d not recv finish", screenHeader.ScreenId)
		return nil
	}
	return screen
}

func (this *InsLiteClient) TestUpdateScreen(data string) *recver.ScreenReceived {
	p := tools.ParseLiteRecvStr(data)
	body := p.Body.(*recver.ScreenDiff)
	s := this.getScreenById(body.ScreenId)
	if s == nil {
		panic("")
	}
	s.DecodeBody.ReadChange(p.Reader, 0, nil)
	log.Info("TestUpdateScreen remain %d", len(p.Reader.PeekRemain()))
	return s
}

func (this *InsLiteClient) TestScreen2Json() string {
	marshal, err := json.Marshal(this.Screen)
	if err != nil {
		panic(err)
	}
	return string(marshal)
}
