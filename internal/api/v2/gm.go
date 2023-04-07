package v2

import (
	"encoding/json"
	"game_slots_vsn/internal/service"
	"game_slots_vsn/internal/service/models"
	"game_slots_vsn/pkg/e"
	"game_slots_vsn/pkg/logger"
	"github.com/gin-gonic/gin"
)

func GetGmInfo(c *gin.Context) {
	appG := e.Gin{C: c}

	gmService := &service.GMService{}
	gmInfo, err := gmService.GetGmInfo()
	if err != nil {
		logger.ErrorF("get gmInfo:%v error:%v", gmInfo, err)
		appG.Fail(e.GmInfoNotFound, gmInfo)
		return
	}

	appG.Success(e.SUCCESS, gmInfo)
	return
}

func UpdateGmInfo(c *gin.Context) {
	appG := e.Gin{C: c}

	buf := make([]byte, 1024)
	n, _ := c.Request.Body.Read(buf)

	gmInfo := &models.GmInfo{}
	err := json.Unmarshal(buf[:n], gmInfo)
	if err != nil {
		logger.ErrorF("Unmarshal gmInfo:%v err:%v", gmInfo, err)
		appG.Fail(e.ParamError, err)
		return
	}

	gmService := &service.GMService{
		GMInfo: gmInfo,
	}
	err = gmService.UpdateGmInfo()
	if err != nil {
		return
	}
	if err != nil {
		logger.ErrorF("Update gmInfo:%v err:%v", gmInfo, err)
		appG.Fail(e.ParamError, err)
		return
	}
	appG.Success(e.SUCCESS, nil)
	return
}
