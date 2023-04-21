package gip

import (
	v2 "game_slots_vsn/internal/api/v2"
	"github.com/gin-gonic/gin"
)

func Register(e *gin.Engine) {
	g1 := e.Group("/v2/ip")
	{
		g1.GET("/c", v2.GetCountry)         //根据id获取用户
		g1.GET("/cc", v2.GetCountryAndCity) //根据id获取用户
	}
	ge1 := e.Group("/:env/v2/ip")
	{
		ge1.GET("/c", v2.GetCountry)         //根据id获取用户
		ge1.GET("/cc", v2.GetCountryAndCity) //根据id获取用户
	}
}
