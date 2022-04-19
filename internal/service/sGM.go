package service

import (
	`encoding/json`
	`game_slots_vsn/internal/consts`
	`game_slots_vsn/internal/module`
	`game_slots_vsn/pkg/redis`
)

type sGM struct {
}

var (
	// sG is the instance of service User.
	sG = sGM{}
)

func GM() *sGM {
	return &sG
}

func (sG sGM) GetGM() *module.GM {
	gmStr, _ := redis.Client.Get(consts.CacheGMKey).Result()
	gmInfo := module.NewGM()
	_ = json.Unmarshal([]byte(gmStr), gmInfo)
	return gmInfo
}

func (sG sGM) AddGM(gm *module.GM) error {
	gmByte, err := json.Marshal(gm)
	if err != nil {
		return err
	}
	_, err = redis.Client.Set(consts.CacheGMKey, gmByte, 0).Result()
	return err
}
