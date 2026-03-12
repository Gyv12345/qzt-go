// Package handler 定义了 HTTP 请求处理器
// 负责处理 HTTP 请求、参数验证和响应返回
package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qzt/backend/internal/service"
)

// AuthHandler 认证处理器
// 处理用户登录、注册和令牌刷新等认证相关的 HTTP 请求
type AuthHandler struct {
	svc service.AuthService // 认证服务接口
}

// NewAuthHandler 创建认证处理器实例
// 参数：
//   - svc: 认证服务接口实例
// 返回值：
//   - *AuthHandler: 认证处理器实例
func NewAuthHandler(svc service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

// LoginReq 登录请求结构体
// 包含用户登录所需的凭证信息
type LoginReq struct {
	Username string `json:"username" binding:"required"` // 用户名（必填）
	Password string `json:"password" binding:"required"` // 密码（必填）
}

// Login 用户登录
// 处理 POST /api/auth/login 请求
// 验证用户凭证并返回访问令牌和刷新令牌
//
// godoc
// @Summary 用户登录
// @Description 用户登录获取 token
// @Tags 认证
// @Accept json
// @Produce json
// @Param body body LoginReq true "登录信息"
// @Success 200 {object} service.LoginResult
// @Failure 401 {object} map[string]string
// @Router /api/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	// 绑定并验证请求参数
	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		// 参数验证失败，返回 400 错误
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 调用认证服务进行登录验证
	result, err := h.svc.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		// 登录失败（用户名或密码错误），返回 401 未授权错误
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// 登录成功，返回令牌信息
	c.JSON(http.StatusOK, result)
}

// Register 用户注册
// 处理 POST /api/auth/register 请求
// 创建新用户账号并返回用户信息
//
// godoc
// @Summary 用户注册
// @Tags 认证
// @Accept json
// @Produce json
// @Param body body service.RegisterReq true "注册信息"
// @Success 201 {object} domain.User
// @Failure 400 {object} map[string]string
// @Router /api/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	// 绑定并验证请求参数
	var req service.RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		// 参数验证失败，返回 400 错误
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 调用认证服务进行用户注册
	user, err := h.svc.Register(c.Request.Context(), &req)
	if err != nil {
		// 注册失败（用户名已存在等），返回 400 错误
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 注册成功，返回 201 状态码和用户信息
	c.JSON(http.StatusCreated, user)
}

// RefreshTokenReq 刷新令牌请求结构体
// 包含用于刷新访问令牌的刷新令牌
type RefreshTokenReq struct {
	RefreshToken string `json:"refreshToken" binding:"required"` // 刷新令牌（必填）
}

// RefreshToken 刷新访问令牌
// 处理 POST /api/auth/refresh 请求
// 使用刷新令牌获取新的访问令牌
//
// godoc
// @Summary 刷新 token
// @Tags 认证
// @Accept json
// @Produce json
// @Param body body RefreshTokenReq true "refresh token"
// @Success 200 {object} service.LoginResult
// @Failure 401 {object} map[string]string
// @Router /api/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	// 绑定并验证请求参数
	var req RefreshTokenReq
	if err := c.ShouldBindJSON(&req); err != nil {
		// 参数验证失败，返回 400 错误
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 调用认证服务刷新令牌
	result, err := h.svc.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		// 刷新失败（令牌无效或已过期），返回 401 未授权错误
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// 刷新成功，返回新的令牌信息
	c.JSON(http.StatusOK, result)
}
