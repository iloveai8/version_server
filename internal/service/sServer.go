package service

import (
	`encoding/json`
	`game_slots_vsn/internal/consts`
	`game_slots_vsn/internal/module`
	`game_slots_vsn/pkg/redis`
)

type sServer struct {
}

var (
	// sV is the instance of service User.
	sV = sServer{}
)

func Server() *sServer {
	return &sV
}

func (sV sServer) GetAllVersion() []*module.Server {
	serverStrList, _ := redis.Client.HVals(consts.CacheServerKey).Result()
	serverInfoList := make([]*module.Server, len(serverStrList), cap(serverStrList))
	for i, serverStr := range serverStrList {
		serverInfo := module.NewServer()
		_ = json.Unmarshal([]byte(serverStr), serverInfo)
		serverInfoList[i] = serverInfo
	}
	return serverInfoList
}

func (sV sServer) GetVersionByKey(Key string) (*module.Server, error) {
	serverStr, err1 := redis.Client.HGet(consts.CacheServerKey, Key).Result()
	if err1 != nil {
		return nil, err1
	}
	serverInfo := module.NewServer()
	_ = json.Unmarshal([]byte(serverStr), serverInfo)
	return serverInfo, nil
}

func (sV sServer) AddVersion(server *module.Server) (bool, error) {
	serverByte, err := json.Marshal(server)
	if err != nil {
		return false, err
	}
	return redis.Client.HSet(consts.CacheServerKey, server.Vsn, serverByte).Result()
}

func (sV sServer) DeleteVersion(key string) (int64, error) {
	return redis.Client.HDel(consts.CacheServerKey, key).Result()
}
