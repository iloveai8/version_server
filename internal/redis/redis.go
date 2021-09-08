package redis
//
//import (
//	`errors`
//	`fmt`
//	`github.com/fsnotify/fsnotify`
//	"github.com/go-redis/redis"
//	`github.com/spf13/viper`
//	`time`
//)
//
//const emptyString = ""
//var (
//	config  RedisConfig
//	RedisDB *redis.Client
//)
//
////RedisConfig 配置信息
//type RedisConfig struct {
//	Name     string
//	hots     string
//	db       int
//	poolSize int
//}
//
//func Init(cfg string) {
//	//RedisDB = redis.NewClient(redisOptions(host, db, poolSize))
//	//pong, err := RedisDB.Ping().Result()
//	//if err != nil {
//	//	fmt.Println(pong, err)
//	//	panic(err)
//	//}
//	err := viper.Unmarshal(&C)
//	if err != nil {
//		t.Fatalf("unable to decode into struct, %v", err)
//	}
//
//	if err := _init(cfg); err != nil {
//		panic(fmt.Errorf("Fatal error config file: %s \n", err))
//	}
//}
//
//func _init(cfg string) error {
//	if cfg == emptyString {
//		return errors.New("config not find")
//	}
//	config = RedisConfig{
//		Name: cfg,
//	}
//	if err := config.initConfig(); err != nil {
//		return err
//	}
//	config.watchConfig()
//	return nil
//}
//
//func (c *RedisConfig) initConfig() error {
//	viper.SetConfigFile(c.Name)
//	viper.SetConfigType("yaml")
//	if err := viper.ReadInConfig(); err != nil {
//		return err
//	}
//	return nil
//}
//
//func (c *RedisConfig) watchConfig() {
//	viper.WatchConfig()
//	viper.OnConfigChange(func(e fsnotify.Event) {
//		fmt.Println("Config file changed:", e.Name)
//	})
//}
//
////Init 初始化配置信息
//func selfInit(cfg string) error {
//	if cfg == emptyString {
//		return errors.New("config not find")
//	}
//	config = RedisConfig{
//		Name: cfg,
//	}
//	if err := config.initConfig(); err != nil {
//		return err
//	}
//	config.watchConfig()
//	return nil
//}
//
//func redisOptions(redisAddr string, db, poolSize int) *redis.Options {
//	return &redis.Options{
//		Addr:               redisAddr,
//		DB:                 db,
//		DialTimeout:        10 * time.Second,
//		ReadTimeout:        30 * time.Second,
//		WriteTimeout:       30 * time.Second,
//		PoolSize:           poolSize,
//		PoolTimeout:        30 * time.Second,
//		IdleTimeout:        500 * time.Millisecond,
//		IdleCheckFrequency: 500 * time.Millisecond,
//	}
//}
