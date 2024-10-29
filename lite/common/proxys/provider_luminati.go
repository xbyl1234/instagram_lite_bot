package proxys

import (
	"CentralizedControl/common/utils"
	"os"
	"strings"
)

type LuminatiProvider struct {
	ipFile       string
	localFile    string
	local2IpMaps map[string]*[]*HttpProxy
}

func convHttpProxy(line string) *HttpProxy {
	sp := strings.Split(line, ":")
	if len(sp) != 4 {
		return nil
	}
	proxy, err := CreateHttpProxy("http://" + sp[2] + ":" + sp[3] + "@" + sp[0] + ":" + sp[1])
	if err != nil {
		return nil
	}
	return proxy
}

func (this *LuminatiProvider) GetProxy(country string) *HttpProxy {
	ps := this.local2IpMaps[country]
	if ps == nil {
		return nil
	}
	return (*ps)[utils.GenNumber(0, len(*ps)-1)]
}

func CreateLuminatiProvider(ipFile string, localFile string) *LuminatiProvider {
	luminati := &LuminatiProvider{
		ipFile:       ipFile,
		localFile:    localFile,
		local2IpMaps: make(map[string]*[]*HttpProxy),
	}
	data, err := os.ReadFile("./luminati/" + localFile)
	lines := strings.Split(string(data), "\n")
	ip2Local := make(map[string]string)
	for _, line := range lines {
		sp := strings.Split(line, ",")
		if len(sp) != 2 {
			continue
		}
		ip2Local[sp[0]] = sp[1]
	}

	data, err = os.ReadFile("./luminati/" + ipFile)
	if err != nil {
		return nil
	}
	lines = strings.Split(string(data), "\n")
	for _, line := range lines {
		p := convHttpProxy(line)
		if p == nil {
			continue
		}
		ip := p.Username[strings.LastIndex(p.Username, "-"):]
		ps := luminati.local2IpMaps[ip]
		if ps == nil {
			tmp := make([]*HttpProxy, 100)
			ps = &tmp
			luminati.local2IpMaps[ip] = ps
		}
		*ps = append(*ps, p)
	}

	return luminati
}
