package common

import "CentralizedControl/ins_lite/proto/types"

type NoData struct {
}

type UnDealData struct {
}

type TextMsg struct {
	Cmd types.VarUInt32String
	Msg types.VarUInt32String
}
