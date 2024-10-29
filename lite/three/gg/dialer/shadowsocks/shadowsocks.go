package shadowsocks

import (
	"encoding/base64"
	"fmt"
	"github.com/mzz2017/gg/common"
	"github.com/mzz2017/gg/dialer"
	"github.com/mzz2017/gg/dialer/transport/simpleobfs"
	"github.com/mzz2017/softwind/protocol"
	"github.com/mzz2017/softwind/protocol/shadowsocks"
	"gopkg.in/yaml.v3"
	"net"
	"net/url"
	"strconv"
	"strings"
)

func init() {
	// Use random salt by default to decrease the boot time
	shadowsocks.DefaultSaltGeneratorType = shadowsocks.RandomSaltGeneratorType

	dialer.FromLinkRegister("shadowsocks", NewShadowsocksFromLink)
	dialer.FromLinkRegister("ss", NewShadowsocksFromLink)
	dialer.FromClashRegister("ss", NewShadowsocksFromClashObj)
}

type Shadowsocks struct {
	Name     string `json:"name"`
	Server   string `json:"server"`
	Port     int    `json:"port"`
	Password string `json:"password"`
	Cipher   string `json:"cipher"`
	Plugin   Sip003 `json:"plugin"`
	UDP      bool   `json:"udp"`
	Protocol string `json:"protocol"`
}

func NewShadowsocksFromLink(link string, opt *dialer.GlobalOption) (*dialer.Dialer, error) {
	s, err := ParseSSURL(link)
	if err != nil {
		return nil, err
	}
	return s.Dialer()
}

func NewShadowsocksFromClashObj(o *yaml.Node, opt *dialer.GlobalOption) (*dialer.Dialer, error) {
	s, err := ParseClash(o)
	if err != nil {
		return nil, err
	}
	return s.Dialer()
}

func (s *Shadowsocks) Dialer() (*dialer.Dialer, error) {
	// FIXME: support plain/none.
	switch s.Cipher {
	case "aes-256-gcm", "aes-128-gcm", "chacha20-poly1305", "chacha20-ietf-poly1305":
	default:
		return nil, fmt.Errorf("unsupported shadowsocks encryption method: %v", s.Cipher)
	}
	var err error
	supportUDP := s.UDP
	d := dialer.SymmetricDirect
	switch s.Plugin.Name {
	case "simple-obfs":
		uSimpleObfs := url.URL{
			Scheme: "simple-obfs",
			Host:   net.JoinHostPort(s.Server, strconv.Itoa(s.Port)),
			RawQuery: url.Values{
				"obfs": []string{s.Plugin.Opts.Obfs},
				"host": []string{s.Plugin.Opts.Host},
				"uri":  []string{s.Plugin.Opts.Path},
			}.Encode(),
		}
		d, err = simpleobfs.NewSimpleObfs(uSimpleObfs.String(), d)
		if err != nil {
			return nil, err
		}
		supportUDP = false
	}
	d, err = protocol.NewDialer("shadowsocks", d, protocol.Header{
		ProxyAddress: net.JoinHostPort(s.Server, strconv.Itoa(s.Port)),
		Cipher:       s.Cipher,
		Password:     s.Password,
		IsClient:     true,
	})
	if err != nil {
		return nil, err
	}
	return dialer.NewDialer(d, supportUDP, s.Name, s.Protocol, s.ExportToURL()), nil
}

func ParseClash(o *yaml.Node) (data *Shadowsocks, err error) {
	type simpleObfsOption struct {
		Mode string `obfs:"mode,omitempty"`
		Host string `obfs:"host,omitempty"`
	}
	type ShadowSocksOption struct {
		Name       string           `yaml:"name"`
		Server     string           `yaml:"server"`
		Port       int              `yaml:"port"`
		Password   string           `yaml:"password"`
		Cipher     string           `yaml:"cipher"`
		UDP        bool             `yaml:"udp,omitempty"`
		Plugin     string           `yaml:"plugin,omitempty"`
		PluginOpts simpleObfsOption `yaml:"plugin-opts,omitempty"`
	}
	var option ShadowSocksOption
	if err = o.Decode(&option); err != nil {
		return nil, err
	}
	data = &Shadowsocks{
		Name:     option.Name,
		Server:   option.Server,
		Port:     option.Port,
		Password: option.Password,
		Cipher:   option.Cipher,
		UDP:      option.UDP,
		Protocol: "shadowsocks",
	}
	if option.Plugin == "obfs" {
		data.Plugin.Name = "simple-obfs"
		data.Plugin.Opts.Obfs = option.PluginOpts.Mode
		data.Plugin.Opts.Host = data.Plugin.Opts.Host
		if data.Plugin.Opts.Host == "" {
			data.Plugin.Opts.Host = "bing.com"
		}
	}
	return data, nil
}

func ParseSSURL(u string) (data *Shadowsocks, err error) {
	// parse attempts to parse ss:// links
	parse := func(content string) (v *Shadowsocks, ok bool) {
		// try to parse in the format of ss://BASE64(method:password)@server:port/?plugin=xxxx#name
		u, err := url.Parse(content)
		if err != nil {
			return nil, false
		}
		username := u.User.String()
		username, _ = common.Base64URLDecode(username)
		arr := strings.SplitN(username, ":", 2)
		if len(arr) != 2 {
			return nil, false
		}
		cipher := arr[0]
		password := arr[1]
		var sip003 Sip003
		plugin := u.Query().Get("plugin")
		if len(plugin) > 0 {
			sip003 = ParseSip003(plugin)
		}
		port, err := strconv.Atoi(u.Port())
		if err != nil {
			return nil, false
		}
		return &Shadowsocks{
			Cipher:   strings.ToLower(cipher),
			Password: password,
			Server:   u.Hostname(),
			Port:     port,
			Name:     u.Fragment,
			Plugin:   sip003,
			UDP:      sip003.Name == "",
			Protocol: "shadowsocks",
		}, true
	}
	var (
		v  *Shadowsocks
		ok bool
	)
	content := u
	// try to parse the ss:// link, if it fails, base64 decode first
	if v, ok = parse(content); !ok {
		// 进行base64解码，并unmarshal到VmessInfo上
		t := content[5:]
		var l, r string
		if ind := strings.Index(t, "#"); ind > -1 {
			l = t[:ind]
			r = t[ind+1:]
		} else {
			l = t
		}
		l, err = common.Base64StdDecode(l)
		if err != nil {
			l, err = common.Base64URLDecode(l)
			if err != nil {
				return
			}
		}
		t = "ss://" + l
		if len(r) > 0 {
			t += "#" + r
		}
		v, ok = parse(t)
	}
	if !ok {
		return nil, fmt.Errorf("%w: unrecognized ss address", dialer.InvalidParameterErr)
	}
	return v, nil
}

type Sip003 struct {
	Name string     `json:"name"`
	Opts Sip003Opts `json:"opts"`
}
type Sip003Opts struct {
	Tls  string `json:"tls"` // for v2ray-plugin
	Obfs string `json:"obfs"`
	Host string `json:"host"`
	Path string `json:"uri"`
}

func ParseSip003Opts(opts string) Sip003Opts {
	var sip003Opts Sip003Opts
	fields := strings.Split(opts, ";")
	for i := range fields {
		a := strings.Split(fields[i], "=")
		if len(a) == 1 {
			// to avoid panic
			a = append(a, "")
		}
		switch a[0] {
		case "tls":
			sip003Opts.Tls = "tls"
		case "obfs", "mode":
			sip003Opts.Obfs = a[1]
		case "obfs-path", "obfs-uri", "path":
			if !strings.HasPrefix(a[1], "/") {
				a[1] += "/"
			}
			sip003Opts.Path = a[1]
		case "obfs-host", "host":
			sip003Opts.Host = a[1]
		}
	}
	return sip003Opts
}
func ParseSip003(plugin string) Sip003 {
	var sip003 Sip003
	fields := strings.SplitN(plugin, ";", 2)
	switch fields[0] {
	case "obfs-local", "simpleobfs":
		sip003.Name = "simple-obfs"
	default:
		sip003.Name = fields[0]
	}
	sip003.Opts = ParseSip003Opts(fields[1])
	return sip003
}

func (s *Sip003) String() string {
	list := []string{s.Name}
	if s.Opts.Obfs != "" {
		list = append(list, "obfs="+s.Opts.Obfs)
	}
	if s.Opts.Host != "" {
		list = append(list, "obfs-host="+s.Opts.Host)
	}
	if s.Opts.Path != "" {
		list = append(list, "obfs-uri="+s.Opts.Path)
	}
	return strings.Join(list, ";")
}

func (s *Shadowsocks) ExportToURL() string {
	// sip002
	u := &url.URL{
		Scheme:   "ss",
		User:     url.User(strings.TrimSuffix(base64.URLEncoding.EncodeToString([]byte(s.Cipher+":"+s.Password)), "=")),
		Host:     net.JoinHostPort(s.Server, strconv.Itoa(s.Port)),
		Fragment: s.Name,
	}
	if s.Plugin.Name != "" {
		q := u.Query()
		q.Set("plugin", s.Plugin.String())
		u.RawQuery = q.Encode()
	}
	return u.String()
}
