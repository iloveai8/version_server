// Package redis 提供Redis客户端的封装
//
// 该包实现了：
// - Redis客户端接口定义
// - Redis命令的封装
// - 错误处理和重试逻辑
//
// 设计原则：
// - 接口抽象：定义清晰的Client接口
// - 可测试性：支持使用miniredis进行测试
// - 错误处理：将Redis错误转换为应用错误
package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"game_slots_vsn/pkg/errcode"
)

// Client Redis客户端接口
//
// 定义了Redis操作的核心方法，支持依赖注入和测试
type Client interface {
	// ==================== 连接操作 ====================

	// Ping 检查Redis连接
	//
	// 参数:
	//   ctx: 上下文
	//
	// 返回:
	//   error: 连接失败时返回错误
	Ping(ctx context.Context) error

	// ==================== String操作 ====================

	// Get 获取字符串值
	Get(ctx context.Context, key string) (string, error)

	// Set 设置字符串值
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error

	// Del 删除键
	Del(ctx context.Context, keys ...string) error

	// Exists 检查键是否存在
	Exists(ctx context.Context, key string) (bool, error)

	// ==================== Hash操作 ====================

	// HGet 获取哈希字段值
	HGet(ctx context.Context, key, field string) (string, error)

	// HSet 设置哈希字段值
	HSet(ctx context.Context, key, field string, value interface{}) error

	// HMGet 批量获取哈希字段值
	HMGet(ctx context.Context, key string, fields ...string) ([]interface{}, error)

	// HMSet 批量设置哈希字段值
	HMSet(ctx context.Context, key string, values map[string]interface{}) error

	// HGetAll 获取哈希所有字段值
	HGetAll(ctx context.Context, key string) (map[string]string, error)

	// HDel 删除哈希字段
	HDel(ctx context.Context, key string, fields ...string) error

	// ==================== Set操作 ====================

	// SAdd 向集合添加成员
	SAdd(ctx context.Context, key string, members ...interface{}) error

	// SRem 从集合移除成员
	SRem(ctx context.Context, key string, members ...interface{}) error

	// SIsMember 检查成员是否在集合中
	SIsMember(ctx context.Context, key string, member interface{}) (bool, error)

	// SMembers 获取集合所有成员
	SMembers(ctx context.Context, key string) ([]string, error)

	// SCard 获取集合成员数量
	SCard(ctx context.Context, key string) (int64, error)

	// ==================== JSON操作 ====================

	// JSONSet 设置JSON值
	JSONSet(ctx context.Context, key string, value interface{}) error

	// JSONGet 获取JSON值
	JSONGet(ctx context.Context, key string) (string, error)
}

// RedisClient Redis客户端实现
//
// 封装了go-redis/v9的客户端
type RedisClient struct {
	client *redis.Client
}

// NewRedisClient 创建Redis客户端
//
// 参数:
//   client: go-redis客户端实例
//
// 返回:
//   Client: Redis客户端接口
func NewRedisClient(client *redis.Client) Client {
	return &RedisClient{client: client}
}

// Ping 检查Redis连接
//
// 参数:
//   ctx: 上下文
//
// 返回:
//   error: 连接失败时返回错误
func (r *RedisClient) Ping(ctx context.Context) error {
	_, err := r.client.Ping(ctx).Result()
	if err != nil {
		return errcode.ErrRedisConnError.WithMessage("Redis连接失败")
	}
	return nil
}

// NewRedisClientFromURL 从URL创建Redis客户端
//
// 参数:
//   url: Redis连接URL
//
// 返回:
//   Client: Redis客户端接口
func NewRedisClientFromURL(url string) Client {
	client := redis.NewClient(&redis.Options{
		Addr: url,
	})
	return &RedisClient{client: client}
}

// ==================== String操作实现 ====================

// Get 获取字符串值
func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", errcode.ErrVersionNotFound.WithDetails(map[string]interface{}{
				"key": key,
			})
		}
		return "", errcode.ErrRedisError.WithMessage(fmt.Sprintf("Redis GET失败: %v", err))
	}
	return val, nil
}

// Set 设置字符串值
func (r *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	err := r.client.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return errcode.ErrRedisError.WithMessage(fmt.Sprintf("Redis SET失败: %v", err))
	}
	return nil
}

// Del 删除键
func (r *RedisClient) Del(ctx context.Context, keys ...string) error {
	_, err := r.client.Del(ctx, keys...).Result()
	if err != nil {
		return errcode.ErrRedisError.WithMessage(fmt.Sprintf("Redis DEL失败: %v", err))
	}
	return nil
}

// Exists 检查键是否存在
func (r *RedisClient) Exists(ctx context.Context, key string) (bool, error) {
	count, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, errcode.ErrRedisError.WithMessage(fmt.Sprintf("Redis EXISTS失败: %v", err))
	}
	return count > 0, nil
}

// ==================== Hash操作实现 ====================

// HGet 获取哈希字段值
func (r *RedisClient) HGet(ctx context.Context, key, field string) (string, error) {
	val, err := r.client.HGet(ctx, key, field).Result()
	if err != nil {
		if err == redis.Nil {
			return "", errcode.ErrVersionNotFound.WithDetails(map[string]interface{}{
				"key":   key,
				"field": field,
			})
		}
		return "", errcode.ErrRedisError.WithMessage(fmt.Sprintf("Redis HGET失败: %v", err))
	}
	return val, nil
}

// HSet 设置哈希字段值
func (r *RedisClient) HSet(ctx context.Context, key, field string, value interface{}) error {
	err := r.client.HSet(ctx, key, field, value).Err()
	if err != nil {
		return errcode.ErrRedisError.WithMessage(fmt.Sprintf("Redis HSET失败: %v", err))
	}
	return nil
}

// HMGet 批量获取哈希字段值
func (r *RedisClient) HMGet(ctx context.Context, key string, fields ...string) ([]interface{}, error) {
	vals, err := r.client.HMGet(ctx, key, fields...).Result()
	if err != nil {
		return nil, errcode.ErrRedisError.WithMessage(fmt.Sprintf("Redis HMGET失败: %v", err))
	}
	return vals, nil
}

// HMSet 批量设置哈希字段值
func (r *RedisClient) HMSet(ctx context.Context, key string, values map[string]interface{}) error {
	err := r.client.HMSet(ctx, key, values).Err()
	if err != nil {
		return errcode.ErrRedisError.WithMessage(fmt.Sprintf("Redis HMSET失败: %v", err))
	}
	return nil
}

// HGetAll 获取哈希所有字段值
func (r *RedisClient) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	vals, err := r.client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, errcode.ErrRedisError.WithMessage(fmt.Sprintf("Redis HGETALL失败: %v", err))
	}
	return vals, nil
}

// HDel 删除哈希字段
func (r *RedisClient) HDel(ctx context.Context, key string, fields ...string) error {
	_, err := r.client.HDel(ctx, key, fields...).Result()
	if err != nil {
		return errcode.ErrRedisError.WithMessage(fmt.Sprintf("Redis HDEL失败: %v", err))
	}
	return nil
}

// ==================== Set操作实现 ====================

// SAdd 向集合添加成员
func (r *RedisClient) SAdd(ctx context.Context, key string, members ...interface{}) error {
	err := r.client.SAdd(ctx, key, members...).Err()
	if err != nil {
		return errcode.ErrRedisError.WithMessage(fmt.Sprintf("Redis SADD失败: %v", err))
	}
	return nil
}

// SRem 从集合移除成员
func (r *RedisClient) SRem(ctx context.Context, key string, members ...interface{}) error {
	err := r.client.SRem(ctx, key, members...).Err()
	if err != nil {
		return errcode.ErrRedisError.WithMessage(fmt.Sprintf("Redis SREM失败: %v", err))
	}
	return nil
}

// SIsMember 检查成员是否在集合中
func (r *RedisClient) SIsMember(ctx context.Context, key string, member interface{}) (bool, error) {
	exists, err := r.client.SIsMember(ctx, key, member).Result()
	if err != nil {
		return false, errcode.ErrRedisError.WithMessage(fmt.Sprintf("Redis SISMEMBER失败: %v", err))
	}
	return exists, nil
}

// SMembers 获取集合所有成员
func (r *RedisClient) SMembers(ctx context.Context, key string) ([]string, error) {
	members, err := r.client.SMembers(ctx, key).Result()
	if err != nil {
		return nil, errcode.ErrRedisError.WithMessage(fmt.Sprintf("Redis SMEMBERS失败: %v", err))
	}
	return members, nil
}

// SCard 获取集合成员数量
func (r *RedisClient) SCard(ctx context.Context, key string) (int64, error) {
	count, err := r.client.SCard(ctx, key).Result()
	if err != nil {
		return 0, errcode.ErrRedisError.WithMessage(fmt.Sprintf("Redis SCARD失败: %v", err))
	}
	return count, nil
}

// ==================== JSON操作实现 ====================

// JSONSet 设置JSON值
func (r *RedisClient) JSONSet(ctx context.Context, key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return errcode.ErrJSONError.WithMessage(fmt.Sprintf("JSON序列化失败: %v", err))
	}

	err = r.client.Set(ctx, key, data, 0).Err()
	if err != nil {
		return errcode.ErrRedisError.WithMessage(fmt.Sprintf("Redis JSON SET失败: %v", err))
	}
	return nil
}

// JSONGet 获取JSON值
func (r *RedisClient) JSONGet(ctx context.Context, key string) (string, error) {
	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", errcode.ErrVersionNotFound.WithDetails(map[string]interface{}{
				"key": key,
			})
		}
		return "", errcode.ErrRedisError.WithMessage(fmt.Sprintf("Redis JSON GET失败: %v", err))
	}
	return data, nil
}
