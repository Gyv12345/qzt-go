package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/qzt/backend/internal/service"
	"github.com/qzt/backend/pkg/response"
)

// ReportHandler 报表处理器
type ReportHandler struct {
	svc service.ReportService
}

func NewReportHandler(svc service.ReportService) *ReportHandler {
	return &ReportHandler{svc: svc}
}

// GetSalesFunnel 获取销售漏斗
func (h *ReportHandler) GetSalesFunnel(c *gin.Context) {
	userID := c.Query("userId")

	report, err := h.svc.GetSalesFunnel(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, response.CodeInternalError)
		return
	}

	response.Success(c, report)
}

// GetPerformance 获取业绩统计
func (h *ReportHandler) GetPerformance(c *gin.Context) {
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
		startDate = time.Now().AddDate(0, -1, 0) // 默认最近一个月
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

	params := service.PerformanceParams{
		UserID:    userID,
		StartDate: startDate,
		EndDate:   endDate,
	}

	report, err := h.svc.GetPerformance(c.Request.Context(), params)
	if err != nil {
		response.Error(c, response.CodeInternalError)
		return
	}

	response.Success(c, report)
}

// GetPaymentStats 获取回款统计
func (h *ReportHandler) GetPaymentStats(c *gin.Context) {
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

	params := service.PaymentStatsParams{
		UserID:    userID,
		StartDate: startDate,
		EndDate:   endDate,
	}

	report, err := h.svc.GetPaymentStats(c.Request.Context(), params)
	if err != nil {
		response.Error(c, response.CodeInternalError)
		return
	}

	response.Success(c, report)
}

// GetCustomerStats 获取客户统计
func (h *ReportHandler) GetCustomerStats(c *gin.Context) {
	userID := c.Query("userId")

	report, err := h.svc.GetCustomerStats(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, response.CodeInternalError)
		return
	}

	response.Success(c, report)
}

// GetDashboard 获取仪表盘
func (h *ReportHandler) GetDashboard(c *gin.Context) {
	userID, _ := c.Get("userId")

	report, err := h.svc.GetDashboard(c.Request.Context(), userID.(string))
	if err != nil {
		response.Error(c, response.CodeInternalError)
		return
	}

	response.Success(c, report)
}
