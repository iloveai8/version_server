package move_vsn

import (
	"encoding/json"
	"fmt"
	"game_slots_vsn/internal/service/models"
	"game_slots_vsn/pkg/consts"
	"game_slots_vsn/pkg/redis"
	"github.com/spf13/viper"
	"testing"
)

var fRdb = &redis.RedDB{}
var tRdb = &redis.RedDB{}

func init() {
	init1()
	//init2()
}

func init1() {
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	viper.SetConfigName("devCluster")
	viper.AutomaticEnv() // read in environment variables that match
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
	fRdb.NewRDB()
}

func init2() {
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	viper.SetConfigName("devCluster")
	viper.AutomaticEnv() // read in environment variables that match
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
	tRdb.NewRDB()
}

func TestRemoveVsn(t *testing.T) {
	platType := "inner"
	delServerInfo(fRdb, platType, "ios244.X")
	m := getServerInfos(fRdb, platType)
	for maxVsn, v := range m {
		fmt.Println("cluster ====================>platType:", platType, " maxVsn:", maxVsn, " v:", v)
	}
}

func TestMoveTOCluster(t *testing.T) {
	doMoveGM(fRdb, tRdb)
	array := [3]string{"ios", "android", "inner"}
	for _, platType := range array {
		doMoveServers(fRdb, tRdb, platType)
	}
}

func TestClusterData(t *testing.T) {
	gmInfo, err := getGmInfo(tRdb)
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

func doMoveGM(fRdb, tRdb *redis.RedDB) {
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

func doMoveServers(fRdb, tRdb *redis.RedDB, platType string) {
	m := getServerInfos(fRdb, platType)
	for maxVsn, v := range m {
		fmt.Println("platType:", platType, " maxVsn:", maxVsn, " v:", v)
		addServerInfo(tRdb, platType, maxVsn, v)
	}
}

func getGmInfo(rdb *redis.RedDB) (*models.GmInfo, error) {
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

func updateGmInfo(rdb *redis.RedDB, g *models.GmInfo) error {
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

func getServerInfos(rdb *redis.RedDB, platType string) map[string]*models.ServerInfo {
	m := rdb.HGetAll(makePlatCacheKey(platType))
	serverInfoMap := make(map[string]*models.ServerInfo, len(m))
	for k, serverInfoStr := range m {
		serverInfo := &models.ServerInfo{}
		_ = json.Unmarshal([]byte(serverInfoStr), serverInfo)
		serverInfoMap[k] = serverInfo
	}
	return serverInfoMap
}

func addServerInfo(rdb *redis.RedDB, platType, maxVsn string, s *models.ServerInfo) bool {
	serverBytes, err := json.Marshal(s)
	if err != nil {
		return false
	}
	return rdb.HSet(makePlatCacheKey(platType), maxVsn, string(serverBytes))
}

func delServerInfo(rdb *redis.RedDB, platType, maxVsn string) bool {
	return rdb.HDel(makePlatCacheKey(platType), maxVsn)
}

func makePlatCacheKey(platType string) string {
	key := consts.CacheServerKey
	if len(platType) > 0 {
		key = fmt.Sprintf("%s%s.", key, platType)
	}
	return key
}
