package handler

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/qzt/backend/internal/service"
	"github.com/qzt/backend/pkg/response"
)

// InvoiceHandler 发票处理器
type InvoiceHandler struct {
	svc service.InvoiceService
}

func NewInvoiceHandler(svc service.InvoiceService) *InvoiceHandler {
	return &InvoiceHandler{svc: svc}
}

// Create 创建发票
func (h *InvoiceHandler) Create(c *gin.Context) {
	var req service.CreateInvoiceReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParams, err.Error())
		return
	}

	invoice, err := h.svc.Create(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, response.CodeInternalError)
		return
	}

	response.Created(c, invoice)
}

// List 获取发票列表
func (h *InvoiceHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	params := service.ListInvoiceParams{
		Page:        page,
		PageSize:    pageSize,
		ContractID:  c.Query("contractId"),
		CustomerID:  c.Query("customerId"),
		Status:      c.Query("status"),
		InvoiceType: c.Query("invoiceType"),
	}

	if startDateStr := c.Query("startDate"); startDateStr != "" {
		if t, err := time.Parse("2006-01-02", startDateStr); err == nil {
			params.StartDate = &t
		}
	}

	if endDateStr := c.Query("endDate"); endDateStr != "" {
		if t, err := time.Parse("2006-01-02", endDateStr); err == nil {
			params.EndDate = &t
		}
	}

	result, err := h.svc.List(c.Request.Context(), params)
	if err != nil {
		response.Error(c, response.CodeInternalError)
		return
	}

	response.SuccessPage(c, result.Items, result.Total, result.Page, result.PageSize)
}

// Get 获取发票详情
func (h *InvoiceHandler) Get(c *gin.Context) {
	id := c.Param("id")

	invoice, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, response.CodeNotFound)
		return
	}

	response.Success(c, invoice)
}

// Update 更新发票
func (h *InvoiceHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req service.UpdateInvoiceReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParams, err.Error())
		return
	}

	invoice, err := h.svc.Update(c.Request.Context(), id, &req)
	if err != nil {
		response.Error(c, response.CodeInternalError)
		return
	}

	response.Success(c, invoice)
}

// Delete 删除发票
func (h *InvoiceHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		response.Error(c, response.CodeInternalError)
		return
	}

	response.NoContent(c)
}

// Submit 提交审批
func (h *InvoiceHandler) Submit(c *gin.Context) {
	id := c.Param("id")

	invoice, err := h.svc.Submit(c.Request.Context(), id)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeInternalError, "提交失败")
		return
	}

	response.Success(c, invoice)
}

// Approve 审批通过
func (h *InvoiceHandler) Approve(c *gin.Context) {
	id := c.Param("id")

	invoice, err := h.svc.Approve(c.Request.Context(), id)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeInternalError, "审批失败")
		return
	}

	response.Success(c, invoice)
}

// Reject 审批拒绝
func (h *InvoiceHandler) Reject(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParams, "请填写拒绝原因")
		return
	}

	invoice, err := h.svc.Reject(c.Request.Context(), id, req.Reason)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeInternalError, "拒绝失败")
		return
	}

	response.Success(c, invoice)
}

// Invoice 已开票
func (h *InvoiceHandler) Invoice(c *gin.Context) {
	id := c.Param("id")

	invoice, err := h.svc.Invoice(c.Request.Context(), id)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeInternalError, "操作失败")
		return
	}

	response.Success(c, invoice)
}

// Receive 已收到
func (h *InvoiceHandler) Receive(c *gin.Context) {
	id := c.Param("id")

	invoice, err := h.svc.Receive(c.Request.Context(), id)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeInternalError, "操作失败")
		return
	}

	response.Success(c, invoice)
}

// GetStats 发票统计
func (h *InvoiceHandler) GetStats(c *gin.Context) {
	userID := c.Query("userId")
	startDateStr := c.Query("startDate")
	endDateStr := c.Query("endDate")

	var startDate, endDate time.Time
	var err error

	if startDateStr != "" {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			response.ErrorWithMsg(c, response.CodeInvalidParams, "startDate 格式错误")
			return
		}
	} else {
		startDate = time.Now().AddDate(0, -1, 0)
	}

	if endDateStr != "" {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			response.ErrorWithMsg(c, response.CodeInvalidParams, "endDate 格式错误")
			return
		}
	} else {
		endDate = time.Now()
	}

	params := service.InvoiceStatsParams{
		UserID:    userID,
		StartDate: startDate,
		EndDate:   endDate,
	}

	report, err := h.svc.GetStats(c.Request.Context(), params)
	if err != nil {
		response.Error(c, response.CodeInternalError)
		return
	}

	response.Success(c, report)
}
