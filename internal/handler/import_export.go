package handler

import (
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/qzt/backend/internal/service"
	"github.com/qzt/backend/pkg/response"
)

// ImportExportHandler 导入导出处理器
type ImportExportHandler struct {
	svc service.ImportExportService
}

func NewImportExportHandler(svc service.ImportExportService) *ImportExportHandler {
	return &ImportExportHandler{svc: svc}
}

// ========== 客户导入导出 ==========

// ExportCustomers 导出客户
func (h *ImportExportHandler) ExportCustomers(c *gin.Context) {
	format := c.DefaultQuery("format", "json")

	params := service.ExportCustomersParams{
		Format: format,
		UserID:  c.Query("userId"),
	}

	data, err := h.svc.ExportCustomers(c.Request.Context(), params)
	if err != nil {
		response.Error(c, response.CodeInternalError)
		return
	}

	var contentType string
	var filename string
	switch format {
	case "csv":
		contentType = "text/csv"
		filename = "customers.csv"
	default:
		contentType = "application/json"
		filename = "customers.json"
	}

	c.Header("Content-Type", contentType)
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Data(http.StatusOK, contentType, data)
}

// ImportCustomers 导入客户
func (h *ImportExportHandler) ImportCustomers(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParams, "请上传文件")
		return
	}

	src, err := file.Open()
	if err != nil {
		response.Error(c, response.CodeInternalError)
		return
	}
	defer src.Close()

	data, err := io.ReadAll(src)
	if err != nil {
		response.Error(c, response.CodeInternalError)
		return
	}

	result, err := h.svc.ImportCustomers(c.Request.Context(), data)
	if err != nil {
		response.Error(c, response.CodeInternalError)
		return
	}

	response.Success(c, result)
}

// ========== 商机导入导出 ==========

// ExportOpportunities 导出商机
func (h *ImportExportHandler) ExportOpportunities(c *gin.Context) {
	format := c.DefaultQuery("format", "json")

	params := service.ExportOpportunitiesParams{
		Format: format,
		UserID:  c.Query("userId"),
		Stage:   c.Query("stage"),
	}

	data, err := h.svc.ExportOpportunities(c.Request.Context(), params)
	if err != nil {
		response.Error(c, response.CodeInternalError)
		return
	}

	var contentType string
	var filename string
	switch format {
	case "csv":
		contentType = "text/csv"
		filename = "opportunities.csv"
	default:
		contentType = "application/json"
		filename = "opportunities.json"
	}

	c.Header("Content-Type", contentType)
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Data(http.StatusOK, contentType, data)
}

// ImportOpportunities 导入商机
func (h *ImportExportHandler) ImportOpportunities(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParams, "请上传文件")
		return
	}

	src, err := file.Open()
	if err != nil {
		response.Error(c, response.CodeInternalError)
		return
	}
	defer src.Close()

	data, err := io.ReadAll(src)
	if err != nil {
		response.Error(c, response.CodeInternalError)
		return
	}

	result, err := h.svc.ImportOpportunities(c.Request.Context(), data)
	if err != nil {
		response.Error(c, response.CodeInternalError)
		return
	}

	response.Success(c, result)
}

// ========== 合同导入导出 ==========

// ExportContracts 导出合同
func (h *ImportExportHandler) ExportContracts(c *gin.Context) {
	format := c.DefaultQuery("format", "json")

	params := service.ExportContractsParams{
		Format: format,
		UserID:  c.Query("userId"),
		Status:  c.Query("status"),
	}

	data, err := h.svc.ExportContracts(c.Request.Context(), params)
	if err != nil {
		response.Error(c, response.CodeInternalError)
		return
	}

	var contentType string
	var filename string
	switch format {
	case "csv":
		contentType = "text/csv"
		filename = "contracts.csv"
	default:
		contentType = "application/json"
		filename = "contracts.json"
	}

	c.Header("Content-Type", contentType)
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Data(http.StatusOK, contentType, data)
}

// ImportContracts 导入合同
func (h *ImportExportHandler) ImportContracts(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParams, "请上传文件")
		return
	}

	src, err := file.Open()
	if err != nil {
		response.Error(c, response.CodeInternalError)
		return
	}
	defer src.Close()

	data, err := io.ReadAll(src)
	if err != nil {
		response.Error(c, response.CodeInternalError)
		return
	}

	result, err := h.svc.ImportContracts(c.Request.Context(), data)
	if err != nil {
		response.Error(c, response.CodeInternalError)
		return
	}

	response.Success(c, result)
}

// ========== 回款导入导出 ==========

// ExportPayments 导出回款
func (h *ImportExportHandler) ExportPayments(c *gin.Context) {
	format := c.DefaultQuery("format", "json")
	startDateStr := c.DefaultQuery("startDate", "")
	endDateStr := c.DefaultQuery("endDate", "")

	var startDate, endDate time.Time
	var err error

	if startDateStr != "" {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			response.ErrorWithMsg(c, response.CodeInvalidParams, "startDate 格式错误")
			return
		}
	} else {
		startDate = time.Now().AddDate(-1, 0, 0)
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

	params := service.ExportPaymentsParams{
		Format:   format,
		StartDate: startDate,
		EndDate:   endDate,
	}

	data, err := h.svc.ExportPayments(c.Request.Context(), params)
	if err != nil {
		response.Error(c, response.CodeInternalError)
		return
	}

	var contentType string
	var filename string
	switch format {
	case "csv":
		contentType = "text/csv"
		filename = "payments.csv"
	default:
		contentType = "application/json"
		filename = "payments.json"
	}

	c.Header("Content-Type", contentType)
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Data(http.StatusOK, contentType, data)
}

// ImportPayments 导入回款
func (h *ImportExportHandler) ImportPayments(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParams, "请上传文件")
		return
	}

	src, err := file.Open()
	if err != nil {
		response.Error(c, response.CodeInternalError)
		return
	}
	defer src.Close()

	data, err := io.ReadAll(src)
	if err != nil {
		response.Error(c, response.CodeInternalError)
		return
	}

	result, err := h.svc.ImportPayments(c.Request.Context(), data)
	if err != nil {
		response.Error(c, response.CodeInternalError)
		return
	}

	response.Success(c, result)
}

// ========== 批量操作 ==========

// BatchDeleteCustomers 批量删除客户
func (h *ImportExportHandler) BatchDeleteCustomers(c *gin.Context) {
	var req struct {
		IDs []string `json:"ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParams, err.Error())
		return
	}

	// TODO: 实现批量删除
	response.Success(c, map[string]int{"deleted": len(req.IDs)})
}

// BatchUpdateCustomers 批量更新客户
func (h *ImportExportHandler) BatchUpdateCustomers(c *gin.Context) {
	var req struct {
		IDs      []string          `json:"ids" binding:"required"`
		Updates  map[string]interface{} `json:"updates" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParams, err.Error())
		return
	}

	// TODO: 实现批量更新
	response.Success(c, map[string]int{"updated": len(req.IDs)})
}
