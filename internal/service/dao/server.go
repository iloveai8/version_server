package dao

import (
	"encoding/json"
	"fmt"
	"game_slots_vsn/internal/service/models"
	"game_slots_vsn/pkg/consts"
	"game_slots_vsn/pkg/redis"
)

func GetServerInfos(platType string) map[string]*models.ServerInfo {
	m := redis.Rdb.HGetAll(makePlatCacheKey(platType))
	serverInfoMap := make(map[string]*models.ServerInfo, len(m))
	for k, serverInfoStr := range m {
		serverInfo := &models.ServerInfo{}
		_ = json.Unmarshal([]byte(serverInfoStr), serverInfo)
		serverInfoMap[k] = serverInfo
	}
	return serverInfoMap
}

func AddServerInfo(platType, maxVsn string, s *models.ServerInfo) bool {
	serverBytes, err := json.Marshal(s)
	if err != nil {
		return false
	}
	return redis.Rdb.HSet(makePlatCacheKey(platType), maxVsn, string(serverBytes))
}

func GetServerInfo(platType, maxVsn string) *models.ServerInfo {
	serverInfoStr := redis.Rdb.HGet(makePlatCacheKey(platType), maxVsn)
	serverInfo := &models.ServerInfo{}
	if serverInfoStr == "" {
		serverInfo.SubServerInfoMap = make(map[string]*models.SubServerInfo, 1)
	} else {
		_ = json.Unmarshal([]byte(serverInfoStr), serverInfo)
		fmt.Println(serverInfo)
	}
	return serverInfo
}

func DeleteServerInfo(playType, maxVsn string) bool {
	return redis.Rdb.HDel(makePlatCacheKey(playType), maxVsn)
}

func makePlatCacheKey(platType string) string {
	key := consts.CacheServerKey
	if len(platType) > 0 {
		key = fmt.Sprintf("%s%s.", key, platType)
	}
	return key
}
