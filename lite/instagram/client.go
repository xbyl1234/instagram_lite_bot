package instagram

import (
	"CentralizedControl/common/proxys"
	http "github.com/bogdanfinn/fhttp"
	"github.com/bytedance/sonic/ast"
	"strconv"
)

type Instagram struct {
	ck         *Cookies
	tmpCk      *TempCookies
	httpClient *http.Client
	proxy      proxys.Proxy
}

func (this *Instagram) GetCookies() *Cookies {
	return this.ck
}

func (this *Instagram) GetSelfUserId() uint64 {
	id, _ := strconv.ParseInt(this.ck.AccountId, 10, 64)
	return uint64(id)
}

func (this *Instagram) newApiRequest(api string, bodyTempKey string) *ApiRequest {
	return newApiRequest(this, api, bodyTempKey)
}

//
//func (this *Instagram) SetProxy(proxy proxys.Proxy) {
//	if proxy == proxys.DebugHttpProxy {
//		http2_helper.DisableHttpSslPinng()(this.httpClient)
//	}
//	http2_helper.HttpSetProxy(proxy)(this.httpClient)
//}
//
//func CreateInstagram(cookies *Cookies) *Instagram {
//	msg := &Instagram{
//		ck: cookies,
//		tmpCk: &TempCookies{
//			BandwidthTotaltimeMs: utils.GenNumber(335, 350),
//			PigeonSessionId:      "UFS-" + utils.GenUUID() + "-0",
//			TimezoneOffset:       fmt.Sprintf("%d", utils.GetTimezoneOffset(cookies.Timezone)),
//		},
//		httpClient: http2_helper.CreateHttp2Client(),
//	}
//	return msg
//}

func (this *Instagram) FollowUserByQrCode(url string) error {
	qrFollow := CreateFollowPeople(this, DecodeShareUrl(url))
	err := qrFollow.NameTagLookupByName()
	if err != nil {
		return err
	}
	return qrFollow.FollowPeople()
}

func (this *Instagram) GetUserInfo(id uint64) (*ast.Node, error) {
	request := this.newApiRequest("/api/v1/users/%v/info/", "")
	request.SetPathParams(id)
	send, err := request.Send()
	if send != nil && send.Json != nil {
		return send.Json, err
	}
	return nil, err
}

//func (this *Instagram) UpdateCookiesValue(key string, value string) {
//	DbCtrl.UpdateCookiesValue(this, key, value)
//}
