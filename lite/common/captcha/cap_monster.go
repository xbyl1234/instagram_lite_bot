package captcha

import (
	"CentralizedControl/common"
	"CentralizedControl/common/http_helper"
	"CentralizedControl/common/log"
	"errors"
	"fmt"
	"strconv"
	"time"
)

type CapMonster struct {
	CaptchaApi
	ServerIp string
}

type ReqCapCreateTask struct {
	ClientKey string `json:"clientKey"`
	Task      struct {
		Type          string `json:"type"`
		WebsiteURL    string `json:"websiteURL"`
		WebsiteKey    string `json:"websiteKey"`
		ProxyType     string `json:"proxyType"`
		ProxyAddress  string `json:"proxyAddress"`
		ProxyPort     int    `json:"proxyPort"`
		ProxyLogin    string `json:"proxyLogin"`
		ProxyPassword string `json:"proxyPassword"`
	} `json:"task"`
}

type RespCapCreateTask struct {
	ErrorId          int    `json:"errorId"`
	ErrorCode        string `json:"errorCode"`
	ErrorDescription string `json:"errorDescription"`
	TaskId           int    `json:"taskId"`
}

func (this *CapMonster) CreateTask() (string, error) {
	resp := &RespCapCreateTask{}

	req := &ReqCapCreateTask{
		ClientKey: this.ClientKey,
		Task: struct {
			Type          string `json:"type"`
			WebsiteURL    string `json:"websiteURL"`
			WebsiteKey    string `json:"websiteKey"`
			ProxyType     string `json:"proxyType"`
			ProxyAddress  string `json:"proxyAddress"`
			ProxyPort     int    `json:"proxyPort"`
			ProxyLogin    string `json:"proxyLogin"`
			ProxyPassword string `json:"proxyPassword"`
		}{
			Type:       "NoCaptchaTask",
			WebsiteURL: this.WebsiteURL,
			WebsiteKey: this.WebsiteKey,
		},
	}

	if this.Proxy != nil {
		proxyType := ""
		if this.Proxy.ProxyType == common.ProxyHttp {
			proxyType = "https"
		} else {
			proxyType = "socket5"
		}
		port, _ := strconv.Atoi(this.Proxy.Port)
		req.Task.ProxyType = proxyType
		req.Task.ProxyAddress = this.Proxy.Ip
		req.Task.ProxyPort = port
		req.Task.ProxyLogin = this.Proxy.Username
		req.Task.ProxyPassword = this.Proxy.Passwd
	}

	err := http_helper.HttpDoJson(this.client, &http_helper.RequestOpt{
		IsPost:   true,
		ReqUrl:   "http://" + this.ServerIp + "/createTask",
		JsonData: req,
	}, resp)
	if err != nil {
		return "", err
	}
	if resp.ErrorId != 0 {
		return "", errors.New(resp.ErrorCode + " " + resp.ErrorDescription)
	}
	return fmt.Sprintf("%d", resp.TaskId), nil
}

type ReqCapGetTaskResult struct {
	ClientKey string `json:"clientKey"`
	TaskId    int    `json:"taskId"`
}
type RespCapGetTaskResult struct {
	ErrorId          int    `json:"errorId"`
	ErrorCode        string `json:"errorCode"`
	ErrorDescription string `json:"errorDescription"`
	Status           string `json:"status"`
	Solution         struct {
		GRecaptchaResponse string `json:"gRecaptchaResponse"`
	} `json:"solution"`
}

func (this *CapMonster) GetTaskResult(taskId string) (string, error) {
	start := time.Now()
	for time.Since(start) < this.RetryTimeoutDuration {
		result, err := this.getTaskResult(taskId)
		if err != nil {
			log.Error("req google code error: %v", err)
		} else {
			if result.ErrorId != 0 {
				log.Error("req google code error: %s", result.ErrorCode+" "+result.ErrorDescription)
				return "", errors.New(result.ErrorCode + " " + result.ErrorDescription)
			}

			if result.Status == "ready" {
				return result.Solution.GRecaptchaResponse, nil
			}
			if result.Status != "processing" {
				log.Error("req google code error: %s", result.ErrorDescription)
			}
		}
		log.Warn("wait for google code...")
		time.Sleep(this.RetryDelayDuration)
	}
	return "", &common.MakeMoneyError{ErrStr: "require google code timeout", ErrType: common.RecvPhoneCodeError}
}

func (this *CapMonster) getTaskResult(taskId string) (*RespCapGetTaskResult, error) {
	resp := &RespCapGetTaskResult{}
	t, _ := strconv.Atoi(taskId)
	err := http_helper.HttpDoJson(this.client, &http_helper.RequestOpt{
		IsPost: true,
		ReqUrl: "http://" + this.ServerIp + "/getTaskResult",
		JsonData: &ReqCapGetTaskResult{
			ClientKey: this.ClientKey,
			TaskId:    t,
		},
	}, resp)
	return resp, err
}
