package handler

import (
	"encoding/json"
	"fmt"
	"game_slots_vsn/internal/mod"
	"game_slots_vsn/internal/respone"
	"game_slots_vsn/pkg/logger"
	"game_slots_vsn/pkg/redis"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
)

const (
	CacheVsnKey    string = "vsn."
	CacheGMConfKey string = "gm.conf."
)

var IPArray = [...]string{
	"103.85.165.146/29",
	"106.120.91.66/28",
	"106.120.91.67/28",
	"106.120.91.68/28",
	"1.202.246.18/28",
	"1.202.246.19/28",
	"1.202.246.20/28",
	"1.202.246.21/28",
	"1.202.246.30/28",
}

type VsnHandler struct {
}

func NewVsnHandler() *VsnHandler {
	return &VsnHandler{}
}

//GetAll 获取所有vsn 信息
func (vh VsnHandler) GetAll(c *gin.Context) {
	result, _ := redis.Client.HVals(CacheVsnKey).Result()

	vsnList := make([]*mod.Vsn, len(result), cap(result))
	for i, vsnStr := range result {
		vsn := mod.NewVsn()
		_ = json.Unmarshal([]byte(vsnStr), vsn)
		vsnList[i] = vsn
	}
	result1, _ := redis.Client.Get(CacheGMConfKey).Result()
	gmConfInfo := mod.NewGmConf()
	_ = json.Unmarshal([]byte(result1), gmConfInfo)
	logger.Logger.Infof(" ==>gmConfInfo:%v vsnInfoList:%v", gmConfInfo, vsnList)
	c.JSON(http.StatusOK, respone.Success(map[string]interface{}{
		"vsnList": vsnList,
		"gmConf":  gmConfInfo,
	}))
}

//Get 根据vsn获取信息
func (vh VsnHandler) Get(c *gin.Context) {
	vsn := c.Query("vsn")
	if vsn == "" {
		logger.Logger.Errorf("get vsn err:%v", vsn)
		c.JSON(http.StatusOK, respone.Fail(respone.ParamsError, map[string]string{
			"vsn": vsn,
		}))
		return
	}
	result, err := redis.Client.HGet(CacheVsnKey, vsn).Result()
	if err != nil {
		logger.Logger.Errorf("get vsn err:%v", err)
		c.JSON(http.StatusOK, respone.Fail(respone.VersionNotFound, map[string]string{
			"vsn": vsn,
		}))
		return
	}
	vsnInfo := mod.NewVsn()
	_ = json.Unmarshal([]byte(result), vsnInfo)

	result, err = redis.Client.Get(CacheGMConfKey).Result()
	gmConf := mod.NewGmConf()
	_ = json.Unmarshal([]byte(result), gmConf)
	isGM := false
	if gmConf.GMEnable {
		ip := c.ClientIP()
		logger.Logger.Warnf(" ==>client ip:%v", ip)
		if MatchIp(ip) {
			logger.Logger.Warnf(" ==>gm enable:%v client ip:%v is in inner", gmConf.GMEnable, ip)
			vsnInfo.SrvUrl = gmConf.GMSrvUrl
			vsnInfo.ResUrl = gmConf.GMResUrl
			isGM = true
		}
		//else if _, ok := IPMap[ip]; ok {
		//	logger.Logger.Warnf(" ==>client ip:%v is in out company white list:%v", ip, IPMap)
		//	vsnInfo.SrvUrl = gmConf.GMSrvUrl
		//	vsnInfo.ResUrl = gmConf.GMResUrl
		//	isGM = true
		//}
	}
	logger.Logger.Infof(" ==>gmConf:%v vsnInfo:%v", gmConf, vsnInfo)
	reply := map[string]interface{}{
		"isGm":   isGM,
		"srvUrl": vsnInfo.SrvUrl,
		"resUrl": vsnInfo.ResUrl,
		"enable": vsnInfo.Enable,
	}
	c.JSON(http.StatusOK, respone.Success(reply))
}

//Insert 存储vsn信息
func (vh VsnHandler) Insert(c *gin.Context) {
	buf := make([]byte, 1024)
	n, _ := c.Request.Body.Read(buf)
	fmt.Println(string(buf[:n]))
	newVsnInfo := mod.NewVsn()
	err := json.Unmarshal(buf[:n], newVsnInfo)
	if err != nil {
		logger.Logger.Errorf("add vsn unmarshal err:%v", err)
		c.JSON(http.StatusOK, respone.Fail(respone.ParamsError, map[string]string{
			"message": string(buf[:n]),
		}))
		return
	}
	if newVsnInfo.Vsn == "" || newVsnInfo.SrvUrl == "" || newVsnInfo.ResUrl == "" {
		logger.Logger.Errorf("add vsn:%v", newVsnInfo)
		c.JSON(http.StatusOK, respone.Fail(respone.ParamsError, newVsnInfo))
		return
	}
	marshal, _ := json.Marshal(newVsnInfo)
	redis.Client.HSet(CacheVsnKey, newVsnInfo.Vsn, marshal)
	logger.Logger.Infof(" ==>add vsnInfo:%v", newVsnInfo)
	c.JSON(http.StatusOK, respone.Success(newVsnInfo))
}

//Update 存储vsn信息
func (vh VsnHandler) Update(c *gin.Context) {
	vsn := c.Query("vsn")

	buf := make([]byte, 1024)
	n, _ := c.Request.Body.Read(buf)
	newVsnInfo := mod.NewVsn()
	err := json.Unmarshal(buf[:n], newVsnInfo)
	if err != nil || vsn == "" || newVsnInfo.SrvUrl == "" || newVsnInfo.ResUrl == "" {
		logger.Logger.Errorf("update vsn:%v", newVsnInfo)
		c.JSON(http.StatusOK, respone.Fail(respone.ParamsError, map[string]string{
			"message": string(buf[:n]),
		}))
		return
	}
	newVsnInfo.Vsn = vsn
	marshal, _ := json.Marshal(newVsnInfo)

	redis.Client.HSet(CacheVsnKey, vsn, marshal)
	logger.Logger.Infof(" ==>update vsnInfo:%v", newVsnInfo)
	c.JSON(http.StatusOK, respone.Success(newVsnInfo))
	return
}

//Delete 存储vsn信息
func (vh VsnHandler) Delete(c *gin.Context) {
	vsn := c.Query("vsn")
	if vsn == "" {
		logger.Logger.Errorf("delete vsn:%v", vsn)
		c.JSON(http.StatusOK, respone.Fail(respone.ParamsError, map[string]string{
			"vsn": vsn,
		}))
		return
	}
	result, _ := redis.Client.HDel(CacheVsnKey, vsn).Result()
	logger.Logger.Infof(" ==>delete vsn:%v result:%v", vsn, result)
	c.JSON(http.StatusOK, respone.Success(map[string]int64{
		"count": result,
	}))
	return
}

//edit gm conf
func (vh VsnHandler) InsertGmConf(c *gin.Context) {
	buf := make([]byte, 1024)
	n, _ := c.Request.Body.Read(buf)
	gmConfInfo := mod.NewGmConf()
	err := json.Unmarshal(buf[:n], gmConfInfo)
	if err != nil || gmConfInfo.GMSrvUrl == "" || gmConfInfo.GMResUrl == "" {
		logger.Logger.Errorf("set global conf err:%v", gmConfInfo)
		c.JSON(http.StatusOK, respone.Fail(respone.ParamsError, map[string]string{
			"message": string(buf[:n]),
		}))
		return
	}
	marshal, _ := json.Marshal(gmConfInfo)
	redis.Client.Set(CacheGMConfKey, marshal, 0)
	logger.Logger.Infof(" ==>set gmConfInfo:%v", gmConfInfo)
	c.JSON(http.StatusOK, respone.Success(gmConfInfo))
}

//"103.85.165.146/29",

//"106.120.91.66/28",
//"106.120.91.67/28",
//"106.120.91.68/28",

//"1.202.246.18/28",
//"1.202.246.19/28",
//"1.202.246.20/28",
//"1.202.246.21/28",
//"1.202.246.30/28",
func MatchIp(IP string) bool {
	for _, network := range IPArray {
		_, subnet, _ := net.ParseCIDR(network)
		if subnet.Contains(net.ParseIP(IP)) {
			fmt.Println("addr:", IP, "in subnet:v%", subnet)
			return true
		}
	}
	return false
}
