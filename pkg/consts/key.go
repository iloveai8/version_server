package consts

const (
	CacheServerKey string = "v2.server."
	CacheGMKey     string = "v2.gm."
	CacheGMIPKey   string = "v2.gmip."

	ServerTypeDefault int = 0
	ServerTypeDEV     int = 1
	ServerTypePRE     int = 2
	ServerTypePRO     int = 3
)

// 配置文件
const (
	ConfigApp    = "app"
	ConfigLogger = "log"
	ConfigServer = "server"
	ConfigRedis  = "redis"
	ConfigMongo  = "mongo"
	ConfigGgeoip = "ggeoip"
)

var IPList = []string{
	"1.202.246.19",
	"40.83.97.197",
	"47.75.45.195",
	"47.75.59.239",
	"49.51.197.144",
	"103.85.165.146",
	"119.81.164.4",
	"129.226.60.247",
	"43.154.154.31",
	"129.226.189.243",
	"94.74.105.98",
	"124.156.132.128",
	"159.138.38.254",
	"169.56.143.199",
	"159.138.154.188",
	"43.129.242.42",
	"103.85.165.146",
	"106.120.91.66",
	"1.202.246.19",
	"47.75.45.195",
	"124.156.183.70",
	"43.129.73.11",
	"43.154.74.107",
	"117.61.200.51",
	"43.154.156.163",
	"117.61.26.1",
	"138.2.71.246",
	"43.134.233.31",
}
