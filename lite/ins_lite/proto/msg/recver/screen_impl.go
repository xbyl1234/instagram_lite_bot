package recver

import "CentralizedControl/ins_lite/proto/io"

type ScreenImpl interface {
	ReadChange(from io.BufferReader, changeType int, parentScreen *SubScreen)
	UpdateScreen(from io.BufferReader)
	InitImpl()
	GetSubScreen() *SubScreenArray
	GetDisplayActionScreenCmdCode() int32
}
