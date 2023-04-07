package server

import (
	v2 "game_slots_vsn/internal/api/v2"
	"github.com/gin-gonic/gin"
)

func Register(e *gin.Engine) {
	g1 := e.Group("/v2/server")
	{
		g1.GET("list", v2.GetServerInfoList) //获取所有用户
		g1.GET("", v2.GetServerInfo)         //根据id获取用户
		g1.POST("", v2.AddServerInfo)        //保存新用户
		g1.PUT("", v2.UpdateServerInfo)      //根据id更新用户
		g1.DELETE("", v2.DeleteServerInfo)   //根据id删除用户
	}
	ge1 := e.Group("/:env/v2/server")
	{
		ge1.GET("list", v2.GetServerInfoList) //获取所有用户
		ge1.GET("", v2.GetServerInfo)         //根据id获取用户
		ge1.POST("", v2.AddServerInfo)        //保存新用户
		ge1.PUT("", v2.UpdateServerInfo)      //根据id更新用户
		ge1.DELETE("", v2.DeleteServerInfo)   //根据id删除用户
	}
}
