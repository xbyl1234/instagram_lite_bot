package facebook

import "gorm.io/gorm"

type TableName struct {
	AccountTable string
}

func InitFacebook(accountDB *gorm.DB, tableName *TableName) {
	InitApiConfig("")
	InitDb(accountDB, tableName)
}
