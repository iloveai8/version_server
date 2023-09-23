package gmip

import (
	v2 "game_slots_vsn/internal/api/v2"
	"github.com/gin-gonic/gin"
)

func Register(e *gin.Engine) {
	g1 := e.Group("/v2/gmip")
	{
		g1.GET("/list", v2.List)  //获取所有GMIP
		g1.GET("/add", v2.Add)    //根据增加IP
		g1.DELETE("/rem", v2.Rem) //删除ip
	}
	g2 := e.Group("/:env/v2/gmip")
	{
		g2.GET("/list", v2.List)  //获取所有GMIP
		g2.GET("/add", v2.Add)    //根据增加IP
		g2.DELETE("/rem", v2.Rem) //删除ip
	}
}
