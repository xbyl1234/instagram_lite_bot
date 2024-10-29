package messenger

import (
	"CentralizedControl/common/utils"
)

type AutoSettingFun func(req *GraphApiRequest) string

var autoSetHeaderFun = map[string]AutoSettingFun{
	//"content-type": func(req *GraphApiRequest) string {
	//	return "application/x-www-form-urlencoded"
	//},
	//"accept-encoding": func(req *GraphApiRequest) string {
	//	return "gzip, deflate"
	//},
	//"user-agent": func(req *GraphApiRequest) string {
	//	return req.msg.ck.UserAgent
	//},
	//"authorization": func(req *GraphApiRequest) string {
	//	return "OAuth " + req.msg.ck.Authorization
	//},
	//"x-fb-ta-logging-ids": func(req *GraphApiRequest) string {
	//	return "graphql:" + utils.GenUUID()
	//},
	//"x-fb-net-hni": func(req *GraphApiRequest) string {
	//	return req.msg.ck.Mcc + req.msg.ck.Mnc
	//},
	//"x-fb-sim-hni": func(req *GraphApiRequest) string {
	//	return req.msg.ck.Mcc + req.msg.ck.Mnc
	//},
}

var autoSetParamsFun = map[string]AutoSettingFun{
	"method":            func(req *GraphApiRequest) string { return "post" },
	"pretty":            func(req *GraphApiRequest) string { return "false" },
	"format":            func(req *GraphApiRequest) string { return "json" },
	"server_timestamps": func(req *GraphApiRequest) string { return "true" },
	"client_trace_id":   func(req *GraphApiRequest) string { return utils.GenUUID() },
	//"locale":            func(req *GraphApiRequest) string { return req.msg.ck.Locale },
	// X.5HI.A03->createClientDocIdForQueryNameHash(0x84b923b7) libcoldstart.so offset: 0xb204d8
	// 固定界面参数
	//"client_doc_id":            func(req *GraphApiRequest) string { return req.apiInfo.ParamsTemplate["client_doc_id"] },
	//"fb_api_req_friendly_name": func(req *GraphApiRequest) string { return req.apiInfo.ParamsTemplate["fb_api_req_friendly_name"] },
	//"fb_api_caller_class":      func(req *GraphApiRequest) string { return req.apiInfo.ParamsTemplate["fb_api_caller_class"] },
	//"fb_api_analytics_tags":    func(req *GraphApiRequest) string { return req.apiInfo.ParamsTemplate["fb_api_analytics_tags"] },
	//"purpose":                  func(req *GraphApiRequest) string { return req.apiInfo.ParamsTemplate["purpose"] },
}

type ApiInfo struct {
	Url               string            `json:"url"`
	Method            string            `json:"method"`
	HeadSeq           []string          `json:"head_seq"`
	ParamsSeq         []string          `json:"params_seq"`
	ParamsTemplate    map[string]string `json:"params_template"`
	HeaderTemplate    map[string]string `json:"header_template"`
	VariablesTemplate string            `json:"variables_template"`

	IsJsonResponse bool `json:"is_json_response"`
	IsBlockSetting bool `json:"is_block_setting"`
}

var apiInfo = map[string]*ApiInfo{}

func InitApiConfig(path string) {
	if path == "" {
		path = "./msg_api_config.json"
	}
	err := utils.LoadJsonFile(path, &apiInfo)
	if err != nil {
		panic(err)
	}
}
