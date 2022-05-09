package controller

import (
	`encoding/json`
	`game_slots_vsn/internal/api/rsp`
	`game_slots_vsn/internal/consts`
	`game_slots_vsn/internal/module`
	`game_slots_vsn/internal/service`
	`game_slots_vsn/internal/utls`
	`game_slots_vsn/pkg/logger`
	`github.com/gin-gonic/gin`
	`net/http`
	`sort`
	`strconv`
	`strings`
)

type cServer struct {
}

func Server() (cV *cServer) {
	return &cServer{}
}

func (cV *cServer) GetServerList(ctx *gin.Context) {
	serverMap := service.Server().GetAllServer()
	subServerList := make([]*module.SubServer, 0)
	for maxVsn, server := range serverMap {
		if strings.IndexAny(maxVsn, ".") == -1 {
			subServerList = append(subServerList, server.SubServer[maxVsn])
		} else {
			subVsnList := server.SubServer
			for _, subServer := range subVsnList {
				subServer.Vsn = joinVsn(maxVsn, subServer.Vsn)
				subServerList = append(subServerList, subServer)
			}
		}
	}
	sort.Slice(subServerList, func(i, j int) bool {
		return subServerList[i].Vsn < subServerList[j].Vsn
	})
	gm := service.GM().GetGM()

	logger.Logger.Infof("get server info list:%v gm:%v", subServerList, gm)
	ctx.JSON(http.StatusOK, rsp.Success(map[string]interface{}{
		"serverList": subServerList,
		"gm":         gm,
	}))
	return
}

func (cV *cServer) GetServer(ctx *gin.Context) {
	env := ctx.Query("env")
	vsn := ctx.Query("vsn")
	if env != "dev" && env != "pro" || strings.IndexAny(vsn, ".") == -1 {
		logger.Logger.Errorf("get vsn err:%v", vsn)
		ctx.JSON(http.StatusOK, rsp.Fail(rsp.ParamsError, map[string]string{
			"vsn": vsn,
			"env": env,
		}))
		return
	}

	isGm := false
	gmInfo := service.GM().GetGM()
	serverMap := service.Server().GetAllServer()
	subServerList := make([]*module.SubServer, 0)
	if env == "dev" {
		isGm = true
		for maxVsn, server := range serverMap {
			if maxVsn == vsn {
				subVsnList := server.SubServer
				for k, v := range subVsnList {
					if v.Type < consts.ServerTypeDEV {
						delete(subVsnList, k)
					}
				}
				if len(subVsnList) > 0 {
					maxSubServer := getMaxServer(subVsnList)
					maxSubServer.Vsn = joinVsn(maxVsn, maxSubServer.Vsn)
					subServerList = append(subServerList, maxSubServer)
				}
			} else {
				subServerList = gmFilter(isGm, maxVsn, server, subServerList)
			}
		}
	} else if env == "pro" {
		isGm = gmInfo.GMEnable && utls.MatchIp(ctx.ClientIP())
		for maxVsn, server := range serverMap {
			if maxVsn == vsn {
				subVsnList := server.SubServer
				for k, v := range subVsnList {
					if isGm {
						if v.Type < consts.ServerTypePRE {
							delete(subVsnList, k)
						}
					} else {
						if v.Type < consts.ServerTypePRO {
							delete(subVsnList, k)
						}
					}
				}
				if len(subVsnList) > 0 {
					maxSubServer := getMaxServer(subVsnList)
					maxSubServer.Vsn = joinVsn(maxVsn, maxSubServer.Vsn)
					subServerList = append(subServerList, maxSubServer)
				}
			} else {
				subServerList = gmFilter(isGm, maxVsn, server, subServerList)
			}
		}
	}
	ctx.JSON(http.StatusOK, rsp.Success(map[string]interface{}{
		"serverList": subServerList,
		"gm":         isGm,
	}))
	return
}

func gmFilter(isGm bool, maxVsn string, server *module.Server, subServerList []*module.SubServer) []*module.SubServer {
	if isGm && strings.IndexAny(maxVsn, ".") == -1 {
		subVsnList := server.SubServer
		for _, subServer := range subVsnList {
			if subServer.Type > consts.ServerTypeDefault {
				subServerList = append(subServerList, subServer)
			}
		}
	}
	//if isGm {
	//	if strings.IndexAny(maxVsn, ".") == -1 {
	//		subVsnList := server.SubServer
	//		for _, subServer := range subVsnList {
	//			if subServer.Type > consts.ServerTypeDefault {
	//				subServerList = append(subServerList, subServer)
	//			}
	//		}
	//	} else {
	//		subVsnList := server.SubServer
	//		for _, subServer := range subVsnList {
	//			if subServer.Type > consts.ServerTypeDefault {
	//				subServer.Vsn = joinVsn(maxVsn, subServer.Vsn)
	//				subServerList = append(subServerList, subServer)
	//			}
	//		}
	//	}
	//}
	return subServerList
}

func (cV *cServer) AddServer(ctx *gin.Context) {
	buf := make([]byte, 1024)
	n, _ := ctx.Request.Body.Read(buf)
	subServerInfo := module.NewSubServer()
	err1 := json.Unmarshal(buf[:n], subServerInfo)
	if err1 != nil {
		logger.Logger.Errorf("add1 server info:%v unmarshal error", buf[:n])
		ctx.JSON(http.StatusOK, rsp.Fail(rsp.ParamsError, map[string]string{}))
		return
	}
	if subServerInfo.Vsn == "" || subServerInfo.SrvUrl == "" || subServerInfo.ResUrl == "" {
		logger.Logger.Errorf("add2 server info:%v error", subServerInfo)
		ctx.JSON(http.StatusOK, rsp.Fail(rsp.ParamsError, subServerInfo))
		return
	}

	if strings.IndexAny(subServerInfo.Vsn, ".") == -1 {
		serverInfo := module.NewServer()
		serverInfo.SubServer[subServerInfo.Vsn] = subServerInfo

		_, _ = service.Server().AddServer(subServerInfo.Vsn, serverInfo)

		logger.Logger.Infof("add server info:%v success", serverInfo)
		ctx.JSON(http.StatusOK, rsp.Success(subServerInfo))
		return
	}

	vsn, subVsn := fmtVsn(subServerInfo.Vsn)
	subServerInfo.Vsn = subVsn

	serverInfo, _ := service.Server().GetServerByKey(vsn)
	serverInfo.SubServer[subVsn] = subServerInfo
	_, err3 := service.Server().AddServer(vsn, serverInfo)
	if err3 != nil {
		logger.Logger.Errorf("add3 server info:%v error:%v", subServerInfo, err3)
		ctx.JSON(http.StatusOK, rsp.Fail(rsp.DataError, subServerInfo))
		return
	}
	logger.Logger.Infof("add server info:%v success", subServerInfo)
	ctx.JSON(http.StatusOK, rsp.Success(subServerInfo))
	return
}

func (cV *cServer) UpdateServer(ctx *gin.Context) {
	key := ctx.Query("vsn")

	buf := make([]byte, 1024)
	n, _ := ctx.Request.Body.Read(buf)
	subServerInfo := module.NewSubServer()
	err := json.Unmarshal(buf[:n], subServerInfo)
	if err != nil || key == "" || subServerInfo.SrvUrl == "" || subServerInfo.ResUrl == "" {
		logger.Logger.Errorf("update server info:%v error", subServerInfo)
		ctx.JSON(http.StatusOK, rsp.Fail(rsp.ParamsError, map[string]string{}))
		return
	}
	if key == subServerInfo.Vsn {
		if strings.IndexAny(subServerInfo.Vsn, ".") == -1 {
			serverInfo := module.NewServer()
			serverInfo.SubServer[key] = subServerInfo
			_, _ = service.Server().AddServer(key, serverInfo)
			if err != nil {
				return
			}
		} else {
			vsn, subVsn := fmtVsn(subServerInfo.Vsn)
			subServerInfo.Vsn = subVsn
			serverInfo, _ := service.Server().GetServerByKey(vsn)
			serverInfo.SubServer[subVsn] = subServerInfo
			_, _ = service.Server().AddServer(vsn, serverInfo)
		}
	} else {
		if strings.IndexAny(key, ".") == -1 {
			_, _ = service.Server().DeleteServer(key)
		} else {
			dVsn, dSubVsn := fmtVsn(key)
			dServerInfo, _ := service.Server().GetServerByKey(dVsn)
			delete(dServerInfo.SubServer, dSubVsn)
			if len(dServerInfo.SubServer) == 0 {
				_, _ = service.Server().DeleteServer(dVsn)
			}
			_, _ = service.Server().AddServer(dVsn, dServerInfo)
		}
		if strings.IndexAny(subServerInfo.Vsn, ".") == -1 {
			serverInfo := module.NewServer()
			serverInfo.SubServer[subServerInfo.Vsn] = subServerInfo
			_, _ = service.Server().AddServer(subServerInfo.Vsn, serverInfo)

			logger.Logger.Infof("update server info:%v success", subServerInfo)
			ctx.JSON(http.StatusOK, rsp.Success(subServerInfo))
		} else {
			vsn, subVsn := fmtVsn(subServerInfo.Vsn)
			subServerInfo.Vsn = subVsn
			serverInfo, _ := service.Server().GetServerByKey(vsn)
			serverInfo.SubServer[subVsn] = subServerInfo
			_, _ = service.Server().AddServer(vsn, serverInfo)
		}
	}
	logger.Logger.Infof("update server info:%v success", subServerInfo)
	ctx.JSON(http.StatusOK, rsp.Success(subServerInfo))
	return
}

func (cV *cServer) DeleteServer(ctx *gin.Context) {
	key := ctx.Query("vsn")
	if strings.IndexAny(key, ".") == -1 {
		count, _ := service.Server().DeleteServer(key)
		logger.Logger.Infof("delete server info:%v success", key)
		ctx.JSON(http.StatusOK, rsp.Success(map[string]int64{
			"count": count,
		}))
		return
	}

	vsn, subVsn := fmtVsn(key)
	serverInfo, _ := service.Server().GetServerByKey(vsn)
	delete(serverInfo.SubServer, subVsn)
	if len(serverInfo.SubServer) == 0 {
		count, _ := service.Server().DeleteServer(vsn)
		logger.Logger.Infof("delete server info:%v success", vsn)
		ctx.JSON(http.StatusOK, rsp.Success(map[string]int64{
			"count": count,
		}))
		return
	}
	_, _ = service.Server().AddServer(vsn, serverInfo)
	ctx.JSON(http.StatusOK, rsp.Success(map[string]int64{
		"count": 1,
	}))
	return
}

func fmtVsn(s string) (string, string) {
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

func getMaxServer(m map[string]*module.SubServer) *module.SubServer {
	max := 0
	for k := range m {
		ikey, err := strconv.Atoi(k)
		if err != nil {
			continue
		}
		if ikey >= max {
			max = ikey
		}
	}
	return m[strconv.Itoa(max)]
}
