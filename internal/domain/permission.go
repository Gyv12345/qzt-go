package domain

import "time"

// Permission 权限模型
// 定义系统中所有可访问的资源权限
type Permission struct {
	ID          string    `gorm:"primaryKey;size:30" json:"id"`
	Name        string    `gorm:"size:50;not null" json:"name"`        // 权限名称
	Code        string    `gorm:"uniqueIndex;size:100;not null" json:"code"` // 权限代码（如：customer:create）
	Type        string    `gorm:"size:20" json:"type"`                 // 权限类型（menu:菜单 api:接口 button:按钮）
	Resource    string    `gorm:"size:100" json:"resource"`            // 资源标识（如：/api/customers）
	Method      string    `gorm:"size:10" json:"method"`               // HTTP方法（GET/POST/PUT/DELETE）
	Description string    `gorm:"size:255" json:"description"`
	Status      int       `gorm:"default:1" json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func (Permission) TableName() string {
	return "permissions"
}

// UserRole 用户角色关联表
type UserRole struct {
	UserID    string    `gorm:"primaryKey;size:30" json:"userId"`
	RoleID    string    `gorm:"primaryKey;size:30" json:"roleId"`
	CreatedAt time.Time `json:"createdAt"`
}

func (UserRole) TableName() string {
	return "user_roles"
}

// RolePermission 角色权限关联表
type RolePermission struct {
	RoleID       string    `gorm:"primaryKey;size:30" json:"roleId"`
	PermissionID string    `gorm:"primaryKey;size:30" json:"permissionId"`
	CreatedAt    time.Time `json:"createdAt"`
}

func (RolePermission) TableName() string {
	return "role_permissions"
}
