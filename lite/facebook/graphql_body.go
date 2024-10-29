package facebook

import (
	"CentralizedControl/common"
	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/ast"
	"strings"
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

type SubParamsType int

const (
	SubParamsTypeParamsNotRealJsonStr SubParamsType = 1
	SubParamsTypeParamsJsonStr        SubParamsType = 2
	SubParamsTypeParamsParamsJsonStr  SubParamsType = 3
	SubParamsTypeParamsNo1Params      SubParamsType = 4
	SubParamsTypeParamsNo2Params      SubParamsType = 5
)

type GraphqlBody struct {
	*FormBody
	req           *ApiRequest
	Root          *ast.Node
	SubParams     *ast.Node
	SubParamsType SubParamsType
}

func createGraphqlBody(req *ApiRequest) *GraphqlBody {
	root, err := sonic.GetFromString(req.apiBodyInfo.FormTemplate["variables"])
	if err != nil {
		panic(err)
	}
	return &GraphqlBody{
		req:      req,
		Root:     &root,
		FormBody: CreateFormBody(req),
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

func CreateGraphqlBody(req *ApiRequest) *GraphqlBody {
	root := createGraphqlBody(req)
	subParams := root.Root.Get("params")
	if subParams.Exists() {
		subParams2 := subParams.Get("params")
		if subParams2.Exists() {
			switch subParams2.Type() {
			case ast.V_STRING: //string
				subStr, _ := subParams2.String()
				if strings.Index(subStr, "{params:") == 0 {
					subStr = subStr[len("{params:") : len(subStr)-len(",}")]
					subJson, err := sonic.GetFromString(subStr)
					if err != nil {
						panic(err)
					}
					root.SubParams = &subJson
					root.SubParamsType = SubParamsTypeParamsNotRealJsonStr
				} else {
					subJson, err := sonic.GetFromString(subStr)
					if err != nil {
						panic(err)
					}
					toMap, err := subJson.Map()
					if err != nil {
						panic(err)
					}
					if _, ok := toMap["params"]; len(toMap) == 1 && ok {
						root.SubParams = subJson.Get("params")
						root.SubParamsType = SubParamsTypeParamsParamsJsonStr
					} else {
						root.SubParams = &subJson
						root.SubParamsType = SubParamsTypeParamsJsonStr
					}
				}
			default:
				panic(common.NerError(""))
			}
		} else {
			root.SubParams = subParams
			root.SubParamsType = SubParamsTypeParamsNo1Params
		}
	} else {
		root.SubParamsType = SubParamsTypeParamsNo2Params
	}

	return root
}

func (this *GraphqlBody) GetRoot() *ast.Node {
	return this.Root
}

func (this *GraphqlBody) GetSubParams() *ast.Node {
	if this.SubParamsType == SubParamsTypeParamsNo2Params {
		panic(common.NerError("GetSubParams SubParamsTypeParamsNo2Params"))
	}
	return this.SubParams
}

func (this *GraphqlBody) GetBlockClientInputParams() *ast.Node {
	if this.SubParams == nil {
		return nil
	}
	return this.SubParams.Get("client_input_params")
}

func (this *GraphqlBody) GetBlockServerParams() *ast.Node {
	if this.SubParams == nil {
		return nil
	}
	return this.SubParams.Get("server_params")
}

func (this *GraphqlBody) AutoSetNtContext(node *ast.Node) {
	autoSetJsonKey(this.req, node, "pixel_ratio")
}

func (this *GraphqlBody) AutoSetScale(node *ast.Node) {
	autoSetJsonKey(this.req, node, "")
}

func (this *GraphqlBody) AutoSetServerParams(node *ast.Node) {
	autoSetJsonKey(this.req, node, "device_id")
	autoSetJsonKey(this.req, node, "waterfall_id")
	autoSetJsonKey(this.req, node, "headers_flow_id")
	autoSetJsonKey(this.req, node, "family_device_id")
}

func (this *GraphqlBody) AutoCommonJson(node *ast.Node) {
	AutoSetJsonBody(this.req, node)
}

func (this *GraphqlBody) AutoSetClientInputParams(node *ast.Node) {
	autoSetJsonKey(this.req, node, "device_id")
	autoSetJsonKey(this.req, node, "machine_id")
	autoSetJsonKey(this.req, node, "family_device_id")
}

func (this *GraphqlBody) AutoSet() {
	switch this.SubParamsType {
	case SubParamsTypeParamsNotRealJsonStr:
		this.AutoSetNtContext(this.Root.Get("nt_context"))
		this.AutoSetClientInputParams(this.SubParams)
		this.AutoSetServerParams(this.SubParams)
		this.AutoSetScale(this.Root)
	case SubParamsTypeParamsJsonStr:
		this.AutoSetNtContext(this.Root.Get("nt_context"))
	case SubParamsTypeParamsParamsJsonStr:
		this.AutoSetNtContext(this.Root.Get("nt_context"))
		this.AutoSetClientInputParams(this.SubParams)
		this.AutoSetServerParams(this.SubParams)
		this.AutoSetScale(this.Root)
	case SubParamsTypeParamsNo1Params:
		this.AutoSetNtContext(this.Root.Get("nt_context"))
	case SubParamsTypeParamsNo2Params:
		this.AutoCommonJson(this.Root)
	}
}

func (this *GraphqlBody) Serialize() string {
	var subStr string
	var err error
	if this.SubParamsType == SubParamsTypeParamsNo2Params {
		subStr, err = sonic.MarshalString(this.GetRoot())
		if err != nil {
			panic(err)
		}
	}
	switch this.SubParamsType {
	case SubParamsTypeParamsNotRealJsonStr:
		this.GetRoot().Get("params").SetString("params", "{params:"+subStr+",}")
	case SubParamsTypeParamsJsonStr:
		var params ast.Node
		params.SetString("params", subStr)
		subStr2, err := sonic.MarshalString(params)
		if err != nil {
			panic(err)
		}
		this.GetRoot().Get("params").SetString("params", subStr2)
	case SubParamsTypeParamsParamsJsonStr:
		this.GetRoot().Get("params").SetString("params", subStr)
	case SubParamsTypeParamsNo1Params:
		this.GetRoot().Set("params", *this.GetSubParams())
	case SubParamsTypeParamsNo2Params:
	}
	marshalString, err := sonic.MarshalString(this.Root)
	if err != nil {
		panic(err)
	}

	this.SetForm("variables", marshalString)

	return this.FormBody.Serialize()
}
