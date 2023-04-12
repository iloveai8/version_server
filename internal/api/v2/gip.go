package v2

import (
	"game_slots_vsn/pkg/e"
	"game_slots_vsn/pkg/ggeoip"
	"github.com/gin-gonic/gin"
)

func GetIP(c *gin.Context) {
	appG := e.Gin{C: c}
	ipStr := c.Query("ip")
	ipInfo := ggeoip.Gip.GetIP(ipStr)
	appG.Success(e.SUCCESS, ipInfo)
	return
}
