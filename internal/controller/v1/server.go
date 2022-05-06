package v1

import (
	"encoding/json"
	"fmt"
	`game_slots_vsn/internal/api/rsp`
	`game_slots_vsn/internal/consts`
	`game_slots_vsn/internal/utls`
	"game_slots_vsn/pkg/logger"
	"game_slots_vsn/pkg/redis"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Vsn struct {
	Vsn    string `json:"vsn"`
	SrvUrl string `json:"srvUrl"`
	ResUrl string `json:"resUrl"`
	Enable bool   `json:"enable"`
}

type GMConf struct {
	GMSrvUrl string `json:"gmSrvUrl"`
	GMResUrl string `json:"gmResUrl"`
	GMEnable bool   `json:"gmEnable"`
}

func newGC() *GMConf {
	return &GMConf{}
}
func newVSN() *Vsn {
	return &Vsn{}
}

type cServer1 struct {
}

func Server1() (cV *cServer1) {
	return &cServer1{}
}

//GetAll 获取所有vsn 信息
func (vh *cServer1) GetAll(c *gin.Context) {
	result, _ := redis.Client.HVals(consts.CacheVsnKey).Result()

	vsnList := make([]*Vsn, len(result), cap(result))
	for i, vsnStr := range result {
		vsn := newVSN()
		_ = json.Unmarshal([]byte(vsnStr), vsn)
		vsnList[i] = vsn
	}
	result1, _ := redis.Client.Get(consts.CacheGMConfKey).Result()
	gmConfInfo := newGC()
	_ = json.Unmarshal([]byte(result1), gmConfInfo)
	logger.Logger.Infof(" ==>gmConfInfo:%v vsnInfoList:%v", gmConfInfo, vsnList)
	c.JSON(http.StatusOK, rsp.Success(map[string]interface{}{
		"vsnList": vsnList,
		"gmConf":  gmConfInfo,
	}))
}

//Get 根据vsn获取信息
func (vh *cServer1) Get(c *gin.Context) {
	vsn := c.Query("vsn")
	if vsn == "" {
		logger.Logger.Errorf("get vsn err:%v", vsn)
		c.JSON(http.StatusOK, rsp.Fail(rsp.ParamsError, map[string]string{
			"vsn": vsn,
		}))
		return
	}
	result, err := redis.Client.HGet(consts.CacheVsnKey, vsn).Result()
	if err != nil {
		logger.Logger.Errorf("get vsn err:%v", err)
		c.JSON(http.StatusOK, rsp.Fail(rsp.ServerInfoNotFound, map[string]string{
			"vsn": vsn,
		}))
		return
	}
	vsnInfo := newVSN()
	_ = json.Unmarshal([]byte(result), vsnInfo)

	result, err = redis.Client.Get(consts.CacheGMConfKey).Result()
	gmConf := newGC()
	_ = json.Unmarshal([]byte(result), gmConf)
	isGM := false
	if gmConf.GMEnable {
		ip := c.ClientIP()
		logger.Logger.Warnf(" ==>client ip:%v", ip)
		if utls.MatchIp(ip) {
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
	c.JSON(http.StatusOK, rsp.Success(reply))
}

//Insert 存储vsn信息
func (vh *cServer1) Insert(c *gin.Context) {
	buf := make([]byte, 1024)
	n, _ := c.Request.Body.Read(buf)
	fmt.Println(string(buf[:n]))
	newVsnInfo := newVSN()
	err := json.Unmarshal(buf[:n], newVsnInfo)
	if err != nil {
		logger.Logger.Errorf("add vsn unmarshal err:%v", err)
		c.JSON(http.StatusOK, rsp.Fail(rsp.ParamsError, map[string]string{
			"message": string(buf[:n]),
		}))
		return
	}
	if newVsnInfo.Vsn == "" || newVsnInfo.SrvUrl == "" || newVsnInfo.ResUrl == "" {
		logger.Logger.Errorf("add vsn:%v", newVsnInfo)
		c.JSON(http.StatusOK, rsp.Fail(rsp.ParamsError, newVsnInfo))
		return
	}
	marshal, _ := json.Marshal(newVsnInfo)
	redis.Client.HSet(consts.CacheVsnKey, newVsnInfo.Vsn, marshal)
	logger.Logger.Infof(" ==>add vsnInfo:%v", newVsnInfo)
	c.JSON(http.StatusOK, rsp.Success(newVsnInfo))
}

//Update 存储vsn信息
func (vh *cServer1) Update(c *gin.Context) {
	vsn := c.Query("vsn")

	buf := make([]byte, 1024)
	n, _ := c.Request.Body.Read(buf)
	newVsnInfo := newVSN()
	err := json.Unmarshal(buf[:n], newVsnInfo)
	if err != nil || vsn == "" || newVsnInfo.SrvUrl == "" || newVsnInfo.ResUrl == "" {
		logger.Logger.Errorf("update vsn:%v", newVsnInfo)
		c.JSON(http.StatusOK, rsp.Fail(rsp.ParamsError, map[string]string{
			"message": string(buf[:n]),
		}))
		return
	}
	newVsnInfo.Vsn = vsn
	marshal, _ := json.Marshal(newVsnInfo)

	redis.Client.HSet(consts.CacheVsnKey, vsn, marshal)
	logger.Logger.Infof(" ==>update vsnInfo:%v", newVsnInfo)
	c.JSON(http.StatusOK, rsp.Success(newVsnInfo))
	return
}

//Delete 存储vsn信息
func (vh *cServer1) Delete(c *gin.Context) {
	vsn := c.Query("vsn")
	if vsn == "" {
		logger.Logger.Errorf("delete vsn:%v", vsn)
		c.JSON(http.StatusOK, rsp.Fail(rsp.ParamsError, map[string]string{
			"vsn": vsn,
		}))
		return
	}
	result, _ := redis.Client.HDel(consts.CacheVsnKey, vsn).Result()
	logger.Logger.Infof(" ==>delete vsn:%v result:%v", vsn, result)
	c.JSON(http.StatusOK, rsp.Success(map[string]int64{
		"count": result,
	}))
	return
}

func (vh *cServer1) InsertGmConf(c *gin.Context) {
	buf := make([]byte, 1024)
	n, _ := c.Request.Body.Read(buf)
	gmConfInfo := newGC()
	err := json.Unmarshal(buf[:n], gmConfInfo)
	if err != nil || gmConfInfo.GMSrvUrl == "" || gmConfInfo.GMResUrl == "" {
		logger.Logger.Errorf("set global conf err:%v", gmConfInfo)
		c.JSON(http.StatusOK, rsp.Fail(rsp.ParamsError, map[string]string{
			"message": string(buf[:n]),
		}))
		return
	}
	marshal, _ := json.Marshal(gmConfInfo)
	redis.Client.Set(consts.CacheGMConfKey, marshal, 0)
	logger.Logger.Infof(" ==>set gmConfInfo:%v", gmConfInfo)
	c.JSON(http.StatusOK, rsp.Success(gmConfInfo))
}
