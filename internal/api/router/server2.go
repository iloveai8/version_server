package router

import (
	`game_slots_vsn/internal/controller`
	`github.com/gin-gonic/gin`
)

func RServer2(e *gin.Engine) {
	v2 := e.Group("/v2")
	v2s := v2.Group("server")
	{
		v2s.GET("list", controller.Server().GetServerList) //获取所有用户
		v2s.GET("", controller.Server().GetServer)         //根据id获取用户
		v2s.POST("", controller.Server().AddServer)        //保存新用户
		v2s.PUT("", controller.Server().UpdateServer)      //根据id更新用户
		v2s.DELETE("", controller.Server().DeleteServer)   //根据id删除用户
	}
	v2g := v2.Group("gm")
	{
		v2g.PUT("", controller.GM().AddGM) //保存GM开关配置
	}

	v2e := e.Group("/:env/v2")
	v2es := v2e.Group("server")
	{
		v2es.GET("list", controller.Server().GetServerList) //获取所有用户
		v2es.GET("", controller.Server().GetServer)         //根据id获取用户
		v2es.POST("", controller.Server().AddServer)        //保存新用户
		v2es.PUT("", controller.Server().UpdateServer)      //根据id更新用户
		v2es.DELETE("", controller.Server().DeleteServer)   //根据id删除用户
	}
	v2eg := v2e.Group("/gm")
	{
		v2eg.PUT("", controller.GM().AddGM) //保存GM开关配置
	}
}
