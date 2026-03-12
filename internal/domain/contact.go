// Package domain 定义了系统的核心领域模型
// 这些模型代表了业务实体和它们之间的关系
package domain

import "time"

// Contact 联系人模型
// 用于存储客户的联系人信息，包括个人基本信息、联系方式和角色属性
// 每个联系人关联一个客户，一个客户可以有多个联系人
type Contact struct {
	// ========== 基本信息 ==========
	ID          string    `gorm:"primaryKey;size:30" json:"id"`                  // 联系人唯一标识符（主键）
	CustomerID  string    `gorm:"size:30;index" json:"customerId"`               // 所属客户ID（外键，带索引）
	Name        string    `gorm:"size:50;not null" json:"name" binding:"required"` // 联系人姓名（必填）
	Position    string    `gorm:"size:50" json:"position"`                       // 职位
	Department  string    `gorm:"size:50" json:"department"`                     // 部门

	// ========== 联系方式 ==========
	Phone       string    `gorm:"size:20" json:"phone"`    // 固定电话
	Mobile      string    `gorm:"size:20" json:"mobile"`   // 手机号码
	Email       string    `gorm:"size:100" json:"email"`   // 电子邮箱
	Wechat      string    `gorm:"size:50" json:"wechat"`   // 微信号
	QQ          string    `gorm:"size:20" json:"qq"`       // QQ号

	// ========== 角色属性 ==========
	IsPrimary   bool      `gorm:"default:false" json:"isPrimary"`  // 是否为主要联系人
	IsDecision  bool      `gorm:"default:false" json:"isDecision"` // 是否为决策者

	// ========== 个人信息 ==========
	Gender      string    `gorm:"size:10" json:"gender"`      // 性别
	Birthday    *time.Time `json:"birthday"`                  // 生日
	Address     string    `gorm:"size:255" json:"address"`    // 联系地址
	Remark      string    `gorm:"type:text" json:"remark"`    // 备注信息
	Status      int       `gorm:"default:1" json:"status"`    // 状态（1:正常 0:禁用）

	// ========== 时间戳 ==========
	CreatedAt   time.Time `json:"createdAt"`  // 创建时间
	UpdatedAt   time.Time `json:"updatedAt"`  // 更新时间

	// ========== 关联关系 ==========
	// 所属客户信息（多对一关联）
	Customer *Customer `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
}

// TableName 指定联系人模型对应的数据库表名
// 返回值：
//   - string: 数据库表名 "contacts"
func (Contact) TableName() string {
	return "contacts"
}
