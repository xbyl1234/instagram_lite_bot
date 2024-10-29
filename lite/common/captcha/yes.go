package captcha

import (
	"CentralizedControl/common"
	"CentralizedControl/common/http_helper"
	"CentralizedControl/common/log"
	"time"
)

type YesCaptcha struct {
	CaptchaApi
}

type ReqYesCreateTask struct {
	ClientKey string `json:"clientKey"`
	Task      struct {
		WebsiteURL string `json:"websiteURL"`
		WebsiteKey string `json:"websiteKey"`
		Type       string `json:"type"`
	} `json:"task"`
}

type RespYesCreateTask struct {
	ErrorId          int    `json:"errorId"`
	ErrorCode        string `json:"errorCode"`
	ErrorDescription string `json:"errorDescription"`
	TaskId           string `json:"taskId"`
}

func (this *YesCaptcha) CreateTask() (string, error) {
	resp := &RespYesCreateTask{}
	err := http_helper.HttpDoJson(this.client, &http_helper.RequestOpt{
		IsPost: true,
		ReqUrl: "https://api.yescaptcha.com/createTask",
		JsonData: &ReqYesCreateTask{
			ClientKey: this.ClientKey,
			Task: struct {
				WebsiteURL string `json:"websiteURL"`
				WebsiteKey string `json:"websiteKey"`
				Type       string `json:"type"`
			}{
				WebsiteURL: this.WebsiteURL,
				WebsiteKey: this.WebsiteKey,
				Type:       this.Type,
			},
		},
	}, resp)
	if err != nil {
		return "", err
	}
	if resp.ErrorCode != "" {
		return "", &common.MakeMoneyError{ErrStr: resp.ErrorDescription}
	}
	return resp.TaskId, nil
}

type ReqYesTaskResult struct {
	ClientKey   string `json:"clientKey"`
	TaskId      string `json:"taskId"`
	CacheRecord string `json:"cacheRecord,omitempty"`
}

type RespYesTaskResult struct {
	ErrorId          int    `json:"errorId"`
	ErrorCode        string `json:"errorCode"`
	ErrorDescription string `json:"errorDescription"`
	Solution         struct {
		GRecaptchaResponse string `json:"gRecaptchaResponse"`
	} `json:"solution"`
	Status string `json:"status"`
}

func (this *YesCaptcha) GetTaskResult(taskId string) (string, error) {
	start := time.Now()
	for time.Since(start) < this.RetryTimeoutDuration {
		result, err := this.getTaskResult(taskId)
		if err != nil {
			log.Error("req google code error: %v", err)
		} else {
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

func (this *YesCaptcha) getTaskResult(taskId string) (*RespYesTaskResult, error) {
	resp := &RespYesTaskResult{}
	err := http_helper.HttpDoJson(this.client, &http_helper.RequestOpt{
		IsPost: true,
		ReqUrl: "https://api.yescaptcha.com/getTaskResult",
		JsonData: &ReqYesTaskResult{
			ClientKey: this.ClientKey,
			TaskId:    taskId,
			//CacheRecord: "",
		},
	}, resp)
	return resp, err
}
