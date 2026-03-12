// Package errors 提供带错误码的错误类型
// 支持错误包装和堆栈追踪
package errors

import (
	"fmt"

	"github.com/qzt/backend/pkg/response"
)

// AppError 应用错误类型
// 包含错误码、消息和原始错误
type AppError struct {
	Code    response.Code // 错误码
	Message string        // 错误消息
	Cause   error         // 原始错误
}

// Error 实现 error 接口
func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%d] %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// Unwrap 实现错误解包
func (e *AppError) Unwrap() error {
	return e.Cause
}

// New 创建新的应用错误
func New(code response.Code, message string) *AppError {
	if message == "" {
		message = response.GetMessage(code)
	}
	return &AppError{
		Code:    code,
		Message: message,
	}
}

// Wrap 包装错误
func Wrap(code response.Code, message string, cause error) *AppError {
	if message == "" {
		message = response.GetMessage(code)
	}
	return &AppError{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}

// WrapErr 使用默认消息包装错误
func WrapErr(code response.Code, cause error) *AppError {
	return &AppError{
		Code:    code,
		Message: response.GetMessage(code),
		Cause:   cause,
	}
}

// Is 检查错误是否为指定的错误码
func Is(err error, code response.Code) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == code
	}
	return false
}

// GetCode 从错误中获取错误码
func GetCode(err error) response.Code {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code
	}
	return response.CodeUnknownError
}

// 预定义常用错误
var (
	ErrInvalidParams   = New(response.CodeInvalidParams, "")
	ErrUnauthorized    = New(response.CodeUnauthorized, "")
	ErrForbidden       = New(response.CodeForbidden, "")
	ErrNotFound        = New(response.CodeNotFound, "")
	ErrInternal        = New(response.CodeInternalError, "")
	ErrTooManyRequests = New(response.CodeTooManyRequests, "")

	// 认证相关
	ErrInvalidCredentials = New(response.CodeInvalidCredentials, "")
	ErrUserInactive       = New(response.CodeUserInactive, "")
	ErrUsernameExists     = New(response.CodeUsernameExists, "")
	ErrInvalidToken       = New(response.CodeInvalidToken, "")
	ErrTokenExpired       = New(response.CodeTokenExpired, "")

	// 客户相关
	ErrCustomerNotFound = New(response.CodeCustomerNotFound, "")
	ErrCustomerExists   = New(response.CodeCustomerExists, "")

	// 联系人相关
	ErrContactNotFound = New(response.CodeContactNotFound, "")

	// 合同相关
	ErrContractNotFound = New(response.CodeContractNotFound, "")

	// API Key 相关
	ErrAPIKeyInvalid = New(response.CodeAPIKeyInvalid, "")
	ErrAPIKeyExpired = New(response.CodeAPIKeyExpired, "")
	ErrAPIKeyRevoked = New(response.CodeAPIKeyRevoked, "")
)
