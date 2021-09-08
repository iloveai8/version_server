package handler
//
//import (
//	"net/http"
//	"oauth-service/internal/service"
//	"oauth-service/pkg/logger"
//
//	"github.com/gin-gonic/gin"
//)
//
////TokenHandler 操作Token的Handler
//type TokenHandler struct {
//	BaseHandler
//	service *service.TokenService
//}
//
////NewTokenHandler 创建Handler
//func NewTokenHandler(service *service.TokenService) *TokenHandler {
//	return &TokenHandler{service: service}
//}
//
////GetToken 获取Token
//func (h *TokenHandler) GetToken(c *gin.Context) {
//	log := logger.FromCtxTracer(c.Request.Context()).WithField("func", "TokenHandler.GetToken")
//	userName := c.Param("userName")
//	passWord := c.Param("passWord")
//	log.Infof("get token param userName:[%s]，passWord:[%s]", userName, passWord)
//	if accessToken, err := h.service.GetToken(c.Request.Context(), userName, passWord); err != nil {
//		log.Warn("get token error", err.Error())
//		c.JSON(http.StatusOK, map[string]string{"message": err.Error()})
//	} else {
//		log.Infof("get token success Token:[%s]", accessToken)
//		c.JSON(http.StatusOK, map[string]string{"token_type": "Bearer", "access_token": accessToken})
//	}
//}
//
////VerifyToken 验证Token
//func (h *TokenHandler) VerifyToken(c *gin.Context) {
//	//fmt.Println(ftutils.Object2JSONNoError(c.Request.URL))
//	log := logger.FromCtxTracer(c.Request.Context()).WithField("func", "TokenHandler.VerifyToken")
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
