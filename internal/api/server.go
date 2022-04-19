package api

import (
	"fmt"
	`game_slots_vsn/internal/controller`
	"game_slots_vsn/pkg/config"
	"game_slots_vsn/pkg/logger"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"time"
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
	router.Use(ginzap.Ginzap(logger.Logger.Desugar(), time.RFC3339, true))
	// Logs all panic to error log
	//   - stack means whether output the stack info.
	router.Use(ginzap.RecoveryWithZap(logger.Logger.Desugar(), true))

	router.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "Not router")
	})
	router.LoadHTMLGlob("static/*")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{"env": os.Getenv("env"), "ver": os.Getenv("ver")})
	})
	router.GET("/:env", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{"env": os.Getenv("env"), "ver": os.Getenv("ver")})
	})
	router.GET("/akamai", func(c *gin.Context) {
		c.HTML(http.StatusOK, "akamai.html", gin.H{})
	})
	router.GET("/:env/akamai", func(c *gin.Context) {
		c.HTML(http.StatusOK, "akamai.html", gin.H{})
	})
	router.GET("/heartbeat", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"vsn":  os.Getenv("ver"),
			"env":  os.Getenv("env"),
			"time": time.Now(),
		})
	})
	router.GET("/:env/heartbeat", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"vsn":  os.Getenv("ver"),
			"env":  os.Getenv("env"),
			"time": time.Now(),
		})
	})
	//vsnHandler := handler.NewVsnHandler()
	//vsn := router.Group("/v2")
	//{
	//	vsn.GET("/vsn/list", vsnHandler.GetServerList) //获取所有用户
	//	vsn.GET("/vsn", vsnHandler.GetVersionByKey)         //根据id获取用户
	//	vsn.POST("/vsn", vsnHandler.Insert)     //保存新用户
	//	vsn.PUT("/vsn", vsnHandler.Update)      //根据id更新用户
	//	vsn.DELETE("/vsn", vsnHandler.Delete)   //根据id删除用户
	//
	//	vsn.POST("/vsn/gmConf", vsnHandler.InsertGmConf) //保存GM开关配置
	//
	//	vsn.GET("/vsn/list", vsnHandler.GetServerList) //获取所有用户
	//
	//}
	//vsn1 := router.Group("/:env/v2")
	//{
	//	vsn1.GET("/vsn/list", vsnHandler.GetServerList) //获取所有用户
	//	vsn1.GET("/vsn", vsnHandler.GetVersionByKey)         //根据id获取用户
	//	vsn1.POST("/vsn", vsnHandler.Insert)     //保存新用户
	//	vsn1.PUT("/vsn", vsnHandler.Update)      //根据id更新用户
	//	vsn1.DELETE("/vsn", vsnHandler.Delete)   //根据id删除用户
	//
	//	vsn1.POST("/vsn/gmConf", vsnHandler.InsertGmConf) //保存GM开关配置
	//}

	v2 := router.Group("/v2")
	v21 := v2.Group("server")
	{
		v21.GET("list", controller.Server().GetServerList) //获取所有用户
		v21.GET("", controller.Server().GetServer)         //根据id获取用户
		v21.POST("", controller.Server().AddServer)        //保存新用户
		v21.PUT("", controller.Server().UpdateServer)      //根据id更新用户
		v21.DELETE("", controller.Server().DeleteServer)   //根据id删除用户
	}
	v22 := v2.Group("gm")
	{
		v22.PUT("", controller.GM().AddGM) //保存GM开关配置
	}

	if err := router.Run(addr); err != nil {
		logger.Logger.Info("service start fail.", err)
		os.Exit(0)
	}
}
