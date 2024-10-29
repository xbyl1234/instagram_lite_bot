package main

import (
	"CentralizedControl/common/email_server"
	"CentralizedControl/common/log"
	"CentralizedControl/common/proxys"
	"CentralizedControl/ctrl"
)

func main() {
	log.InitDefaultLog("ctrl", true, true)
	log.Info("ctrl server start...")
	ctrl.InitMysql()
	http := ctrl.RunHttpServer(
		email_server.CreateEmailServer(""),
		proxys.CreateProxysManage(),
	)
	http.Run()
}
