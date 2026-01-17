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

// GMHandler GM配置处理器
//
// 处理GM后台的GM配置请求
type GMHandler struct {
	gmService service.GMService
}

// NewGMHandler 创建GM配置处理器
//
// 参数:
//   gmService: GM配置服务
//
// 返回:
//   *GMHandler: GM配置处理器
func NewGMHandler(gmService service.GMService) *GMHandler {
	return &GMHandler{
		gmService: gmService,
	}
}

// GetGMConfig 获取GM配置
//
// @Summary 获取GM配置
// @Description 获取当前的GM配置
// @Tags GM
// @Produce json
// @Success 200 {object} map[string]interface{} "data": {...}
// @Failure 404 {object} map[string]interface{} "error": {...}
// @Router /admin/v1/gm/config [get]
func (h *GMHandler) GetGMConfig(c *gin.Context) {
	// 调用Service层
	config, err := h.gmService.GetGMConfig(c.Request.Context())
	if err != nil {
		handler.RespondWithError(c, err)
		return
	}

	handler.RespondWithSuccess(c, config)
}

// UpdateGMConfig 更新GM配置
//
// @Summary 更新GM配置
// @Description 更新GM配置
// @Tags GM
// @Accept json
// @Produce json
// @Param body body domain.GMConfig true "GM配置"
// @Success 200 {object} map[string]interface{} "data": {...}
// @Failure 400 {object} map[string]interface{} "error": {...}
// @Router /admin/v1/gm/config [put]
func (h *GMHandler) UpdateGMConfig(c *gin.Context) {
	var config domain.GMConfig
	if !handler.ShouldBindJSON(c, &config) {
		return
	}

	// 调用Service层
	if err := h.gmService.UpdateGMConfig(c.Request.Context(), &config); err != nil {
		handler.RespondWithError(c, err)
		return
	}

	handler.RespondWithSuccess(c, gin.H{
		"gmEnable": config.GMEnable,
		"block":    config.Block,
	})
}

// ToggleGM 切换GM功能开关
//
// @Summary 切换GM功能
// @Description 切换GM功能开关
// @Tags GM
// @Produce json
// @Success 200 {object} map[string]interface{} "data": {...}
// @Failure 500 {object} map[string]interface{} "error": {...}
// @Router /admin/v1/gm/toggle [post]
func (h *GMHandler) ToggleGM(c *gin.Context) {
	// 调用Service层
	state, err := h.gmService.ToggleGM(c.Request.Context())
	if err != nil {
		handler.RespondWithError(c, err)
		return
	}

	handler.RespondWithSuccess(c, gin.H{
		"gmEnable": state,
	})
}

// ToggleBlock 切换封锁状态
//
// @Summary 切换封锁状态
// @Description 切换封锁状态
// @Tags GM
// @Produce json
// @Success 200 {object} map[string]interface{} "data": {...}
// @Failure 500 {object} map[string]interface{} "error": {...}
// @Router /admin/v1/gm/block/toggle [post]
func (h *GMHandler) ToggleBlock(c *gin.Context) {
	// 调用Service层
	state, err := h.gmService.ToggleBlock(c.Request.Context())
	if err != nil {
		handler.RespondWithError(c, err)
		return
	}

	handler.RespondWithSuccess(c, gin.H{
		"block": state,
	})
}

// InitializeGMConfig 初始化GM配置
//
// @Summary 初始化GM配置
// @Description 初始化GM配置（如果不存在）
// @Tags GM
// @Produce json
// @Success 200 {object} map[string]interface{} "data": {...}
// @Failure 500 {object} map[string]interface{} "error": {...}
// @Router /admin/v1/gm/initialize [post]
func (h *GMHandler) InitializeGMConfig(c *gin.Context) {
	// 调用Service层
	if err := h.gmService.InitializeGMConfig(c.Request.Context()); err != nil {
		handler.RespondWithError(c, err)
		return
	}

	handler.RespondWithSuccess(c, gin.H{
		"message": "GM配置已初始化",
	})
}
