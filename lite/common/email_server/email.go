package email_server

import (
	"CentralizedControl/common"
	"CentralizedControl/common/email_server/base"
	"CentralizedControl/common/email_server/imap"
	"CentralizedControl/common/email_server/temp"
	"CentralizedControl/common/fastjson"
	"CentralizedControl/common/log"
	"CentralizedControl/common/proxys"
	"errors"
	"fmt"
	"sync"
	"time"
)

const (
	pendingStateJustGet = iota
	pendingStateRequesting
	pendingStateFinish
	pendingStateError
)

const (
	ProjectFacebook      = "facebook"
	ProjectFacebookLite  = "facebook_lite"
	ProjectInstagram     = "instagram"
	ProjectInstagramLite = "instagram_lite"
	ProjectMessenger     = "messenger"
)

type PendingAccount struct {
	lock         sync.Mutex
	Email        base.EmailClient `json:"-"`
	EmailAccount string           `json:"email"`
	State        int              `json:"state"`
	Code         string           `json:"code"`
	Error        error            `json:"error"`
	Project      string           `json:"project"`
	Type         string           `json:"type"`
	StartTime    time.Time        `json:"startTime"`
}

type EmailServer struct {
	pendingEmail     map[string]*PendingAccount
	pendingLock      sync.Mutex
	curProvider      base.Provider
	extractEmailFunc map[string]FuncWaitCodeFunc
	//cache            *EmailCache
}

func CreateEmailServer(providerName string) *EmailServer {
	log.Info("email server start...")
	es := &EmailServer{
		pendingEmail: make(map[string]*PendingAccount),
		extractEmailFunc: map[string]FuncWaitCodeFunc{
			//ProjectFacebookLite:  WaitForCodeFacebook,
			//ProjectMessenger:     WaitForCodeMessenger,
			ProjectInstagramLite: WaitForInstagram,
		},
	}
	es.curProvider = CreateEmailProvider(providerName, nil)
	return es
}

func CreateEmailProvider(providerName string, p proxys.Proxy) base.Provider {
	provider := imap.CreateImapProvider(providerName, p)
	if provider != nil {
		return provider
	}
	provider = temp.CreateProviderWebTempEmail(providerName, p)
	if provider == nil {
		panic(errors.New(fmt.Sprintf("unknow email provider %s", providerName)))
	}
	return provider
}

//func (this *EmailServer) SetDB(DB *gorm.DB) {
//	this.cache = CreateDbCacheEmailProvider(DB)
//}

func (this *EmailServer) CreatePendingEmail(email base.EmailClient, project string) {
	username := email.GetEmailAddr()
	var pending = &PendingAccount{
		Email:        email,
		EmailAccount: email.GetEmailAddr(),
		Project:      project,
		Type:         this.curProvider.GetType(),
		State:        pendingStateJustGet,
		StartTime:    time.Now(),
	}
	this.pendingLock.Lock()
	defer this.pendingLock.Unlock()
	this.pendingEmail[username] = pending
}

func (this *EmailServer) GetServerStatus() (string, error) {
	this.pendingLock.Lock()
	defer this.pendingLock.Unlock()
	resp := fastjson.MustParse("[]")
	idx := 0
	for _, v := range this.pendingEmail {
		item := fastjson.MustParse("{}")
		item.Set("state", fastjson.AutoParse(v.State))
		item.Set("code", fastjson.AutoParse(v.Code))
		item.Set("error", fastjson.AutoParse(v.Error))
		item.Set("project", fastjson.AutoParse(v.Project))
		item.Set("type", fastjson.AutoParse(v.Type))
		item.Set("startTime", fastjson.AutoParse(v.StartTime))
		item.Set("email", fastjson.AutoParse(v.Email.GetEmailAddr()))
		item.Set("passwd", fastjson.AutoParse(v.Email.GetEmailAddr()))
		item.Set("login", fastjson.AutoParse(v.Email.GetStatus()))
		resp.SetArrayItem(idx, item)
		idx++
	}
	return resp.String(), nil
}

func (this *EmailServer) GetEmail(project string, useCache bool) (base.EmailClient, error) {
	var _email base.EmailClient
	var err error
	//if useCache && this.cache != nil {
	//	_email, err = this.cache.GetEmailByProject(project)
	//	if err != nil {
	//		log.Error("get email form cache error: %v", err)
	//	}
	//}

	//if _email == nil {
	//	_email, err = this.curProvider.GetEmail()
	//if err == nil && this.cache != nil {
	//	err = this.cache.SaveNewEmailDb(_email, project, this.curProvider.GetType())
	//	if err != nil {
	//		log.Error("save email error: %v", err)
	//	}
	//} else {
	//	log.Error("get email form imap_provider %s error: %v", this.curProvider.GetType(), err)
	//}
	//}

	_email, err = this.curProvider.GetEmail()
	if _email == nil {
		log.Error("get email failed!")
		return nil, err
	}

	this.CreatePendingEmail(_email, project)
	return _email, nil
}

func (this *EmailServer) doGetCode(pending *PendingAccount, project string) {
	pending.State = pendingStateRequesting
	log.Info("go get code coro: %s", pending.EmailAccount)
	if this.extractEmailFunc[project] != nil {
		pending.Code, pending.Error = this.extractEmailFunc[project](pending.Email)
		if pending.Error != nil {
			log.Error("email: %s, get code error: %v", pending.EmailAccount, pending.Error)
		} else {
			log.Info("email: %s, get code: %v", pending.EmailAccount, pending.Code)
		}
	} else {
		pending.Error = common.NerError("extractEmailFunc not exist")
	}

	pending.lock.Lock()
	if pending.Code != "" {
		pending.State = pendingStateFinish
	}
	if pending.Error != nil {
		pending.State = pendingStateError
	}
	pending.lock.Unlock()
}

func (this *EmailServer) SyncGetCode(email string, project string) (string, error) {
	this.pendingLock.Lock()
	var pending = this.pendingEmail[email]
	this.pendingLock.Unlock()
	if pending == nil {
		return "", common.NerError("not in pending")
	}

	switch pending.State {
	case pendingStateJustGet:
		pending.lock.Lock()
		pending.State = pendingStateRequesting
		pending.lock.Unlock()
		this.doGetCode(pending, project)
		if pending.Error != nil {
			log.Error("email: %s, project: %s, get code error: %v",
				email, project, pending.Error)
		}
		this.releasePendingEmail(pending, project)
		return pending.Code, pending.Error
	case pendingStateRequesting:
		return "", common.NerError("rein")
	case pendingStateFinish, pendingStateError:
		return pending.Code, pending.Error
	}
	return "", common.NerError("wtf?")
}

func (this *EmailServer) AsyncGetCode(email string, project string) (string, error) {
	this.pendingLock.Lock()
	defer this.pendingLock.Unlock()

	var pending = this.pendingEmail[email]
	if pending == nil {
		return "", common.NerError("not in pending")
	}

	pending.lock.Lock()
	defer pending.lock.Unlock()

	switch pending.State {
	case pendingStateJustGet:
		pending.State = pendingStateRequesting
		go this.doGetCode(pending, project)
		return "", common.NerError("commit")
	case pendingStateRequesting:
		return "", common.NerError("commit")
	case pendingStateFinish, pendingStateError:
		if pending.Error != nil {
			log.Error("email: %s, project: %s, get code error: %v",
				email, project, pending.Error)
		}
		this.releasePendingEmail(pending, project)
		return pending.Code, pending.Error
	}
	return "", common.NerError("wtf?")
}

func (this *EmailServer) ReleaseEmail(email string, project string, opt string) error {
	this.pendingLock.Lock()
	defer this.pendingLock.Unlock()
	this.releasePendingEmail(this.pendingEmail[email], project)
	//if this.cache != nil {
	//	err := this.cache.UpdateEmailProject(email, project, opt)
	//	if err != nil {
	//	}
	//}
	return nil
}

func (this *EmailServer) releasePendingEmail(pending *PendingAccount, project string) {
	if pending != nil {
		delete(this.pendingEmail, pending.EmailAccount)
		pending.Email.Close()
	}
}
