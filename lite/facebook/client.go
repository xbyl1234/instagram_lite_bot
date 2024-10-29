package facebook

import (
	"CentralizedControl/common/proxys"
	http "github.com/bogdanfinn/fhttp"
)

type Facebook struct {
	ck         *Cookies
	tmpCk      *TempCookies
	httpClient *http.Client
	proxy      proxys.Proxy
}

//
//func (this *Facebook) GetCookies() *Cookies {
//	return this.ck
//}
//
//func (this *Facebook) GetSelfUserId() uint64 {
//	id, _ := strconv.ParseInt(this.ck.AccountId, 10, 64)
//	return uint64(id)
//}
//
//func (this *Facebook) newApiRequest(api string, bodyTempKey string) *ApiRequest {
//	return newApiRequest(this, api, bodyTempKey)
//}
//
//func (this *Facebook) SetProxy(proxy proxys.Proxy) {
//	if proxy == proxys.DebugHttpProxy {
//		http_helper.DisableHttpSslPinng()(this.httpClient)
//	}
//	http_helper.HttpSetProxy(proxy)(this.httpClient)
//}
//
//func CreateFacebook(cookies *Cookies) *Facebook {
//	msg := &Facebook{
//		ck:         cookies,
//		tmpCk:      &TempCookies{},
//		httpClient: http_helper.CreateHttp2Client(http_helper.DisableRedirect()),
//	}
//	return msg
//}

//func (this *Facebook) UpdateCookiesValue(key string, value string) {
//	DbCtrl.UpdateCookiesValue(this, key, value)
//}
