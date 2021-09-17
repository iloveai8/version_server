package api

import (
	`fmt`
	"game_slots_vsn/internal/handler"
	`game_slots_vsn/pkg/config`
	"game_slots_vsn/pkg/logger"
	`github.com/gin-gonic/gin`
	"net/http"
)

func StartServer(web *config.WebConfig) {
	var addr string
	if web.IP == "" {
		addr = fmt.Sprintf(":%d", web.Port)
	} else {
		addr = fmt.Sprintf("%s:%d", web.IP, web.Port)
	}

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(
		logger.GinLogger(),
		gin.Recovery(),
	)
	router.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "Not router")
	})
	router.GET("/favicon.ico", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	vsnHandler := handler.NewVsnHandler()
	vsn := router.Group("v1")
	{
		vsn.GET("/vsn/list", vsnHandler.GetAll) //获取所有用户
		vsn.GET("/vsn", vsnHandler.Get)         //根据id获取用户
		vsn.POST("/vsn", vsnHandler.Insert)     //保存新用户
		vsn.PUT("/vsn", vsnHandler.Update)      //根据id更新用户
		vsn.DELETE("/vsn", vsnHandler.Delete)   //根据id删除用户
	}
	if err := router.Run(addr); err != nil {
		logger.Logger.Info("service start fail.")
	}
}
