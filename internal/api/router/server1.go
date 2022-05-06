package router

import (
	`game_slots_vsn/internal/controller/v1`
	`github.com/gin-gonic/gin`
)

func RServer1(e *gin.Engine) {
	vsn := e.Group("/v1")
	{
		vsn.GET("/vsn/list", v1.Server1().GetAll) //获取所有用户
		vsn.GET("/vsn", v1.Server1().Get)         //根据id获取用户
		vsn.POST("/vsn", v1.Server1().Insert)     //保存新用户
		vsn.PUT("/vsn", v1.Server1().Update)      //根据id更新用户
		vsn.DELETE("/vsn", v1.Server1().Delete)   //根据id删除用户

		vsn.POST("/vsn/gmConf", v1.Server1().InsertGmConf) //保存GM开关配置
	}
	vsn1 := e.Group("/:env/v1")
	{
		vsn1.GET("/vsn/list", v1.Server1().GetAll) //获取所有用户
		vsn1.GET("/vsn", v1.Server1().Get)         //根据id获取用户
		vsn1.POST("/vsn", v1.Server1().Insert)     //保存新用户
		vsn1.PUT("/vsn", v1.Server1().Update)      //根据id更新用户
		vsn1.DELETE("/vsn", v1.Server1().Delete)   //根据id删除用户

		vsn1.POST("/vsn/gmConf", v1.Server1().InsertGmConf) //保存GM开关配置
	}
}
