package gm

import (
	v2 "game_slots_vsn/internal/api/v2"
	"github.com/gin-gonic/gin"
)

func Register(e *gin.Engine) {
	g1 := e.Group("/v2/gm")
	{
		g1.GET("", v2.GetGmInfo)
		g1.PUT("", v2.UpdateGmInfo)
	}
	g2 := e.Group("/:env/v2/gm")
	{
		g2.GET("", v2.GetGmInfo)
		g2.PUT("", v2.UpdateGmInfo)
	}
}
