package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	`github.com/spf13/viper`
	`os`
)

var (
	G = &GlobalConfig{}
)

type GlobalConfig struct {
	Web   *WebConfig   `yaml:"web"`
	Log   *LogConfig   `yaml:"log"`
	Redis *RedisConfig `yaml:"redis"`
}

type WebConfig struct {
	IP   string `yaml:"ip,omitempty"`
	Port int    `yaml:"port"`
}

type LogConfig struct {
	FileEnable bool   `yaml:"fileEnable"`
	FileName   string `yaml:"fileName"`
	FileLevel  string `yaml:"-,fileLevel"`

	ConsoleEnable bool   `yaml:"consoleEnable"`
	ConsoleLevel  string `yaml:"-,consoleLevel"`

	MaxSize    int  `yaml:"maxSize"`
	MaxBackups int  `yaml:"maxBackups"`
	MaxAges    int  `yaml:"maxAges"`
	Compress   bool `yaml:"compress"`
	JsonEnable bool `yaml:"jsonEnable"`
}

type RedisConfig struct {
	DriverName string `yaml:"driveName"`
	Host       string `yaml:"addr"`

	MasterName string   `yaml:"masterName"`
	Hosts      []string `yaml:"hosts"`

	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
	PoolSize int    `yaml:"poolSize"`
}

//init AppConfig eg:mysql redis.....
func InitConfig() {
	if err := viper.ReadInConfig(); err != nil {
		cobra.CheckErr(err)
	}
	if err := viper.Unmarshal(G); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Using Config File:", viper.ConfigFileUsed())
		_, _ = fmt.Fprintln(os.Stderr, "config modify fail.", err)
		os.Exit(0)
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		if err := viper.Unmarshal(G); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "config modify fail.", err)
			os.Exit(0)
		}
	})
}
