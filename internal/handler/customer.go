// Package handler 定义了 HTTP 请求处理器
package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/qzt/backend/internal/repository"
	"github.com/qzt/backend/internal/service"
	"github.com/qzt/backend/pkg/response"
)

// CustomerHandler 客户处理器
type CustomerHandler struct {
	svc service.CustomerService
}

// NewCustomerHandler 创建客户处理器实例
func NewCustomerHandler(svc service.CustomerService) *CustomerHandler {
	return &CustomerHandler{svc: svc}
}

// List 获取客户列表
// GET /api/customers
func (h *CustomerHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	var status *int
	if s := c.Query("status"); s != "" {
		v, _ := strconv.Atoi(s)
		status = &v
	}

	params := repository.ListParams{
		Page:     page,
		PageSize: pageSize,
		Keyword:  c.Query("keyword"),
		Level:    c.Query("level"),
		Status:   status,
	}

	result, err := h.svc.List(c.Request.Context(), params)
	if err != nil {
		response.Error(c, response.CodeInternalError)
		return
	}

	response.SuccessPage(c, result.Items, result.Total, result.Page, result.PageSize)
}

// Create 创建客户
// POST /api/customers
func (h *CustomerHandler) Create(c *gin.Context) {
	var req service.CreateCustomerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParams, err.Error())
		return
	}

	customer, err := h.svc.Create(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, response.CodeInternalError)
		return
	}

	response.Created(c, customer)
}

// Get 获取客户详情
// GET /api/customers/:id
func (h *CustomerHandler) Get(c *gin.Context) {
	id := c.Param("id")

	customer, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, response.CodeCustomerNotFound)
		return
	}

	response.Success(c, customer)
}

// Update 更新客户
// PUT /api/customers/:id
func (h *CustomerHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req service.UpdateCustomerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParams, err.Error())
		return
	}

	customer, err := h.svc.Update(c.Request.Context(), id, &req)
	if err != nil {
		response.Error(c, response.CodeInternalError)
		return
	}

	response.Success(c, customer)
}

// Delete 删除客户
// DELETE /api/customers/:id
func (h *CustomerHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		response.Error(c, response.CodeInternalError)
		return
	}

	response.NoContent(c)
}
