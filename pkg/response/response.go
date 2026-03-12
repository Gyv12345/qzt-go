// Package response 提供统一的 API 响应格式
// 标准化所有 HTTP 响应的结构和错误码处理
package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Code 错误码类型
type Code int

// 定义标准错误码
// 格式：模块(2位) + 具体错误(3位)
const (
	// 通用错误码 100xx
	CodeSuccess         Code = 0     // 成功
	CodeUnknownError    Code = 10001 // 未知错误
	CodeInvalidParams   Code = 10002 // 参数错误
	CodeUnauthorized    Code = 10003 // 未授权
	CodeForbidden       Code = 10004 // 禁止访问
	CodeNotFound        Code = 10005 // 资源不存在
	CodeInternalError   Code = 10006 // 内部错误
	CodeTooManyRequests Code = 10007 // 请求过于频繁

	// 认证模块错误码 101xx
	CodeInvalidCredentials Code = 10101 // 用户名或密码错误
	CodeUserInactive       Code = 10102 // 用户未激活
	CodeUsernameExists     Code = 10103 // 用户名已存在
	CodeInvalidToken       Code = 10104 // Token 无效
	CodeTokenExpired       Code = 10105 // Token 已过期
	CodeRefreshFailed      Code = 10106 // Token 刷新失败

	// 客户模块错误码 102xx
	CodeCustomerNotFound Code = 10201 // 客户不存在
	CodeCustomerExists   Code = 10202 // 客户已存在
	CodeCustomerInvalid  Code = 10203 // 客户信息无效

	// 联系人模块错误码 103xx
	CodeContactNotFound Code = 10301 // 联系人不存在

	// 合同模块错误码 104xx
	CodeContractNotFound Code = 10401 // 合同不存在

	// API Key 模块错误码 105xx
	CodeAPIKeyInvalid Code = 10501 // API Key 无效
	CodeAPIKeyExpired Code = 10502 // API Key 已过期
	CodeAPIKeyRevoked Code = 10503 // API Key 已撤销
)

// 错误码消息映射
var codeMessages = map[Code]string{
	CodeSuccess:            "成功",
	CodeUnknownError:       "未知错误",
	CodeInvalidParams:      "参数错误",
	CodeUnauthorized:       "未授权",
	CodeForbidden:          "禁止访问",
	CodeNotFound:           "资源不存在",
	CodeInternalError:      "内部错误",
	CodeTooManyRequests:    "请求过于频繁",
	CodeInvalidCredentials: "用户名或密码错误",
	CodeUserInactive:       "用户未激活",
	CodeUsernameExists:     "用户名已存在",
	CodeInvalidToken:       "Token 无效",
	CodeTokenExpired:       "Token 已过期",
	CodeRefreshFailed:      "Token 刷新失败",
	CodeCustomerNotFound:   "客户不存在",
	CodeCustomerExists:     "客户已存在",
	CodeCustomerInvalid:    "客户信息无效",
	CodeContactNotFound:    "联系人不存在",
	CodeContractNotFound:   "合同不存在",
	CodeAPIKeyInvalid:      "API Key 无效",
	CodeAPIKeyExpired:      "API Key 已过期",
	CodeAPIKeyRevoked:      "API Key 已撤销",
}

// GetMessage 获取错误码对应的默认消息
func GetMessage(code Code) string {
	if msg, ok := codeMessages[code]; ok {
		return msg
	}
	return "未知错误"
}

// Response 统一响应结构
type Response struct {
	Code    Code        `json:"code"`    // 业务状态码
	Message string      `json:"message"` // 提示信息
	Data    interface{} `json:"data"`    // 响应数据
}

// PageData 分页数据结构
type PageData struct {
	Items    interface{} `json:"items"`    // 数据列表
	Total    int64       `json:"total"`    // 总记录数
	Page     int         `json:"page"`     // 当前页码
	PageSize int         `json:"pageSize"` // 每页数量
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: GetMessage(CodeSuccess),
		Data:    data,
	})
}

// SuccessWithMsg 成功响应（自定义消息）
func SuccessWithMsg(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: message,
		Data:    data,
	})
}

// SuccessPage 分页成功响应
func SuccessPage(c *gin.Context, items interface{}, total int64, page, pageSize int) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: GetMessage(CodeSuccess),
		Data: PageData{
			Items:    items,
			Total:    total,
			Page:     page,
			PageSize: pageSize,
		},
	})
}

// Created 创建成功响应
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Code:    CodeSuccess,
		Message: "创建成功",
		Data:    data,
	})
}

// NoContent 无内容响应
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// Error 业务错误响应
func Error(c *gin.Context, code Code) {
	httpStatus := getHTTPStatus(code)
	c.JSON(httpStatus, Response{
		Code:    code,
		Message: GetMessage(code),
		Data:    nil,
	})
}

// ErrorWithMsg 业务错误响应（自定义消息）
func ErrorWithMsg(c *gin.Context, code Code, message string) {
	httpStatus := getHTTPStatus(code)
	c.JSON(httpStatus, Response{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

// ErrorWithData 业务错误响应（带数据）
func ErrorWithData(c *gin.Context, code Code, data interface{}) {
	httpStatus := getHTTPStatus(code)
	c.JSON(httpStatus, Response{
		Code:    code,
		Message: GetMessage(code),
		Data:    data,
	})
}

// getHTTPStatus 根据业务错误码获取 HTTP 状态码
func getHTTPStatus(code Code) int {
	switch code {
	case CodeSuccess:
		return http.StatusOK
	case CodeInvalidParams:
		return http.StatusBadRequest
	case CodeUnauthorized, CodeInvalidToken, CodeTokenExpired:
		return http.StatusUnauthorized
	case CodeForbidden:
		return http.StatusForbidden
	case CodeNotFound, CodeCustomerNotFound, CodeContactNotFound, CodeContractNotFound:
		return http.StatusNotFound
	case CodeTooManyRequests:
		return http.StatusTooManyRequests
	default:
		return http.StatusInternalServerError
	}
}
