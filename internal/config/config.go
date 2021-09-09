package config

import (
	`fmt`
	`github.com/spf13/viper`
	`os`
)

var (
	C = &Config{}
)

type Config struct {
	Mysql  MySQLConfig
	Redis  RedisConfig
	Server ServerConfig
}

type ServerConfig struct {
	IP   string
	Port int
}

type MySQLConfig struct {
	IP       string
	Port     int
	User     string
	Password string
	Database string
}

type RedisConfig struct {
	Host     string
	DB       int
	PoolSize int
}

//init Config eg:mysql redis.....
func Init() {
	err := viper.Unmarshal(C)
	if err != nil {
		os.Exit(0)
	}
	fmt.Println("server:", C.Server)
	fmt.Println("redis:", C.Redis)
	fmt.Println("mysql:", C.Mysql)
}
