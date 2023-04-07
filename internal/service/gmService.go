package service

import (
	"game_slots_vsn/internal/service/dao"
	"game_slots_vsn/internal/service/models"
)

type GMService struct {
	GMInfo *models.GmInfo
}

func (gs *GMService) GetGmInfo() (*models.GmInfo, error) {
	return dao.GetGmInfo()
}

func (gs *GMService) UpdateGmInfo() error {
	return dao.UpdateGmInfo(gs.GMInfo)
}
