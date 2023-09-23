package v2

import (
	"game_slots_vsn/internal/service"
	"game_slots_vsn/pkg/e"
	"github.com/gin-gonic/gin"
)

func GetGmIPList(c *gin.Context) {
	appG := e.Gin{C: c}
	ipService := &service.GmIPService{}
	ipList := ipService.GetGmIPList()
	appG.Success(e.SUCCESS, map[string]interface{}{
		"ipList": ipList,
	})
	return
}

func AddGmIP(c *gin.Context) {
	appG := e.Gin{C: c}
	ip := c.Query("gmIP")
	ipService := &service.GmIPService{}
	add := ipService.AddGmIP(ip)
	appG.Success(e.SUCCESS, add)
	return
}
func RemGmIP(c *gin.Context) {
	appG := e.Gin{C: c}
	ip := c.Query("gmIP")
	ipService := &service.GmIPService{}
	rem := ipService.RemGmIP(ip)
	appG.Success(e.SUCCESS, rem)
	return
}
