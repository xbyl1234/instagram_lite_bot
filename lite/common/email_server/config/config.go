package config

import (
	"CentralizedControl/common/email_server/base"
	"time"
)

var DefaultRetryConfig = &base.RetryConfig{
	RetryTimeoutDuration: 3 * time.Minute,
	RetryDelayDuration:   2 * time.Second,
}

var DefaultOutlookConfig = &base.ServerConfig{
	Server: "outlook.office365.com",
	Port:   "993",
}
