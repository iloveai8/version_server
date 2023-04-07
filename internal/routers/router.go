package routers

import (
	"game_slots_vsn/internal/routers/gm"
	"game_slots_vsn/internal/routers/server"
	"game_slots_vsn/pkg/logger"
	"game_slots_vsn/pkg/setting"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"time"
)

type Router interface {
	register(e *gin.Engine)
}

func Register() *gin.Engine {
	gin.SetMode(setting.SrvSetting.RunMode)

	e := gin.New()
	_ = e.SetTrustedProxies([]string{"127.0.0.1", "localhost"})
	//e.Use(logger.GinZap())
	e.Use(logger.GinZapWithSkipPaths())
	e.Use(logger.RecoverZap())

	e.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "Not engine")
	})
	e.LoadHTMLGlob("static/*")
	e.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{"env": os.Getenv("env"), "ver": os.Getenv("ver")})
	})
	e.GET("/akamai", func(c *gin.Context) {
		c.HTML(http.StatusOK, "akamai.html", gin.H{})
	})
	e.GET("/:env/akamai", func(c *gin.Context) {
		c.HTML(http.StatusOK, "akamai.html", gin.H{})
	})
	e.GET("/heartbeat", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"vsn":  os.Getenv("ver"),
			"env":  os.Getenv("env"),
			"time": time.Now(),
		})
	})
	e.GET("/:env/heartbeat", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"vsn":  os.Getenv("ver"),
			"env":  os.Getenv("env"),
			"time": time.Now(),
		})
	})
	gm.Register(e)
	server.Register(e)

	return e
}
