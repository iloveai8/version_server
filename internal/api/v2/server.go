package v2

import (
	"encoding/json"
	"game_slots_vsn/internal/service"
	"game_slots_vsn/internal/service/models"
	"game_slots_vsn/pkg/e"
	"game_slots_vsn/pkg/ggeoip"
	"game_slots_vsn/pkg/logger"
	"game_slots_vsn/pkg/utils"
	"github.com/gin-gonic/gin"
	"sort"
	"strings"
)

func GetServerInfoList(c *gin.Context) {
	appG := e.Gin{C: c}
	platType := c.Query("platType")

	serverService := &service.ServerService{
		PlatType: platType,
	}
	serverInfoMap := serverService.GetServerInfos()

	subServerInfoList := make([]*models.SubServerInfo, 0)
	for maxVsn, serverInfo := range serverInfoMap {
		if strings.IndexAny(maxVsn, ".") == -1 {
			subServerInfoList = append(subServerInfoList, serverInfo.SubServerInfoMap[maxVsn])
		} else {
			subServerInfoMap := serverInfo.SubServerInfoMap
			if len(subServerInfoMap) == 0 {
				serverService.MaxVsn = maxVsn
				serverService.DeleteServerInfo()
			}
			for _, subServer := range subServerInfoMap {
				subServer.Vsn = joinVsn(maxVsn, subServer.Vsn)
				subServerInfoList = append(subServerInfoList, subServer)
			}
		}
	}
	sort.Slice(subServerInfoList, func(i, j int) bool {
		return subServerInfoList[i].Vsn < subServerInfoList[j].Vsn
	})

	gmService := &service.GMService{}
	gmInfo, err := gmService.GetGmInfo()
	if err != nil {
		logger.ErrorF("get gmInfo:%v error:%v", *gmInfo, err)
	}
	logger.InfoF("get serverInfoList:%v gm:%v", subServerInfoList, *gmInfo)
	appG.Success(e.SUCCESS, map[string]interface{}{
		"serverList": subServerInfoList,
		"gm":         gmInfo,
	})
	return
}

func GetServerInfo(c *gin.Context) {
	appG := e.Gin{C: c}

	platType := c.Query("platType")
	env := c.Query("env")
	vsn := c.Query("vsn")
	if env != "dev" && env != "pro" || strings.IndexAny(vsn, ".") == -1 {
		logger.ErrorF("get server info platType:% v env:%v vsn:%v err", platType, env, vsn)
		appG.Fail(e.ParamError, map[string]string{
			"platType": platType,
			"vsn":      vsn,
			"env":      env,
		})
		return
	}
	gmService := &service.GMService{}
	gmInfo, err := gmService.GetGmInfo()
	if err != nil {
		logger.ErrorF("get gmInfo:%v error:%v", gmInfo, err)
	}
	isGm := gmInfo.GMEnable
	isBlock := gmInfo.Block
	if env == "pro" {
		logger.InfoF("X-Forwarded-For:%s  X-Real-IP:%s cli:%s", c.Request.Header.Get("X-Forwarded-For"), c.Request.Header.Get("X-Real-IP"), c.ClientIP())
		isGm = gmInfo.GMEnable && utils.IsGmIP(c.ClientIP())
	}
	serverService := &service.ServerService{
		PlatType: platType,
		Env:      env,
		Vsn:      vsn,
		IsGm:     isGm,
	}
	country := ggeoip.Gip.GetCountryByIP(c.ClientIP())
	subServerList := serverService.GetServerInfo()
	appG.Success(e.SUCCESS, map[string]interface{}{
		"serverList": subServerList,
		"gm":         isGm,
		"block":      isBlock,
		"country":    country,
		"ip":         c.ClientIP(),
	})
	return
}

func AddServerInfo(c *gin.Context) {
	appG := e.Gin{C: c}

	platType := c.Query("platType")

	buf := make([]byte, 1024)
	n, _ := c.Request.Body.Read(buf)
	subServerInfo := &models.SubServerInfo{}
	err := json.Unmarshal(buf[:n], subServerInfo)

	if err != nil {
		logger.ErrorF("add serverInfo:%v unmarshal err:%v", buf[:n], err)
		appG.Fail(e.ParamError, map[string]string{})
		return
	}

	if subServerInfo.Vsn == "" || subServerInfo.SrvUrl == "" || subServerInfo.ResUrl == "" {
		logger.ErrorF("add serverInfo:%v fail", subServerInfo)
		appG.Fail(e.ParamError, map[string]string{})
		return
	}

	serverService := &service.ServerService{
		PlatType: platType,
	}
	if strings.IndexAny(subServerInfo.Vsn, ".") == -1 {
		subServerInfoMap := make(map[string]*models.SubServerInfo, 1)
		serverInfo := &models.ServerInfo{
			SubServerInfoMap: subServerInfoMap,
		}
		serverInfo.SubServerInfoMap[subServerInfo.Vsn] = subServerInfo
		serverService.MaxVsn = subServerInfo.Vsn
		hSet := serverService.AddServerInfo(serverInfo)
		logger.InfoF("add serverInfo:%v hSet:~%v", serverInfo, hSet)
		appG.Success(e.SUCCESS, serverInfo)
		return
	} else {
		maxVsn, subVsn := splitVsn(subServerInfo.Vsn)
		subServerInfo.Vsn = subVsn

		serverService.MaxVsn = maxVsn
		serverInfo := serverService.GetServerInfoByMaxVsn()
		serverInfo.SubServerInfoMap[subServerInfo.Vsn] = subServerInfo

		hSet := serverService.AddServerInfo(serverInfo)
		logger.InfoF("add serverInfo:%v hSet:~%v", serverInfo, hSet)

		appG.Success(e.SUCCESS, serverInfo)
		return
	}
}

func UpdateServerInfo(c *gin.Context) {
	appG := e.Gin{C: c}

	platType := c.Query("platType")
	key := c.Query("vsn")

	buf := make([]byte, 1024)
	n, _ := c.Request.Body.Read(buf)

	subServerInfo := &models.SubServerInfo{}
	err := json.Unmarshal(buf[:n], subServerInfo)
	if err != nil || key == "" || subServerInfo.SrvUrl == "" || subServerInfo.ResUrl == "" {
		logger.ErrorF("update serverInfo:%v err:%v", subServerInfo, "error")
		appG.Fail(e.ParamError, map[string]string{})
		return
	}

	serverService := &service.ServerService{
		PlatType: platType,
	}
	if key != subServerInfo.Vsn {
		if strings.IndexAny(key, ".") == -1 {
			serverService.MaxVsn = key
			serverService.DeleteServerInfo()
		} else {
			dMaxVsn, dSubVsn := splitVsn(key)
			serverService.MaxVsn = dMaxVsn

			dServerInfo := serverService.GetServerInfoByMaxVsn()
			delete(dServerInfo.SubServerInfoMap, dSubVsn)

			if len(dServerInfo.SubServerInfoMap) == 0 {
				serverService.DeleteServerInfo()
			} else {
				serverService.AddServerInfo(dServerInfo)
			}
		}
	}
	if strings.IndexAny(subServerInfo.Vsn, ".") == -1 {
		subServerInfoMap := make(map[string]*models.SubServerInfo, 1)
		serverInfo := &models.ServerInfo{
			SubServerInfoMap: subServerInfoMap,
		}
		serverInfo.SubServerInfoMap[subServerInfo.Vsn] = subServerInfo

		serverService.MaxVsn = subServerInfo.Vsn
		hSet := serverService.AddServerInfo(serverInfo)

		logger.InfoF("update serverInfo:%v hSet:~%v", serverInfo, hSet)
		appG.Success(e.SUCCESS, serverInfo)
	} else {
		maxVsn, subVsn := splitVsn(subServerInfo.Vsn)
		subServerInfo.Vsn = subVsn

		serverService.MaxVsn = maxVsn
		serverInfo := serverService.GetServerInfoByMaxVsn()

		serverInfo.SubServerInfoMap[subServerInfo.Vsn] = subServerInfo
		hSet := serverService.AddServerInfo(serverInfo)
		logger.InfoF("update serverInfo:%v hSet:~%v", serverInfo, hSet)
		appG.Success(e.SUCCESS, serverInfo)
	}
	return
}

func DeleteServerInfo(c *gin.Context) {
	appG := e.Gin{C: c}

	platType := c.Query("platType")
	key := c.Query("vsn")
	serverService := &service.ServerService{
		PlatType: platType,
	}
	if strings.IndexAny(key, ".") == -1 {
		serverService.MaxVsn = key
		serverService.DeleteServerInfo()
		logger.InfoF("delete serverInfo:%v success", key)
	} else {
		dMaxVsn, dSubVsn := splitVsn(key)
		serverService.MaxVsn = dMaxVsn

		dServerInfo := serverService.GetServerInfoByMaxVsn()
		delete(dServerInfo.SubServerInfoMap, dSubVsn)

		if len(dServerInfo.SubServerInfoMap) == 0 {
			serverService.DeleteServerInfo()
		} else {
			serverService.AddServerInfo(dServerInfo)
		}
		logger.InfoF("delete serverInfo:%v success", key)
	}
	appG.Success(e.SUCCESS, map[string]int64{})
	return
}

func splitVsn(s string) (string, string) {
	ss := strings.Split(strings.Replace(s, " ", "", -1), ".")
	if len(ss) == 3 {
		return strings.Join(ss, "."), "0"
	} else if len(ss) == 4 {
		return strings.Join(ss[:3], "."), ss[3]
	}
	return s, "0"
}

func joinVsn(max, min string) string {
	return max + "." + min
}
