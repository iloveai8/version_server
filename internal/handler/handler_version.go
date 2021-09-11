package handler

import (
	`encoding/json`
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
	//var vsnList []*mod.Vsn
	//for _, vsnStr := range result {
	//	vsn := mod.NewVsn()
	//	_ = json.Unmarshal([]byte(vsnStr), vsn)
	//	vsnList = append(vsnList, vsn)
	//}
	//eg2
	vsnList1 := make([]*mod.Vsn, len(result), cap(result))
	for i, vsnStr := range result {
		vsn := mod.NewVsn()
		_ = json.Unmarshal([]byte(vsnStr), vsn)
		vsnList1[i] = vsn
	}

	c.JSON(http.StatusOK, respone.Success(vsnList))
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
	vsn := c.PostForm("vsn")
	srvUrl := c.PostForm("srvUrl")
	resUrl := c.PostForm("resUrl")
	//c.MustBindWith()
	//c.BindJSON(mod.Vsn{})
	if vsn == "" || srvUrl == "" || resUrl == "" {
		c.JSON(http.StatusOK, respone.Fail(respone.ParamsError, map[string]string{
			"vsn":    vsn,
			"srvUrl": srvUrl,
			"resUrl": resUrl,
		}))
		return
	}
	newVsn := mod.NewVsn()
	newVsn.Vsn = vsn
	newVsn.SrvUrl = srvUrl
	newVsn.ResUrl = resUrl
	marshal, _ := json.Marshal(newVsn)
	vh.Client.HSet(CacheVsnKey, vsn, marshal)
	c.JSON(http.StatusOK, respone.Success(newVsn))
}

//Update 存储vsn信息
func (vh VsnHandler) Update(c *gin.Context) {
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

	srvUrl := c.DefaultPostForm("srvUrl", vsnInfo.SrvUrl)
	resUrl := c.DefaultPostForm("resUrl", vsnInfo.ResUrl)

	vsnInfo.SrvUrl = srvUrl
	vsnInfo.ResUrl = resUrl
	marshal, _ := json.Marshal(vsnInfo)
	vh.Client.HSet(CacheVsnKey, vsn, marshal)
	c.JSON(http.StatusOK, respone.Success(vsnInfo))
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
