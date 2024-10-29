package proxys

import (
	"CentralizedControl/common"
	"CentralizedControl/common/log"
	"CentralizedControl/common/utils"
	"fmt"
	ssr "github.com/v2rayA/shadowsocksR/client"
	"golang.org/x/net/proxy"
	"net/url"
	"strconv"
	"strings"
)

type SsrConfig struct {
	Method        string `json:"method"`
	Passwd        string `json:"password"`
	Host          string `json:"serverAddress"`
	Port          int    `json:"serverPort"`
	Obfs          string `json:"obfs"`
	ObfsParam     string `json:"obfsParam"`
	Protocol      string `json:"protocol"`
	ProtocolParam string `json:"protocolParam"`
	Group         string `json:"group"`
	Remarks       string `json:"name"`
	UdpPort       string
	UOT           string
}

type SsrProxy struct {
	*ProxyInfo
	*SsrConfig
}

func (this *SsrProxy) GetDialer() Dialer {
	s, err := convertDialerURL(this.SsrConfig)
	if err != nil {
	}
	newSSR, err := ssr.NewSSR(s, proxy.Direct, nil)
	if err != nil {
		return nil
	}
	return newSSR
}

func CreateSsrClient(lineData string) (client *SsrProxy, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Error("CreateSsrClient recover error: %v", err)
			client = nil
			err = r.(error)
		}
	}()
	client = &SsrProxy{}
	client.SsrConfig, err = parseSSR(lineData)
	if err != nil {
		return nil, err
	}
	client.ProxyInfo = createProxyInfo("ssr", client.Remarks, client.SsrConfig, client)
	return client, nil
}

func convertDialerURL(params *SsrConfig) (s string, err error) {
	u, err := url.Parse(fmt.Sprintf(
		"ssr://%v:%v@%v:%v",
		params.Method,
		params.Passwd,
		params.Host,
		params.Port,
	))
	if err != nil {
		return
	}
	q := u.Query()
	if len(strings.TrimSpace(params.Obfs)) <= 0 {
		params.Obfs = "plain"
	}
	if len(strings.TrimSpace(params.Protocol)) <= 0 {
		params.Protocol = "origin"
	}
	q.Set("obfs", params.Obfs)
	q.Set("obfs_param", params.ObfsParam)
	q.Set("protocol", params.Protocol)
	q.Set("protocol_param", params.ProtocolParam)
	u.RawQuery = q.Encode()
	return u.String(), nil
}

func parseSSR(lineData string) (*SsrConfig, error) {
	var params = &SsrConfig{}
	if len(lineData) < 7 {
		return nil, common.NerError("ssr link too short")
	}
	lineData = lineData[6:]
	lineData = strings.ReplaceAll(lineData, "_", "/")
	lineData = strings.ReplaceAll(lineData, "-", "+")
	dl, _ := utils.DecodeBase64(lineData)
	sp := strings.Split(string(dl), "?")
	if len(sp) != 2 {
		return nil, common.NerError("link ftm error!")
	}
	spFornt := strings.Split(strings.ReplaceAll(sp[0], "/", ""), ":")
	if len(spFornt) != 6 {
		return nil, common.NerError("link ftm error!")
	}
	params.Host = spFornt[0]
	port, _ := strconv.ParseInt(spFornt[1], 10, 64)
	params.Port = int(port)
	params.Protocol = spFornt[2]
	params.Method = spFornt[3]
	params.Obfs = spFornt[4]
	passwd, _ := utils.DecodeBase64(spFornt[5])
	params.Passwd = string(passwd)
	parse := url.URL{RawQuery: sp[1]}

	var param, _ = utils.DecodeBase64(parse.Query().Get("obfsparam"))
	params.ObfsParam = string(param)
	param, _ = utils.DecodeBase64(parse.Query().Get("protoparam"))
	params.ProtocolParam = string(param)
	param, _ = utils.DecodeBase64(parse.Query().Get("group"))
	params.Group = string(param)
	param, _ = utils.DecodeBase64(parse.Query().Get("remarks"))
	params.Remarks = string(param)

	params.UdpPort = parse.Query().Get("udpport")
	params.UOT = parse.Query().Get("uot")
	return params, nil
}
