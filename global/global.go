package global

import (
	"gorm.io/gorm"
	"sql_bank/config"
)

var (
	DB           *gorm.DB
	ServerConfig config.ServerConfig
)
