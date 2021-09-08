package config

var config *Config

type Config struct {
	MySQL MySQLConfig
	REDIS RedisConfig
}

type MySQLConfig struct {
	IP       string
	Port     int
	User     string
	Password string
	Database string
}

type RedisConfig struct {
	HOST     string
	DB       int
	POOLSIZE int
}

//init config eg:mysql redis.....
func Init(cfg string) {
	//config := &Config{}
	//err := config.initConfig(cfg)
	//if err != nil {
	//	panic(fmt.Errorf("Fatal error config file: %s \n", err))
	//	return
	//}
}

func (c Config) initConfig(cfg string) error {
	//viper.SetConfigFile(c.Name)
	//viper.SetConfigType("yaml")
	//config.watchConfig()
	//if err := viper.ReadInConfig(); err != nil {
	//	return err
	//}
	//return nil
	return nil
}

//
//func (c *Config) watchConfig() {
//	viper.WatchConfig()
//	viper.OnConfigChange(func(e fsnotify.Event) {
//		if c.LevelKey == emptyString {
//			return
//		}
//		logger.ChangeLevel(viper.GetString(c.LevelKey))
//	})
//}
