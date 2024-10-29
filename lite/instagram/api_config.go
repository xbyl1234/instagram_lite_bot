package instagram

import (
	"CentralizedControl/common/utils"
	"fmt"
	"strings"
	"time"
)

type ApiRequestType int

const (
	ApiRequestTypeSignedBody = 0
	ApiRequestTypeJsonBody   = 1
	ApiRequestTypeFormBody   = 2
	ApiRequestTypeNoBody     = 3
)

type ApiType int

const (
	ApiTypeInsV1    = 0
	ApiTypeInsGraph = 1
	ApiTypeInsWebFb = 2
)

type NoType struct {
}

var PassHeader = map[string]*NoType{
	"Content-Length": nil,
}

type AutoSettingHeaderFun func(req *ApiRequest) string
type AutoSettingJsonFun func(req *ApiRequest) interface{}

var autoSetHeaderFun = map[string]AutoSettingHeaderFun{
	//"X-Bloks-Version-Id": "8ca96ca267e30c02cf90888d91eeff09627f0e3fd2bd9df472278c9a6c022cbb",
	//"X-Bloks-Is-Layout-Rtl": "false",
	//"X-Ig-Capabilities": "3brTv10=",
	//"X-Ig-App-Id": "567067343352427",
	//"Priority": "u=3",
	//"Content-Type": "application/x-www-form-urlencoded; charset=UTF-8",
	//"Accept-Encoding": "gzip, deflate",
	//"X-Fb-Http-Engine": "Liger",
	//"X-Fb-Client-Ip": "True",
	//"X-Fb-Server-Cluster": "True",
	//"X-Ig-Transfer-Encoding": "chunked",
	//"X-Ig-Prefetch-Request": "foreground",
	//"X-Fb-Friendly-Name": "CrossPostingContentCompatibilityConfig",
	//"X-Root-Field-Name": "fxcal_accounts",
	//"X-Graphql-Client-Library": "graphservice",
	//"X-Tigon-Is-Retry": "False",
	//"X-Ads-Opt-Out": "0",
	//"X-Fb": "1",
	//"X-Messenger": "1",
	//X-Tigon-Is-Retry:  是否重试
	"X-Ig-Nav-Chain": func(req *ApiRequest) string {
		return req.GetNavChain().Serialize()
	},
	"X-Cm-Latency": func(req *ApiRequest) string {
		return fmt.Sprintf("%.03f", float64(utils.GenNumber(1000, 3000))/1000.0)
	},
	"X-Cm-Bandwidth-Kbps": func(req *ApiRequest) string {
		return fmt.Sprintf("%.03f", float64(req.tempData.RandomBytes)/float64(req.ins.tmpCk.BandwidthTotaltimeMs))
	},
	"X-Ig-Bandwidth-Speed-Kbps": func(req *ApiRequest) string {
		return fmt.Sprintf("%d.000", req.tempData.RandomBytes/req.ins.tmpCk.BandwidthTotaltimeMs)
	},
	"X-Ig-Bandwidth-Totalbytes-B": func(req *ApiRequest) string { return fmt.Sprintf("%d", req.tempData.RandomBytes) },
	"X-Ig-Bandwidth-Totaltime-Ms": func(req *ApiRequest) string { return fmt.Sprintf("%d", req.ins.tmpCk.BandwidthTotaltimeMs) },
	"X-Ig-Timezone-Offset":        func(req *ApiRequest) string { return req.ins.tmpCk.TimezoneOffset },
	"X-Ig-App-Locale": func(req *ApiRequest) string {
		return strings.ToLower(req.ins.ck.Language) + "_" + strings.ToUpper(req.ins.ck.Country)
	},
	"X-Ig-Device-Locale": func(req *ApiRequest) string {
		return strings.ToLower(req.ins.ck.Language) + "_" + strings.ToUpper(req.ins.ck.Country)
	},
	"X-Ig-Mapped-Locale": func(req *ApiRequest) string {
		return strings.ToLower(req.ins.ck.Language) + "_" + strings.ToUpper(req.ins.ck.Country)
	},
	"Accept-Language": func(req *ApiRequest) string {
		return strings.ToLower(req.ins.ck.Language) + "-" + strings.ToUpper(req.ins.ck.Country)
	},
	"X-Pigeon-Session-Id": func(req *ApiRequest) string { return req.ins.tmpCk.PigeonSessionId },
	"X-Pigeon-Rawclienttime": func(req *ApiRequest) string {
		return fmt.Sprintf("%.03f", float64(time.Now().UnixMilli())/1000.0)
	},
	"X-Fb-Connection-Type": func(req *ApiRequest) string { return "WIFI" },
	"X-Ig-Connection-Type": func(req *ApiRequest) string { return "WIFI" },
	"User-Agent":           func(req *ApiRequest) string { return req.ins.ck.UserAgent },

	"X-Ig-Device-Id":        func(req *ApiRequest) string { return req.ins.ck.DeviceId },
	"X-Device-Id":           func(req *ApiRequest) string { return req.ins.ck.DeviceId },
	"X-Ig-Family-Device-Id": func(req *ApiRequest) string { return req.ins.ck.FamilyId },
	"X-Ig-Android-Id":       func(req *ApiRequest) string { return "android-" + req.ins.ck.AndroidId },

	"Ig-U-Ig-Direct-Region-Hint": func(req *ApiRequest) string { return req.ins.ck.DirectRegionHint },
	"X-Ig-Www-Claim":             func(req *ApiRequest) string { return req.ins.ck.Claim },
	"X-Mid":                      func(req *ApiRequest) string { return req.ins.ck.Mid }, //DEVICE_HEADER_ID
	"Ig-U-Rur":                   func(req *ApiRequest) string { return req.ins.ck.Rur },
	"Ig-U-Ds-User-Id":            func(req *ApiRequest) string { return req.ins.ck.AccountId },
	"Ig-Intended-User-Id":        func(req *ApiRequest) string { return req.ins.ck.AccountId },
	"Authorization":              func(req *ApiRequest) string { return req.ins.ck.Authorization },
}

var autoSetJsonFun = map[string]AutoSettingJsonFun{
	"radio_type": func(req *ApiRequest) interface{} { return "wifi-none" },
	"_uid":       func(req *ApiRequest) interface{} { return req.ins.ck.AccountId },
	"device_id": func(req *ApiRequest) interface{} {
		if strings.Contains(req.GetJsonBody().GetString("device_id"), "android-") {
			return "android-" + req.ins.ck.AndroidId
		} else {
			return req.ins.ck.DeviceId
		}
	},
	"_uuid": func(req *ApiRequest) interface{} { return req.ins.ck.DeviceId },
}

func AutoSetJsonBody(req *ApiRequest) {
	var json = req.GetJsonBody()
	for k, v := range autoSetJsonFun {
		if json.Get(k).Exists() {
			json.SetAny(k, v(req))
		}
	}
}

type ApiBodyInfo struct {
	BodyType       string `json:"body_type"`
	JsonTemplate   string `json:"json_template"`
	HadNavChain    bool   `json:"had_nav_chain"`
	lowCaseHeadSeq []string
	HeadSeq        []string          `json:"head_seq"`
	HeaderTemplate map[string]string `json:"header_template"`
	Query          map[string]string `json:"query"`
	QuerySeq       []string          `json:"query_seq"`
	FormSeq        []string          `json:"form_seq"`
	FormTemplate   map[string]string `json:"form_template"`

	NavChain    *NavChain
	requestType ApiRequestType
}

type ApiInfo struct {
	Host           string                  `json:"host"`
	Path           string                  `json:"path"`
	Method         string                  `json:"method"`
	IsJsonResponse bool                    `json:"is_json_response"`
	Body           map[string]*ApiBodyInfo `json:"body"`
	ApiType        ApiType
}

var apiInfo = map[string]*ApiInfo{}

func InitApiConfig(path string) {
	if path == "" {
		path = "./ins_api_config.json"
	}
	err := utils.LoadJsonFile(path, &apiInfo)
	if err != nil {
		panic(err)
	}

	for _, apiValue := range apiInfo {
		if apiValue.Host == "i.instagram.com" {
			apiValue.ApiType = ApiTypeInsV1
		} else if apiValue.Host == "graph.instagram.com" {
			apiValue.ApiType = ApiTypeInsGraph
		} else if apiValue.Host == "web.facebook.com" {
			apiValue.ApiType = ApiTypeInsWebFb
		}

		for _, item := range apiValue.Body {
			item.lowCaseHeadSeq = make([]string, len(item.HeadSeq))
			for idx := range item.HeadSeq {
				item.lowCaseHeadSeq[idx] = strings.ToLower(item.HeadSeq[idx])
			}

			if item.HadNavChain {
				item.NavChain = CreateNavChain(item.HeaderTemplate["X-Ig-Nav-Chain"])
			}
			switch item.BodyType {
			case "signed_body":
				item.requestType = ApiRequestTypeSignedBody
			case "form_body":
				item.requestType = ApiRequestTypeFormBody
			case "json_body":
				item.requestType = ApiRequestTypeJsonBody
			case "":
				item.requestType = ApiRequestTypeNoBody
			}
		}
	}
}

func (this *Instagram) onResponse(resp *Response) {
	for k, v := range resp.resp.Header {
		if strings.Index(k, "Ig-Set") == -1 {
			continue
		}
		if len(v) == 0 || len(v[0]) == 0 {
			continue
		}
		switch k {
		case "X-Ig-Set-Www-Claim":
			if v[0] != this.ck.Claim {
				DbCtrl.UpdateCookiesValue(this, "claim", v[0])
			}
			this.ck.Claim = v[0]
			break
		case "Ig-Set-Ig-U-Ig-Direct-Region-Hint":
			if v[0] != this.ck.DirectRegionHint {
				DbCtrl.UpdateCookiesValue(this, "direct_region_hint", v[0])
			}
			this.ck.DirectRegionHint = v[0]
			break
		case "Ig-Set-Ig-U-Shbid":
			if v[0] != this.ck.IgUShbid {
				DbCtrl.UpdateCookiesValue(this, "ig_u_shbid", v[0])
			}
			this.ck.IgUShbid = v[0]
			break
		case "Ig-Set-Ig-U-Shbts":
			if v[0] != this.ck.IgUShbts {
				DbCtrl.UpdateCookiesValue(this, "ig_u_shbts", v[0])
			}
			this.ck.IgUShbts = v[0]
			break
		case "Ig-Set-Ig-U-Rur":
			if v[0] != this.ck.Rur {
				DbCtrl.UpdateCookiesValue(this, "rur", v[0])
			}
			this.ck.Rur = v[0]
			break
		}
	}
}
