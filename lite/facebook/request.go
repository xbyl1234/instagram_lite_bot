package facebook

import (
	"CentralizedControl/common"
	"CentralizedControl/common/log"
	"CentralizedControl/common/utils"
	"bytes"
	"fmt"
	http "github.com/bogdanfinn/fhttp"
	"github.com/bytedance/sonic/ast"
	"github.com/klauspost/compress/gzip"
	"github.com/klauspost/compress/zstd"
	"io"
	"strings"
)

type TempData struct {
	RandomBytes int
	Path        string
	Url         string
	LoggingId   string
}

type ApiRequest struct {
	fb            *Facebook
	tempData      TempData
	bodyTempKey   string
	userSetHeader http.Header
	userSetQuery  map[string]string
	body          Body

	apiInfo     *ApiInfo
	apiBodyInfo *ApiBodyInfo
}

func newEmptyApiRequest() *ApiRequest {
	req := &ApiRequest{
		userSetHeader: http.Header{},
		tempData: TempData{
			RandomBytes: utils.GenNumber(1000000, 1500000),
			LoggingId:   strings.ToLower(utils.GenUUID()),
		},
	}
	return req
}

func newApiRequest(msg *Facebook, api string, bodyTempKey string) *ApiRequest {
	req := newEmptyApiRequest()
	req.fb = msg
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

	switch req.apiBodyInfo.requestType {
	case ApiRequestTypeFormBody:
		req.body = CreateFormBody(req)
	case ApiRequestTypeJsonBody:
		req.body = CreateJsonBody(req)
	case ApiRequestTypeGraphqlBody:
		req.body = CreateGraphqlBody(req)
	case ApiRequestTypeNoBody:
		req.body = nil
	}
	return req
}

func (this *ApiRequest) GetFormBody() *FormBody {
	if this.apiBodyInfo.requestType == ApiRequestTypeFormBody {
		return this.body.(*FormBody)
	}
	panic("error req type!")
}

func (this *ApiRequest) GetJsonBody() *JsonBody {
	if this.apiBodyInfo.requestType == ApiRequestTypeJsonBody {
		return this.body.(*JsonBody)
	}
	panic("error req type!")
}

func (this *ApiRequest) GetGraphqlBody() *GraphqlBody {
	if this.apiBodyInfo.requestType == ApiRequestTypeGraphqlBody {
		return this.body.(*GraphqlBody)
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
		query := "?" + SerializeForm(this.userSetQuery, this.apiBodyInfo.QuerySeq,
			this.apiBodyInfo.Query, utils.EscapeEncodeNone)
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

func (this *ApiRequest) OnCheckResponse(response *Response) error {
	//if this.apiInfo.IsJsonResponse {
	//	jsonResp, err := sonic.Get(respData)
	//	if err != nil {
	//		return result, err
	//	}
	//	result.Json = &jsonResp
	//	this.OnCheckResponse()
	//}
	if response.resp.StatusCode >= 400 {
		return common.NerError(response.resp.Status)
	}
	return nil
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

	err = this.OnCheckResponse(result)
	if err != nil {
		return result, err
	}
	this.fb.onResponse(result)
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

	resp, err := this.fb.httpClient.Do(req)
	if err != nil {
		return nil, nil, err
	}

	defer resp.Body.Close()
	encoding := resp.Header.Get("Content-Encoding")
	var respBody io.Reader
	switch encoding {
	case "gzip":
		respBody, err = gzip.NewReader(resp.Body)
		break
	case "zstd":
		respBody, err = zstd.NewReader(resp.Body)
		break
	case "deflate":
	case "":
		respBody = resp.Body
		break
	}
	if err != nil {
		return nil, nil, err
	}

	readBytes, err := io.ReadAll(respBody)
	if err != nil {
		return nil, nil, err
	}
	return resp, readBytes, nil
}
