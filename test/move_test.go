package test

import (
	"encoding/json"
	"fmt"
	"game_slots_vsn/internal/service/models"
	"game_slots_vsn/pkg/consts"
	"github.com/spf13/viper"
	"testing"
)

var fRdb = &rdb{}
var tRdb = &rdb{}

func init() {
	init1(fRdb)
	init2(tRdb)
}

func init1(fRdb *rdb) {
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	viper.SetConfigName("dev1")
	viper.AutomaticEnv() // read in environment variables that match
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
	rs := &setting{}
	err := viper.UnmarshalKey(consts.ConfigRedis, rs)
	if err != nil {
		panic(err)
	}
	fRdb.setup(rs)
}

func init2(tRdb *rdb) {
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	viper.SetConfigName("dev2")
	viper.AutomaticEnv() // read in environment variables that match
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
	rs := &setting{}
	err := viper.UnmarshalKey(consts.ConfigRedis, rs)
	if err != nil {
		panic(err)
	}
	tRdb.setup(rs)
}

func TestMoveTOCluster(t *testing.T) {
	doMoveGM(fRdb, tRdb)
	array := [3]string{"ios", "android", "inner"}
	for _, platType := range array {
		doMoveServers(fRdb, tRdb, platType)
	}
}

func doMoveGM(fRdb, tRdb *rdb) {
	gmInfo, err := getGmInfo(fRdb)
	if err != nil {
		return
	}
	err = updateGmInfo(tRdb, gmInfo)
	if err != nil {
		return
	}
	return
}

func doMoveServers(fRdb, tRdb *rdb, platType string) {
	m := getServerInfos(fRdb, platType)
	for maxVsn, v := range m {
		fmt.Println("platType:", platType, " maxVsn:", maxVsn, " v:", v)
		addServerInfo(tRdb, platType, maxVsn, v)
	}
}

func TestClusterData(t *testing.T) {
	gmInfo, err := getGmInfo(fRdb)
	if err != nil {
		return
	}
	fmt.Println("cluster ====================>gmInfo:", gmInfo)
	array := [3]string{"ios", "android", "inner"}
	for _, platType := range array {
		m := getServerInfos(tRdb, platType)
		for maxVsn, v := range m {
			fmt.Println("cluster ====================>platType:", platType, " maxVsn:", maxVsn, " v:", v)
		}
	}
}

func getGmInfo(rdb *rdb) (*models.GmInfo, error) {
	gmInfo := &models.GmInfo{}
	gmStr, err := rdb.Get(consts.CacheGMKey)
	if err != nil {
		return gmInfo, err
	}
	err = json.Unmarshal([]byte(gmStr), gmInfo)
	if err != nil {
		return gmInfo, err
	}
	return gmInfo, nil
}

func updateGmInfo(rdb *rdb, g *models.GmInfo) error {
	gmByte, err := json.Marshal(g)
	if err != nil {
		return err
	}
	err = rdb.Set(consts.CacheGMKey, string(gmByte), 0)
	if err != nil {
		return err
	}
	return nil
}

func getServerInfos(rdb *rdb, platType string) map[string]*models.ServerInfo {
	m := rdb.HGetAll(makePlatCacheKey(platType))
	serverInfoMap := make(map[string]*models.ServerInfo, len(m))
	for k, serverInfoStr := range m {
		serverInfo := &models.ServerInfo{}
		_ = json.Unmarshal([]byte(serverInfoStr), serverInfo)
		serverInfoMap[k] = serverInfo
	}
	return serverInfoMap
}

func addServerInfo(rdb *rdb, platType, maxVsn string, s *models.ServerInfo) bool {
	serverBytes, err := json.Marshal(s)
	if err != nil {
		return false
	}
	return rdb.HSet(makePlatCacheKey(platType), maxVsn, string(serverBytes))
}

func makePlatCacheKey(platType string) string {
	key := consts.CacheServerKey
	if len(platType) > 0 {
		key = fmt.Sprintf("%s%s.", key, platType)
	}
	return key
}
