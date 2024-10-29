package ins_lite

import (
	"CentralizedControl/ins_lite/proto"
	"CentralizedControl/ins_lite/proto/msg/sender"
)

type Model struct {
	InitAppMsg proto.MessageC[sender.AppInitMsg]
}

func (this *InsLiteClient) loadModelFromHex(path string) *Model {
	var model Model

	err := model.InitAppMsg.ReadFromHexFile(path + "/1.txt")
	if err != nil {
		panic(err)
	}
	model.InitAppMsg.ClientId = -1
	model.InitAppMsg.Magic = 0xcf3

	return &model
}

func (this *InsLiteClient) loadModelFromJson(path string) *Model {
	var model Model

	err := model.InitAppMsg.ReadFromJsonFile(path + "/1.json")
	if err != nil {
		panic(err)
	}
	model.InitAppMsg.ClientId = -1
	model.InitAppMsg.Magic = 0xcf3

	return &model
}
