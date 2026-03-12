package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/qzt/backend/internal/repository"
	"github.com/qzt/backend/internal/service"
	"github.com/qzt/backend/pkg/response"
)

// OpportunityHandler 商机处理器
type OpportunityHandler struct {
	svc service.OpportunityService
}

func NewOpportunityHandler(svc service.OpportunityService) *OpportunityHandler {
	return &OpportunityHandler{svc: svc}
}

// Create 创建商机
func (h *OpportunityHandler) Create(c *gin.Context) {
	var req service.CreateOpportunityReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParams, err.Error())
		return
	}

	opp, err := h.svc.Create(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, response.CodeInternalError)
		return
	}

	response.Created(c, opp)
}

// List 获取商机列表
func (h *OpportunityHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	params := repository.ListOpportunityParams{
		Page:     page,
		PageSize: pageSize,
		UserID:   c.Query("userId"),
		Stage:    c.Query("stage"),
		Status:   c.Query("status"),
	}

	result, err := h.svc.List(c.Request.Context(), params)
	if err != nil {
		response.Error(c, response.CodeInternalError)
		return
	}

	response.SuccessPage(c, result.Items, result.Total, result.Page, result.PageSize)
}

// Get 获取商机详情
func (h *OpportunityHandler) Get(c *gin.Context) {
	id := c.Param("id")

	opp, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, response.CodeNotFound)
		return
	}

	response.Success(c, opp)
}

// Update 更新商机
func (h *OpportunityHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req service.UpdateOpportunityReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParams, err.Error())
		return
	}

	opp, err := h.svc.Update(c.Request.Context(), id, &req)
	if err != nil {
		response.Error(c, response.CodeInternalError)
		return
	}

	response.Success(c, opp)
}

// Delete 删除商机
func (h *OpportunityHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		response.Error(c, response.CodeInternalError)
		return
	}

	response.NoContent(c)
}

// MoveToNextStage 推进到下一阶段
func (h *OpportunityHandler) MoveToNextStage(c *gin.Context) {
	id := c.Param("id")

	opp, err := h.svc.MoveToNextStage(c.Request.Context(), id)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeInternalError, "无法推进到下一阶段")
		return
	}

	response.Success(c, opp)
}

// MoveToPrevStage 回退到上一阶段
func (h *OpportunityHandler) MoveToPrevStage(c *gin.Context) {
	id := c.Param("id")

	opp, err := h.svc.MoveToPrevStage(c.Request.Context(), id)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeInternalError, "无法回退到上一阶段")
		return
	}

	response.Success(c, opp)
}

// ConvertToContract 商机转化为合同
func (h *OpportunityHandler) ConvertToContract(c *gin.Context) {
	id := c.Param("id")

	contract, err := h.svc.ConvertToContract(c.Request.Context(), id)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeInternalError, "转化失败")
		return
	}

	response.Created(c, contract)
}

// GetStageStats 获取商机阶段统计
func (h *OpportunityHandler) GetStageStats(c *gin.Context) {
	userID := c.Query("userId")

	stats, err := h.svc.GetStageStats(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, response.CodeInternalError)
		return
	}

	response.Success(c, stats)
}

// PaymentHandler 回款处理器
type PaymentHandler struct {
	svc service.PaymentService
}

func NewPaymentHandler(svc service.PaymentService) *PaymentHandler {
	return &PaymentHandler{svc: svc}
}

type CreatePaymentReq struct {
	ContractID  string  `json:"contractId" binding:"required"`
	CustomerID  string  `json:"customerId" binding:"required"`
	Amount      float64 `json:"amount" binding:"required"`
	PaymentDate string  `json:"paymentDate" binding:"required"`
	PaymentType string  `json:"paymentType"`
	Remark      string  `json:"remark"`
}

// Create 创建回款记录
func (h *PaymentHandler) Create(c *gin.Context) {
	var req CreatePaymentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParams, err.Error())
		return
	}

	paymentDate, err := time.Parse("2006-01-02", req.PaymentDate)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParams, "日期格式错误")
		return
	}

	payment := &service.CreatePaymentReq{
		ContractID:  req.ContractID,
		CustomerID:  req.CustomerID,
		Amount:      req.Amount,
		PaymentDate: paymentDate,
		PaymentType: req.PaymentType,
		Remark:      req.Remark,
	}

	result, err := h.svc.Create(c.Request.Context(), payment)
	if err != nil {
		response.Error(c, response.CodeInternalError)
		return
	}

	response.Created(c, result)
}

// List 获取回款列表
func (h *PaymentHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	params := repository.ListPaymentParams{
		Page:       page,
		PageSize:   pageSize,
		ContractID: c.Query("contractId"),
		CustomerID: c.Query("customerId"),
		Status:     c.Query("status"),
	}

	result, err := h.svc.List(c.Request.Context(), params)
	if err != nil {
		response.Error(c, response.CodeInternalError)
		return
	}

	response.SuccessPage(c, result.Items, result.Total, result.Page, result.PageSize)
}

// Confirm 确认回款
func (h *PaymentHandler) Confirm(c *gin.Context) {
	id := c.Param("id")

	payment, err := h.svc.Confirm(c.Request.Context(), id)
	if err != nil {
		response.Error(c, response.CodeInternalError)
		return
	}

	response.Success(c, payment)
}
