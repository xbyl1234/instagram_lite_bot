package messenger

import (
	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/ast"
)

type blockSettingParams struct {
	Params            string `json:"params"`
	BloksVersioningId string `json:"bloks_versioning_id"`
	AppId             string `json:"app_id"`
}

type blockSettingContext struct {
	StylesId         string `json:"styles_id"`
	UsingWhiteNavbar bool   `json:"using_white_navbar"`
	PixelRatio       int    `json:"pixel_ratio"`
	IsPushOn         bool   `json:"is_push_on"`
	BloksVersion     string `json:"bloks_version"`
}

type BlockSettingVariables struct {
	Params    blockSettingParams  `json:"params"`
	Scale     string              `json:"scale"`
	NtContext blockSettingContext `json:"nt_context"`
}

type JsonVariables struct {
	msg         *Messenger
	Root        *ast.Node
	BlockParams *ast.Node
}

func CreateJsonVariables(msg *Messenger, temp string) *JsonVariables {
	root, err := sonic.GetFromString(temp)
	if err != nil {
		panic(err)
	}
	return &JsonVariables{
		msg:  msg,
		Root: &root,
	}
}

func JsonIfSetUInt64(node *ast.Node, key string, value string) {
	if node.Get(key) != nil {
		node.Set(key, ast.NewNumber(value))
	}
}

func JsonIfSetString(node *ast.Node, key string, value string) {
	if node.Get(key) != nil {
		node.Set(key, ast.NewString(value))
	}
}

func CreateBlockSettingVariables(msg *Messenger, temp string) *JsonVariables {
	root := CreateJsonVariables(msg, temp)
	subParams, _ := root.Root.Get("params").Get("params").String()
	subParams = subParams[len("{params:") : len(subParams)-len(",}")]
	blockParams, err := sonic.GetFromString(subParams)
	if err != nil {
		panic(err)
	}
	root.BlockParams = &blockParams

	clientInputParams := root.GetBlockClientInputParams()
	if clientInputParams != nil {
		JsonIfSetString(clientInputParams, "device_id", root.msg.ck.DeviceId)
		JsonIfSetString(clientInputParams, "machine_id", root.msg.ck.MachineId)
		JsonIfSetString(clientInputParams, "family_device_id", root.msg.ck.FamilyDeviceId)
	}
	serverParams := root.GetBlockServerParams()
	if serverParams != nil {
		JsonIfSetUInt64(serverParams, "account_type", root.msg.ck.AccountType)
		JsonIfSetUInt64(serverParams, "account_id", root.msg.ck.AccountId)
	}
	return root
}

func (this *JsonVariables) GetRoot() *ast.Node {
	return this.Root
}

func (this *JsonVariables) GetBlock() *ast.Node {
	return this.BlockParams
}

func (this *JsonVariables) GetBlockClientInputParams() *ast.Node {
	return this.BlockParams.Get("client_input_params")
}

func (this *JsonVariables) GetBlockServerParams() *ast.Node {
	return this.BlockParams.Get("server_params")
}

func (this *JsonVariables) Serialize() string {
	if this.BlockParams != nil {
		data, err := this.BlockParams.MarshalJSON()
		if err != nil {
			panic(err)
		}
		this.Root.Get("params").Set("params", ast.NewString("{params:"+string(data)+",}"))
	}
	data, err := this.Root.MarshalJSON()
	if err != nil {
		panic(err)
	}
	return string(data)
}
