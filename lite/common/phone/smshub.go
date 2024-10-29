package phone

import (
	"CentralizedControl/common/android"
	"CentralizedControl/common/http_helper"
	"CentralizedControl/common/log"
	"errors"
	"net/http"
	"strings"
	"time"
)

const host = "https://smshub.org/stubs/handler_api.php"
const apiKey = "213294Uaa48dc6027b1b8d6488935ce0b6d3ac0"
const (
	England = "16"
	Hk      = "14"
	India   = "22"
	Us      = "12"
	Sg      = "196"
)
const (
	Instagram = "ig"
	Facebook  = "fb"
)

var CountryCode2Number = map[string]string{
	"gb": England,
	"hk": Hk,
	"in": India,
	"us": Us,
	"sg": Sg,
}

var DefaultRetryConfig = &RetryConfig{
	RetryTimeoutDuration: 5 * time.Minute,
	RetryDelayDuration:   5 * time.Second,
}

type RetryConfig struct {
	LastReqTime          time.Time
	RetryTimeoutDuration time.Duration
	RetryDelayDuration   time.Duration
}

type SmsHub struct {
	client      *http.Client
	RetryConfig *RetryConfig
	ApiKey      string `json:"api_key"`
	MaxPrice    string `json:"max_price"`
	Country     string `json:"country"`
	Operator    string `json:"operator"`
	Service     string `json:"service"`
}

type Number struct {
	client    *SmsHub
	StartTime time.Time
	Number    string
	AreaCode  string
	Id        string
}

func (this *Number) SyncGetCode() (string, error) {
	start := time.Now()
	for time.Since(start) < this.client.RetryConfig.RetryTimeoutDuration {
		code, err := this.client.GetCode(this)
		if code != "" {
			return code, nil
		}
		if err != nil {
			return "", err
		}
		time.Sleep(this.client.RetryConfig.RetryDelayDuration)
	}
	return "", errors.New("timeout")
}

func CreateSmsHub(service string, country string) *SmsHub {
	return &SmsHub{
		client:      &http.Client{},
		RetryConfig: DefaultRetryConfig,
		ApiKey:      apiKey,
		MaxPrice:    "15",
		Country:     country,
		Operator:    "",
		Service:     service,
	}
}

func getNumber(country string, number string) string {
	areaCode := android.GetAreaCode(country)
	return number[len(areaCode)-1:]
}

func getAreaCode(country string, number string) string {
	areaCode := android.GetAreaCode(country)
	return areaCode[1:]
}

func (this *SmsHub) GetPhoneNumber() *Number {
	resp, err := http_helper.HttpDo(this.client, &http_helper.RequestOpt{
		Params: map[string]string{
			"api_key":  this.ApiKey,
			"action":   "getNumber",
			"service":  this.Service,
			"operator": this.Operator,
			"country":  CountryCode2Number[this.Country],
			"maxPrice": this.MaxPrice,
		},
		IsPost: false,
		ReqUrl: host,
	})
	if err != nil {
		log.Error("GetPhoneNumber error: %v", err)
		return nil
	}
	if strings.Contains(resp, "ACCESS_NUMBER") {
		sp := strings.Split(resp, ":")
		if len(sp) != 3 {
			log.Error("GetPhoneNumber error: %v", resp)
			return nil
		}
		number := &Number{
			client:    this,
			StartTime: time.Now(),
			Number:    getNumber(this.Country, sp[2]),
			AreaCode:  getAreaCode(this.Country, sp[2]),
			Id:        sp[1],
		}
		log.Info("get phone number: %s-%s, id: %s", number.AreaCode, number.Number, number.Id)
		return number
	} else {
		log.Error("GetPhoneNumber error: %v", resp)
		return nil
	}
}

func (this *SmsHub) SetStatus(number *Number, status string) (string, error) {
	return http_helper.HttpDo(this.client, &http_helper.RequestOpt{
		Params: map[string]string{
			"api_key": this.ApiKey,
			"action":  "setStatus",
			"status":  status,
			"id":      number.Id,
		},
		IsPost: false,
		ReqUrl: host,
	})
}

func (this *SmsHub) GetStatus(number *Number) (string, error) {
	resp, err := http_helper.HttpDo(this.client, &http_helper.RequestOpt{
		Params: map[string]string{
			"api_key": this.ApiKey,
			"action":  "getStatus",
			"id":      number.Id,
		},
		IsPost: false,
		ReqUrl: host,
	})
	if err != nil {
		return "", err
	}
	return resp, err
}

func (this *SmsHub) Release(number *Number) {
	this.SetStatus(number, "8")
}

func (this *SmsHub) Continue(number *Number) (err error) {
	defer func() {
		if err != nil {
			log.Error("Continue number: %s, error: %v", number.Number, err)
		}
	}()
	var resp string
	resp, err = this.GetStatus(number)
	if err != nil {
		return err
	}
	if strings.Contains(resp, "STATUS_OK") {
		resp, err = this.SetStatus(number, "3")
		if err != nil {
			return err
		}
		if resp != "ACCESS_RETRY_GET" {
			return errors.New(resp)
		}
		return nil
	} else if resp == "STATUS_WAIT_CODE" {
		return nil
	} else {
		return errors.New(resp)
	}
}

func (this *SmsHub) GetCode(number *Number) (string, error) {
	status, err := this.GetStatus(number)
	log.Info("%v %v", status, err)
	if status == "STATUS_WAIT_CODE" {
		log.Debug("wait for code...")
		return "", nil
	}
	if strings.Contains(status, "STATUS_OK") {
		sp := strings.Split(status, ":")
		log.Info("get code: %s", sp[1])
		return sp[1], nil
	}
	log.Error("GetCode number: %s, error: %v", number.Number, err)
	return "", errors.New(status)
}
