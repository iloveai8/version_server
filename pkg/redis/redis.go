package redis

import (
	"context"
	"errors"
	"fmt"
	"game_slots_vsn/pkg/setting"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

type Redis struct {
	ctx    context.Context
	client redis.UniversalClient
}

func Setup(c *setting.RedisSetting) (*Redis, error) {
	ctx := context.Background()
	client := redis.NewUniversalClient(
		&redis.UniversalOptions{
			MasterName:   c.MasterName,
			Password:     c.Password,
			DB:           c.DB,
			Addrs:        c.Hosts,
			PoolSize:     c.PoolSize,
			DialTimeout:  5 * time.Second,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
			PoolTimeout:  5 * time.Second,

			OnConnect: func(ctx context.Context, cn *redis.Conn) error {
				fmt.Println("redis ping")
				return cn.Ping(ctx).Err()
			},
		},
	)
	err := client.Ping(ctx).Err()
	if err != nil {
		return nil, err
	}
	return &Redis{
		ctx:    ctx,
		client: client,
	}, nil
}

func (r *Redis) HGet(key, field string) string {
	value, err := r.client.HGet(r.ctx, key, field).Result()
	if err != nil {
		return ""
	}
	return value
}

func (r *Redis) HSet(key, field, value string) bool {
	if _, err := r.client.HSet(r.ctx, key, field, value).Result(); err != nil {
		return false
	}
	return true
}

func (r *Redis) HMGet(key string, fields ...string) []interface{} {
	value, err := r.client.HMGet(r.ctx, key, fields...).Result()
	if err != nil {
		return nil
	}
	return value
}

func (r *Redis) HMSet(key string, data map[string]interface{}) bool {
	_, err := r.client.HMSet(r.ctx, key, data).Result()
	if err != nil {
		return false
	}
	return true
}

func (r *Redis) HDel(key string, fields ...string) bool {
	_, err := r.client.HDel(r.ctx, key, fields...).Result()
	if err != nil {
		return false
	}
	return true
}

func (r *Redis) HGetAll(key string) map[string]string {
	value, err := r.client.HGetAll(r.ctx, key).Result()
	if err != nil {
		return nil
	}
	return value
}

func (r *Redis) HIncrBy(key, field string, incr int64) int64 {
	value, err := r.client.HIncrBy(r.ctx, key, field, incr).Result()
	if err != nil {
		return 0
	}
	return value
}

func (r *Redis) SAdd(key string, values []string) bool {
	if _, err := r.client.SAdd(r.ctx, key, values).Result(); err != nil {
		return false
	}
	return true
}

func (r *Redis) SRem(key string, values []string) bool {
	if _, err := r.client.SRem(r.ctx, key, values).Result(); err != nil {
		return false
	}
	return true
}

func (r *Redis) SIsMember(key string, value string) bool {
	if _, err := r.client.SIsMember(r.ctx, key, value).Result(); err != nil {
		return false
	}
	return true
}

func (r *Redis) SCard(key string) (count int64, isSuccess bool) {
	var err error
	if count, err = r.client.SCard(r.ctx, key).Result(); err != nil {
		return 0, false
	}
	return count, true
}

func (r *Redis) SMembers(key string) []string {
	if values, err := r.client.SMembers(r.ctx, key).Result(); err == nil {
		return values
	}
	return nil
}

func (r *Redis) Del(key ...string) bool {
	if _, err := r.client.Del(r.ctx, key...).Result(); err != nil {
		return false
	}
	return true
}

func (r *Redis) IncrVersion(key string) (version int64, err error) {
	if version, err = r.client.Incr(r.ctx, key).Result(); err != nil {
		return
	}
	return
}

func (r *Redis) GetIncrVersion(key string) (version int64, err error) {
	var versionStr string
	if versionStr, err = r.client.Get(r.ctx, key).Result(); err != nil && !errors.Is(err, redis.Nil) {
		return
	}
	if versionStr == "" {
		return 0, nil
	}
	version, err = strconv.ParseInt(versionStr, 10, 64)
	return
}

func (r *Redis) Get(key string) (value string, err error) {
	if value, err = r.client.Get(r.ctx, key).Result(); err != nil && !errors.Is(err, redis.Nil) {
		return
	}
	return
}

func (r *Redis) Set(key, value string, expire time.Duration) error {
	if _, err := r.client.Set(r.ctx, key, value, expire).Result(); err != nil {
		return err
	}
	return nil
}

func (r *Redis) GetSet(key string) (values []string) {
	var index uint64
	var tempValues []string
	var err error
	for {
		tempValues, index, err = r.client.SScan(r.ctx, key, index, "", -1).Result()
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

func (r *Redis) Subscribe(key string, fun func(message string)) {
	subscribe := r.client.(*redis.Client).Subscribe(r.ctx, key)
	for {
		msg := <-subscribe.Channel()
		go fun(msg.Payload)
	}
}

func (r *Redis) Publish(key, msg string) bool {
	if _, err := r.client.Publish(r.ctx, key, msg).Result(); err != nil {
		return false
	}
	return true
}

func (r *Redis) RPop(key string) (value string, isSuccess bool) {
	var err error
	if value, err = r.client.RPop(r.ctx, key).Result(); err != nil {
		return value, false
	}
	return value, true
}

func (r *Redis) LPush(key string, values []interface{}) (count int64, isSuccess bool) {
	var err error
	if count, err = r.client.LPush(r.ctx, key, values...).Result(); err != nil {
		return count, false
	}
	return count, true
}

func (r *Redis) SetNX(key string, value interface{}, expiration time.Duration) (isSuccess bool) {
	_, err := r.client.SetNX(r.ctx, key, value, expiration).Result()
	if err != nil {
		return
	}
	return
}

func (r *Redis) Expire(key string, expiration time.Duration) (isSuccess bool) {
	_, err := r.client.Expire(r.ctx, key, expiration).Result()
	if err != nil {
		return
	}
	return
}

func (r *Redis) Exists(key string) bool {
	return r.client.Exists(r.ctx, key).Val() > 0
}

func (r *Redis) Pipeline() redis.Pipeliner {
	return r.client.Pipeline()
}

func (r *Redis) GetCtx() context.Context {
	return r.ctx
}
