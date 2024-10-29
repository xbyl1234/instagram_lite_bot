package common

import (
	"CentralizedControl/common/utils"
	"time"
)

type RandomConfig struct {
	Rate        float64
	Range       float64
	ClickBtn    time.Duration
	InputText   time.Duration
	NewViewShow time.Duration
}

var DefaultRandomConfig = RandomConfig{
	Rate:        1,
	Range:       1,
	ClickBtn:    2 * time.Second,
	InputText:   10 * time.Second,
	NewViewShow: 5 * time.Second,
}

type RandomSleep struct {
	Config RandomConfig
}

func CreateRandomSleep(config RandomConfig) *RandomSleep {
	return &RandomSleep{
		Config: config,
	}
}

func (this *RandomSleep) SleepForClickBtn() {
	this.Sleep(this.Config.ClickBtn)
}

func (this *RandomSleep) SleepForInputText() {
	this.Sleep(this.Config.InputText)
}

func (this *RandomSleep) SleepForRecvNewView() {
	this.Sleep(this.Config.NewViewShow)
}

func (this *RandomSleep) Sleep(base time.Duration) {
	time.Sleep(time.Duration(float64(base) * this.Config.Range * utils.GenFloat(1-this.Config.Range, 1+this.Config.Range)))
}
