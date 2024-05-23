package initialize

import (
	"fmt"
	"github.com/spf13/viper"
	"sql_bank/global"
)

func InitConfig() {
	configFilePrefix := "config"
	configFileName := fmt.Sprintf("%s.yaml", configFilePrefix)

	v := viper.New()
	v.SetConfigFile(configFileName)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := v.Unmarshal(&global.ServerConfig); err != nil {
		panic(err)
	}
}
