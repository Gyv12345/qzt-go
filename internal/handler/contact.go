// Package handler 定义了 HTTP 请求处理器
// 负责处理 HTTP 请求、参数验证和响应返回
package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/qzt/backend/internal/repository"
	"github.com/qzt/backend/internal/service"
)

// ContactHandler 联系人处理器
// 处理联系人相关的 HTTP 请求（CRUD 操作）
type ContactHandler struct {
	svc service.ContactService // 联系人服务接口
}

// NewContactHandler 创建联系人处理器实例
// 参数：
//   - svc: 联系人服务接口实例
// 返回值：
//   - *ContactHandler: 联系人处理器实例
func NewContactHandler(svc service.ContactService) *ContactHandler {
	return &ContactHandler{svc: svc}
}

// List 获取联系人列表
// 处理 GET /api/contacts 请求
// 支持分页、按客户ID过滤和搜索关键词
//
// godoc
// @Summary 获取联系人列表
// @Tags 联系人管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Param customerId query string false "客户ID"
// @Param keyword query string false "搜索关键词"
// @Success 200 {object} service.ContactListResult
// @Router /api/contacts [get]
func (h *ContactHandler) List(c *gin.Context) {
	// 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	// 构建查询参数
	params := repository.ContactListParams{
		Page:       page,
		PageSize:   pageSize,
		CustomerID: c.Query("customerId"),
		Keyword:    c.Query("keyword"),
	}

	// 调用服务层获取联系人列表
	result, err := h.svc.List(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回结果
	c.JSON(http.StatusOK, result)
}

// Create 创建联系人
// 处理 POST /api/contacts 请求
// 创建新的联系人记录
//
// godoc
// @Summary 创建联系人
// @Tags 联系人管理
// @Accept json
// @Produce json
// @Param body body service.CreateContactReq true "联系人信息"
// @Success 201 {object} domain.Contact
// @Router /api/contacts [post]
func (h *ContactHandler) Create(c *gin.Context) {
	// 绑定并验证请求参数
	var req service.CreateContactReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 调用服务层创建联系人
	contact, err := h.svc.Create(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回 201 状态码和创建的联系人信息
	c.JSON(http.StatusCreated, contact)
}

// Get 获取联系人详情
// 处理 GET /api/contacts/{id} 请求
// 根据联系人 ID 获取单个联系人的详细信息
//
// godoc
// @Summary 获取联系人详情
// @Tags 联系人管理
// @Accept json
// @Produce json
// @Param id path string true "联系人ID"
// @Success 200 {object} domain.Contact
// @Failure 404 {object} map[string]string
// @Router /api/contacts/{id} [get]
func (h *ContactHandler) Get(c *gin.Context) {
	// 从 URL 路径参数中获取联系人 ID
	id := c.Param("id")

	// 调用服务层获取联系人详情
	contact, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		// 联系人不存在，返回 404 错误
		c.JSON(http.StatusNotFound, gin.H{"error": "联系人不存在"})
		return
	}

	// 返回联系人信息
	c.JSON(http.StatusOK, contact)
}

// Update 更新联系人
// 处理 PUT /api/contacts/{id} 请求
// 更新指定联系人的信息
//
// godoc
// @Summary 更新联系人
// @Tags 联系人管理
// @Accept json
// @Produce json
// @Param id path string true "联系人ID"
// @Param body body service.UpdateContactReq true "更新信息"
// @Success 200 {object} domain.Contact
// @Router /api/contacts/{id} [put]
func (h *ContactHandler) Update(c *gin.Context) {
	// 从 URL 路径参数中获取联系人 ID
	id := c.Param("id")

	// 绑定并验证请求参数
	var req service.UpdateContactReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 调用服务层更新联系人
	contact, err := h.svc.Update(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回更新后的联系人信息
	c.JSON(http.StatusOK, contact)
}

// Delete 删除联系人
// 处理 DELETE /api/contacts/{id} 请求
// 删除指定的联系人记录
//
// godoc
// @Summary 删除联系人
// @Tags 联系人管理
// @Accept json
// @Produce json
// @Param id path string true "联系人ID"
// @Success 204
// @Router /api/contacts/{id} [delete]
func (h *ContactHandler) Delete(c *gin.Context) {
	// 从 URL 路径参数中获取联系人 ID
	id := c.Param("id")

	// 调用服务层删除联系人
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回 204 No Content 状态码
	c.Status(http.StatusNoContent)
}
