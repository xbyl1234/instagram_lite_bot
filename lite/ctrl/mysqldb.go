package ctrl

import (
	"CentralizedControl/common/log"
	"CentralizedControl/facebook"
	"CentralizedControl/instagram"
	"CentralizedControl/messenger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var accountDB *gorm.DB

func InitMysql() {
	var err error
	accountDB, err = gorm.Open(mysql.Open("root:xbylAbc123++..@tcp(127.0.0.1:3306)/account"),
		&gorm.Config{})
	if err != nil {
		panic(err)
	}
	log.Info("db connected...")
	err = accountDB.Table(MessengerAccountTable).AutoMigrate(&messenger.Cookies{})
	if err != nil {
		panic(err)
	}
	err = accountDB.Table(InstagramAccountTable).AutoMigrate(&instagram.Cookies{})
	if err != nil {
		panic(err)
	}
	err = accountDB.Table(FacebookAccountTable).AutoMigrate(&facebook.Cookies{})
	if err != nil {
		panic(err)
	}
}
