package facebook

import (
	"CentralizedControl/common"
	"CentralizedControl/common/log"
	"gorm.io/gorm"
)

type DbCtrlInterface interface {
	SaveCookies(ck *Cookies) error
	UpdateCookiesValue(ins *Facebook, key string, value interface{}) error
	LoadAccount()
}

type _DbCtrl struct {
	accountDB *gorm.DB
	tableName *TableName
}

func (this *_DbCtrl) SaveCookies(ck *Cookies) error {
	if this.accountDB == nil {
		log.Warn("no init data base")
		return common.NerError("no init data base")
	}
	result := this.accountDB.Table(this.tableName.AccountTable).Create(ck)
	if result.Error != nil {
		log.Error("update messenger db error: %v", result.Error)
	}
	return result.Error
}

func (this *_DbCtrl) UpdateCookiesValue(ins *Facebook, key string, value interface{}) error {
	if this.accountDB == nil {
		log.Warn("no init data base")
		return common.NerError("no init data base")
	}
	result := this.accountDB.Table(this.tableName.AccountTable).Where("email = ?", ins.ck.Id).
		Update(key, value)
	if result.Error != nil {
		log.Error("UpdateCookiesValue %d error: %v", ins.ck.Id, result.Error)
	}
	return result.Error
}

var DbCtrl *_DbCtrl

func InitDb(accountDB *gorm.DB, tableName *TableName) {
	DbCtrl = &_DbCtrl{
		accountDB: accountDB,
		tableName: tableName,
	}
}
