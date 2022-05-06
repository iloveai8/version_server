package router

import (
	`github.com/gin-gonic/gin`
	`net/http`
	`os`
	`time`
)

func ROther(router *gin.Engine) {
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
}
