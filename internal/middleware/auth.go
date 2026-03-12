// Package middleware 定义了 HTTP 中间件
// 用于处理请求的横切关注点，如认证、授权、日志、CORS 等
package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/qzt/backend/internal/config"
	"github.com/qzt/backend/internal/service"
)

// JWTAuth JWT 认证中间件
// 验证请求中的 JWT Token，并将用户信息注入到 Context 中
// 参数：
//   - cfg: 应用配置对象（包含 JWT 密钥等）
// 返回值：
//   - gin.HandlerFunc: Gin 中间件函数
func JWTAuth(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// ========== 第一步：提取 Token ==========
		// 从请求头或查询参数中提取 JWT Token
		tokenString := extractToken(c)
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供认证 token"})
			c.Abort()
			return
		}

		// ========== 第二步：解析 Token ==========
		// 解析 JWT Token 并验证签名
		claims := &service.JWTClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.JWT.Secret), nil
		})

		// 验证 Token 是否有效
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的 token"})
			c.Abort()
			return
		}

		// ========== 第三步：检查过期时间 ==========
		// 检查 Token 是否已过期
		if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token 已过期"})
			c.Abort()
			return
		}

		// ========== 第四步：注入用户信息 ==========
		// 将解析出的用户信息注入到 Context 中，供后续处理器使用
		c.Set("userId", claims.UserID)
		c.Set("username", claims.Username)

		// 继续执行后续处理器
		c.Next()
	}
}

// extractToken 从请求中提取 JWT Token
// 支持两种提取方式：
//   1. Authorization Header: Bearer <token>
//   2. Query Parameter: ?token=<token>
// 参数：
//   - c: Gin Context 对象
// 返回值：
//   - string: 提取到的 Token 字符串，如果未找到则返回空字符串
func extractToken(c *gin.Context) string {
	// 方式1: 从 Authorization Header 提取
	auth := c.GetHeader("Authorization")
	if strings.HasPrefix(auth, "Bearer ") {
		return strings.TrimPrefix(auth, "Bearer ")
	}

	// 方式2: 从 Query Parameter 提取
	return c.Query("token")
}
