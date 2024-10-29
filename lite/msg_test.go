package main

import (
	"CentralizedControl/common/proxys"
	"testing"
)

func TestIp(t *testing.T) {
	ip, err := proxys.GetIpInfoByIp("50.114.59.9", nil)
	println(ip)
	println(err)
}

//func TestMsg(t *testing.T) {
//	messenger.InitApiConfig("")
//	ck, err := os.ReadFile("./test_ck.json")
//	if err != nil {
//		panic(err)
//	}
//	msg := messenger.CreateMessenger(messenger.ConvDeviceFile2Cookies(string(ck)))
//	msg.SetProxy(proxys.DebugHttpProxy)
//	msg.Set2FA()
//	_ = msg
//}

//func TestEncodeMsg(t *testing.T) {
//	//data := "client_doc_id=10537346114073748047743187847&method=post&locale=zh_CN&pretty=false&format=json&variables={\"params\":{\"params\":\"{params:{\\\"client_input_params\\\":{\\\"machine_id\\\":\\\"NPK8ZAUINXPyoTBGL55DJmj4\\\"},\\\"server_params\\\":{\\\"should_show_done_button\\\":0,\\\"account_type\\\":0,\\\"account_id\\\":100094758223469,\\\"INTERNAL_INFRA_screen_id\\\":\\\"0\\\"}},}\",\"bloks_versioning_id\":\"e4ac5d6944a7e6c15b7353672e454568871b99f82a8810e33351bd7aa0bd97ea\",\"app_id\":\"com.bloks.www.fx.settings.security.two_factor.select_method\"},\"scale\":\"3\",\"nt_context\":{\"styles_id\":\"48aa9bafd01f10f7bc81ee10dd0b81aa\",\"using_white_navbar\":true,\"pixel_ratio\":3,\"is_push_on\":true,\"bloks_version\":\"e4ac5d6944a7e6c15b7353672e454568871b99f82a8810e33351bd7aa0bd97ea\"}}&fb_api_req_friendly_name=FbBloksAppRootQuery-com.bloks.www.fx.settings.security.two_factor.select_method&fb_api_caller_class=graphservice&fb_api_analytics_tags=[\"GraphServices\"]&client_trace_id=767209f7-a58a-4023-8d9b-e9bcf0c56dd1&server_timestamps=true&purpose=fetch"
//	data := ",:"
//	//params := map[string]string{
//	//	"client_doc_id":"10537346114073748047743187847",
//	//	"method":"post",
//	//	"locale":"zh_CN",
//	//	"pretty":"false",
//	//	"format":"json",
//	//	"variables":
//	//	"fb_api_req_friendly_name":"FbBloksAppRootQuery-com.bloks.www.fx.settings.security.two_factor.select_method
//	//	"fb_api_caller_class":"graphservice"
//	//	"fb_api_analytics_tags":
//	//	"client_trace_id":"767209f7-a58a-4023-8d9b-e9bcf0c56dd1"
//	//	"server_timestamps":"true"
//	//	"purpose":"fetch"
//	//}
//	//
//
//	println(common.Escape(data, common.EscapeEncodePath))
//	println(common.Escape(data, common.EscapeEncodePathSegment))
//	println(common.Escape(data, common.EscapeEncodeHost))
//	println(common.Escape(data, common.EscapeEncodeZone))
//	println(common.Escape(data, common.EscapeEncodeUserPassword))
//	println(common.Escape(data, common.EscapeEncodeQueryComponent))
//	println(common.Escape(data, common.EscapeEncodeFragment))
//}

//func TestMyHttp(t *testing.T) {
//	client := http2_helper.CreateHttp2Client(http2_helper.DisableHttpSslPinng())
//	http2_helper.HttpSetProxy(proxys.DebugS5Proxy)(client)
//
//	//proxys.DebugS5Proxy.GetProxy()(client)
//
//	req, err := http.NewRequest(http.MethodGet, "https://tls.peet.ws/api/all", nil)
//	if err != nil {
//		log.Println(err)
//		return
//	}
//
//	req.Header = http.Header{
//		"accept":          {"*/*"},
//		"accept-language": {"de-DE,de;q=0.9,en-US;q=0.8,en;q=0.7"},
//		"user-agent":      {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.75 Safari/537.36"},
//		http.HeaderOrderKey: {
//			"accept",
//			"accept-language",
//			"user-agent",
//		},
//	}
//
//	resp, err := client.Do(req)
//	if err != nil {
//		log.Println(err)
//		return
//	}
//
//	defer resp.Body.Close()
//
//	log.Println(fmt.Sprintf("status code: %d", resp.StatusCode))
//
//	readBytes, err := io.ReadAll(resp.Body)
//	if err != nil {
//		log.Println(err)
//		return
//	}
//
//	log.Println(string(readBytes))
//}
