// Package middleware 定义了 HTTP 中间件
// 用于处理请求的横切关注点，如认证、授权、日志、CORS 等
package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/qzt/backend/internal/config"
)

// CORS 跨域资源共享中间件
// 处理跨域请求，允许指定的源、方法和头部访问 API
// 参数：
//   - cfg: CORS 配置对象
// 返回值：
//   - gin.HandlerFunc: Gin 中间件函数
func CORS(cfg *config.CORSConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取请求的源（Origin）
		origin := c.Request.Header.Get("Origin")
		
		// ========== 第一步：检查 Origin 是否允许 ==========
		// 验证请求的源是否在允许的源列表中
		allowed := false
		for _, o := range cfg.AllowOrigins {
			// 匹配指定的源或通配符 "*"
			if o == origin || o == "*" {
				allowed = true
				break
			}
		}

		// 如果源允许，设置 Access-Control-Allow-Origin 响应头
		if allowed {
			c.Header("Access-Control-Allow-Origin", origin)
		}

		// ========== 第二步：设置允许的 HTTP 方法 ==========
		// 设置允许的 HTTP 方法（GET、POST、PUT、DELETE 等）
		c.Header("Access-Control-Allow-Methods", joinStrings(cfg.AllowMethods, ", "))

		// ========== 第三步：设置允许的请求头 ==========
		// 设置允许的请求头（如 Content-Type、Authorization 等）
		c.Header("Access-Control-Allow-Headers", joinStrings(cfg.AllowHeaders, ", "))

		// ========== 第四步：设置暴露的响应头 ==========
		// 设置允许浏览器访问的响应头
		c.Header("Access-Control-Expose-Headers", joinStrings(cfg.ExposeHeaders, ", "))
		
		// ========== 第五步：设置允许携带凭证 ==========
		// 允许请求携带 Cookie 等凭证信息
		if cfg.AllowCredentials {
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		
		// ========== 第六步：设置预检请求缓存时间 ==========
		// 设置预检请求（OPTIONS）的缓存时间（秒）
		if cfg.MaxAge > 0 {
			c.Header("Access-Control-Max-Age", string(rune(cfg.MaxAge)))
		}

		// ========== 第七步：处理预检请求 ==========
		// 对于浏览器的预检请求（OPTIONS 方法），直接返回 204 状态码
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		// 继续执行后续处理器
		c.Next()
	}
}

// joinStrings 将字符串数组连接成一个字符串
// 参数：
//   - strs: 字符串数组
//   - sep: 分隔符
// 返回值：
//   - string: 连接后的字符串
func joinStrings(strs []string, sep string) string {
	result := ""
	for i, s := range strs {
		if i > 0 {
			result += sep
		}
		result += s
	}
	return result
}
