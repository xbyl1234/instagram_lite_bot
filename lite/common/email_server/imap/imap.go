package imap

import (
	"CentralizedControl/common/email_server/base"
	"CentralizedControl/common/email_server/config"
	"CentralizedControl/common/log"
	"CentralizedControl/common/proxys"
	"CentralizedControl/common/utils"
	"errors"
	"fmt"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"io/ioutil"
	"net/http"
	"net/mail"
	"strings"
	"sync"
	"time"
)

type ImapClientData struct {
	base.EmailClientData
	base.ServerConfig
	imapClient *client.Client
}

func (this *ImapClientData) SetProxy(p proxys.Proxy) {
	this.Proxy = p
}

func (this *ImapClientData) Init() error {
	var err error
	defer func() {
		if err != nil {
			log.Error("email %s login error:%v", this.EmailAddr, err)
		} else {
			log.Info("email %s login success", this.EmailAddr)
		}
	}()
	if this.imapClient != nil {
		if this.Status == base.StatusNoLogin {
			this.imapClient.Close()
		}
		if this.Status == base.StatusLogin {
			return nil
		}
	}
	this.imapClient, err = client.DialTLS(this.Server+":"+this.Port, nil)
	if err != nil {
		log.Error("imap_provider create imapClient error: %v", err)
		this.Status = base.StatusNoLogin
		return nil
	}
	err = this.imapClient.Login(this.EmailAddr, this.EmailPasswd)
	if err != nil {
		return this.CheckError(err)
	}
	this.Status = base.StatusLogin
	return nil
}

func (this *ImapClientData) Close() {
	this.imapClient.Close()
}

func (this *ImapClientData) GetEmail(from string) ([]base.RespEmail, error) {
	email, err := this.requireEmail(from, false)
	if err != nil {
		log.Error("imap_provider %s error: %v", this.EmailAddr, err)
		return nil, this.CheckError(err)
	}
	if email == nil {
		return nil, nil
	}
	return email, nil
}

func (this *ImapClientData) CheckError(err error) error {
	if err == nil {
		return nil
	}
	switch err.Error() {
	case client.ErrAlreadyLoggedIn.Error():
		this.Status = base.StatusLogin
		return nil
	case "LOGIN failed.":
		this.Status = base.StatusPasswdError
		return base.NeedPasswdError
	case "User is authenticated but not connected.":
		this.Status = base.StatusBan
		return base.NeedBanError
	case "connection closed":
		this.Status = base.StatusNoLogin
		return base.NeedLoginError
	default:
		if strings.Contains(err.Error(), "An established connection was aborted by the software in your host machine.") {
			this.Status = base.StatusNoLogin
			return base.NeedLoginError
		}
		return nil
	}
}

func (this *ImapClientData) requireEmail(from string, fetchSeen bool) ([]base.RespEmail, error) {
	_, err := this.imapClient.Select("INBOX", false)
	if err != nil {
		return nil, err
	}

	criteria := imap.NewSearchCriteria()
	if fetchSeen {
		criteria.WithFlags = []string{imap.SeenFlag}
	} else {
		criteria.WithFlags = []string{imap.RecentFlag}
	}
	//criteria.Text = []string{to}
	//criteria.WithoutFlags = []string{imap_provider.RecentFlag}
	//criteria.Header.Add("Delivered-To", to)
	//criteria.Header.Add("To", to)
	criteria.Header.Add("From", from)
	ids, err := this.imapClient.Search(criteria)
	if err != nil {
		return nil, err
	}
	if len(ids) > 0 {
		seqset := &imap.SeqSet{}
		seqset.AddNum(ids[len(ids)-1])
		//seqset.AddNum(ids[0])

		messages := make(chan *imap.Message, 2)
		section := &imap.BodySectionName{}
		err = this.imapClient.Fetch(seqset, []imap.FetchItem{section.FetchItem()}, messages)
		if err != nil {
			return nil, err
		}
		var emails []base.RespEmail
		for msg := range messages {
			r := msg.GetBody(section)
			m, err := mail.ReadMessage(r)
			if err != nil {
				return nil, err
			}
			body, err := ioutil.ReadAll(m.Body)
			if err != nil {
				return nil, err
			}

			emails = append(emails, base.RespEmail{
				To:     utils.GetMidString(m.Header.Get("To"), "<", ">"),
				Header: m.Header,
				Body:   string(body),
			})
		}
		return emails, nil
	}

	return nil, nil
}

func (this *ImapClientData) Test() error {
	this.Lock.Lock()
	defer this.Lock.Unlock()
	for true {
		switch this.Status {
		case base.StatusLogin:
			_, _ = this.GetEmail("")
			if this.Status == base.StatusLogin {
				return nil
			}
		case base.StatusNoLogin:
			this.Init()
		case base.StatusBan:
			return base.NeedBanError
		case base.StatusPasswdError:
			return base.NeedPasswdError
		}
	}
	return nil
}

func (this *ImapClientData) WaitForEmail(self base.EmailClient, from string, callback base.ExtractEmailCallback) (string, error) {
	this.Lock.Lock()
	defer this.Lock.Unlock()

	retry := 0
	start := time.Now()
	for time.Since(start) < config.DefaultRetryConfig.RetryTimeoutDuration {
		var emails []base.RespEmail
		var err error

		switch this.Status {
		case base.StatusLogin:
			emails, err = self.GetEmail(from)
			if emails != nil {
				for _, item := range emails {
					var code string
					code, err = callback(self, &item)
					if err == nil && code != "" {
						return code, nil
					}
				}
				return "", errors.New("no one accept")
			}
			retry = 0
		case base.StatusNoLogin:
			err = this.Init()
			if err != nil {
				log.Warn("email: %s login error: %v", this.EmailAddr, err)
				if retry > 3 {
					log.Warn("email: %s login retry more than 3 so exit!", this.EmailAddr)
					return "", errors.New(fmt.Sprintf("email: %s login retry more than 3 so exit!", this.EmailAddr))
				}
			}
			retry++
		case base.StatusBan:
			return "", base.NeedBanError
		case base.StatusPasswdError:
			return "", base.NeedPasswdError
		}
		log.Warn("wait for %s code...error: %v", this.EmailAddr, err)
		time.Sleep(config.DefaultRetryConfig.RetryDelayDuration)
	}
	return "", errors.New("require code timeout")
}

func CreateImapEmail(server base.ServerConfig, emailAddr, emailPasswd string) *ImapClientData {
	imap := &ImapClientData{
		EmailClientData: base.EmailClientData{
			EmailAddr:   emailAddr,
			EmailPasswd: emailPasswd,
			Status:      base.StatusNoLogin,
			Type:        "imap_provider",
			Lock:        sync.Mutex{},
		},
		ServerConfig: server,
		imapClient:   nil,
	}
	go func() {
		imap.Lock.Lock()
		defer imap.Lock.Unlock()
		err := imap.Init()
		if err != nil {
			log.Error("email %s login error: %v", imap.EmailPasswd, err)
		}
	}()
	return imap
}

type ImapProviderData struct {
	base.ServerConfig
	Client http.Client
	Proxy  proxys.Proxy
	_type  string
}

func (this *ImapProviderData) GetType() string {
	return this._type
}

var yx1024Url = "https://api.yx1024.cc/getAccountApi.aspx?uid=55436&type=69&token=17b521b1e6c66e4fca6c6a006bb82a56&count=1"
var youxiang555Url = "https://www.youxiang555.com/api/pub/email.html?type=33&token=f822649b9102bb58cfc3d0506e55a58c&count=1"

func CreateImapProvider(providerName string, p proxys.Proxy) base.Provider {
	switch providerName {
	case base.ProviderYX1024Name:
		return CreateEmailProviderYX1024(yx1024Url, p)
	case base.ProviderYouXiang555Name:
		return CreateEmailProviderYouXiang555(youxiang555Url, p)
	}
	return nil
}
