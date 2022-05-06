package redis

import (
	"errors"
	`fmt`
	"game_slots_vsn/pkg/logger"
	"github.com/go-redis/redis"
	`github.com/spf13/viper`
	`os`
	"time"
)

type RedisConf struct {
	DriverName string `yaml:"driveName"`
	Host       string `yaml:"addr"`

	MasterName string   `yaml:"masterName"`
	Hosts      []string `yaml:"hosts"`

	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
	PoolSize int    `yaml:"poolSize"`
}

const (
	DriverRedis    string = "redis"
	DriverSentinel string = "redisSentinel"
)

var (
	c      = &RedisConf{}
	Client *redis.Client
)

func Init() {
	if err := viper.UnmarshalKey("redis", c); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "config modify fail.", err)
		os.Exit(0)
	}
	switch c.DriverName {
	case DriverRedis:
		logger.Logger.Infof("redis match driver:%v.", DriverRedis)
		Client = redis.NewClient(redisOptions(c))
	case DriverSentinel:
		logger.Logger.Infof("redis match driver:%v.", DriverSentinel)
		Client = redis.NewFailoverClient(sentinelOptions(c))
	default:
		panic(errors.New("connection not available"))
	}

	if _, err := Client.Ping().Result(); err != nil {
		panic(errors.New("init redis error"))
	}
	logger.Logger.Infof("redis:%v start success.", c)
}

// redis option
func redisOptions(c *RedisConf) *redis.Options {
	return &redis.Options{
		Addr:               c.Host,
		DB:                 c.DB,
		PoolSize:           c.PoolSize,
		DialTimeout:        10 * time.Second,
		ReadTimeout:        30 * time.Second,
		WriteTimeout:       30 * time.Second,
		PoolTimeout:        30 * time.Second,
		IdleTimeout:        500 * time.Millisecond,
		IdleCheckFrequency: 500 * time.Millisecond,
	}
}

// redis sentinel option
func sentinelOptions(c *RedisConf) *redis.FailoverOptions {
	return &redis.FailoverOptions{
		MasterName:         c.MasterName,
		SentinelAddrs:      c.Hosts,
		Password:           c.Password,
		DB:                 c.DB,
		PoolSize:           c.PoolSize,
		DialTimeout:        10 * time.Second,
		ReadTimeout:        30 * time.Second,
		WriteTimeout:       30 * time.Second,
		PoolTimeout:        30 * time.Second,
		IdleTimeout:        500 * time.Millisecond,
		IdleCheckFrequency: 500 * time.Millisecond,
	}
}
