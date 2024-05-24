package config

type SqlConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     string `mapstructure:"port" json:"port"`
	Db       string `mapstructure:"db" json:"db"`
	User     string `mapstructure:"user" json:"user"`
	Password string `mapstructure:"password" json:"password"`
}
type redisConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     string `mapstructure:"port" json:"port"`
	DB       int    `mapstructure:"db" json:"db"`
	Password string `mapstructure:"password" json:"password"`
}

type ServerConfig struct {
	SqlConfig   SqlConfig   `mapstructure:"mysql" json:"mysql"`
	RedisConfig redisConfig `mapstructure:"redis" json:"redis"`
}
