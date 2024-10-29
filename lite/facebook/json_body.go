package facebook

import (
	"CentralizedControl/common/utils"
	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/ast"
)

type JsonBody struct {
	BaseBody
	*ast.Node
}

func (this *JsonBody) Serialize() string {
	json, _ := this.MarshalJSON()
	return utils.Escape(string(json), utils.EscapeEncodeQueryComponent)
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
		Node: &root,
	}
}
