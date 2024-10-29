package facebook

import (
	"CentralizedControl/common/log"
	"CentralizedControl/common/utils"
	"fmt"
)

type FormBody struct {
	BaseBody
	userSetForm map[string]string
}

func (this *FormBody) SetForm(k, v string) {
	this.userSetForm[k] = v
}

func (this *FormBody) Serialize() string {
	return SerializeForm(this.userSetForm, this.req.apiBodyInfo.FormSeq,
		this.req.apiBodyInfo.FormTemplate, utils.EscapeEncodeQueryComponent)
}

func (this *FormBody) AutoSetForm() {
	for _, key := range this.req.apiBodyInfo.FormSeq {
		if _, ok := this.userSetForm[key]; ok {
			continue
		}
		tempValue, hasTemp := this.req.apiBodyInfo.FormTemplate[key]
		autoFunc := autoSetJsonFun[key]
		if autoFunc == nil && !hasTemp {
			log.Error("not find auto params func: %s", key)
			continue
		}
		if hasTemp {
			this.SetForm(key, tempValue)
		}
		if autoFunc != nil {
			this.SetForm(key, fmt.Sprintf("%v", autoFunc(this.req)))
		}
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
