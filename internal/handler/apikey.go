// Package handler 定义了 HTTP 请求处理器
// 负责处理 HTTP 请求、参数验证和响应返回
package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qzt/backend/internal/service"
)

// APIKeyHandler API Key 处理器
// 处理 API Key 相关的 HTTP 请求（创建、列表、撤销）
type APIKeyHandler struct {
	svc service.APIKeyService // API Key 服务接口
}

// NewAPIKeyHandler 创建 API Key 处理器实例
// 参数：
//   - svc: API Key 服务接口实例
// 返回值：
//   - *APIKeyHandler: API Key 处理器实例
func NewAPIKeyHandler(svc service.APIKeyService) *APIKeyHandler {
	return &APIKeyHandler{svc: svc}
}

// Create 创建 API Key
// 处理 POST /api/apikeys 请求
// 为当前用户创建新的 API Key
//
// godoc
// @Summary 创建 API Key
// @Tags API Key 管理
// @Accept json
// @Produce json
// @Param body body service.CreateAPIKeyReq true "API Key 信息"
// @Success 201 {object} service.CreateAPIKeyResult
// @Router /api/apikeys [post]
func (h *APIKeyHandler) Create(c *gin.Context) {
	// 绑定并验证请求参数
	var req service.CreateAPIKeyReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 调用服务层创建 API Key
	result, err := h.svc.Create(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回 201 状态码和创建的 API Key 信息
	c.JSON(http.StatusCreated, result)
}

// List 获取 API Key 列表
// 处理 GET /api/apikeys 请求
// 获取当前用户的所有 API Key 列表
//
// godoc
// @Summary 获取 API Key 列表
// @Tags API Key 管理
// @Accept json
// @Produce json
// @Success 200 {array} domain.APIKey
// @Router /api/apikeys [get]
func (h *APIKeyHandler) List(c *gin.Context) {
	// 从 context 中获取当前用户 ID（由 JWT 中间件设置）
	userID, exists := c.Get("userId")
	if !exists {
		// 未获取到用户 ID，返回未授权错误
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 调用服务层获取当前用户的 API Key 列表
	keys, err := h.svc.List(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回 API Key 列表
	c.JSON(http.StatusOK, keys)
}

// Revoke 撤销 API Key
// 处理 POST /api/apikeys/{id}/revoke 请求
// 撤销（禁用）指定的 API Key
//
// godoc
// @Summary 撤销 API Key
// @Tags API Key 管理
// @Accept json
// @Produce json
// @Param id path string true "API Key ID"
// @Success 204
// @Router /api/apikeys/{id}/revoke [post]
func (h *APIKeyHandler) Revoke(c *gin.Context) {
	// 从 URL 路径参数中获取 API Key ID
	id := c.Param("id")

	// 调用服务层撤销 API Key
	if err := h.svc.Revoke(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回 204 No Content 状态码
	c.Status(http.StatusNoContent)
}
