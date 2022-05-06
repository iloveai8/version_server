package api

import (
	`game_slots_vsn/internal/api/router`
	`github.com/gin-gonic/gin`
)

func Init(e *gin.Engine) {
	router.ROther(e)
	router.RServer1(e)
	router.RServer2(e)
}
