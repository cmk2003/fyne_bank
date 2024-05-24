package initialize

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"sql_bank/global"
)

// var rdb *redis.Client
//var ctx = context.Background()

func InitRedis() {
	addr := fmt.Sprintf("%s:%s", global.ServerConfig.RedisConfig.Host, global.ServerConfig.RedisConfig.Port)
	global.RDB = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: global.ServerConfig.RedisConfig.Password,
		DB:       global.ServerConfig.RedisConfig.DB,
	})
}
