package captcha

import (
	"CentralizedControl/common/proxys"
	"net/http"
	"time"
)

type GoogleCaptchaConfig struct {
	CapServerIp  string `json:"cap_server_ip"`
	YesClientKey string `json:"yes_clientKey"`
	WebsiteURL   string `json:"websiteURL"`
	WebsiteKey   string `json:"websiteKey"`
	Type         string `json:"type"`
	RetryTimeout int    `json:"retry_timeout"`
	RetryDelay   int    `json:"retry_delay"`
}

type GoogleCaptcha interface {
	CreateTask() (string, error)
	GetTaskResult(taskId string) (string, error)
}

type CaptchaApi struct {
	ClientKey            string
	WebsiteURL           string
	WebsiteKey           string
	Type                 string
	client               *http.Client
	RetryTimeoutDuration time.Duration
	RetryDelayDuration   time.Duration
	Proxy                proxys.Proxy
}

func NewGoogleCaptcha(config *GoogleCaptchaConfig, usedType string, _proxy proxys.Proxy) GoogleCaptcha {
	composite := &Composite{
		used: usedType,
		yes: &YesCaptcha{
			CaptchaApi: CaptchaApi{
				ClientKey:            config.YesClientKey,
				WebsiteURL:           config.WebsiteURL,
				WebsiteKey:           config.WebsiteKey,
				Type:                 config.Type,
				client:               &http.Client{},
				RetryTimeoutDuration: time.Duration(config.RetryTimeout) * time.Second,
				RetryDelayDuration:   time.Duration(config.RetryDelay) * time.Second,
			},
		},
		cap: &CapMonster{
			CaptchaApi: CaptchaApi{
				ClientKey:            "",
				WebsiteURL:           config.WebsiteURL,
				WebsiteKey:           config.WebsiteKey,
				Type:                 config.Type,
				client:               &http.Client{},
				RetryTimeoutDuration: time.Duration(config.RetryTimeout) * time.Second,
				RetryDelayDuration:   time.Duration(config.RetryDelay) * time.Second,
				Proxy:                _proxy,
			},
			ServerIp: config.CapServerIp,
		},
	}
	return composite
}
