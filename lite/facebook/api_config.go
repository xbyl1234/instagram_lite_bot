package facebook

import (
	"CentralizedControl/common/utils"
	"github.com/bytedance/sonic/ast"
	"strings"
)

type ApiRequestType int

const (
	ApiRequestTypeJsonBody    = 0
	ApiRequestTypeFormBody    = 1
	ApiRequestTypeNoBody      = 2
	ApiRequestTypeGraphqlBody = 3
)

type HostType int

const (
	ApiTypeBGraphFacebook = 0
	ApiTypeGraphFacebook  = 1
)

type NoType struct {
}

var PassHeader = map[string]*NoType{
	"Content-Length": nil,
}

type AutoSettingHeaderFun func(req *ApiRequest) string

type AutoSettingJsonFun func(req *ApiRequest) interface{}

var autoSetHeaderFun = map[string]AutoSettingHeaderFun{
	//"X-Fb-Request-Analytics-Tags": func(req *ApiRequest) string {},
	//"X-Fb-Background-State": func(req *ApiRequest) string {},
	//"X-Graphql-Request-Purpose": func(req *ApiRequest) string {},
	//"X-Fb-Device-Group": func(req *ApiRequest) string {},
	//"Priority": func(req *ApiRequest) string {},
	"X-Fb-Ta-Logging-Ids": func(req *ApiRequest) string {
		return "graphql:" + req.tempData.LoggingId
	},
	"X-Fb-Sim-Hni": func(req *ApiRequest) string {
		return req.fb.ck.Mcc + req.fb.ck.Mnc
	},
	"X-Fb-Net-Hni": func(req *ApiRequest) string {
		return req.fb.ck.Mcc + req.fb.ck.Mnc
	},
	"X-Fb-Connection-Type": func(req *ApiRequest) string {
		return "WIFI"
	},
	"User-Agent": func(req *ApiRequest) string {
		return req.fb.ck.UserAgent
	},
	"Authorization": func(req *ApiRequest) string {
		return req.fb.ck.Authorization
	},
	"X-Tigon-Is-Retry": func(req *ApiRequest) string {
		return "False"
	},
}

var autoSetJsonFun = map[string]AutoSettingJsonFun{
	"family_device_id": func(req *ApiRequest) interface{} {
		return req.fb.ck.FamilyDeviceId
	},
	"device_id": func(req *ApiRequest) interface{} {
		return req.fb.ck.DeviceId
	},
	"waterfall_id": func(req *ApiRequest) interface{} {
		return req.fb.tmpCk.WaterfallId
	},
	"headers_flow_id": func(req *ApiRequest) interface{} {
		return ""
	},
	"machine_id": func(req *ApiRequest) interface{} {
		return req.fb.ck.MachineId
	},
	"pixel_ratio": func(req *ApiRequest) interface{} {
		return req.fb.ck.Density
	},
	"scale": func(req *ApiRequest) interface{} {
		return req.fb.ck.Density
	},
	"appid": func(req *ApiRequest) interface{} {
		return "350685531728"
	},
	"app_scoped_id": func(req *ApiRequest) interface{} {
		return req.fb.ck.DeviceId
	},
	"locale": func(req *ApiRequest) interface{} {
		return req.fb.ck.Locale
	},
	"client_trace_id": func(req *ApiRequest) interface{} {
		return req.tempData.LoggingId
	},
	"access_token": func(req *ApiRequest) interface{} {
		return req.fb.ck.AccessToken
	},
	"client_country_code": func(req *ApiRequest) interface{} {
		return strings.ToUpper(req.fb.ck.Country)
	},
}

func AutoSetJsonBody(req *ApiRequest, commonJson *ast.Node) {
	for k, v := range autoSetJsonFun {
		if commonJson.Get(k).Exists() {
			commonJson.SetAny(k, v(req))
		}
	}
}

func autoSetJsonKey(req *ApiRequest, commonJson *ast.Node, key string) {
	if commonJson.Get(key).Exists() {
		commonJson.SetAny(key, autoSetJsonFun[key](req))
	}
}

type ApiBodyInfo struct {
	BodyType     string `json:"body_type"`
	JsonTemplate string `json:"json_template"`
	//HadNavChain    bool   `json:"had_nav_chain"`
	HeadSeq        []string          `json:"head_seq"`
	HeaderTemplate map[string]string `json:"header_template"`
	Query          map[string]string `json:"query"`
	QuerySeq       []string          `json:"query_seq"`
	FormSeq        []string          `json:"form_seq"`
	FormTemplate   map[string]string `json:"form_template"`

	lowCaseHeadSeq    []string
	variablesTemplate string
	requestType       ApiRequestType
}

type ApiInfo struct {
	Host           string                  `json:"host"`
	Path           string                  `json:"path"`
	Method         string                  `json:"method"`
	IsJsonResponse bool                    `json:"is_json_response"`
	Body           map[string]*ApiBodyInfo `json:"body"`
	HostType       HostType
}

var apiInfo = map[string]*ApiInfo{}

func InitApiConfig(path string) {
	if path == "" {
		path = "./fb_api_config.json"
	}
	err := utils.LoadJsonFile(path, &apiInfo)
	if err != nil {
		panic(err)
	}

	for _, apiValue := range apiInfo {
		if apiValue.Host == "b-graph.facebook.com" {
			apiValue.HostType = ApiTypeBGraphFacebook
		} else if apiValue.Host == "graph.facebook.com" {
			apiValue.HostType = ApiTypeGraphFacebook
		}

		for _, item := range apiValue.Body {
			item.lowCaseHeadSeq = make([]string, len(item.HeadSeq))
			for idx := range item.HeadSeq {
				item.lowCaseHeadSeq[idx] = strings.ToLower(item.HeadSeq[idx])
			}

			switch item.BodyType {
			case "form_body":
				item.requestType = ApiRequestTypeFormBody
			case "json_body":
				item.requestType = ApiRequestTypeJsonBody
			case "graphql_body":
				item.requestType = ApiRequestTypeGraphqlBody
			case "":
				item.requestType = ApiRequestTypeNoBody
			}
		}
	}
}

func (this *Facebook) onResponse(resp *Response) {
	for k, v := range resp.resp.Header {
		if strings.Index(k, "Ig-Set") == -1 {
			continue
		}
		if len(v) == 0 || len(v[0]) == 0 {
			continue
		}
		//switch k {
		//case "X-Ig-Set-Www-Claim":
		//	if v[0] != this.ck.Claim {
		//		DbCtrl.UpdateCookiesValue(this, "claim", v[0])
		//	}
		//	this.ck.Claim = v[0]
		//	break
		//case "Ig-Set-Ig-U-Ig-Direct-Region-Hint":
		//	if v[0] != this.ck.DirectRegionHint {
		//		DbCtrl.UpdateCookiesValue(this, "direct_region_hint", v[0])
		//	}
		//	this.ck.DirectRegionHint = v[0]
		//	break
		//case "Ig-Set-Ig-U-Shbid":
		//	if v[0] != this.ck.IgUShbid {
		//		DbCtrl.UpdateCookiesValue(this, "ig_u_shbid", v[0])
		//	}
		//	this.ck.IgUShbid = v[0]
		//	break
		//case "Ig-Set-Ig-U-Shbts":
		//	if v[0] != this.ck.IgUShbts {
		//		DbCtrl.UpdateCookiesValue(this, "ig_u_shbts", v[0])
		//	}
		//	this.ck.IgUShbts = v[0]
		//	break
		//case "Ig-Set-Ig-U-Rur":
		//	if v[0] != this.ck.Rur {
		//		DbCtrl.UpdateCookiesValue(this, "rur", v[0])
		//	}
		//	this.ck.Rur = v[0]
		//	break
		//}
	}
}
