package global

import (
	"context"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"sql_bank/config"
)

var (
	DB           *gorm.DB
	ServerConfig config.ServerConfig
	RDB          *redis.Client
	Ctx          = context.Background() // 定义全局的上下文
)
