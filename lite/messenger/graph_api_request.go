package messenger

import (
	"CentralizedControl/common/log"
	"CentralizedControl/common/utils"
	"bytes"
	http "github.com/bogdanfinn/fhttp"
	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/ast"

	"io"
)

type GraphApiRequest struct {
	msg           *Messenger
	userSetHeader http.Header
	bodyKv        map[string]string
	jsonVariables *JsonVariables
	apiInfo       *ApiInfo
}

func newEmptyGraphApi() *GraphApiRequest {
	req := &GraphApiRequest{
		userSetHeader: http.Header{},
		bodyKv:        map[string]string{},
	}
	return req
}

func newJsonGraphApi(msg *Messenger, api string) *GraphApiRequest {
	req := newEmptyGraphApi()
	req.msg = msg
	req.apiInfo = apiInfo[api]
	if req.apiInfo == nil {
		panic(api)
	}
	if req.apiInfo.IsBlockSetting {
		req.jsonVariables = CreateBlockSettingVariables(msg, req.apiInfo.VariablesTemplate)
	} else {
		req.jsonVariables = CreateJsonVariables(msg, req.apiInfo.VariablesTemplate)
	}
	return req
}

func (this *GraphApiRequest) GetVariables() *JsonVariables {
	return this.jsonVariables
}

func (this *GraphApiRequest) Set(k string, v string) {
	this.bodyKv[k] = v
}

func (this *GraphApiRequest) SetHeader(k string, v string) {
	this.userSetHeader.Set(k, v)
}

func (this *GraphApiRequest) autoSetParams() {
	for _, key := range this.apiInfo.ParamsSeq {
		if _, ok := this.bodyKv[key]; ok {
			continue
		}
		tempValue, hasTemp := this.apiInfo.ParamsTemplate[key]
		autoFunc := autoSetParamsFun[key]
		if autoFunc == nil && !hasTemp {
			log.Error("not find auto params func: %s", key)
			continue
		}
		if hasTemp {
			this.Set(key, tempValue)
		}
		if autoFunc != nil {
			this.Set(key, autoFunc(this))
		}
	}
}

func (this *GraphApiRequest) autoSetHeader() {
	for _, key := range this.apiInfo.HeadSeq {
		if _, ok := this.userSetHeader[key]; ok {
			continue
		}
		tempValue, hasTemp := this.apiInfo.HeaderTemplate[key]
		autoFunc := autoSetHeaderFun[key]
		if autoFunc == nil && !hasTemp {
			log.Error("not find auto header func: %s", key)
			continue
		}
		if hasTemp {
			this.SetHeader(key, tempValue)
		}
		if autoFunc != nil {
			this.SetHeader(key, autoFunc(this))
		}
	}
}

func (this *GraphApiRequest) Serialize() string {
	if this.jsonVariables != nil {
		this.Set("variables", this.jsonVariables.Serialize())
	}
	this.autoSetParams()
	data := ""
	for _, key := range this.apiInfo.ParamsSeq {
		data += key + "="
		data += utils.Escape(this.bodyKv[key], utils.EscapeEncodeQueryComponent) + "&"
	}
	return data[:len(data)-1]
}

type Response struct {
	Err  error
	Json *ast.Node
	Row  []byte
}

func (this *GraphApiRequest) Send() (*Response, error) {
	this.autoSetHeader()
	body := this.Serialize()
	if DebugPkg {
		log.Debug("%s", body)
	}
	response, err := this.sendRequest([]byte(body))
	if err != nil {
		return nil, err
	}
	result := &Response{}
	result.Row = response
	if this.apiInfo.IsJsonResponse {
		jsonResp, err := sonic.Get(response)
		if err != nil {
			return nil, err
		}
		result.Json = &jsonResp
	}
	return result, nil
}

func (this *GraphApiRequest) sendRequest(body []byte) ([]byte, error) {
	req, err := http.NewRequest(this.apiInfo.Method, this.apiInfo.Url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header = this.userSetHeader
	req.Header[http.HeaderOrderKey] = this.apiInfo.HeadSeq

	log.Info("send request: %s, friend name: %s",
		this.apiInfo.Url,
		req.Header.Get("x-fb-friendly-name"))

	resp, err := this.msg.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	readBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return readBytes, nil
}
