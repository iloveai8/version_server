package handler

import (
	`encoding/json`
	`fmt`
	`game_slots_vsn/internal/mod`
	`game_slots_vsn/internal/respone`
	"game_slots_vsn/pkg/logger"
	"game_slots_vsn/pkg/redis"
	`github.com/gin-gonic/gin`
	"go.uber.org/zap"
	`net/http`
)

const CacheVsnKey string = "vsn."

type VsnHandler struct {
}

func NewVsnHandler() *VsnHandler {
	return &VsnHandler{}
}

//GetAll 获取所有vsn 信息
func (vh VsnHandler) GetAll(c *gin.Context) {
	result, _ := redis.Client.HVals(CacheVsnKey).Result()

	vsnInfoList := make([]*mod.Vsn, len(result), cap(result))
	for i, vsnStr := range result {
		vsn := mod.NewVsn()
		_ = json.Unmarshal([]byte(vsnStr), vsn)
		vsnInfoList[i] = vsn
	}
	logger.Logger.Info("vsn list:", vsnInfoList)
	c.JSON(http.StatusOK, respone.Success(vsnInfoList))
}

//Get 根据vsn获取信息
func (vh VsnHandler) Get(c *gin.Context) {
	vsn := c.Query("vsn")
	if vsn == "" {
		logger.Logger.Error("get vsn err:", vsn)
		c.JSON(http.StatusOK, respone.Fail(respone.ParamsError, map[string]string{
			"vsn": vsn,
		}))
		return
	}
	result, err := redis.Client.HGet(CacheVsnKey, vsn).Result()
	if err != nil {
		logger.Logger.Error("get vsn err:", err)
		c.JSON(http.StatusOK, respone.Fail(respone.ResultNotFound, map[string]string{
			"vsn": vsn,
		}))
		return
	}
	vsnInfo := mod.NewVsn()
	_ = json.Unmarshal([]byte(result), vsnInfo)
	logger.Logger.Info(zap.Object("vsn", vsnInfo))
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
		logger.Logger.Error("add vsn unmarshal err:", err)
		c.JSON(http.StatusOK, respone.Fail(respone.ParamsError, map[string]string{
			"message": string(buf[:n]),
		}))
	}
	if newVsn.Vsn == "" || newVsn.SrvUrl == "" || newVsn.ResUrl == "" {
		logger.Logger.Error("add vsn:", newVsn)
		c.JSON(http.StatusOK, respone.Fail(respone.ParamsError, newVsn))
		return
	}
	marshal, _ := json.Marshal(newVsn)
	redis.Client.HSet(CacheVsnKey, newVsn.Vsn, marshal)

	logger.Logger.Info("add vsn:", newVsn)
	c.JSON(http.StatusOK, respone.Success(newVsn))
}

//Update 存储vsn信息
func (vh VsnHandler) Update(c *gin.Context) {
	vsn := c.Query("vsn")

	buf := make([]byte, 1024)
	n, _ := c.Request.Body.Read(buf)
	fmt.Println(string(buf[:n]))
	newVsn := mod.NewVsn()
	err := json.Unmarshal(buf[:n], newVsn)
	if err != nil || vsn == "" || newVsn.SrvUrl == "" || newVsn.ResUrl == "" {
		logger.Logger.Error("update vsn:", newVsn)
		c.JSON(http.StatusOK, respone.Fail(respone.ParamsError, map[string]string{
			"message": string(buf[:n]),
		}))
	}
	newVsn.Vsn = vsn
	marshal, _ := json.Marshal(newVsn)

	redis.Client.HSet(CacheVsnKey, vsn, marshal)
	logger.Logger.Info("update vsn:", newVsn)
	c.JSON(http.StatusOK, respone.Success(newVsn))
}

//Delete 存储vsn信息
func (vh VsnHandler) Delete(c *gin.Context) {
	vsn := c.Query("vsn")
	if vsn == "" {
		logger.Logger.Error("delete vsn:", vsn)
		c.JSON(http.StatusOK, respone.Fail(respone.ParamsError, map[string]string{
			"vsn": vsn,
		}))
	} else {
		result, _ := redis.Client.HDel(CacheVsnKey, vsn).Result()

		logger.Logger.Info("delete vsn:", vsn)
		c.JSON(http.StatusOK, respone.Success(map[string]int64{
			"count": result,
		}))
	}
}
