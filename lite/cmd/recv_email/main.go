package main

import (
	"CentralizedControl/common/email_server"
	"CentralizedControl/common/email_server/base"
	"CentralizedControl/common/email_server/config"
	"CentralizedControl/common/log"
	"CentralizedControl/common/proxys"
	"time"
)

func main() {
	log.DisAbleDebugLog()
	log.InitDefaultLog("", true, true)
	config.DefaultRetryConfig = &base.RetryConfig{
		RetryTimeoutDuration: 10 * time.Minute,
		RetryDelayDuration:   1 * time.Second,
	}

	testProxy, err := proxys.CreateSocks5Proxy("socks5://zP5Y4JL1:8ZTQmMIH@45.204.208.90:9001")
	email := email_server.CreateEmailProvider(base.ProviderTempMailIo, testProxy)
	if email == nil {
		return
	}
	getEmail, err := email.GetEmail()
	if err != nil {
		return
	}
	log.Info("email: %s", getEmail.GetEmailAddr())
	instagram, err := email_server.WaitForInstagram(getEmail)
	if err != nil {
		return
	}
	log.Info("%s", instagram)
}
