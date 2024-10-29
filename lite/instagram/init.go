package instagram

import "gorm.io/gorm"

type TableName struct {
	AccountTable string
}

func InitInstagram(accountDB *gorm.DB, tableName *TableName) {
	InitApiConfig("")
	InitDb(accountDB, tableName)
}
