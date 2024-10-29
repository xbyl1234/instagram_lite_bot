package messenger

import (
	"CentralizedControl/common/proxys"
	http "github.com/bogdanfinn/fhttp"
)

type Messenger struct {
	ck         *Cookies
	tmpCk      *TempCookies
	httpClient *http.Client
	proxy      proxys.Proxy
}

func CreateMessenger(cookies *Cookies) *Messenger {
	msg := &Messenger{
		ck: cookies,
		//httpClient: http2_helper.CreateHttp2Client(),
	}
	return msg
}

func (this *Messenger) newJsonGraphApi(api string) *GraphApiRequest {
	return newJsonGraphApi(this, api)
}

func (this *Messenger) sendGraphApi() {

}

//func (this *Messenger) SetProxy(proxy proxys.Proxy) {
//	if proxy == proxys.DebugHttpProxy {
//		http2_helper.DisableHttpSslPinng()(this.httpClient)
//	}
//	http2_helper.HttpSetProxy(proxy)(this.httpClient)
//}
