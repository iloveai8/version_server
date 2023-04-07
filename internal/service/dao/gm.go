package dao

import (
	"encoding/json"
	"game_slots_vsn/internal/service/models"
	"game_slots_vsn/pkg/consts"
)

func GetGmInfo() (*models.GmInfo, error) {
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

func UpdateGmInfo(g *models.GmInfo) error {
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
