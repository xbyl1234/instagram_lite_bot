package temp

import (
	"CentralizedControl/common/email_server/base"
	"CentralizedControl/common/http_helper"
	"CentralizedControl/common/utils"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

type LingshiMail struct {
	WebTempEmailData
	PHPSESSID             string
	logglytrackingsession string
	tmail                 string
}

func (this *LingshiMail) Init() error {
	this.Type = base.ProviderLinShiName
	http_helper.NeedJar()(&this.Client)
	this.EmailAddr = utils.GenString(utils.CharSet_abc, 10) + "@youxiang.dev"
	this.PHPSESSID = utils.GenString(utils.CharSet_All, 26)
	this.logglytrackingsession = utils.GenUUID()
	this.tmail = utils.Escape(this.EmailAddr, utils.EscapeEncodePathSegment)

	_url, _ := url.Parse("https://linshiyou.com/")
	ck := make([]*http.Cookie, 3)
	ck[0] = &http.Cookie{
		Name:  "PHPSESSID",
		Value: this.PHPSESSID,
	}
	ck[1] = &http.Cookie{
		Name:  "logglytrackingsession",
		Value: this.logglytrackingsession,
	}
	ck[2] = &http.Cookie{
		Name:  "tmail-emails",
		Value: utils.EncodeQueryPath(fmt.Sprintf("a:1:{i:0;s:20:\"%s\";}", this.EmailAddr)),
	}
	this.Client.Jar.SetCookies(_url, ck)

	do, err := http_helper.HttpDo(&this.Client, &http_helper.RequestOpt{
		Params: map[string]string{
			"user": this.EmailAddr,
		},
		Header: map[string]string{
			"referer":          "https://linshiyou.com/",
			"x-requested-with": "XMLHttpRequest",
			"accept":           "*/*",
		},
		IsPost: false,
		ReqUrl: "https://linshiyou.com/user.php",
	})
	if err != nil {
		return err
	}
	if do != this.EmailAddr {
		return errors.New("not match email username")
	}
	this.Status = base.StatusLogin
	return nil
}

func (this *LingshiMail) GetEmail(from string) ([]base.RespEmail, error) {
	do, err := http_helper.HttpDo(&this.Client, &http_helper.RequestOpt{
		Params: map[string]string{
			"unseen": "1",
		},
		Header: map[string]string{
			"referer":          "https://linshiyou.com/",
			"x-requested-with": "XMLHttpRequest",
			"accept":           "*/*",
		},
		IsPost: false,
		ReqUrl: "https://linshiyou.com/mail.php",
	})
	if err != nil {
		return nil, err
	}
	if len(do) == 0 {
		return nil, nil
	}
	if do == "DIE" {
		return nil, errors.New("email die")
	}
	return []base.RespEmail{{
		Header: nil,
		Body:   do,
	}}, nil
}
