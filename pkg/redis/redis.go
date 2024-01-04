package redis

import (
	"context"
	"errors"
	"fmt"
	"game_slots_vsn/pkg/consts"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"strconv"
	"time"
)

type setting struct {
	MasterName string   `yaml:"masterName,option,omitempty"`
	Password   string   `yaml:"password,option,omitempty"`
	Hosts      []string `yaml:"hosts"`
	DB         int      `yaml:"db"`
	PoolSize   int      `yaml:"poolSize"`
}

type RedDB struct {
	ctx    context.Context
	client redis.UniversalClient
	S      *setting
}

var Rdb = &RedDB{}

func SetUp() {
	rs := &setting{}
	err := viper.UnmarshalKey(consts.ConfigRedis, rs)
	if err != nil {
		_ = fmt.Errorf("error redis: %v", err)
		return
	}
	Rdb.setup(rs)
}

func (rdb *RedDB) NewRDB() {
	rs := &setting{}
	err := viper.UnmarshalKey(consts.ConfigRedis, rs)
	if err != nil {
		panic(err)
	}
	rdb.setup(rs)
}

func (rdb *RedDB) setup(s *setting) {
	fmt.Printf("redis setting:%v\n", *s)
	ctx := context.Background()
	client := redis.NewUniversalClient(
		&redis.UniversalOptions{
			MasterName:   s.MasterName,
			Password:     s.Password,
			DB:           s.DB,
			Addrs:        s.Hosts,
			PoolSize:     s.PoolSize,
			DialTimeout:  5 * time.Second,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
			PoolTimeout:  5 * time.Second,
			OnConnect: func(ctx context.Context, cn *redis.Conn) error {
				return cn.Ping(ctx).Err()
			},
		},
	)
	err := client.Ping(ctx).Err()
	if err != nil {
		panic(err)
	}
	rdb.S = s
	rdb.ctx = ctx
	rdb.client = client
}

func (rdb *RedDB) HGet(key, field string) string {
	value, err := rdb.client.HGet(rdb.ctx, key, field).Result()
	if err != nil {
		return ""
	}
	return value
}

func (rdb *RedDB) HSet(key, field, value string) bool {
	if _, err := rdb.client.HSet(rdb.ctx, key, field, value).Result(); err != nil {
		return false
	}
	return true
}

func (rdb *RedDB) HMGet(key string, fields ...string) []interface{} {
	value, err := rdb.client.HMGet(rdb.ctx, key, fields...).Result()
	if err != nil {
		return nil
	}
	return value
}

func (rdb *RedDB) HMSet(key string, data map[string]interface{}) bool {
	_, err := rdb.client.HMSet(rdb.ctx, key, data).Result()
	if err != nil {
		return false
	}
	return true
}

func (rdb *RedDB) HDel(key string, fields ...string) bool {
	_, err := rdb.client.HDel(rdb.ctx, key, fields...).Result()
	if err != nil {
		return false
	}
	return true
}

func (rdb *RedDB) HGetAll(key string) map[string]string {
	value, err := rdb.client.HGetAll(rdb.ctx, key).Result()
	if err != nil {
		return nil
	}
	return value
}

func (rdb *RedDB) HIncrBy(key, field string, incr int64) int64 {
	value, err := rdb.client.HIncrBy(rdb.ctx, key, field, incr).Result()
	if err != nil {
		return 0
	}
	return value
}

func (rdb *RedDB) SAdd(key string, values []string) bool {
	if _, err := rdb.client.SAdd(rdb.ctx, key, values).Result(); err != nil {
		return false
	}
	return true
}

func (rdb *RedDB) SRem(key string, values []string) bool {
	if _, err := rdb.client.SRem(rdb.ctx, key, values).Result(); err != nil {
		return false
	}
	return true
}

func (rdb *RedDB) SIsMember(key string, value string) bool {
	//if _, err := rdb.client.SIsMember(rdb.ctx, key, value).Result(); err != nil {
	//	return false
	//}
	//return true
	b, err := rdb.client.SIsMember(rdb.ctx, key, value).Result()
	if err != nil {
		return false
	}
	return b
}

func (rdb *RedDB) SCard(key string) (count int64, isSuccess bool) {
	var err error
	if count, err = rdb.client.SCard(rdb.ctx, key).Result(); err != nil {
		return 0, false
	}
	return count, true
}

func (rdb *RedDB) SMembers(key string) []string {
	if values, err := rdb.client.SMembers(rdb.ctx, key).Result(); err == nil {
		return values
	}
	return nil
}

func (rdb *RedDB) Del(key ...string) bool {
	if _, err := rdb.client.Del(rdb.ctx, key...).Result(); err != nil {
		return false
	}
	return true
}

func (rdb *RedDB) IncrVersion(key string) (version int64, err error) {
	if version, err = rdb.client.Incr(rdb.ctx, key).Result(); err != nil {
		return
	}
	return
}

func (rdb *RedDB) GetIncrVersion(key string) (version int64, err error) {
	var versionStr string
	if versionStr, err = rdb.client.Get(rdb.ctx, key).Result(); err != nil && !errors.Is(err, redis.Nil) {
		return
	}
	if versionStr == "" {
		return 0, nil
	}
	version, err = strconv.ParseInt(versionStr, 10, 64)
	return
}

func (rdb *RedDB) Get(key string) (value string, err error) {
	if value, err = rdb.client.Get(rdb.ctx, key).Result(); err != nil && !errors.Is(err, redis.Nil) {
		return
	}
	return
}

func (rdb *RedDB) Set(key, value string, expire time.Duration) error {
	if _, err := rdb.client.Set(rdb.ctx, key, value, expire).Result(); err != nil {
		return err
	}
	return nil
}

func (rdb *RedDB) GetSet(key string) (values []string) {
	var index uint64
	var tempValues []string
	var err error
	for {
		tempValues, index, err = rdb.client.SScan(rdb.ctx, key, index, "", -1).Result()
		if err != nil {
			break
		}
		values = append(values, tempValues...)
		if index == 0 {
			break
		}
	}
	return
}

func (rdb *RedDB) Subscribe(key string, fun func(message string)) {
	subscribe := rdb.client.(*redis.Client).Subscribe(rdb.ctx, key)
	for {
		msg := <-subscribe.Channel()
		go fun(msg.Payload)
	}
}

func (rdb *RedDB) Publish(key, msg string) bool {
	if _, err := rdb.client.Publish(rdb.ctx, key, msg).Result(); err != nil {
		return false
	}
	return true
}

func (rdb *RedDB) RPop(key string) (value string, isSuccess bool) {
	var err error
	if value, err = rdb.client.RPop(rdb.ctx, key).Result(); err != nil {
		return value, false
	}
	return value, true
}

func (rdb *RedDB) LPush(key string, values []interface{}) (count int64, isSuccess bool) {
	var err error
	if count, err = rdb.client.LPush(rdb.ctx, key, values...).Result(); err != nil {
		return count, false
	}
	return count, true
}

func (rdb *RedDB) SetNX(key string, value interface{}, expiration time.Duration) (isSuccess bool) {
	_, err := rdb.client.SetNX(rdb.ctx, key, value, expiration).Result()
	if err != nil {
		return
	}
	return
}

func (rdb *RedDB) Expire(key string, expiration time.Duration) (isSuccess bool) {
	_, err := rdb.client.Expire(rdb.ctx, key, expiration).Result()
	if err != nil {
		return
	}
	return
}

func (rdb *RedDB) Exists(key string) bool {
	return rdb.client.Exists(rdb.ctx, key).Val() > 0
}

func (rdb *RedDB) Pipeline() redis.Pipeliner {
	return rdb.client.Pipeline()
}

func (rdb *RedDB) GetCtx() context.Context {
	return rdb.ctx
}
