package dao

import (
	"game_slots_vsn/pkg/redis"
	"game_slots_vsn/pkg/setting"
)

var (
	rdb *redis.Redis
)

func Setup(c *setting.RedisSetting) {
	r, err := redis.Setup(c)
	if err != nil {
		panic(err)
	}
	rdb = r
}
