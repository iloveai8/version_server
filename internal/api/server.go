package api

import (
	"fmt"
	`game_slots_vsn/internal/consts`
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	`net/http`
	`os`
)

func ServerStart() {
	gin.SetMode(gin.ReleaseMode)
	addr := fmt.Sprintf("%s:%s", viper.GetString(consts.AppServerIp), viper.GetString(consts.AppServerPort))
	fmt.Println("add:", addr)
	router := setupRouter()
	if err := router.Run(addr); err != nil {
		//logger.Logger.Error(consts.MessageServiceFail, err)
		os.Exit(1)
	}
}

//配置路由
func setupRouter() *gin.Engine {
	router := gin.Default()
	v1 := router.Group("/vsn")
	{
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})
	}
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	return router
}
