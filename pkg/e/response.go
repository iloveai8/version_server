package e

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Gin struct {
	C *gin.Context
}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

func (g *Gin) Success(errCode int, data interface{}) {
	g.C.JSON(http.StatusOK, Response{
		Code:    errCode,
		Message: GetMsg(errCode),
		Data:    data,
	})
	return
}

func (g *Gin) Fail(errCode int, error interface{}) {
	g.C.JSON(http.StatusOK, Response{
		Code:    errCode,
		Message: GetMsg(errCode),
		Errors:  error,
	})
	return
}
