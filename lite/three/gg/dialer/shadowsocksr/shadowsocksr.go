package shadowsocksr

import (
	"encoding/base64"
	"fmt"
	"github.com/mzz2017/gg/common"
	"github.com/mzz2017/gg/dialer"
	ssr "github.com/v2rayA/shadowsocksR/client"
	"gopkg.in/yaml.v3"
	"net"
	"net/url"
	"strconv"
	"strings"
)

func init() {
	dialer.FromLinkRegister("shadowsocksr", NewShadowsocksR)
	dialer.FromLinkRegister("ssr", NewShadowsocksR)
	dialer.FromClashRegister("ssr", NewShadowsocksRFromClashObj)
}

type ShadowsocksR struct {
	Name       string `json:"name"`
	Server     string `json:"server"`
	Port       int    `json:"port"`
	Password   string `json:"password"`
	Cipher     string `json:"cipher"`
	Proto      string `json:"proto"`
	ProtoParam string `json:"protoParam"`
	Obfs       string `json:"obfs"`
	ObfsParam  string `json:"obfsParam"`
	Protocol   string `json:"protocol"`
}

func NewShadowsocksR(link string, opt *dialer.GlobalOption) (*dialer.Dialer, error) {
	s, err := ParseSSRURL(link)
	if err != nil {
		return nil, err
	}
	return s.Dialer()
}

func NewShadowsocksRFromClashObj(o *yaml.Node, opt *dialer.GlobalOption) (*dialer.Dialer, error) {
	s, err := ParseClash(o)
	if err != nil {
		return nil, err
	}
	return s.Dialer()
}

func (s *ShadowsocksR) Dialer() (*dialer.Dialer, error) {
	u := url.URL{
		Scheme: "ssr",
		User:   url.UserPassword(s.Cipher, s.Password),
		Host:   net.JoinHostPort(s.Server, strconv.Itoa(s.Port)),
		RawQuery: url.Values{
			"protocol":       []string{s.Proto},
			"protocol_param": []string{s.ProtoParam},
			"obfs":           []string{s.Obfs},
			"obfs_param":     []string{s.ObfsParam},
		}.Encode(),
	}
	d, err := ssr.NewSSR(u.String(), dialer.SymmetricDirect, nil)
	if err != nil {
		return nil, err
	}
	return dialer.NewDialer(d, false, s.Name, s.Protocol, s.ExportToURL()), nil
}

func ParseClash(o *yaml.Node) (data *ShadowsocksR, err error) {
	type ShadowSocksROption struct {
		Name          string `yaml:"name"`
		Server        string `yaml:"server"`
		Port          int    `yaml:"port"`
		Password      string `yaml:"password"`
		Cipher        string `yaml:"cipher"`
		Obfs          string `yaml:"obfs"`
		ObfsParam     string `yaml:"obfs-param,omitempty"`
		Protocol      string `yaml:"protocol"`
		ProtocolParam string `yaml:"protocol-param,omitempty"`
		UDP           bool   `yaml:"udp,omitempty"`
	}
	var option ShadowSocksROption
	if err = o.Decode(&option); err != nil {
		return nil, err
	}
	return &ShadowsocksR{
		Name:       option.Name,
		Server:     option.Server,
		Port:       option.Port,
		Password:   option.Password,
		Cipher:     option.Cipher,
		Proto:      option.Protocol,
		ProtoParam: option.ProtocolParam,
		Obfs:       option.Obfs,
		ObfsParam:  option.ObfsParam,
		Protocol:   "shadowsocksr",
	}, nil
}

func ParseSSRURL(u string) (data *ShadowsocksR, err error) {
	// parse attempts to parse ss:// links
	parse := func(content string) (v ShadowsocksR, ok bool) {
		arr := strings.Split(content, "/?")
		if strings.Contains(content, ":") && len(arr) < 2 {
			content += "/?remarks=&protoparam=&obfsparam="
			arr = strings.Split(content, "/?")
		} else if len(arr) != 2 {
			return v, false
		}
		pre := strings.Split(arr[0], ":")
		if len(pre) > 6 {
			//if the length is more than 6, it means that the host contains the characters:,
			//re-merge the first few groups into the host
			pre[len(pre)-6] = strings.Join(pre[:len(pre)-5], ":")
			pre = pre[len(pre)-6:]
		} else if len(pre) < 6 {
			return v, false
		}
		q, err := url.ParseQuery(arr[1])
		if err != nil {
			return v, false
		}
		pswd, _ := common.Base64URLDecode(pre[5])
		add, _ := common.Base64URLDecode(pre[0])
		remarks, _ := common.Base64URLDecode(q.Get("remarks"))
		protoparam, _ := common.Base64URLDecode(q.Get("protoparam"))
		obfsparam, _ := common.Base64URLDecode(q.Get("obfsparam"))
		port, err := strconv.Atoi(pre[1])
		if err != nil {
			return v, false
		}
		v = ShadowsocksR{
			Name:       remarks,
			Server:     add,
			Port:       port,
			Password:   pswd,
			Cipher:     pre[3],
			Proto:      pre[2],
			ProtoParam: protoparam,
			Obfs:       pre[4],
			ObfsParam:  obfsparam,
			Protocol:   "shadowsocksr",
		}
		return v, true
	}
	content := u[6:]
	var (
		info ShadowsocksR
		ok   bool
	)
	// try parsing the ssr:// link, if it fails, base64 decode first
	if info, ok = parse(content); !ok {
		// perform base64 decoding and parse again
		content, err = common.Base64StdDecode(content)
		if err != nil {
			content, err = common.Base64URLDecode(content)
			if err != nil {
				return
			}
		}
		info, ok = parse(content)
	}
	if !ok {
		err = fmt.Errorf("%w: unrecognized ssr address", dialer.InvalidParameterErr)
		return
	}
	return &info, nil
}

func (s *ShadowsocksR) ExportToURL() string {
	/* ssr://server:port:proto:method:obfs:URLBASE64(password)/?remarks=URLBASE64(remarks)&protoparam=URLBASE64(protoparam)&obfsparam=URLBASE64(obfsparam)) */
	return fmt.Sprintf("ssr://%v", strings.TrimSuffix(base64.URLEncoding.EncodeToString([]byte(
		fmt.Sprintf(
			"%v:%v:%v:%v:%v/?remarks=%v&protoparam=%v&obfsparam=%v",
			net.JoinHostPort(s.Server, strconv.Itoa(s.Port)),
			s.Proto,
			s.Cipher,
			s.Obfs,
			base64.URLEncoding.EncodeToString([]byte(s.Password)),
			base64.URLEncoding.EncodeToString([]byte(s.Name)),
			base64.URLEncoding.EncodeToString([]byte(s.ProtoParam)),
			base64.URLEncoding.EncodeToString([]byte(s.ObfsParam)),
		),
	)), "="))
}
