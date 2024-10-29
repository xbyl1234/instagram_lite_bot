package temp

import (
	"CentralizedControl/common/email_server/base"
	"CentralizedControl/common/http_helper"
	"encoding/json"
	"fmt"
	"github.com/bytedance/sonic"
)

type TempMailIo struct {
	WebTempEmailData
	token string
}

type TempMailIoGetEmailResp struct {
	Id   string `json:"id"`
	From string `json:"from"`
	To   string `json:"to"`
	//Cc          interface{}   `json:"cc"`
	Subject  string `json:"subject"`
	BodyText string `json:"body_text"`
	BodyHtml string `json:"body_html"`
	//CreatedAt time.Time `json:"created_at"`
	//Attachments []interface{} `json:"attachments"`
}

func (this *TempMailIo) GetEmail(from string) ([]base.RespEmail, error) {
	do, err := http_helper.HttpDo(&this.Client, &http_helper.RequestOpt{
		Header: map[string]string{
			"pragma":        "no-cache",
			"cache-control": "no-cache",
			"referer":       "https://temp-mail.io/",
			"origin":        "https://temp-mail.io",
			"accept":        "*/*",
			"content-type":  "application/json;charset=UTF-8",
			"user-agent":    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36",
		},
		IsPost: false,
		ReqUrl: fmt.Sprintf("https://api.internal.temp-mail.io/api/v3/email/%s/messages", this.EmailAddr),
	})
	if err != nil {
		return nil, err
	}
	var resp []TempMailIoGetEmailResp
	err = json.Unmarshal([]byte(do), &resp)
	if err != nil {
		return nil, err
	}
	if len(resp) == 0 {
		return nil, nil
	}
	var result []base.RespEmail
	for _, item := range resp {
		result = append(result,
			base.RespEmail{
				To: item.To,
				Header: map[string][]string{
					"Subject": {item.Subject},
				},
				Body: item.BodyText,
			})
	}
	return result, nil
}

func (this *TempMailIo) Init() error {
	this.Type = base.ProviderTempMailIo
	do, err := http_helper.HttpDoRetry(&this.Client, &http_helper.RequestOpt{
		Header: map[string]string{
			"referer":      "https://temp-mail.io/",
			"origin":       "https://temp-mail.io",
			"accept":       "*/*",
			"content-type": "application/json;charset=UTF-8",
			"user-agent":   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36",
		},
		IsPost: true,
		ReqUrl: "https://api.internal.temp-mail.io/api/v3/email/new",
		JsonData: struct {
			MinNameLength int `json:"min_name_length"`
			MaxNameLength int `json:"max_name_length"`
		}{
			MinNameLength: 10,
			MaxNameLength: 10,
		},
	})
	if err != nil {
		return err
	}
	resp, err := sonic.GetFromString(do)
	if err != nil {
		return err
	}
	this.EmailAddr = resp.GetString("email")
	this.token = resp.GetString("token")
	return nil
}
