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

// ContractHandler 合同处理器
// 处理合同相关的 HTTP 请求（CRUD 操作）
type ContractHandler struct {
	svc service.ContractService // 合同服务接口
}

// NewContractHandler 创建合同处理器实例
// 参数：
//   - svc: 合同服务接口实例
// 返回值：
//   - *ContractHandler: 合同处理器实例
func NewContractHandler(svc service.ContractService) *ContractHandler {
	return &ContractHandler{svc: svc}
}

// List 获取合同列表
// 处理 GET /api/contracts 请求
// 支持分页、按客户ID过滤、按状态过滤和搜索关键词
//
// godoc
// @Summary 获取合同列表
// @Tags 合同管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Param customerId query string false "客户ID"
// @Param status query string false "合同状态"
// @Param keyword query string false "搜索关键词"
// @Success 200 {object} service.ContractListResult
// @Router /api/contracts [get]
func (h *ContractHandler) List(c *gin.Context) {
	// 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	// 构建查询参数
	params := repository.ContractListParams{
		Page:       page,
		PageSize:   pageSize,
		CustomerID: c.Query("customerId"),
		Status:     c.Query("status"),
		Keyword:    c.Query("keyword"),
	}

	// 调用服务层获取合同列表
	result, err := h.svc.List(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回结果
	c.JSON(http.StatusOK, result)
}

// Create 创建合同
// 处理 POST /api/contracts 请求
// 创建新的合同记录
//
// godoc
// @Summary 创建合同
// @Tags 合同管理
// @Accept json
// @Produce json
// @Param body body service.CreateContractReq true "合同信息"
// @Success 201 {object} domain.Contract
// @Router /api/contracts [post]
func (h *ContractHandler) Create(c *gin.Context) {
	// 绑定并验证请求参数
	var req service.CreateContractReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 调用服务层创建合同
	contract, err := h.svc.Create(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回 201 状态码和创建的合同信息
	c.JSON(http.StatusCreated, contract)
}

// Get 获取合同详情
// 处理 GET /api/contracts/{id} 请求
// 根据合同 ID 获取单个合同的详细信息
//
// godoc
// @Summary 获取合同详情
// @Tags 合同管理
// @Accept json
// @Produce json
// @Param id path string true "合同ID"
// @Success 200 {object} domain.Contract
// @Failure 404 {object} map[string]string
// @Router /api/contracts/{id} [get]
func (h *ContractHandler) Get(c *gin.Context) {
	// 从 URL 路径参数中获取合同 ID
	id := c.Param("id")

	// 调用服务层获取合同详情
	contract, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		// 合同不存在，返回 404 错误
		c.JSON(http.StatusNotFound, gin.H{"error": "合同不存在"})
		return
	}

	// 返回合同信息
	c.JSON(http.StatusOK, contract)
}

// Update 更新合同
// 处理 PUT /api/contracts/{id} 请求
// 更新指定合同的信息
//
// godoc
// @Summary 更新合同
// @Tags 合同管理
// @Accept json
// @Produce json
// @Param id path string true "合同ID"
// @Param body body service.UpdateContractReq true "更新信息"
// @Success 200 {object} domain.Contract
// @Router /api/contracts/{id} [put]
func (h *ContractHandler) Update(c *gin.Context) {
	// 从 URL 路径参数中获取合同 ID
	id := c.Param("id")

	// 绑定并验证请求参数
	var req service.UpdateContractReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 调用服务层更新合同
	contract, err := h.svc.Update(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回更新后的合同信息
	c.JSON(http.StatusOK, contract)
}

// Delete 删除合同
// 处理 DELETE /api/contracts/{id} 请求
// 删除指定的合同记录
//
// godoc
// @Summary 删除合同
// @Tags 合同管理
// @Accept json
// @Produce json
// @Param id path string true "合同ID"
// @Success 204
// @Router /api/contracts/{id} [delete]
func (h *ContractHandler) Delete(c *gin.Context) {
	// 从 URL 路径参数中获取合同 ID
	id := c.Param("id")

	// 调用服务层删除合同
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回 204 No Content 状态码
	c.Status(http.StatusNoContent)
}
