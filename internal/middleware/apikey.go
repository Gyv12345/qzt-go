// Package middleware 定义了 HTTP 中间件
// 用于处理请求的横切关注点，如认证、授权、日志、CORS 等
package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/qzt/backend/internal/domain"
	"github.com/qzt/backend/internal/repository"
	"github.com/qzt/backend/pkg/apikey"
)

// APIKeyMiddleware API Key 认证中间件
// 用于第三方系统调用 API 时的身份验证和权限控制
type APIKeyMiddleware struct {
	repo repository.APIKeyRepository // API Key 仓储接口
}

// NewAPIKeyMiddleware 创建 API Key 中间件实例
// 参数：
//   - repo: API Key 仓储接口
// 返回值：
//   - *APIKeyMiddleware: API Key 中间件实例
func NewAPIKeyMiddleware(repo repository.APIKeyRepository) *APIKeyMiddleware {
	return &APIKeyMiddleware{repo: repo}
}

// Handler 返回 Gin 中间件处理函数
// 验证 API Key 的有效性、权限和 IP 白名单限制
// 返回值：
//   - gin.HandlerFunc: Gin 中间件函数
func (m *APIKeyMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// ========== 第一步：提取 API Key ==========
		// 从请求中提取 API Key
		apiKeyStr := extractAPIKey(c)
		if apiKeyStr == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "API key required",
			})
			c.Abort()
			return
		}

		// ========== 第二步：验证 API Key 格式 ==========
		// 检查 API Key 格式是否正确（以 "qzt_" 开头）
		if !apikey.Validate(apiKeyStr) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid API key format",
			})
			c.Abort()
			return
		}

		// ========== 第三步：查询数据库验证 ==========
		// 计算密钥哈希并查询数据库
		keyHash := apikey.Hash(apiKeyStr)
		keyInfo, err := m.repo.GetByHash(c.Request.Context(), keyHash)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid API key",
			})
			c.Abort()
			return
		}

		// ========== 第四步：检查状态 ==========
		// 检查 API Key 是否被禁用
		if keyInfo.Status != 1 {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "API key is disabled",
			})
			c.Abort()
			return
		}

		// ========== 第五步：检查过期时间 ==========
		// 检查 API Key 是否已过期
		if keyInfo.ExpiresAt != nil && keyInfo.ExpiresAt.Before(time.Now()) {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "API key has expired",
			})
			c.Abort()
			return
		}

		// ========== 第六步：检查 IP 白名单 ==========
		// 如果配置了 IP 白名单，验证客户端 IP 是否在白名单中
		if keyInfo.IPWhitelist != "" && keyInfo.IPWhitelist != "[]" {
			clientIP := c.ClientIP()
			if !isIPAllowed(clientIP, keyInfo.IPWhitelist) {
				c.JSON(http.StatusForbidden, gin.H{
					"error": "IP not allowed",
				})
				c.Abort()
				return
			}
		}

		// ========== 第七步：检查权限 ==========
		// 检查 API Key 是否有访问当前接口的权限
		requiredPerm := getRequiredPermission(c)
		if requiredPerm != "" && !hasPermission(keyInfo.Permissions, requiredPerm) {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Permission denied",
			})
			c.Abort()
			return
		}

		// ========== 第八步：注入信息到 Context ==========
		// 将 API Key 信息和用户 ID 注入到 Context
		c.Set("apiKey", keyInfo)
		c.Set("userId", keyInfo.UserID)

		// ========== 第九步：更新最后使用时间 ==========
		// 异步更新 API Key 的最后使用时间
		go m.repo.UpdateLastUsed(keyInfo.ID)

		// 继续执行后续处理器
		c.Next()
	}
}

// extractAPIKey 从请求中提取 API Key
// 支持三种提取方式：
//   1. Authorization Header: Bearer <key>
//   2. X-API-Key Header
//   3. Query Parameter: ?api_key=<key>
// 参数：
//   - c: Gin Context 对象
// 返回值：
//   - string: 提取到的 API Key 字符串，如果未找到则返回空字符串
func extractAPIKey(c *gin.Context) string {
	// 方式1: 从 Authorization Header 提取
	auth := c.GetHeader("Authorization")
	if strings.HasPrefix(auth, "Bearer ") {
		return strings.TrimPrefix(auth, "Bearer ")
	}

	// 方式2: 从 X-API-Key Header 提取（推荐方式）
	key := c.GetHeader("X-API-Key")
	if key != "" {
		return key
	}

	// 方式3: 从 Query Parameter 提取（不推荐，仅用于测试）
	return c.Query("api_key")
}

// getRequiredPermission 根据请求方法和路径获取所需的权限
// 参数：
//   - c: Gin Context 对象
// 返回值：
//   - string: 所需权限字符串（如 "customer:read"），如果不需要权限则返回空字符串
func getRequiredPermission(c *gin.Context) string {
	method := c.Request.Method
	path := c.Request.URL.Path

	// 客户管理接口的权限映射
	if strings.HasPrefix(path, "/api/v1/customers") {
		switch method {
		case "GET":
			return "customer:read"
		case "POST", "PUT", "DELETE":
			return "customer:write"
		}
	}

	// 联系人管理接口的权限映射
	if strings.HasPrefix(path, "/api/v1/contacts") {
		switch method {
		case "GET":
			return "contact:read"
		case "POST", "PUT", "DELETE":
			return "contact:write"
		}
	}

	// 合同管理接口的权限映射
	if strings.HasPrefix(path, "/api/v1/contracts") {
		switch method {
		case "GET":
			return "contract:read"
		case "POST", "PUT", "DELETE":
			return "contract:write"
		}
	}

	return ""
}

// hasPermission 检查权限列表中是否包含所需权限
// 参数：
//   - permissionsJSON: 权限列表的 JSON 字符串
//   - required: 所需权限
// 返回值：
//   - bool: 如果有权限返回 true，否则返回 false
func hasPermission(permissionsJSON string, required string) bool {
	// 如果权限为空或包含通配符 *，允许所有访问
	if permissionsJSON == "" || permissionsJSON == "[]" || permissionsJSON == `["*"]` {
		return true
	}

	// 简单的字符串匹配（实际应用中应解析 JSON 后精确匹配）
	return strings.Contains(permissionsJSON, required) || strings.Contains(permissionsJSON, `"*"`)
}

// isIPAllowed 检查客户端 IP 是否在白名单中
// 参数：
//   - clientIP: 客户端 IP 地址
//   - whitelistJSON: 白名单的 JSON 字符串
// 返回值：
//   - bool: 如果 IP 在白名单中返回 true，否则返回 false
func isIPAllowed(clientIP string, whitelistJSON string) bool {
	// 简单实现：检查 IP 是否在白名单 JSON 字符串中
	// TODO: 支持更完善的 CIDR 网段匹配
	return strings.Contains(whitelistJSON, clientIP)
}

// GetAPIKey 从 Context 中获取 API Key 信息
// 供后续处理器使用
// 参数：
//   - c: Gin Context 对象
// 返回值：
//   - *domain.APIKey: API Key 对象，如果不存在则返回 nil
func GetAPIKey(c *gin.Context) *domain.APIKey {
	if key, exists := c.Get("apiKey"); exists {
		return key.(*domain.APIKey)
	}
	return nil
}
