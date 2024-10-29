package instagram

import (
	"CentralizedControl/common"
	"CentralizedControl/common/log"
	"CentralizedControl/common/utils"
	"bytes"
	"fmt"
	http "github.com/bogdanfinn/fhttp"
	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/ast"
	"io"
)

type TempData struct {
	RandomBytes int
	Path        string
	Url         string
}

type ApiRequest struct {
	ins           *Instagram
	tempData      TempData
	bodyTempKey   string
	userSetHeader http.Header
	userSetQuery  map[string]string
	body          Body

	apiInfo     *ApiInfo
	apiBodyInfo *ApiBodyInfo
	navChain    *NavChain
}

func newEmptyApiRequest() *ApiRequest {
	req := &ApiRequest{
		userSetHeader: http.Header{},
		tempData: TempData{
			RandomBytes: utils.GenNumber(1000000, 1500000),
		},
	}
	return req
}

func newApiRequest(msg *Instagram, api string, bodyTempKey string) *ApiRequest {
	req := newEmptyApiRequest()
	req.ins = msg
	req.apiInfo = apiInfo[api]
	if req.apiInfo == nil {
		panic(api)
	}

	if bodyTempKey == "" {
		for key := range req.apiInfo.Body {
			bodyTempKey = key
			break
		}
	}
	req.bodyTempKey = bodyTempKey
	req.apiBodyInfo = req.apiInfo.Body[bodyTempKey]

	if req.apiBodyInfo.HadNavChain {
		req.navChain = req.apiBodyInfo.NavChain
	}

	switch req.apiBodyInfo.requestType {
	case ApiRequestTypeSignedBody:
		req.body = CreateSignedBody(req)
	case ApiRequestTypeFormBody:
		req.body = CreateFormBody(req)
	case ApiRequestTypeJsonBody:
		req.body = CreateJsonBody(req)
	case ApiRequestTypeNoBody:
		req.body = nil
	}
	return req
}

func (this *ApiRequest) GetNavChain() *NavChain {
	return this.navChain
}

func (this *ApiRequest) NewNavChain(fakeTime bool) *NavChain {
	this.navChain = &NavChain{
		chain:    make([]*NavChainItem, 0),
		fakeTime: fakeTime,
	}
	return this.navChain
}

func (this *ApiRequest) GetFormBody() *FormBody {
	return this.body.(*FormBody)
}

func (this *ApiRequest) GetJsonBody() *ast.Node {
	if this.apiBodyInfo.requestType == ApiRequestTypeJsonBody {
		return this.body.(*JsonBody).json
	} else if this.apiBodyInfo.requestType == ApiRequestTypeSignedBody {
		return this.body.(*SignedBody).json
	}
	panic("error req type!")
}

func (this *ApiRequest) getUrl() string {
	if this.tempData.Url != "" {
		return this.tempData.Url
	}
	this.tempData.Url = "https://" + this.apiInfo.Host
	if this.tempData.Path != "" {
		this.tempData.Url += this.tempData.Path
	} else {
		this.tempData.Url += this.apiInfo.Path
	}
	if this.apiBodyInfo.QuerySeq != nil && len(this.apiBodyInfo.QuerySeq) > 0 {
		query := "?" + SerializeForm(this.userSetQuery, this.apiBodyInfo.QuerySeq, this.apiBodyInfo.Query)
		this.tempData.Url += query
	}
	return this.tempData.Url
}

func (this *ApiRequest) SetHeader(k string, v string) {
	this.userSetHeader.Set(k, v)
}

func (this *ApiRequest) SetQuery(k string, v string) {
	this.userSetQuery[k] = v
}

func (this *ApiRequest) SetPathParams(params ...interface{}) {
	this.tempData.Path = fmt.Sprintf(this.apiInfo.Path, params...)
}

func (this *ApiRequest) AutoTempJson() {
	AutoSetJsonBody(this)
}

func (this *ApiRequest) checkPassHeader(header string) bool {
	_, ok := PassHeader[header]
	return ok
}

func (this *ApiRequest) autoSetHeader() {
	for _, key := range this.apiBodyInfo.HeadSeq {
		if _, ok := this.userSetHeader[key]; ok {
			continue
		}
		tempValue, hasTemp := this.apiBodyInfo.HeaderTemplate[key]
		autoFunc := autoSetHeaderFun[key]
		if autoFunc == nil && !hasTemp {
			log.Error("not find auto header func or temp: %s", key)
			continue
		}
		if hasTemp && !this.checkPassHeader(key) {
			this.SetHeader(key, tempValue)
		}
		if autoFunc != nil {
			this.SetHeader(key, autoFunc(this))
		}
	}
}

func (this *ApiRequest) Serialize() string {
	if this.body != nil {
		return this.body.Serialize()
	}
	return ""
}

type Response struct {
	Json *ast.Node
	Row  []byte
	resp *http.Response
}

func (this *ApiRequest) Send() (*Response, error) {
	this.autoSetHeader()
	body := this.Serialize()
	if DebugPkg {
		log.Debug("%s", body)
	}
	resp, respData, err := this.sendRequest([]byte(body))
	if err != nil {
		return nil, err
	}
	result := &Response{}
	result.Row = respData
	result.resp = resp

	if this.apiInfo.IsJsonResponse {
		jsonResp, err := sonic.Get(respData)
		if err != nil {
			return result, err
		}
		result.Json = &jsonResp
		if this.apiInfo.ApiType == ApiTypeInsV1 {
			if jsonResp.GetString("status") != "ok" {
				errstr := "api error"
				if jsonResp.Get("message").Exists() {
					errstr = jsonResp.GetString("message")
				}
				return result, common.NerError(errstr)
			}
		}
	}

	if resp.StatusCode >= 400 {
		return result, common.NerError(resp.Status)
	}

	this.ins.onResponse(result)
	return result, nil
}

func (this *ApiRequest) sendRequest(body []byte) (*http.Response, []byte, error) {
	req, err := http.NewRequest(this.apiInfo.Method, this.getUrl(), bytes.NewBuffer(body))
	if err != nil {
		return nil, nil, err
	}

	req.Header = this.userSetHeader
	req.Header[http.HeaderOrderKey] = this.apiBodyInfo.lowCaseHeadSeq

	log.Info("send request: %s", this.getUrl())

	resp, err := this.ins.httpClient.Do(req)
	if err != nil {
		return nil, nil, err
	}

	defer resp.Body.Close()

	readBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	return resp, readBytes, nil
}
