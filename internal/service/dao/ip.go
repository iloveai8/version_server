package dao

import (
	"game_slots_vsn/pkg/consts"
	"game_slots_vsn/pkg/redis"
)

func GetGmIPList() []string {
	return redis.Rdb.SMembers(consts.CacheGMIPKey)
}

func AddGmIP(ips []string) bool {
	return redis.Rdb.SAdd(consts.CacheGMIPKey, ips)
}

func RemGmIP(ips []string) bool {
	return redis.Rdb.SRem(consts.CacheGMIPKey, ips)
}

func IsGmIP(ip string) bool {
	return redis.Rdb.SIsMember(consts.CacheGMIPKey, ip)
}
