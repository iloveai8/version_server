package api

import (
	`errors`
	`fmt`
	`game_slots_vsn/internal/config`
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

	router.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "Not router")
	})
	router.GET("/favicon.ico", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})
	vsnHandler := handler.NewVersionHandler()
	vsn := router.Group("/vsn")
	{
		vsn.GET("/all", vsnHandler.Get)
		vsn.GET("/get", vsnHandler.Get)
		vsn.DELETE("/del", vsnHandler.Get)
		vsn.POST("/update", vsnHandler.Get)
	}
	return router
}
