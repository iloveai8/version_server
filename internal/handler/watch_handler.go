package handler
//
//import (
//	"context"
//	"oauth-service/internal/service"
//	"oauth-service/pkg/logger"
//	"gitlab.ftsview.com/OpenPlatform/open-go-lib/ftutils"
//	"github.com/gin-gonic/gin"
//)
//
////var cache map[string]*atomic.Uint64
//
////WatchHandler 处理监控Handler
//type WatchHandler struct {
//	service *service.WatchService
//	second  int
//}
//
////NewWatchHandler 创建监控服务，second为0时为同步保存QPS
//func NewWatchHandler(service *service.WatchService, second int) *WatchHandler {
//	//cache = make(map[string]*atomic.Uint64)
//	handler := &WatchHandler{service: service, second: second}
//	//if second > 0 {
//	//	go func() {
//	//		interval := time.Duration(second) * time.Second
//	//		timer := time.NewTimer(interval)
//	//		for range timer.C {
//	//			handler.sync()
//	//			timer.Reset(interval)
//	//		}
//	//	}()
//	//}
//	return handler
//}
//
////Watch 获取Token
//func (h *WatchHandler) Watch(c *gin.Context) {
//	//yyyyMMddHH := time.Now().Format(consts.DateFormatYYYYMMDDHH)
//	//if v, ok := cache[yyyyMMddHH]; ok {
//	//	v.Inc()
//	//} else {
//	//	cache[yyyyMMddHH] = atomic.NewUint64(1)
//	//}
//	uuID := ftutils.UUID()
//	c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), logger.TraceKey, uuID))
//	c.Next()
//}
//
////func (h *WatchHandler) sync() {
////	fmt.Println(cache)
////}
