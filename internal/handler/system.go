package handler

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/qzt/backend/internal/repository"
	"github.com/qzt/backend/internal/service"
	"github.com/qzt/backend/pkg/response"
)

// NotificationHandler 通知处理器
type NotificationHandler struct {
	svc service.NotificationService
}

func NewNotificationHandler(svc service.NotificationService) *NotificationHandler {
	return &NotificationHandler{svc: svc}
}

// List 获取通知列表
func (h *NotificationHandler) List(c *gin.Context) {
	userID, _ := c.Get("userId")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	params := repository.ListNotificationParams{
		Page:     page,
		PageSize: pageSize,
		Type:     c.Query("type"),
		Category: c.Query("category"),
	}

	if isReadStr := c.Query("isRead"); isReadStr != "" {
		isRead := isReadStr == "true"
		params.IsRead = &isRead
	}

	result, err := h.svc.List(c.Request.Context(), userID.(string), params)
	if err != nil {
		response.Error(c, response.CodeInternalError)
		return
	}

	response.SuccessPage(c, result.Items, result.Total, result.Page, result.PageSize)
}

// GetUnreadCount 获取未读数量
func (h *NotificationHandler) GetUnreadCount(c *gin.Context) {
	userID, _ := c.Get("userId")

	count, err := h.svc.GetUnreadCount(c.Request.Context(), userID.(string))
	if err != nil {
		response.Error(c, response.CodeInternalError)
		return
	}

	response.Success(c, map[string]int64{"count": count})
}

// MarkAsRead 标记为已读
func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	id := c.Param("id")

	if err := h.svc.MarkAsRead(c.Request.Context(), id); err != nil {
		response.Error(c, response.CodeInternalError)
		return
	}

	response.Success(c, nil)
}

// MarkAllAsRead 全部标记为已读
func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	userID, _ := c.Get("userId")

	if err := h.svc.MarkAllAsRead(c.Request.Context(), userID.(string)); err != nil {
		response.Error(c, response.CodeInternalError)
		return
	}

	response.Success(c, nil)
}

// Delete 删除通知
func (h *NotificationHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		response.Error(c, response.CodeInternalError)
		return
	}

	response.NoContent(c)
}

// OperationLogHandler 操作日志处理器
type OperationLogHandler struct {
	svc service.OperationLogService
}

func NewOperationLogHandler(svc service.OperationLogService) *OperationLogHandler {
	return &OperationLogHandler{svc: svc}
}

// List 获取操作日志列表
func (h *OperationLogHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	params := service.ListOperationLogParams{
		Page:     page,
		PageSize: pageSize,
		UserID:   c.Query("userId"),
		Method:   c.Query("method"),
		Path:     c.Query("path"),
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
