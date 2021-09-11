package redis

import (
	`errors`
	`game_slots_vsn/internal/config`
	"github.com/go-redis/redis"
	`time`
)

var (
	_client      *redis.Client
	ErrInitRedis = errors.New("connection not available")
)

func Init(rConfig *config.Config) {
	o := redisOptions(rConfig.Redis.Host, rConfig.Redis.DB, rConfig.Redis.PoolSize)
	_client = redis.NewClient(o)
	if _, err := _client.Ping().Result(); err != nil {
		panic(ErrInitRedis)
	}
}

func redisOptions(redisAddr string, db, poolSize int) *redis.Options {
	return &redis.Options{
		Addr:               redisAddr,
		DB:                 db,
		DialTimeout:        10 * time.Second,
		ReadTimeout:        30 * time.Second,
		WriteTimeout:       30 * time.Second,
		PoolSize:           poolSize,
		PoolTimeout:        30 * time.Second,
		IdleTimeout:        500 * time.Millisecond,
		IdleCheckFrequency: 500 * time.Millisecond,
	}
}

func GetClient() *redis.Client {
	return _client
}
