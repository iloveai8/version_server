package gmip

import (
	v2 "game_slots_vsn/internal/api/v2"
	"github.com/gin-gonic/gin"
)

func Register(e *gin.Engine) {
	g1 := e.Group("/v2/gmIP")
	{
		g1.GET("", v2.GetGmIPList) //获取所有GMIP
		g1.PUT("", v2.AddGmIP)     //根据增加IP
		g1.DELETE("", v2.RemGmIP)  //删除ip
	}
	g2 := e.Group("/:env/v2/gmIP")
	{
		g2.GET("", v2.GetGmIPList) //获取所有GMIP
		g2.PUT("", v2.AddGmIP)     //根据增加IP
		g2.DELETE("", v2.RemGmIP)  //删除ip
	}
}
