package service

import (
	`encoding/json`
	`fmt`
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

func (sV sServer) GetAllServer(platType string) map[string]*module.Server {
	m, _ := redis.Client.HGetAll(makeServerCacheKey(platType)).Result()
	serverMap := make(map[string]*module.Server, len(m))
	for key, serverStr := range m {
		serverInfo := module.NewServer()
		_ = json.Unmarshal([]byte(serverStr), serverInfo)
		serverMap[key] = serverInfo
	}
	return serverMap
}

func (sV sServer) GetServerByKey(platType, maxVsn string) (*module.Server, error) {
	serverStr, err1 := redis.Client.HGet(makeServerCacheKey(platType), maxVsn).Result()
	if err1 != nil {
		serverInfo := module.NewServer()
		return serverInfo, nil
	}
	serverInfo := module.NewServer()
	_ = json.Unmarshal([]byte(serverStr), serverInfo)
	return serverInfo, nil
}

func (sV sServer) AddServer(platType, maxVsn string, server *module.Server) (bool, error) {
	serverByte, err := json.Marshal(server)
	if err != nil {
		return false, err
	}
	return redis.Client.HSet(makeServerCacheKey(platType), maxVsn, serverByte).Result()
}

func (sV sServer) DeleteServer(platType, maxVsn string) (int64, error) {
	return redis.Client.HDel(makeServerCacheKey(platType), maxVsn).Result()
}

func makeServerCacheKey(platType string) string {
	key := consts.CacheServerKey
	if len(platType) > 0 {
		key = fmt.Sprintf("%s%s.", key, platType)
	}
	return key
}
