package handler

import (
	`game_slots_vsn/internal/db/redis`
	`github.com/gin-gonic/gin`
	`net/http`
)

//VersionHandler 操作Token的Handler
type VersionHandler struct {
}

//NewVersionHandler 创建Handler
func NewVersionHandler() *VersionHandler {
	return &VersionHandler{}
}

//Get 获取ServerInfo
func (vsn VersionHandler)Get(c *gin.Context) {
	redis.GetClient().Set("aa", "hello", 0)

	c.JSON(http.StatusOK, gin.H{
		"message":redis.GetClient().Get("aa"),
	})
	//redis.GetClient().Get("aa")
	//userName := c.Param("userName")
	//passWord := c.Param("passWord")

	//if accessToken, err := h.service.GetToken(c.Request.Context(), userName, passWord); err != nil {
	//	//log.Warn("get token error", err.Error())
	//	c.JSON(http.StatusOK, map[string]string{"message": err.Error()})
	//} else {
	//	log.Infof("get token success Token:[%s]", accessToken)
	//	c.JSON(http.StatusOK, map[string]string{"token_type": "Bearer", "access_token": accessToken})
	//}
}
func (vsn *VersionHandler) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message":"pong11",
	})
	//userName := c.Param("userName")
	//passWord := c.Param("passWord")

	//if accessToken, err := h.service.GetToken(c.Request.Context(), userName, passWord); err != nil {
	//	//log.Warn("get token error", err.Error())
	//	c.JSON(http.StatusOK, map[string]string{"message": err.Error()})
	//} else {
	//	log.Infof("get token success Token:[%s]", accessToken)
	//	c.JSON(http.StatusOK, map[string]string{"token_type": "Bearer", "access_token": accessToken})
	//}
}

////VerifyToken 验证Token
//func (h *VersionHandler) VerifyToken(c *gin.Context) {
//	//fmt.Println(ftutils.Object2JSONNoError(c.Request.URL))
//	log := logger.FromCtxTracer(c.Request.Context()).WithField("func", "VersionHandler.VerifyToken")
//	accessToken := c.Param("access_token")
//	log.Infof("verify token param access_token:[%s]", accessToken)
//	if claims, err := h.service.VerifyToken(c.Request.Context(), accessToken); err != nil {
//		log.Warn("verify token error", err.Error())
//		c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
//	} else {
//		log.Infof("get token success Claims:[%+v]", claims)
//		c.JSON(http.StatusOK, claims)
//	}
//}
