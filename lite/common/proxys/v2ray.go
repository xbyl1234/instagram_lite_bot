package proxys

import (
	"CentralizedControl/common/log"
	"errors"
	"fmt"
	"github.com/mzz2017/gg/dialer"
	"github.com/mzz2017/gg/dialer/v2ray"
	_ "github.com/mzz2017/softwind/protocol/vless"
	_ "github.com/mzz2017/softwind/protocol/vmess"
	"strconv"
	"strings"
)

type VmessInfo struct {
	AlterId                          int    `json:"alterId"`
	ExperimentalAuthenticatedLength  bool   `json:"experimentalAuthenticatedLength"`
	ExperimentalNoTerminationSignal  bool   `json:"experimentalNoTerminationSignal"`
	AllowInsecure                    bool   `json:"allowInsecure"`
	Alpn                             string `json:"alpn"`
	Certificates                     string `json:"certificates"`
	EarlyDataHeaderName              string `json:"earlyDataHeaderName"`
	Encryption                       string `json:"encryption"`
	GrpcServiceName                  string `json:"grpcServiceName"`
	HeaderType                       string `json:"headerType"`
	Host                             string `json:"host"`
	MKcpSeed                         string `json:"mKcpSeed"`
	PacketEncoding                   int    `json:"packetEncoding"`
	Path                             string `json:"path"`
	PinnedPeerCertificateChainSha256 string `json:"pinnedPeerCertificateChainSha256"`
	QuicKey                          string `json:"quicKey"`
	QuicSecurity                     string `json:"quicSecurity"`
	Security                         string `json:"security"`
	Sni                              string `json:"sni"`
	Type                             string `json:"type"`
	UtlsFingerprint                  string `json:"utlsFingerprint"`
	Uuid                             string `json:"uuid"`
	WsMaxEarlyData                   int    `json:"wsMaxEarlyData"`
	WsUseBrowserForwarder            bool   `json:"wsUseBrowserForwarder"`
	ExtraType                        int    `json:"extraType"`
	Group                            string `json:"group"`
	Name                             string `json:"name"`
	ProfileId                        string `json:"profileId"`
	ServerAddress                    string `json:"serverAddress"`
	ServerPort                       int    `json:"serverPort"`
}

func Conv2VmessInfo(vmess *v2ray.V2Ray) *VmessInfo {
	prot, _ := strconv.Atoi(vmess.Port)
	aid, _ := strconv.Atoi(vmess.Aid)
	return &VmessInfo{
		AlterId:                          aid,
		ExperimentalAuthenticatedLength:  false,
		ExperimentalNoTerminationSignal:  false,
		AllowInsecure:                    false,
		Alpn:                             "",
		Certificates:                     "",
		EarlyDataHeaderName:              "",
		Encryption:                       "",
		GrpcServiceName:                  "",
		HeaderType:                       "",
		Host:                             "",
		MKcpSeed:                         "",
		PacketEncoding:                   0,
		Path:                             "",
		PinnedPeerCertificateChainSha256: "",
		QuicKey:                          "",
		QuicSecurity:                     "",
		Security:                         "",
		Sni:                              "",
		Type:                             vmess.Net,
		UtlsFingerprint:                  "",
		Uuid:                             vmess.ID,
		WsMaxEarlyData:                   0,
		WsUseBrowserForwarder:            false,
		ExtraType:                        0,
		Group:                            "",
		Name:                             vmess.Ps,
		ProfileId:                        "",
		ServerAddress:                    vmess.Add,
		ServerPort:                       prot,
	}
}

type V2RayProxy struct {
	*ProxyInfo
	*v2ray.V2Ray
}

func CreateV2rayClient(link string) (client *V2RayProxy, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Error("CreateV2rayClient recover error: %v", err)
			client = nil
			err = r.(error)
		}
	}()
	var protocol string
	client = &V2RayProxy{}
	switch {
	case strings.HasPrefix(link, "vmess://"):
		client.V2Ray, err = v2ray.ParseVmessURL(link)
		if err != nil {
			return nil, err
		}
		if client.Aid != "0" && client.Aid != "" {
			return nil, fmt.Errorf("%w: aid: %v, we only support AEAD encryption", dialer.UnexpectedFieldErr, client.Aid)
		}
		protocol = "vmess"
	case strings.HasPrefix(link, "vless://"):
		client.V2Ray, err = v2ray.ParseVlessURL(link)
		if err != nil {
			return nil, err
		}
		protocol = "vless"
	default:
		return nil, errors.New("unknow protocol")
	}
	client.AllowInsecure = true
	client.ProxyInfo = createProxyInfo(protocol, client.Ps, Conv2VmessInfo(client.V2Ray), client)
	return client, nil
}

func (this *V2RayProxy) GetDialer() Dialer {
	d, err := this.V2Ray.Dialer()
	if err != nil {
		log.Warn("V2RayProxy GetDialer error: %v", err)
	}
	return d
}
