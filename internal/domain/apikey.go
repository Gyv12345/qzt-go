// Package domain 定义了系统的核心领域模型
// 这些模型代表了业务实体和它们之间的关系
package domain

import "time"

// APIKey API密钥模型
// 用于存储第三方系统调用本系统API所需的认证密钥信息
// 包括密钥本身、权限配置、使用限制和安全设置等
type APIKey struct {
	// ========== 基本信息 ==========
	ID          string     `gorm:"primaryKey;size:30" json:"id"`              // API Key 唯一标识符（主键）
	Name        string     `gorm:"size:100;not null" json:"name"`             // API Key 名称（便于识别）
	KeyHash     string     `gorm:"uniqueIndex;size:64;not null" json:"-"`     // 密钥哈希值（SHA256，不对外暴露）
	Prefix      string     `gorm:"size:10;not null" json:"prefix"`            // 密钥前缀（用于识别，如 "qzt_"）
	UserID      string     `gorm:"size:30;not null" json:"userId"`            // 所属用户ID

	// ========== 权限与安全配置 ==========
	Permissions string     `gorm:"type:text" json:"permissions"` // 权限列表（JSON格式：["customer:read", "contact:read"]）
	IPWhitelist string     `gorm:"type:text" json:"ipWhitelist"` // IP白名单（JSON格式：["192.168.1.0/24"]）
	RateLimit   int        `gorm:"default:100" json:"rateLimit"` // 速率限制（每分钟最大请求数）

	// ========== 时间信息 ==========
	ExpiresAt   *time.Time `json:"expiresAt"`  // 过期时间（nil表示永不过期）
	LastUsedAt  *time.Time `json:"lastUsedAt"` // 最后使用时间

	// ========== 状态信息 ==========
	Status      int        `gorm:"default:1" json:"status"`   // 状态（1:启用 0:禁用/已撤销）
	CreatedAt   time.Time  `json:"createdAt"`                 // 创建时间
	UpdatedAt   time.Time  `json:"updatedAt"`                 // 更新时间

	// ========== 关联关系 ==========
	// 所属用户信息（多对一关联）
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName 指定 API Key 模型对应的数据库表名
// 返回值：
//   - string: 数据库表名 "api_keys"
func (APIKey) TableName() string {
	return "api_keys"
}
