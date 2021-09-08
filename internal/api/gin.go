package api
//
//import (
//	"github.com/gin-gonic/gin"
//	`log`
//	"net/http"
//)
//
////Start 启动Gin服务
//func Start() {
//	//addr := fmt.Sprintf("%s:%s", conf.GlobalConfig.ServiceIP, conf.GlobalConfig.ServicePort)
//	router := initRouter()
//	router.GET("/ping", func(c *gin.Context) {
//		c.JSON(http.StatusOK, gin.H{
//			"message": "pong",
//		})
//	})
//	vsnGroup := router.Group("/vsn")
//
//	vsnGroup.POST()
//
//	if err := router.Run(addr); err != nil {
//
//	}
//}
//
//func initRouter() *gin.Engine {
//	gin.SetMode(gin.ReleaseMode)
//	router := gin.Default()
//	vsnGroup := router.Group("/vsn")
//	vsnGroup.GET("/version", func(c *gin.Context) {
//		c.JSON(http.StatusOK, map[string]string{"version": conf.GlobalConfig.ImageVersion})
//	})
//	return router
//}
