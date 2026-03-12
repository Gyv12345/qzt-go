package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/qzt/backend/internal/domain"
	"github.com/qzt/backend/internal/service"
)

// OperationLogger 操作日志中间件
func OperationLogger(logSvc service.OperationLogService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录开始时间
		startTime := time.Now()

		// 读取请求体
		var requestBody string
		if c.Request.Body != nil {
			bodyBytes, _ := io.ReadAll(c.Request.Body)
			requestBody = string(bodyBytes)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		// 处理请求
		c.Next()

		// 计算耗时
		latency := time.Since(startTime).Milliseconds()

		// 获取用户信息
		var userID, username string
		if val, exists := c.Get("userId"); exists {
			userID = val.(string)
		}
		if val, exists := c.Get("username"); exists {
			username = val.(string)
		}

		// 只记录需要审计的操作（POST/PUT/DELETE）
		if c.Request.Method == "GET" || c.Request.Method == "OPTIONS" {
			return
		}

		// 限制请求体大小，避免记录过大的数据
		if len(requestBody) > 2000 {
			requestBody = requestBody[:2000] + "..."
		}

		// 创建日志记录
		log := &domain.OperationLog{
			UserID:     userID,
			Username:   username,
			Method:     c.Request.Method,
			Path:       c.Request.URL.Path,
			Query:      c.Request.URL.RawQuery,
			Body:       requestBody,
			IP:         c.ClientIP(),
			UserAgent:  c.Request.UserAgent(),
			StatusCode: c.Writer.Status(),
			Latency:    latency,
		}

		// 记录错误信息
		if len(c.Errors) > 0 {
			log.Error = c.Errors.String()
		}

		// 异步保存日志
		go func() {
			_ = logSvc.Log(c.Request.Context(), log)
		}()
	}
}

// responseWriter 用于捕获响应体
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w *responseWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

// NotificationHelper 通知辅助函数
type NotificationHelper struct {
	svc service.NotificationService
}

func NewNotificationHelper(svc service.NotificationService) *NotificationHelper {
	return &NotificationHelper{svc: svc}
}

// NotifyCustomerAssigned 客户分配通知
func (h *NotificationHelper) NotifyCustomerAssigned(ctx *gin.Context, userID, customerName string) {
	_ = h.svc.Send(ctx.Request.Context(), &service.SendNotificationReq{
		UserID:   userID,
		Title:    "新客户分配",
		Content:  "您有新的客户：" + customerName + "，请及时跟进",
		Type:     "BUSINESS",
		Category: "客户",
		Priority: 1,
	})
}

// NotifyOpportunityStageChange 商机阶段变更通知
func (h *NotificationHelper) NotifyOpportunityStageChange(ctx *gin.Context, userID, oppName, stage string) {
	_ = h.svc.Send(ctx.Request.Context(), &service.SendNotificationReq{
		UserID:   userID,
		Title:    "商机阶段更新",
		Content:  "商机「" + oppName + "」已进入「" + stage + "」阶段",
		Type:     "BUSINESS",
		Category: "商机",
		Priority: 0,
	})
}

// NotifyPaymentReceived 回款通知
func (h *NotificationHelper) NotifyPaymentReceived(ctx *gin.Context, userID, customerName string, amount float64) {
	_ = h.svc.Send(ctx.Request.Context(), &service.SendNotificationReq{
		UserID:   userID,
		Title:    "收到回款",
		Content:  "客户「" + customerName + "」已回款，金额：" + formatMoney(amount),
		Type:     "BUSINESS",
		Category: "回款",
		Priority: 1,
	})
}

// NotifyContractExpiring 合同到期提醒
func (h *NotificationHelper) NotifyContractExpiring(ctx *gin.Context, userID, contractName string, days int) {
	_ = h.svc.Send(ctx.Request.Context(), &service.SendNotificationReq{
		UserID:   userID,
		Title:    "合同即将到期",
		Content:  "合同「" + contractName + "」将在 " + string(rune(days)) + " 天后到期，请及时处理",
		Type:     "REMINDER",
		Category: "合同",
		Priority: 2,
	})
}

func formatMoney(amount float64) string {
	// 简单的金额格式化
	return "¥" + string(rune(int(amount)))
}
