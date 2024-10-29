package instagram

import (
	"CentralizedControl/common/utils"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type NavChainItem struct {
	ViewName  string
	SubName   string
	Index     int
	ActName   string
	StartTime string
	Unknow1   string
	Unknow2   string
}

func CreateNavChainItem(temp string) *NavChainItem {
	result := &NavChainItem{}
	sp := strings.Split(temp, ":")
	result.ViewName = sp[0]
	result.SubName = sp[1]
	result.Index, _ = strconv.Atoi(sp[2])
	result.ActName = sp[3]
	result.StartTime = sp[4]
	result.Unknow1 = sp[5]
	result.Unknow2 = sp[6]
	return result
}

type NavChain struct {
	chain          []*NavChainItem
	fakeTime       bool
	serializeCache string
}

func (this *NavChain) Push(viewName, subName, actName string) *NavChainItem {
	index := 0
	if len(this.chain) > 0 {
		index = this.chain[len(this.chain)-1].Index + 1
	}
	item := &NavChainItem{
		ViewName:  viewName,
		SubName:   subName,
		Index:     index,
		ActName:   actName,
		StartTime: fmt.Sprintf("%.03f", float64(time.Now().UnixMilli())/1000.0),
		Unknow1:   "",
		Unknow2:   "",
	}
	this.chain = append(this.chain, item)
	return item
}

func (this *NavChain) Serialize() string {
	if this.serializeCache != "" {
		return this.serializeCache
	}
	now := time.Now().UnixMilli()
	del := 0.0
	if this.fakeTime {
		for idx := range this.chain {
			this.chain[len(this.chain)-idx-1].StartTime = fmt.Sprintf("%.03f", float64(now)/1000.0-del)
			del += float64(utils.GenNumber(1000, 10000)) / 1000.0
		}
	}
	result := ""
	for _, item := range this.chain {
		result += item.ViewName + ":"
		result += item.SubName + ":"
		result += fmt.Sprintf("%d", item.Index) + ":"
		result += item.ActName + ":"
		result += item.StartTime + ":"
		result += item.Unknow1 + ":"
		result += item.Unknow2
		result += ","
	}
	this.serializeCache = result[:len(result)-1]
	return this.serializeCache
}

func CreateNavChain(temp string) *NavChain {
	navChain := &NavChain{
		chain:    nil,
		fakeTime: true,
	}
	sp := strings.Split(temp, ",")
	for _, item := range sp {
		navChain.chain = append(navChain.chain, CreateNavChainItem(item))
	}
	return navChain
}
