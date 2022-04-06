package redis

import (
	"errors"
	"game_slots_vsn/pkg/config"
	"game_slots_vsn/pkg/logger"
	"github.com/go-redis/redis"
	"time"
)

const (
	DriverRedis    string = "redis"
	DriverSentinel string = "redisSentinel"
)

var Client *redis.Client

// init redis
func InitRedis(c *config.RedisConfig) {
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
func redisOptions(c *config.RedisConfig) *redis.Options {
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
func sentinelOptions(c *config.RedisConfig) *redis.FailoverOptions {
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
