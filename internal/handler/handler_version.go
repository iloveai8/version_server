package handler

import (
	`encoding/json`
	`fmt`
	`game_slots_vsn/internal/mod`
	`game_slots_vsn/internal/respone`
	"game_slots_vsn/pkg/logger"
	"game_slots_vsn/pkg/redis"
	`github.com/gin-gonic/gin`
	`net/http`
)

const (
	CacheVsnKey    string = "vsn."
	CacheGMConfKey string = "gm.conf."
)

var IPMap = map[string]int{
	"127.0.0.1":       1,
	"40.83.97.197":    1,
	"129.226.60.247":  1,
	"47.75.45.195":    1,
	"129.226.189.243": 1,
	"47.75.59.239":    1,
	"119.81.164.4":    1,
	"1.202.246.19":    1,
	"106.120.91.66":   1,
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
	logger.Logger.Infof("vsnInfo list:%v", vsnList)

	result1, _ := redis.Client.Get(CacheGMConfKey).Result()
	gmConf := mod.NewGmConf()
	_ = json.Unmarshal([]byte(result1), gmConf)
	logger.Logger.Infof("gmConf:%v", gmConf)

	c.JSON(http.StatusOK, respone.Success(map[string]interface{}{
		"vsnList": vsnList,
		"gmConf":  gmConf,
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
	logger.Logger.Infof("gmConf:%v", gmConf)
	if gmConf.GMEnable {
		ip := c.ClientIP()
		if _, ok := IPMap[ip]; ok {
			logger.Logger.Infof(" client ip:%v is in white list:%v", ip, IPMap)

			vsnInfo.SrvUrl = gmConf.GMSrvUrl
			vsnInfo.ResUrl = gmConf.GMResUrl
		}
	}
	logger.Logger.Infof("vsnInfo:%v", vsnInfo)
	c.JSON(http.StatusOK, respone.Success(vsnInfo))
}

//Insert 存储vsn信息
func (vh VsnHandler) Insert(c *gin.Context) {
	buf := make([]byte, 1024)
	n, _ := c.Request.Body.Read(buf)
	fmt.Println(string(buf[:n]))
	newVsn := mod.NewVsn()
	err := json.Unmarshal(buf[:n], newVsn)
	if err != nil {
		logger.Logger.Errorf("add vsn unmarshal err:%v", err)
		c.JSON(http.StatusOK, respone.Fail(respone.ParamsError, map[string]string{
			"message": string(buf[:n]),
		}))
		return
	}
	if newVsn.Vsn == "" || newVsn.SrvUrl == "" || newVsn.ResUrl == "" {
		logger.Logger.Errorf("add vsn:%v", newVsn)
		c.JSON(http.StatusOK, respone.Fail(respone.ParamsError, newVsn))
		return
	}
	marshal, _ := json.Marshal(newVsn)
	redis.Client.HSet(CacheVsnKey, newVsn.Vsn, marshal)

	logger.Logger.Infof("add vsn:%v", newVsn)
	c.JSON(http.StatusOK, respone.Success(newVsn))
}

//Update 存储vsn信息
func (vh VsnHandler) Update(c *gin.Context) {
	vsn := c.Query("vsn")

	buf := make([]byte, 1024)
	n, _ := c.Request.Body.Read(buf)
	newVsn := mod.NewVsn()
	err := json.Unmarshal(buf[:n], newVsn)
	if err != nil || vsn == "" || newVsn.SrvUrl == "" || newVsn.ResUrl == "" {
		logger.Logger.Errorf("update vsn:%v", newVsn)
		c.JSON(http.StatusOK, respone.Fail(respone.ParamsError, map[string]string{
			"message": string(buf[:n]),
		}))
		return
	}
	newVsn.Vsn = vsn
	marshal, _ := json.Marshal(newVsn)

	redis.Client.HSet(CacheVsnKey, vsn, marshal)
	logger.Logger.Infof("update vsn:%v", newVsn)
	c.JSON(http.StatusOK, respone.Success(newVsn))
}

//Delete 存储vsn信息
func (vh VsnHandler) Delete(c *gin.Context) {
	vsn := c.Query("vsn")
	if vsn == "" {
		logger.Logger.Errorf("delete vsn:%v", vsn)
		c.JSON(http.StatusOK, respone.Fail(respone.ParamsError, map[string]string{
			"vsn": vsn,
		}))
	} else {
		result, _ := redis.Client.HDel(CacheVsnKey, vsn).Result()

		logger.Logger.Infof("delete vsn:%v", vsn)
		c.JSON(http.StatusOK, respone.Success(map[string]int64{
			"count": result,
		}))
	}
}

//edit gm conf
func (vh VsnHandler) InsertGmConf(c *gin.Context) {
	buf := make([]byte, 1024)
	n, _ := c.Request.Body.Read(buf)
	globalConf := mod.NewGmConf()
	err := json.Unmarshal(buf[:n], globalConf)
	if err != nil || globalConf.GMSrvUrl == "" || globalConf.GMResUrl == "" {
		logger.Logger.Errorf("set global conf err:%v", globalConf)
		c.JSON(http.StatusOK, respone.Fail(respone.ParamsError, map[string]string{
			"message": string(buf[:n]),
		}))
		return
	}
	marshal, _ := json.Marshal(globalConf)
	redis.Client.Set(CacheGMConfKey, marshal, 0)
	logger.Logger.Infof("set global conf:%v", globalConf)
	c.JSON(http.StatusOK, respone.Success(globalConf))
}
