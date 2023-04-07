package setting

import (
	"game_slots_vsn/pkg/consts"
	"github.com/spf13/viper"
	"time"
)

type AppSetting struct {
}

type ServeSetting struct {
	RunMode string `yaml:"runMode"`
	IP      string `yaml:"ip"`
	Port    int    `yaml:"port"`
}

type LoggerSetting struct {
	ConsoleEnable bool   `yaml:"consoleEnable"`
	ConsoleLevel  string `yaml:"consoleLevel"`

	FileEnable     bool   `yaml:"fileEnable"`
	FileName       string `yaml:"fileName"`
	FileLevel      string `yaml:"fileLevel"`
	FileMaxSize    int    `yaml:"maxSize"`
	FileMaxBackups int    `yaml:"maxBackups"`
	FileMaxAges    int    `yaml:"maxAges"`
	FileCompress   bool   `yaml:"compress"`
	FileJsonEnable bool   `yaml:"jsonEnable"`
}

type RedisSetting struct {
	DriverName  string        `yaml:"driverName"`
	MasterName  string        `yaml:"masterName,option,omitempty"`
	Password    string        `yaml:"password,option,omitempty"`
	Hosts       []string      `yaml:"hosts"`
	DB          int           `yaml:"db"`
	PoolSize    int           `yaml:"poolSize"`
	MaxIdle     int           `yaml:"maxIdle"`
	MaxActive   int           `yaml:"maxActive"`
	IdleTimeout time.Duration `yaml:"idleTimeout"`
}

var (
	ASetting   = &AppSetting{}
	LogSetting = &LoggerSetting{}
	SrvSetting = &ServeSetting{}
	RisSetting = &RedisSetting{}
)

func Setup() {
	m := make(map[string]interface{})
	m[consts.APP] = ASetting
	m[consts.SERVER] = SrvSetting
	m[consts.Logger] = LogSetting
	m[consts.REDIS] = RisSetting
	loadSetting(m)
}

func loadSetting(settings map[string]interface{}) {
	for k, v := range settings {
		viper.UnmarshalKey(k, v)
	}
}
