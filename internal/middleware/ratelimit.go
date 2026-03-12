// Package middleware 提供 HTTP 中间件
package middleware

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/qzt/backend/pkg/response"
)

// RateLimiter 限流器接口
type RateLimiter interface {
	Allow(key string) bool
}

// TokenBucketLimiter 令牌桶限流器
type TokenBucketLimiter struct {
	rate       int           // 令牌生成速率（个/秒）
	capacity   int           // 桶容量
	tokens     map[string]int // 每个 key 的令牌数
	lastUpdate map[string]time.Time
	mu         sync.Mutex
}

// NewTokenBucketLimiter 创建令牌桶限流器
// rate: 每秒生成的令牌数
// capacity: 桶的最大容量
func NewTokenBucketLimiter(rate, capacity int) *TokenBucketLimiter {
	return &TokenBucketLimiter{
		rate:       rate,
		capacity:   capacity,
		tokens:     make(map[string]int),
		lastUpdate: make(map[string]time.Time),
	}
}

// Allow 检查是否允许请求
func (l *TokenBucketLimiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	
	// 初始化
	if _, exists := l.tokens[key]; !exists {
		l.tokens[key] = l.capacity
		l.lastUpdate[key] = now
		return true
	}

	// 计算需要补充的令牌
	elapsed := now.Sub(l.lastUpdate[key])
	tokensToAdd := int(elapsed.Seconds()) * l.rate
	
	// 补充令牌
	l.tokens[key] += tokensToAdd
	if l.tokens[key] > l.capacity {
		l.tokens[key] = l.capacity
	}
	l.lastUpdate[key] = now

	// 检查是否有令牌
	if l.tokens[key] > 0 {
		l.tokens[key]--
		return true
	}

	return false
}

// SlidingWindowLimiter 滑动窗口限流器
type SlidingWindowLimiter struct {
	limit    int           // 窗口内最大请求数
	window   time.Duration // 窗口大小
	requests map[string][]time.Time
	mu       sync.Mutex
}

// NewSlidingWindowLimiter 创建滑动窗口限流器
func NewSlidingWindowLimiter(limit int, window time.Duration) *SlidingWindowLimiter {
	return &SlidingWindowLimiter{
		limit:    limit,
		window:   window,
		requests: make(map[string][]time.Time),
	}
}

// Allow 检查是否允许请求
func (l *SlidingWindowLimiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-l.window)

	// 获取该 key 的请求记录
	requests, exists := l.requests[key]
	if !exists {
		l.requests[key] = []time.Time{now}
		return true
	}

	// 过滤掉窗口外的请求
	validRequests := []time.Time{}
	for _, t := range requests {
		if t.After(windowStart) {
			validRequests = append(validRequests, t)
		}
	}

	// 检查是否超过限制
	if len(validRequests) >= l.limit {
		l.requests[key] = validRequests
		return false
	}

	// 添加当前请求
	validRequests = append(validRequests, now)
	l.requests[key] = validRequests
	return true
}

// RateLimit 限流中间件
// 使用 IP 作为限流 key
func RateLimit(limiter RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.ClientIP()
		
		if !limiter.Allow(key) {
			response.Error(c, response.CodeTooManyRequests)
			c.Abort()
			return
		}
		
		c.Next()
	}
}

// RateLimitByUser 基于用户的限流中间件
// 使用用户 ID 作为限流 key，未登录用户使用 IP
func RateLimitByUser(limiter RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		var key string
		
		// 尝试获取用户 ID
		if userID, exists := c.Get("userId"); exists {
			key = userID.(string)
		} else {
			key = c.ClientIP()
		}
		
		if !limiter.Allow(key) {
			response.Error(c, response.CodeTooManyRequests)
			c.Abort()
			return
		}
		
		c.Next()
	}
}

// RateLimitByAPIKey 基于 API Key 的限流中间件
func RateLimitByAPIKey(limiter RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		var key string
		
		// 尝试获取 API Key ID
		if apiKeyID, exists := c.Get("apiKeyId"); exists {
			key = apiKeyID.(string)
		} else {
			key = c.ClientIP()
		}
		
		if !limiter.Allow(key) {
			response.Error(c, response.CodeTooManyRequests)
			c.Abort()
			return
		}
		
		c.Next()
	}
}
