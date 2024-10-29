package base

import (
	"CentralizedControl/common"
	"CentralizedControl/common/proxys"
	"net/mail"
	"sync"
	"time"
)

type RetryConfig struct {
	LastReqTime          time.Time
	RetryTimeoutDuration time.Duration
	RetryDelayDuration   time.Duration
}

type ServerConfig struct {
	Server string `json:"server"`
	Port   string `json:"port"`
}

type RespEmail struct {
	To     string
	Header mail.Header
	Body   string
}

const (
	StatusNoLogin = iota
	StatusLogin
	StatusPasswdError
	StatusBan
)

var (
	NeedLoginError  = common.NerError("need login")
	NeedPasswdError = common.NerError("Passwd Error")
	NeedBanError    = common.NerError("Ban Error")
)

type ExtractEmailCallback func(e EmailClient, resp *RespEmail) (string, error)

type EmailClient interface {
	Init() error
	Close()
	GetStatus() int
	GetType() string
	GetEmailAddr() string
	GetEmailPasswd() string
	GetEmail(from string) ([]RespEmail, error)
	WaitForEmail(self EmailClient, from string, callback ExtractEmailCallback) (string, error)
	SetProxy(p proxys.Proxy)
}

type EmailClientData struct {
	EmailAddr   string
	EmailPasswd string
	Status      int
	Type        string
	Proxy       proxys.Proxy
	Lock        sync.Mutex
}

func (this *EmailClientData) GetStatus() int {
	return this.Status
}

func (this *EmailClientData) GetType() string {
	return this.Type
}

func (this *EmailClientData) GetEmailAddr() string {
	return this.EmailAddr
}

func (this *EmailClientData) GetEmailPasswd() string {
	return this.EmailPasswd
}
