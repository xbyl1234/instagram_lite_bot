package instagram

import (
	"CentralizedControl/common/utils"
	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/ast"
)

func SerializeForm(userSet map[string]string, seq []string, temp map[string]string) string {
	data := ""
	for _, key := range seq {
		data += key + "="
		if value, ok := userSet[key]; ok {
			data += utils.Escape(value, utils.EscapeEncodeQueryComponent)
		} else {
			data += utils.Escape(temp[key], utils.EscapeEncodeQueryComponent)
		}
		data += "&"
	}
	return data[:len(data)-1]
}

type Body interface {
	Serialize() string
}

type BaseBody struct {
	req *ApiRequest
}

type SignedBody struct {
	BaseBody
	json *ast.Node
}

func (this *SignedBody) GetJson() *ast.Node {
	return this.json
}

func (this *SignedBody) Serialize() string {
	json, _ := this.json.MarshalJSON()
	return "signed_body=SIGNATURE." + utils.Escape(string(json), utils.EscapeEncodeQueryComponent)
}

type FormBody struct {
	BaseBody
	userSetForm map[string]string
}

func (this *FormBody) SetForm(k, v string) {
	this.userSetForm[k] = v
}

func (this *FormBody) Serialize() string {
	return SerializeForm(this.userSetForm, this.req.apiBodyInfo.FormSeq, this.req.apiBodyInfo.FormTemplate)
}

type JsonBody struct {
	BaseBody
	json *ast.Node
}

func (this *JsonBody) Serialize() string {
	json, _ := this.json.MarshalJSON()
	return utils.Escape(string(json), utils.EscapeEncodeQueryComponent)
}

func CreateSignedBody(req *ApiRequest) *SignedBody {
	root, err := sonic.GetFromString(req.apiBodyInfo.JsonTemplate)
	if err != nil {
		panic(err)
	}
	return &SignedBody{
		BaseBody: BaseBody{
			req: req,
		},
		json: &root,
	}
}

func CreateFormBody(req *ApiRequest) *FormBody {
	return &FormBody{
		BaseBody: BaseBody{
			req: req,
		},
		userSetForm: map[string]string{},
	}
}

func CreateJsonBody(req *ApiRequest) *JsonBody {
	root, err := sonic.GetFromString(req.apiBodyInfo.JsonTemplate)
	if err != nil {
		panic(err)
	}
	return &JsonBody{
		BaseBody: BaseBody{
			req: req,
		},
		json: &root,
	}
}
