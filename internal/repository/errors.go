// Package repository 定义了数据访问层
package repository

import "errors"

// 仓储层通用错误
var (
	// ErrNotFound 资源不存在
	ErrNotFound = errors.New("资源不存在")
	
	// ErrDuplicate 重复记录
	ErrDuplicate = errors.New("记录已存在")
	
	// ErrInvalidID 无效 ID
	ErrInvalidID = errors.New("无效的 ID")
)
