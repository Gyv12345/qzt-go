// Package middleware 提供 HTTP 中间件
package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	// RequestIDHeader HTTP Header 名称
	RequestIDHeader = "X-Request-ID"
	// RequestIDKey Gin Context 中的键名
	RequestIDKey = "requestId"
)

// RequestID 请求 ID 中间件
// 为每个请求生成唯一的 Request ID，用于日志追踪
//
// 工作流程：
//  1. 检查请求头中是否已有 X-Request-ID
//  2. 如果有，使用该值；如果没有，生成新的 UUID
//  3. 将 Request ID 存入 Gin Context 和响应头
//
// 使用示例：
//
//	r := gin.New()
//	r.Use(middleware.RequestID())
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 尝试从请求头获取 Request ID
		requestID := c.GetHeader(RequestIDHeader)
		
		// 如果没有，生成新的 UUID
		if requestID == "" {
			requestID = uuid.New().String()
		}
		
		// 存入 Context
		c.Set(RequestIDKey, requestID)
		
		// 设置响应头
		c.Header(RequestIDHeader, requestID)
		
		c.Next()
	}
}

// GetRequestID 从 Context 获取 Request ID
func GetRequestID(c *gin.Context) string {
	if requestID, exists := c.Get(RequestIDKey); exists {
		if id, ok := requestID.(string); ok {
			return id
		}
	}
	return ""
}
