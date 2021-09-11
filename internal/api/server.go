package api

import (
	`errors`
	`fmt`
	`game_slots_vsn/internal/config`
	`game_slots_vsn/internal/db/redis`
	`game_slots_vsn/internal/handler`
	`github.com/gin-gonic/gin`
	`net/http`
)

var (
	ErrInitServer = errors.New("server start error")
)

func Init(sConfig *config.Config) {
	addr := fmt.Sprintf("%s:%d", sConfig.Server.IP, sConfig.Server.Port)
	router := generateRouter()
	err := router.Run(addr)
	if err != nil {
		panic(ErrInitServer)
	}
}

func generateRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()
	defaultRouter(router)
	vsnRouter(router)
	return router
}

func defaultRouter(router *gin.Engine) {
	router.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "Not router")
	})
	router.GET("/favicon.ico", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})
}

func vsnRouter(router *gin.Engine) {
	vsnHandler := handler.NewVsnHandler(redis.GetClient())
	vsn := router.Group("v1")
	{
		vsn.GET("/vsn/list", vsnHandler.GetAll) //获取所有用户
		vsn.GET("/vsn", vsnHandler.Get)         //根据id获取用户
		vsn.POST("/vsn", vsnHandler.Insert)     //保存新用户
		vsn.PUT("/vsn", vsnHandler.Update)      //根据id更新用户
		vsn.DELETE("/vsn", vsnHandler.Delete)   //根据id删除用户
	}
}
