package controller

import (
	`encoding/json`
	`game_slots_vsn/internal/api/rsp`
	`game_slots_vsn/internal/module`
	`game_slots_vsn/internal/service`
	`game_slots_vsn/internal/utls`
	`game_slots_vsn/pkg/logger`
	`github.com/gin-gonic/gin`
	`net/http`
	`sort`
)

type cServer struct {
}

func Server() (cV *cServer) {
	return &cServer{}
}

func (cV *cServer) GetServerList(ctx *gin.Context) {
	serverList := service.Server().GetAllVersion()
	sort.Slice(serverList, func(i, j int) bool {
		return serverList[i].Vsn < serverList[j].Vsn
	})
	gm := service.GM().GetGM()

	logger.Logger.Infof("get server info list:%v gm:%v", serverList, gm)
	ctx.JSON(http.StatusOK, rsp.Success(map[string]interface{}{
		"serverList": serverList,
		"gm":         gm,
	}))
	return
}

func (cV *cServer) GetServer(ctx *gin.Context) {
	vsn := ctx.Query("vsn")
	if vsn == "" {
		logger.Logger.Errorf("get server info by vsn:%v err", vsn)
		ctx.JSON(http.StatusOK, rsp.Fail(rsp.ParamsError, map[string]string{
			"vsn": vsn,
		}))
		return
	}
	versionInfo, err1 := service.Server().GetVersionByKey(vsn)
	if err1 != nil {
		logger.Logger.Errorf("get server info:%v error:%v", vsn, err1)
		ctx.JSON(http.StatusOK, rsp.Fail(rsp.ServerInfoNotFound, map[string]string{}))
		return
	}
	logger.Logger.Infof("get server info:%v ResultMap:%v", vsn, versionInfo)
	gmInfo := service.GM().GetGM()
	isGM := false
	if gmInfo.GMEnable {
		if utls.MatchIp(ctx.ClientIP()) {
			versionInfo.SrvUrl = gmInfo.GMSrvUrl
			versionInfo.ResUrl = gmInfo.GMResUrl
			isGM = true
		}
	}

	ctx.JSON(http.StatusOK, rsp.Success(map[string]interface{}{
		"isGm":   isGM,
		"srvUrl": versionInfo.SrvUrl,
		"resUrl": versionInfo.ResUrl,
		"enable": versionInfo.Enable,
	}))
	return
}

func (cV *cServer) AddServer(ctx *gin.Context) {
	buf := make([]byte, 1024)
	n, _ := ctx.Request.Body.Read(buf)
	versionInfo := module.NewServer()
	err1 := json.Unmarshal(buf[:n], versionInfo)
	if err1 != nil {
		logger.Logger.Errorf("add1 server info:%v unmarshal error", buf[:n])
		ctx.JSON(http.StatusOK, rsp.Fail(rsp.ParamsError, map[string]string{}))
		return
	}
	if versionInfo.Vsn == "" || versionInfo.SrvUrl == "" || versionInfo.ResUrl == "" {
		logger.Logger.Errorf("add2 server info:%v error", versionInfo)
		ctx.JSON(http.StatusOK, rsp.Fail(rsp.ParamsError, versionInfo))
		return
	}
	_, err3 := service.Server().AddVersion(versionInfo)
	if err3 != nil {
		logger.Logger.Errorf("add4 server info:%v error:%v", versionInfo, err3)
		ctx.JSON(http.StatusOK, rsp.Fail(rsp.DataError, versionInfo))
		return
	}
	logger.Logger.Infof("add server info:%v success", versionInfo)
	ctx.JSON(http.StatusOK, rsp.Success(versionInfo))
	return
}

func (cV *cServer) UpdateServer(ctx *gin.Context) {
	key := ctx.Query("vsn")

	buf := make([]byte, 1024)
	n, _ := ctx.Request.Body.Read(buf)
	versionInfo := module.NewServer()
	err := json.Unmarshal(buf[:n], versionInfo)
	if err != nil || key == "" || versionInfo.SrvUrl == "" || versionInfo.ResUrl == "" {
		logger.Logger.Errorf("update server info:%v error", versionInfo)
		ctx.JSON(http.StatusOK, rsp.Fail(rsp.ParamsError, map[string]string{}))
		return
	}
	if key == versionInfo.Vsn {
		logger.Logger.Infof("update server info:%v is same", versionInfo.Vsn)
	} else {
		logger.Logger.Infof("update server info:%v is not same", versionInfo.Vsn)
		service.Server().DeleteVersion(key)
	}
	_, err2 := service.Server().AddVersion(versionInfo)
	if err2 != nil {
		logger.Logger.Errorf("update server info error:%v", err2)
		ctx.JSON(http.StatusOK, rsp.Fail(rsp.ParamsError, map[string]string{
			"message": err.Error(),
		}))
		return
	}
	logger.Logger.Infof("update server info:%v success", versionInfo)
	ctx.JSON(http.StatusOK, rsp.Success(versionInfo))
}

func (cV *cServer) DeleteServer(ctx *gin.Context) {
	key := ctx.Query("vsn")
	count, err := service.Server().DeleteVersion(key)
	if err != nil {
		logger.Logger.Errorf("delete server info:%v error", key)
		ctx.JSON(http.StatusOK, rsp.Fail(rsp.ParamsError, map[string]string{
			"message": err.Error(),
		}))
		return
	}
	logger.Logger.Infof("delete server info:%v success", key)
	ctx.JSON(http.StatusOK, rsp.Success(map[string]int64{
		"count": count,
	}))
	return
}
