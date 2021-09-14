package handler

import (
	`encoding/json`
	`fmt`
	`game_slots_vsn/internal/mod`
	`game_slots_vsn/internal/respone`
	`github.com/gin-gonic/gin`
	`github.com/go-redis/redis`
	`net/http`
)

const CacheVsnKey string = "vsn."

type VsnHandler struct {
	Client *redis.Client
}

func NewVsnHandler(redis *redis.Client) *VsnHandler {
	return &VsnHandler{Client: redis}
}

//GetAll 获取所有vsn 信息
func (vh VsnHandler) GetAll(c *gin.Context) {
	result, _ := vh.Client.HVals(CacheVsnKey).Result()

	//eg1
	//var vsnInfoList []*mod.Vsn
	//for _, vsnStr := range result {
	//	vsn := mod.NewVsn()
	//	_ = json.Unmarshal([]byte(vsnStr), vsn)
	//	vsnInfoList = append(vsnInfoList, vsn)
	//}
	//eg2
	vsnInfoList := make([]*mod.Vsn, len(result), cap(result))
	for i, vsnStr := range result {
		vsn := mod.NewVsn()
		_ = json.Unmarshal([]byte(vsnStr), vsn)
		vsnInfoList[i] = vsn
	}

	c.JSON(http.StatusOK, respone.Success(vsnInfoList))
}

//Get 根据vsn获取信息
func (vh VsnHandler) Get(c *gin.Context) {
	vsn := c.Query("vsn")
	if vsn == "" {
		c.JSON(http.StatusOK, respone.Fail(respone.ParamsError, map[string]string{
			"vsn": vsn,
		}))
		return
	}
	result, _ := vh.Client.HGet(CacheVsnKey, vsn).Result()
	vsnInfo := mod.NewVsn()
	_ = json.Unmarshal([]byte(result), vsnInfo)
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
		c.JSON(http.StatusOK, respone.Fail(respone.ParamsError, map[string]string{
			"message": string(buf[:n]),
		}))
	}
	if newVsn.Vsn == "" || newVsn.SrvUrl == "" || newVsn.ResUrl == "" {
		c.JSON(http.StatusOK, respone.Fail(respone.ParamsError, newVsn))
		return
	}
	marshal, _ := json.Marshal(newVsn)
	vh.Client.HSet(CacheVsnKey, newVsn.Vsn, marshal)
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
		c.JSON(http.StatusOK, respone.Fail(respone.ParamsError, map[string]string{
			"message": string(buf[:n]),
		}))
	}
	newVsn.Vsn = vsn
	marshal, _ := json.Marshal(newVsn)
	vh.Client.HSet(CacheVsnKey, vsn, marshal)
	c.JSON(http.StatusOK, respone.Success(newVsn))
}

//Delete 存储vsn信息
func (vh VsnHandler) Delete(c *gin.Context) {
	vsn := c.Query("vsn")
	if vsn == "" {
		c.JSON(http.StatusOK, respone.Fail(respone.ParamsError, map[string]string{
			"vsn": vsn,
		}))
	} else {
		result, _ := vh.Client.HDel(CacheVsnKey, vsn).Result()
		c.JSON(http.StatusOK, respone.Success(map[string]int64{
			"count": result,
		}))
	}
}
