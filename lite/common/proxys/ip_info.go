package proxys

import (
	"CentralizedControl/common/http_helper"
	"CentralizedControl/common/log"
)

type IpInfo struct {
	Status        string  `json:"status"`
	Continent     string  `json:"continent"`
	ContinentCode string  `json:"continentCode"`
	Country       string  `json:"country"`
	CountryCode   string  `json:"countryCode"`
	Region        string  `json:"region"`
	RegionName    string  `json:"regionName"`
	City          string  `json:"city"`
	District      string  `json:"district"`
	Zip           string  `json:"zip"`
	Lat           float64 `json:"lat"`
	Lon           float64 `json:"lon"`
	Timezone      string  `json:"timezone"`
	Offset        int     `json:"offset"`
	Currency      string  `json:"currency"`
	Isp           string  `json:"isp"`
	Org           string  `json:"org"`
	As            string  `json:"as"`
	Asname        string  `json:"asname"`
	Mobile        bool    `json:"mobile"`
	Proxy         bool    `json:"proxy"`
	Hosting       bool    `json:"hosting"`
	Query         string  `json:"query"`
}

func GetIpInfoByIp(ip string, proxy Proxy) (*IpInfo, error) {
	client := http_helper.CreateGoHttpClient(http_helper.EnableHttp2(),
		http_helper.HttpTimeout(20))
	if proxy != nil {
		proxy.GetProxy()(client)
	}
	var info IpInfo
	err := http_helper.HttpDoJson(client, &http_helper.RequestOpt{
		Header: map[string]string{
			"Origin":  "https://ip-api.com",
			"Referer": "https://ip-api.com/",
		},
		IsPost: false,
		ReqUrl: "https://demo.ip-api.com/json/" + ip + "?fields=66842623&lang=en",
	}, &info)
	if err != nil {
		log.Error("GetIpInfoByIp: %s error: %v", ip, err)
	}
	return &info, err
}
