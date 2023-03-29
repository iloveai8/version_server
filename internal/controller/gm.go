package controller

import (
	"encoding/json"
	"game_slots_vsn/internal/api/rsp"
	"game_slots_vsn/internal/module"
	"game_slots_vsn/internal/service"
	"game_slots_vsn/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

type cGm struct {
}

func GM() (cGM *cGm) {
	return &cGm{}
}

func (cGM *cGm) AddGM(ctx *gin.Context) {
	buf := make([]byte, 1024)
	n, _ := ctx.Request.Body.Read(buf)
	gmInfo := module.NewGM()
	err := json.Unmarshal(buf[:n], gmInfo)
	logger.Logger.Infof("gmInfo:%v", gmInfo)
	//if err != nil || gmInfo.GMSrvUrl == "" || gmInfo.GMResUrl == "" {
	if err != nil {
		logger.Logger.Errorf("add1 gm info error:%v", gmInfo)
		ctx.JSON(http.StatusOK, rsp.Fail(rsp.ParamsError, map[string]string{}))
		return
	}

	err = service.GM().AddGM(gmInfo)
	if err != nil {
		logger.Logger.Errorf("add2 gm info error:%v", gmInfo)
		ctx.JSON(http.StatusOK, rsp.Fail(rsp.ParamsError, map[string]string{}))
		return
	}

	logger.Logger.Infof("add gm info:%v success", gmInfo)
	ctx.JSON(http.StatusOK, rsp.Success(gmInfo))
	return
}
