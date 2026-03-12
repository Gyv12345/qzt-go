package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/qzt/backend/internal/repository"
	"github.com/qzt/backend/pkg/response"
)

// PermissionMiddleware 权限中间件
// 自动检查用户是否有访问权限
func PermissionMiddleware(repos *repository.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取用户 ID
		userID, exists := c.Get("userId")
		if !exists {
			response.Error(c, response.CodeUnauthorized)
			c.Abort()
			return
		}

		// 获取请求路径和方法
		path := c.Request.URL.Path
		method := c.Request.Method

		// 跳过不需要权限检查的路径
		if skipPermissionCheck(path) {
			c.Next()
			return
		}

		// 构建权限代码（如：customer:create）
		permissionCode := buildPermissionCode(path, method)
		if permissionCode == "" {
			c.Next()
			return
		}

		// 检查用户是否有该权限
		permissions, err := repos.User.GetPermissions(c.Request.Context(), userID.(string))
		if err != nil {
			response.Error(c, response.CodeInternalError)
			c.Abort()
			return
		}

		// 超级管理员跳过权限检查
		if hasPermission(permissions, "*") {
			c.Next()
			return
		}

		// 检查具体权限
		if !hasPermission(permissions, permissionCode) {
			response.Error(c, response.CodeForbidden)
			c.Abort()
			return
		}

		c.Next()
	}
}

// skipPermissionCheck 跳过权限检查的路径
func skipPermissionCheck(path string) bool {
	skipPaths := []string{
		"/api/auth/login",
		"/api/auth/register",
		"/api/auth/refresh",
		"/health",
	}

	for _, p := range skipPaths {
		if strings.HasPrefix(path, p) {
			return true
		}
	}
	return false
}

// buildPermissionCode 构建权限代码
// 例如：POST /api/customers -> customer:create
func buildPermissionCode(path, method string) string {
	// 解析路径
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) < 2 {
		return ""
	}

	// 获取资源名称（如：customers）
	resource := parts[1]
	if strings.HasPrefix(resource, ":") {
		return ""
	}

	// 去掉复数形式
	if strings.HasSuffix(resource, "s") {
		resource = resource[:len(resource)-1]
	}

	// 根据方法确定操作
	var action string
	switch method {
	case "GET":
		if len(parts) > 2 && !strings.HasPrefix(parts[2], ":") {
			action = "read"
		} else {
			action = "list"
		}
	case "POST":
		action = "create"
	case "PUT", "PATCH":
		action = "update"
	case "DELETE":
		action = "delete"
	default:
		return ""
	}

	return resource + ":" + action
}

// hasPermission 检查权限列表中是否包含指定权限
func hasPermission(permissions []interface{}, code string) bool {
	for _, perm := range permissions {
		if p, ok := perm.(map[string]interface{}); ok {
			if p["code"] == code {
				return true
			}
		}
	}
	return false
}

// DataScopeMiddleware 数据权限中间件
// 限制用户只能访问自己或部门的数据
func DataScopeMiddleware(repos *repository.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userId")
		// deptID, _ := c.Get("deptId")

		// 将用户信息存入上下文，供后续使用
		c.Set("currentUserId", userID)

		// TODO: 实现数据权限逻辑
		// 1. 全部数据权限
		// 2. 本部门数据权限
		// 3. 本部门及以下数据权限
		// 4. 仅本人数据权限

		c.Next()
	}
}
