// Package domain 定义了系统的核心领域模型
// 这些模型代表了业务实体和它们之间的关系
package domain

import "time"

// UserStatus 用户状态类型定义
type UserStatus int

// 用户状态常量定义
const (
	// UserStatusInactive 未激活状态 - 用户已创建但未激活
	UserStatusInactive UserStatus = 0
	// UserStatusActive 激活状态 - 用户正常使用中
	UserStatusActive UserStatus = 1
)

// User 用户模型
// 用于存储系统用户的基本信息、认证信息和组织架构关联
// 这是系统权限管理和身份认证的核心实体
type User struct {
	// ========== 基本信息 ==========
	ID           string     `gorm:"primaryKey;size:30" json:"id"`          // 用户唯一标识符（主键）
	Username     string     `gorm:"uniqueIndex;size:50;not null" json:"username"` // 用户名（唯一索引，必填）
	Password     string     `gorm:"size:255;not null" json:"-"`            // 密码哈希值（不对外暴露）
	Email        string     `gorm:"size:100" json:"email"`                 // 电子邮箱
	Mobile       string     `gorm:"size:20" json:"mobile"`                 // 手机号码
	RealName     string     `gorm:"size:50" json:"realName"`               // 真实姓名
	Avatar       string     `gorm:"size:255" json:"avatar"`                // 头像URL

	// ========== 组织架构 ==========
	DepartmentID string     `gorm:"size:30" json:"departmentId"` // 所属部门ID

	// ========== 状态与时间 ==========
	Status       UserStatus `gorm:"default:1" json:"status"`     // 用户状态（1:激活 0:未激活）
	LastLoginAt  *time.Time `json:"lastLoginAt"`                 // 最后登录时间
	CreatedAt    time.Time  `json:"createdAt"`                   // 创建时间
	UpdatedAt    time.Time  `json:"updatedAt"`                   // 更新时间

	// ========== 关联关系 ==========
	// 所属部门信息（多对一关联）
	Department *Department `gorm:"foreignKey:DepartmentID" json:"department,omitempty"`
	// 用户拥有的角色（多对多关联）
	Roles      []Role      `gorm:"many2many:user_roles" json:"roles,omitempty"`
}

// TableName 指定用户模型对应的数据库表名
// 返回值：
//   - string: 数据库表名 "users"
func (User) TableName() string {
	return "users"
}

// Role 角色模型
// 用于定义系统角色，实现基于角色的访问控制（RBAC）
// 每个角色可以关联多个菜单权限
type Role struct {
	// ========== 基本信息 ==========
	ID          string    `gorm:"primaryKey;size:30" json:"id"`              // 角色唯一标识符（主键）
	Name        string    `gorm:"uniqueIndex;size:50;not null" json:"name"`  // 角色名称（唯一索引，必填）
	Code        string    `gorm:"uniqueIndex;size:50;not null" json:"code"`  // 角色代码（唯一索引，必填）
	Description string    `gorm:"size:255" json:"description"`               // 角色描述

	// ========== 状态与时间 ==========
	Status      int       `gorm:"default:1" json:"status"`   // 状态（1:启用 0:禁用）
	CreatedAt   time.Time `json:"createdAt"`                 // 创建时间
	UpdatedAt   time.Time `json:"updatedAt"`                 // 更新时间

	// ========== 关联关系 ==========
	// 角色拥有的菜单权限（多对多关联）
	Menus []Menu `gorm:"many2many:role_menus" json:"menus,omitempty"`
}

// TableName 指定角色模型对应的数据库表名
// 返回值：
//   - string: 数据库表名 "roles"
func (Role) TableName() string {
	return "roles"
}

// Department 部门模型
// 用于存储组织架构中的部门信息
// 支持树形结构，可以表示部门的上下级关系
type Department struct {
	// ========== 基本信息 ==========
	ID       string       `gorm:"primaryKey;size:30" json:"id"`           // 部门唯一标识符（主键）
	Name     string       `gorm:"size:50;not null" json:"name"`           // 部门名称（必填）
	Code     string       `gorm:"uniqueIndex;size:50" json:"code"`        // 部门代码（唯一索引）
	ParentID string       `gorm:"size:30" json:"parentId"`                // 父部门ID（用于构建树形结构）
	Sort     int          `gorm:"default:0" json:"sort"`                  // 排序号（同级部门排序）
	Status   int          `gorm:"default:1" json:"status"`                // 状态（1:启用 0:禁用）

	// ========== 关联关系 ==========
	// 子部门列表（一对多关联，用于树形结构）
	Children []Department `gorm:"foreignKey:ParentID" json:"children,omitempty"`
}

// TableName 指定部门模型对应的数据库表名
// 返回值：
//   - string: 数据库表名 "departments"
func (Department) TableName() string {
	return "departments"
}

// Menu 菜单模型
// 用于存储系统的菜单和按钮权限配置
// 支持树形结构，可以表示菜单的层级关系
type Menu struct {
	// ========== 基本信息 ==========
	ID       string    `gorm:"primaryKey;size:30" json:"id"`       // 菜单唯一标识符（主键）
	Name     string    `gorm:"size:50;not null" json:"name"`       // 菜单名称（必填）
	Code     string    `gorm:"uniqueIndex;size:50" json:"code"`    // 菜单代码（唯一索引）
	ParentID string    `gorm:"size:30" json:"parentId"`            // 父菜单ID（用于构建树形结构）

	// ========== 显示配置 ==========
	Path     string    `gorm:"size:255" json:"path"`               // 路由路径
	Icon     string    `gorm:"size:50" json:"icon"`                // 图标名称
	Sort     int       `gorm:"default:0" json:"sort"`              // 排序号
	Type     string    `gorm:"size:20" json:"type"`                // 类型（menu:菜单 button:按钮）
	Status   int       `gorm:"default:1" json:"status"`            // 状态（1:启用 0:禁用）
}

// TableName 指定菜单模型对应的数据库表名
// 返回值：
//   - string: 数据库表名 "menus"
func (Menu) TableName() string {
	return "menus"
}
