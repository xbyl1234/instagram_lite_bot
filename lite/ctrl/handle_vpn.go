package ctrl

import (
	"CentralizedControl/common/proxys"
	"github.com/bytedance/sonic/ast"
	"net/http"
	"os"
	"strings"
)

func (this *HttpServer) HttpHandleVpnGetSetting(req *http.Request, reqJson *ast.Node) (interface{}, error) {
	file, err := os.ReadFile("./proxys_setting.json")
	if err != nil {
		return nil, err
	}
	return string(file), nil
}

func (this *HttpServer) HttpHandleVpnGetProxysFile(req *http.Request, reqJson *ast.Node) (interface{}, error) {
	file, err := os.ReadFile("./proxys_file.json")
	if err != nil {
		return nil, err
	}
	return string(file), nil
}

type ProxyItem struct {
	GroupId       int    `json:"groupId"`
	Id            int    `json:"id"`
	Ping          int    `json:"ping"`
	Rx            int    `json:"rx"`
	SocksBean     any    `json:"socksBean"`
	HttpBean      any    `json:"httpBean"`
	SsBean        any    `json:"ssBean"`
	SsrBean       any    `json:"ssrBean"`
	VmessBean     any    `json:"vmessBean"`
	Status        int    `json:"status"`
	Tx            int    `json:"tx"`
	Type          int    `json:"type"`
	UserOrder     int    `json:"userOrder"`
	Uuid          string `json:"uuid"`
	CtrlProxyId   string `json:"ctrl_proxy_id"`
	CtrlProxyName string `json:"ctrl_proxy_name"`
}

type ProxyResp struct {
	Proxys []*ProxyItem `json:"proxys"`
}

func getIp(ipPort string) string {
	sp := strings.Split(ipPort, ":")
	return sp[0]
}

func (this *HttpServer) HttpHandleVpnGetProxys(req *http.Request, reqJson *ast.Node) (interface{}, error) {
	var _proxy proxys.Proxy
	switch ConfigProxyUse {
	case ConfigProxyTypeSubscribe:
		_proxy = this.proxyManage.GetSubscribeProxy()
	case ConfigProxyTypeLuminati:
		_proxy = this.proxyManage.GetLuminatiProxy("hk")
	case ConfigProxyTypeOxylabs:
		_proxy = this.proxyManage.GetOxylabsProxy("hk")
	case ConfigProxyTypeZenoo:
		_proxy = this.proxyManage.GetZenooProxy(getIp(req.RemoteAddr), "hk")
	}
	resp := ProxyResp{
		Proxys: make([]*ProxyItem, 1),
	}
	resp.Proxys[0] = &ProxyItem{
		Type:          _proxy.GetTypeInt(),
		CtrlProxyId:   _proxy.GetId(),
		CtrlProxyName: _proxy.GetName(),
	}
	switch _proxy.GetTypeInt() {
	case proxys.ProxySocks:
		resp.Proxys[0].SocksBean = _proxy.GetData()
	case proxys.ProxyHttp:
		resp.Proxys[0].HttpBean = _proxy.GetData()
	case proxys.ProxySs:
		resp.Proxys[0].SsBean = _proxy.GetData()
	case proxys.ProxySsr:
		resp.Proxys[0].SsrBean = _proxy.GetData()
	case proxys.ProxyVmess:
		resp.Proxys[0].VmessBean = _proxy.GetData()
	}

	return resp, nil
}

type IpInfos struct {
	Ip      string `json:"ip"`
	Asn     string `json:"asn"`
	Country string `json:"country"`
}

func (this *HttpServer) HttpHandleGetIpInfos(req *http.Request, reqJson *ast.Node) (interface{}, error) {
	ipRaw, err := proxys.GetIpInfoByIp(reqJson.GetString("ip"), nil)
	if err != nil {
		return nil, err
	}

	return &IpInfos{
		Ip:      reqJson.GetString("ip"),
		Asn:     ipRaw.Asname,
		Country: strings.ToLower(ipRaw.CountryCode),
	}, nil
}
