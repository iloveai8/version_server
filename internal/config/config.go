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
	Server ServerConfig
	Redis  RedisConfig
}

type ServerConfig struct {
	IP   string
	Port int
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
}
