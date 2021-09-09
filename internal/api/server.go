package api

import (
	`fmt`
	`game_slots_vsn/internal/config`
	`github.com/gin-gonic/gin`
	`log`
	`net/http`
)

//web start
func ServerStart() {
	addr := fmt.Sprintf("%s:%d", config.C.Server.IP, config.C.Server.Port)
	router := generateRouter()
	if err := router.Run(addr); err != nil {
		log.Fatal("server start:", err)
	}
}

func generateRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.GET("/heartbeat", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	return router
}
