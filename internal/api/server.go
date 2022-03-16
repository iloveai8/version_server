package api

import (
	`fmt`
	"game_slots_vsn/internal/handler"
	`game_slots_vsn/pkg/config`
	"game_slots_vsn/pkg/logger"
	ginzap "github.com/gin-contrib/zap"
	`github.com/gin-gonic/gin`
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
	router.GET("/akamai", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})
	router.GET("/:env/akamai", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})
	router.GET("/favicon.ico", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})
	router.GET("/:env/favicon.ico", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
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
	vsnHandler := handler.NewVsnHandler()
	vsn := router.Group("/v1")
	{
		vsn.GET("/vsn/list", vsnHandler.GetAll) //获取所有用户
		vsn.GET("/vsn", vsnHandler.Get)         //根据id获取用户
		vsn.POST("/vsn", vsnHandler.Insert)     //保存新用户
		vsn.PUT("/vsn", vsnHandler.Update)      //根据id更新用户
		vsn.DELETE("/vsn", vsnHandler.Delete)   //根据id删除用户

		vsn.POST("/vsn/gmConf", vsnHandler.InsertGmConf) //保存GM开关配置
	}
	vsn1 := router.Group("/:env/v1")
	{
		vsn1.GET("/vsn/list", vsnHandler.GetAll) //获取所有用户
		vsn1.GET("/vsn", vsnHandler.Get)         //根据id获取用户
		vsn1.POST("/vsn", vsnHandler.Insert)     //保存新用户
		vsn1.PUT("/vsn", vsnHandler.Update)      //根据id更新用户
		vsn1.DELETE("/vsn", vsnHandler.Delete)   //根据id删除用户

		vsn1.POST("/vsn/gmConf", vsnHandler.InsertGmConf) //保存GM开关配置
	}
	if err := router.Run(addr); err != nil {
		logger.Logger.Info("service start fail.", err)
		os.Exit(0)
	}
}
