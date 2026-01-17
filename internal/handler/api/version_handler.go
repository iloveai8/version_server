// Package api 提供对外API的HTTP处理器
//
// 该包实现了游戏客户端访问的API接口
package api

import (
	"github.com/gin-gonic/gin"

	"game_slots_vsn/internal/domain"
	"game_slots_vsn/internal/handler"
	"game_slots_vsn/internal/service"
	"game_slots_vsn/pkg/errcode"
)

// VersionHandler 版本配置处理器
//
// 处理游戏客户端的版本配置请求
type VersionHandler struct {
	versionService service.VersionService
}

// NewVersionHandler 创建版本配置处理器
//
// 参数:
//   versionService: 版本配置服务
//
// 返回:
//   *VersionHandler: 版本配置处理器
func NewVersionHandler(versionService service.VersionService) *VersionHandler {
	return &VersionHandler{
		versionService: versionService,
	}
}

// GetVersion 获取版本配置
//
// @Summary 获取版本配置
// @Description 根据版本号和环境获取版本配置
// @Tags Version
// @Accept json
// @Produce json
// @Param vsn query string true "版本号"
// @Param env query string true "环境 (dev/pro)"
// @Success 200 {object} map[string]interface{} "data": {...}"
// @Failure 400 {object} map[string]interface{} "error": {...}
// @Failure 404 {object} map[string]interface{} "error": {...}
// @Router /api/v1/version [get]
func (h *VersionHandler) GetVersion(c *gin.Context) {
	var req handler.GetVersionRequest
	if !handler.ShouldBindQuery(c, &req) {
		return
	}

	// 调用Service层
	version, err := h.versionService.GetVersion(c.Request.Context(), req.Vsn, req.Env)
	if err != nil {
		handler.RespondWithError(c, err)
		return
	}

	handler.RespondWithSuccess(c, version)
}

// ==================== 未实现的Handler方法 ====================
//
// 以下方法保留用于未来扩展，当前版本不提供这些API

// CreateVersion 创建版本（仅用于测试，生产环境不提供）
//
// @Summary 创建版本
// @Description 创建新的版本配置
// @Tags Version
// @Accept json
// @Produce json
// @Param env query string true "环境 (dev/pro)"
// @Param body body domain.Version true "版本配置"
// @Success 201 {object} map[string]interface{} "data": {...}
// @Failure 400 {object} map[string]interface{} "error": {...}"
// @Router /api/v1/version [post]
func (h *VersionHandler) CreateVersion(c *gin.Context) {
	var req handler.CreateVersionRequest
	if !handler.ShouldBindQuery(c, &req) {
		return
	}

	var version domain.Version
	if !handler.ShouldBindJSON(c, &version) {
		return
	}

	// 调用Service层
	if err := h.versionService.CreateVersion(c.Request.Context(), &version, req.Env); err != nil {
		handler.RespondWithError(c, err)
		return
	}

	handler.RespondWithCreated(c, gin.H{
		"vsn": version.Vsn,
	})
}

// UpdateVersion 更新版本（仅用于测试，生产环境不提供）
//
// @Summary 更新版本
// @Description 更新版本配置
// @Tags Version
// @Accept json
// @Produce json
// @Param vsn path string true "版本号"
// @Param env query string true "环境 (dev/pro)"
// @Param body body domain.Version true "版本配置"
// @Success 200 {object} map[string]interface{} "data": {...}
// @Failure 400 {object} map[string]interface{} "error": {...}
// @Failure 404 {object} map[string]interface{} "error": {...}
// @Router /api/v1/version/{vsn} [put]
func (h *VersionHandler) UpdateVersion(c *gin.Context) {
	vsn := c.Param("vsn")
	if vsn == "" {
		handler.RespondWithError(c, errcode.ErrInvalidParams.WithDetails(map[string]interface{}{
			"field": "vsn",
			"reason": "版本号不能为空",
		}))
		return
	}

	var req handler.UpdateVersionRequest
	if !handler.ShouldBindQuery(c, &req) {
		return
	}

	var version domain.Version
	if !handler.ShouldBindJSON(c, &version) {
		return
	}

	// 设置版本号
	version.Vsn = vsn

	// 调用Service层
	if err := h.versionService.UpdateVersion(c.Request.Context(), &version, req.Env); err != nil {
		handler.RespondWithError(c, err)
		return
	}

	handler.RespondWithSuccess(c, gin.H{
		"vsn": version.Vsn,
	})
}

// DeleteVersion 删除版本（仅用于测试，生产环境不提供）
//
// @Summary 删除版本
// @Description 删除版本配置
// @Tags Version
// @Produce json
// @Param vsn path string true "版本号"
// @Param env query string true "环境 (dev/pro)"
// @Success 204
// @Failure 400 {object} map[string]interface{} "error": {...}
// @Failure 404 {object} map[string]interface{} "error": {...}
// @Router /api/v1/version/{vsn} [delete]
func (h *VersionHandler) DeleteVersion(c *gin.Context) {
	vsn := c.Param("vsn")
	if vsn == "" {
		handler.RespondWithError(c, errcode.ErrInvalidParams.WithDetails(map[string]interface{}{
			"field": "vsn",
			"reason": "版本号不能为空",
		}))
		return
	}

	var req handler.DeleteVersionRequest
	if !handler.ShouldBindQuery(c, &req) {
		return
	}

	// 调用Service层
	if err := h.versionService.DeleteVersion(c.Request.Context(), vsn, req.Env); err != nil {
		handler.RespondWithError(c, err)
		return
	}

	handler.RespondWithNoContent(c)
}

// ListVersions 列出版本（仅用于测试，生产环境不提供）
//
// @Summary 列出版本
// @Description 获取所有版本列表
// @Tags Version
// @Produce json
// @Param env query string true "环境 (dev/pro)"
// @Success 200 {object} map[string]interface{} "data": {...}}
// @Failure 400 {object} map[string]interface{} "error": {...}
// @Router /api/v1/versions [get]
func (h *VersionHandler) ListVersions(c *gin.Context) {
	var req handler.ListVersionsRequest
	if !handler.ShouldBindQuery(c, &req) {
		return
	}

	// 调用Service层
	versions, err := h.versionService.ListVersions(c.Request.Context(), req.Env)
	if err != nil {
		handler.RespondWithError(c, err)
		return
	}

	handler.RespondWithSuccess(c, gin.H{
		"versions": versions,
		"count":    len(versions),
	})
}

// AddSubServer 添加子服务器（仅用于测试，生产环境不提供）
//
// @Summary 添加子服务器
// @Description 向版本添加子服务器
// @Tags Version
// @Accept json
// @Produce json
// @Param vsn path string true "版本号"
// @Param env query string true "环境 (dev/pro)"
// @Param key query string true "子服务器键"
// @Param body object{vsn string, srvUrl string, resUrl string, type int} true "子服务器配置"
// @Success 200 {object} map[string]interface{} "data": {...}}
// @Failure 400 {object} map[string]interface{} "error": {...}
// @Failure 404 {object} map[string]interface{} "error": {...}"
// @Router /api/v1/version/{vsn}/subserver [post]
func (h *VersionHandler) AddSubServer(c *gin.Context) {
	vsn := c.Param("vsn")
	if vsn == "" {
		handler.RespondWithError(c, errcode.ErrInvalidParams.WithDetails(map[string]interface{}{
			"field": "vsn",
			"reason": "版本号不能为空",
		}))
		return
	}

	var req handler.AddSubServerRequest
	if !handler.ShouldBindQuery(c, &req) {
		return
	}

	// 构建子服务器对象
	sub := domain.NewSubServer(req.Vsn, req.URL, req.Res, req.Type)

	// 调用Service层
	if err := h.versionService.AddSubServer(c.Request.Context(), vsn, req.Env, req.Key, sub); err != nil {
		handler.RespondWithError(c, err)
		return
	}

	handler.RespondWithSuccess(c, gin.H{
		"key": req.Key,
	})
}

// RemoveSubServer 移除子服务器（仅用于测试，生产环境不提供）
//
// @Summary 移除子服务器
// @Description 从版本移除子服务器
// @Tags Version
// @Produce json
// @Param vsn path string true "版本号"
// @Param env query string true "环境 (dev/pro)"
// @Param key query string true "子服务器键"
// @Success 204
// @Failure 400 {object} map[string]interface{} "error": {...}
// @Failure 404 {object} map[string]interface{} "error": {...}"
// @Router /api/v1/version/{vsn}/subserver [delete]
func (h *VersionHandler) RemoveSubServer(c *gin.Context) {
	vsn := c.Param("vsn")
	if vsn == "" {
		handler.RespondWithError(c, errcode.ErrInvalidParams.WithDetails(map[string]interface{}{
			"field": "vsn",
			"reason": "版本号不能为空",
		}))
		return
	}

	var req handler.RemoveSubServerRequest
	if !handler.ShouldBindQuery(c, &req) {
		return
	}

	// 调用Service层
	if err := h.versionService.RemoveSubServer(c.Request.Context(), vsn, req.Env, req.Key); err != nil {
		handler.RespondWithError(c, err)
		return
	}

	handler.RespondWithNoContent(c)
}
