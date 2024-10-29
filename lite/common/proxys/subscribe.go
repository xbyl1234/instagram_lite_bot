package proxys

import (
	"CentralizedControl/common"
	"CentralizedControl/common/http_helper"
	"CentralizedControl/common/log"
	"CentralizedControl/common/utils"
	"strings"
)

type SubscribeProvider struct {
	Subscribe string
	Type      string
}

func (this *SubscribeProvider) GetLink() string {
	return this.Subscribe
}

func (this *SubscribeProvider) Update() ([]Proxy, error) {
	c := http_helper.CreateGoHttpClient(http_helper.HttpTimeout(20))
	resp, err := http_helper.HttpDoRetry(c, &http_helper.RequestOpt{
		ReqUrl: this.Subscribe,
	})
	if err != nil {
		log.Error("update subscribe %s error: %v", this.Subscribe, err)
		return nil, err
	}
	var data string
	this.Type, data, err = this.GuessProvider(resp)
	if err != nil {
		log.Error("guess subscribe %s error: %v", this.Subscribe, err)
		return nil, err
	}

	switch this.Type {
	case "ss":
		return this.CreateMultipleProxy(data, this.CreateSsProxyInfo)
	case "ssr":
		return this.CreateMultipleProxy(data, this.CreateSsrProxyInfo)
	case "vless", "vmess":
		return this.CreateMultipleProxy(data, this.CreateV2RayProxyInfo)
	default:
		log.Error("unknown subscribe %s", this.Subscribe)
		return nil, common.NerError("unknow imap_provider")
	}
}

func (this *SubscribeProvider) GuessProvider(data string) (string, string, error) {
	decode, _ := utils.DecodeBase64(data)
	data = string(decode)
	if strings.HasPrefix(data, "ss://") {
		return "ss", data, nil
	} else if strings.HasPrefix(data, "ssr://") {
		return "ssr", data, nil
	} else if strings.HasPrefix(data, "trojan://") {
		return "trojan", data, nil
	} else if strings.HasPrefix(data, "vmess://") {
		return "vmess", data, nil
	} else if strings.HasPrefix(data, "vless://") {
		return "vless", data, nil
	}
	return "", "", common.NerError("unknown protocol")
}

func (this *SubscribeProvider) CreateMultipleProxy(data string, create func(data string) (Proxy, error)) ([]Proxy, error) {
	sp := strings.Split(string(data), "\n")
	result := make([]Proxy, len(sp))
	count := 0
	for _, line := range sp {
		if len(line) < 5 {
			continue
		}
		ssr, err := create(line)
		if err != nil {
			continue
		}
		result[count] = ssr
		count++
	}
	return result[:count], nil
}

func (this *SubscribeProvider) CreateSsrProxyInfo(data string) (Proxy, error) {
	return CreateSsrClient(data)
}

func (this *SubscribeProvider) CreateV2RayProxyInfo(data string) (Proxy, error) {
	return CreateV2rayClient(data)
}
func (this *SubscribeProvider) CreateSsProxyInfo(data string) (Proxy, error) {
	return CreateSsClient(data)
}

func CreateSubscribeProvider(link string) *SubscribeProvider {
	sub := &SubscribeProvider{
		Subscribe: link,
	}
	return sub
}
