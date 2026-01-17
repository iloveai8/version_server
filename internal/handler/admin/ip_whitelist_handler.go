// Package admin 提供管理后台API的HTTP处理器
//
// 该包实现了GM后台访问的管理接口
package admin

import (
	"github.com/gin-gonic/gin"

	"game_slots_vsn/internal/domain"
	"game_slots_vsn/internal/handler"
	"game_slots_vsn/internal/service"
)

// IPWhitelistHandler IP白名单处理器
//
// 处理GM后台的IP白名单请求
type IPWhitelistHandler struct {
	whitelistService service.IPWhitelistService
}

// NewIPWhitelistHandler 创建IP白名单处理器
//
// 参数:
//   whitelistService: IP白名单服务
//
// 返回:
//   *IPWhitelistHandler: IP白名单处理器
func NewIPWhitelistHandler(whitelistService service.IPWhitelistService) *IPWhitelistHandler {
	return &IPWhitelistHandler{
		whitelistService: whitelistService,
	}
}

// GetWhitelist 获取IP白名单
//
// @Summary 获取IP白名单
// @Description 获取当前的IP白名单
// @Tags IPWhitelist
// @Produce json
// @Success 200 {object} map[string]interface{} "data": {...}
// @Router /admin/v1/ip-whitelist [get]
func (h *IPWhitelistHandler) GetWhitelist(c *gin.Context) {
	// 调用Service层
	whitelist, err := h.whitelistService.GetWhitelist(c.Request.Context())
	if err != nil {
		handler.RespondWithError(c, err)
		return
	}

	handler.RespondWithSuccess(c, whitelist)
}

// UpdateWhitelist 更新IP白名单
//
// @Summary 更新IP白名单
// @Description 更新IP白名单
// @Tags IPWhitelist
// @Accept json
// @Produce json
// @Param body body domain.IPWhitelist true "IP白名单"
// @Success 200 {object} map[string]interface{} "data": {...}
// @Failure 400 {object} map[string]interface{} "error": {...}
// @Router /admin/v1/ip-whitelist [put]
func (h *IPWhitelistHandler) UpdateWhitelist(c *gin.Context) {
	var whitelist domain.IPWhitelist
	if !handler.ShouldBindJSON(c, &whitelist) {
		return
	}

	// 调用Service层
	if err := h.whitelistService.UpdateWhitelist(c.Request.Context(), &whitelist); err != nil {
		handler.RespondWithError(c, err)
		return
	}

	handler.RespondWithSuccess(c, gin.H{
		"count": whitelist.Count(),
		"ips":   whitelist.GetIPs(),
	})
}

// AddIP 添加IP到白名单
//
// @Summary 添加IP
// @Description 添加IP到白名单
// @Tags IPWhitelist
// @Accept json
// @Produce json
// @Param body body object{ip string} true "IP地址"
// @Success 200 {object} map[string]interface{} "data": {...}
// @Failure 400 {object} map[string]interface{} "error": {...}
// @Router /admin/v1/ip-whitelist/ip [post]
func (h *IPWhitelistHandler) AddIP(c *gin.Context) {
	var req handler.AddIPRequest
	if !handler.ShouldBindJSON(c, &req) {
		return
	}

	// 调用Service层
	if err := h.whitelistService.AddIP(c.Request.Context(), req.IP); err != nil {
		handler.RespondWithError(c, err)
		return
	}

	handler.RespondWithSuccess(c, gin.H{
		"ip": req.IP,
	})
}

// RemoveIP 从白名单移除IP
//
// @Summary 移除IP
// @Description 从白名单移除IP
// @Tags IPWhitelist
// @Accept json
// @Produce json
// @Param body body object{ip string} true "IP地址"
// @Success 204
// @Failure 400 {object} map[string]interface{} "error": {...}
// @Router /admin/v1/ip-whitelist/ip [delete]
func (h *IPWhitelistHandler) RemoveIP(c *gin.Context) {
	var req handler.RemoveIPRequest
	if !handler.ShouldBindJSON(c, &req) {
		return
	}

	// 调用Service层
	if err := h.whitelistService.RemoveIP(c.Request.Context(), req.IP); err != nil {
		handler.RespondWithError(c, err)
		return
	}

	handler.RespondWithNoContent(c)
}

// CheckIP 检查IP是否在白名单中
//
// @Summary 检查IP
// @Description 检查IP是否在白名单中
// @Tags IPWhitelist
// @Accept json
// @Produce json
// @Param body body object{ip string} true "IP地址"
// @Success 200 {object} map[string]interface{} "data": {...}
// @Failure 400 {object} map[string]interface{} "error": {...}
// @Router /admin/v1/ip-whitelist/check [post]
func (h *IPWhitelistHandler) CheckIP(c *gin.Context) {
	var req handler.AddIPRequest
	if !handler.ShouldBindJSON(c, &req) {
		return
	}

	// 调用Service层
	contains, err := h.whitelistService.CheckIP(c.Request.Context(), req.IP)
	if err != nil {
		handler.RespondWithError(c, err)
		return
	}

	handler.RespondWithSuccess(c, gin.H{
		"contains": contains,
	})
}

// BatchAddIPs 批量添加IP
//
// @Summary 批量添加IP
// @Description 批量添加IP到白名单
// @Tags IPWhitelist
// @Accept json
// @Produce json
// @Param body body object{ips []string} true "IP列表"
// @Success 200 {object} map[string]interface{} "data": {...}
// @Failure 400 {object} map[string]interface{} "error": {...}
// @Router /admin/v1/ip-whitelist/batch [post]
func (h *IPWhitelistHandler) BatchAddIPs(c *gin.Context) {
	var req handler.BatchAddIPsRequest
	if !handler.ShouldBindJSON(c, &req) {
		return
	}

	// 调用Service层
	if err := h.whitelistService.BatchAddIPs(c.Request.Context(), req.IPs); err != nil {
		handler.RespondWithError(c, err)
		return
	}

	handler.RespondWithSuccess(c, gin.H{
		"added": len(req.IPs),
	})
}

// BatchRemoveIPs 批量移除IP
//
// @Summary 批量移除IP
// @Description 批量从白名单移除IP
// @Tags IPWhitelist
// @Accept json
// @Produce json
// @Param body body object{ips []string} true "IP列表"
// @Success 200 {object} map[string]interface{} "data": {...}
// @Failure 400 {object} map[string]interface{} "error": {...}
// @Router /admin/v1/ip-whitelist/batch [delete]
func (h *IPWhitelistHandler) BatchRemoveIPs(c *gin.Context) {
	var req handler.BatchRemoveIPsRequest
	if !handler.ShouldBindJSON(c, &req) {
		return
	}

	// 调用Service层
	if err := h.whitelistService.BatchRemoveIPs(c.Request.Context(), req.IPs); err != nil {
		handler.RespondWithError(c, err)
		return
	}

	handler.RespondWithSuccess(c, gin.H{
		"removed": len(req.IPs),
	})
}

// ClearWhitelist 清空白名单
//
// @Summary 清空白名单
// @Description 清空IP白名单
// @Tags IPWhitelist
// @Produce json
// @Success 200 {object} map[string]interface{} "data": {...}
// @Router /admin/v1/ip-whitelist/clear [post]
func (h *IPWhitelistHandler) ClearWhitelist(c *gin.Context) {
	// 调用Service层
	if err := h.whitelistService.ClearWhitelist(c.Request.Context()); err != nil {
		handler.RespondWithError(c, err)
		return
	}

	handler.RespondWithSuccess(c, gin.H{
		"message": "IP白名单已清空",
	})
}
