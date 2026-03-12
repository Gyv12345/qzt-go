// Package service 定义了业务逻辑层
// 实现了核心业务规则和流程编排
package service

import "errors"

// 定义业务错误常量
// 这些错误用于表示业务操作中的常见错误情况

var (
	// ErrInvalidCredentials 凭据无效错误
	// 表示用户名或密码错误，登录认证失败
	ErrInvalidCredentials = errors.New("用户名或密码错误")
	
	// ErrUserInactive 用户未激活错误
	// 表示用户账户已被禁用，无法进行登录等操作
	ErrUserInactive = errors.New("用户已被禁用")
	
	// ErrUsernameExists 用户名已存在错误
	// 表示注册时用户名已被其他用户使用
	ErrUsernameExists = errors.New("用户名已存在")
	
	// ErrInvalidToken 令牌无效错误
	// 表示 JWT Token 格式错误或签名验证失败
	ErrInvalidToken = errors.New("无效的 token")
	
	// ErrPermissionDenied 权限不足错误
	// 表示用户没有执行某项操作的权限
	ErrPermissionDenied = errors.New("权限不足")
	
	// ErrNotFound 资源不存在错误
	// 表示请求的资源在系统中不存在
	ErrNotFound = errors.New("资源不存在")
)
