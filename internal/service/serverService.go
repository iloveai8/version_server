package service

import (
	"game_slots_vsn/internal/service/dao"
	"game_slots_vsn/internal/service/models"
	"game_slots_vsn/pkg/consts"
	"strconv"
)

type ServerService struct {
	PlatType string
	MaxVsn   string
	Env      string
	Vsn      string
	IsGm     bool
}

func (gs *ServerService) GetServerInfo() []*models.SubServerInfo {
	serverInfoList := make([]*models.SubServerInfo, 0)

	serverInfoMap := dao.GetServerInfos(gs.PlatType)

	if gs.Env == "pro" {
		for maxVsn, serverInfo := range serverInfoMap {
			if maxVsn == gs.Vsn {
				subServerInfoMap := serverInfo.SubServerInfoMap
				for subVsn, subServerInfo := range subServerInfoMap {
					if gs.IsGm {
						if subServerInfo.Type < consts.ServerTypePRE {
							delete(subServerInfoMap, subVsn)
						}
					} else {
						if subServerInfo.Type != consts.ServerTypePRO && subServerInfo.Type != consts.ServerTypeDefault {
							delete(subServerInfoMap, subVsn)
						}
					}
				}
				if len(subServerInfoMap) > 0 {
					maxSubServer := getMaxSubServerInfo(subServerInfoMap)
					maxSubServer.Vsn = joinVsn(maxVsn, maxSubServer.Vsn)
					serverInfoList = append(serverInfoList, maxSubServer)
				}
			}
		}
		if gs.IsGm {
			innerServerMap := dao.GetServerInfos("inner")
			for _, server := range innerServerMap {
				subVsnMap := server.SubServerInfoMap
				for _, subServer := range subVsnMap {
					serverInfoList = append(serverInfoList, subServer)
				}
			}
		}
	} else {
		innerServerMap := dao.GetServerInfos("inner")
		for _, server := range innerServerMap {
			subVsnMap := server.SubServerInfoMap
			for _, subServer := range subVsnMap {
				serverInfoList = append(serverInfoList, subServer)
			}
		}
	}
	return serverInfoList
}

func joinVsn(max, min string) string {
	return max + "." + min
}

func getMaxSubServerInfo(m map[string]*models.SubServerInfo) *models.SubServerInfo {
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

func (gs *ServerService) GetServerInfos() map[string]*models.ServerInfo {
	return dao.GetServerInfos(gs.PlatType)
}

func (gs *ServerService) GetServerInfoByMaxVsn() *models.ServerInfo {
	return dao.GetServerInfo(gs.PlatType, gs.MaxVsn)
}

func (gs *ServerService) AddServerInfo(server *models.ServerInfo) bool {
	return dao.AddServerInfo(gs.PlatType, gs.MaxVsn, server)
}

func (gs *ServerService) DeleteServerInfo() bool {
	return dao.DeleteServerInfo(gs.PlatType, gs.MaxVsn)
}
