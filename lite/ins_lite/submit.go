package ins_lite

import (
	"CentralizedControl/ins_lite/proto"
	"CentralizedControl/ins_lite/proto/msg/recver"
	"CentralizedControl/ins_lite/proto/msg/sender"
	"CentralizedControl/ins_lite/proto/types"
	"time"
)

type SubmitParams struct {
	event                     *ScreenEventData
	cmdData                   *recver.ScreenCmd
	isCode83                  bool
	screenCmdHasExternData    bool
	likeNotIsCode83           bool
	likeZeroOrScreenCmdResult int32
	instanceKeyParam          int32
	likeActionResourceId      int32
	submitData                *ScreenSubmitData
	unusedData                []sender.SendActionUnUsedData
}

func (this *InsLiteClient) submitWindowData(action *SubmitParams) {
	var dealActionData []sender.SendSubmitData
	submit := action.submitData.Get()
	dealActionData = make([]sender.SendSubmitData, len(submit))
	for idx, item := range submit {
		if item.Data != "" && item.NeedEncrypt {
			dealActionData[idx].Data = this.MakePassword(item.Data)
		} else {
			dealActionData[idx].Data = item.Data
		}
		dealActionData[idx].Type = submit[idx].Type
	}

	var instanceKey int32
	var markerId int32
	if action.likeNotIsCode83 {
		markerId = 0x25C0004
		if action.isCode83 {
			instanceKey = action.instanceKeyParam
		} else {
			instanceKey = 0
		}
	} else {
		markerId = 0x2400001
		instanceKey = action.instanceKeyParam
	}
	trackingState := sender.CreateTrackingState(true, markerId, true, instanceKey, instanceKey != -1, false)

	var flag uint32 = 0
	if action.screenCmdHasExternData {
		flag = 2
	}
	if action.unusedData != nil || len(action.unusedData) > 0 {
		flag = flag | 1 | 32
	}
	if action.likeZeroOrScreenCmdResult != 0 {
		flag |= 8
	}
	if !trackingState.IsDisableTracking() {
		flag |= 64
	}

	msgBody := &sender.ActionMsg{
		FromScreenId:              action.event.Screen.Header.ScreenId,
		ToScreenId:                action.cmdData.ToScreenId,
		ResourceId:                action.cmdData.ResourceId,
		LikeActionResourceId:      int16(action.likeActionResourceId),
		RespMsgData:               *types.CreateListValue[sender.SendSubmitData, int16](dealActionData),
		Flag1:                     *types.CreateVarUInt32(flag),
		UnusedData:                *types.CreateListValue[sender.SendActionUnUsedData, int16](action.unusedData),
		LikeZeroOrScreenCmdResult: action.likeZeroOrScreenCmdResult,
		InstanceKey:               action.instanceKeyParam,
		Const:                     -1,
		Time:                      time.Now().UnixMilli(),
		TrackingState:             *trackingState,
	}

	if action.isCode83 {
		msg := &proto.Message[sender.PassivityActionMsg]{}
		msg.Body.ActionMsg = *msgBody
		this.SendMsg(msg)
	} else {
		msg := &proto.Message[sender.InitiativeActionMsg]{}
		msg.Body.ActionMsg = *msgBody
		this.SendMsg(msg)
	}
}
