// Package domain 定义了系统的核心领域模型
// 这些模型代表了业务实体和它们之间的关系
package domain

import (
	"time"
)

// CustomerLevel 客户等级类型定义
// 用于表示客户在销售漏斗中的不同阶段
type CustomerLevel string

// 客户等级常量定义
const (
	// CustomerLevelLead 线索客户 - 初步接触，尚未明确需求
	CustomerLevelLead CustomerLevel = "LEAD"
	// CustomerLevelProspect 潜在客户 - 已表达兴趣，正在跟进中
	CustomerLevelProspect CustomerLevel = "PROSPECT"
	// CustomerLevelCustomer 正式客户 - 已签约或已购买产品/服务
	CustomerLevelCustomer CustomerLevel = "CUSTOMER"
	// CustomerLevelVIP VIP客户 - 重要客户，享受特殊待遇
	CustomerLevelVIP CustomerLevel = "VIP"
)

// Customer 客户模型
// 用于存储企业客户的基本信息、跟进状态和关联数据
// 这是 CRM 系统的核心实体之一
type Customer struct {
	// ========== 基本信息 ==========
	ID               string        `gorm:"primaryKey;size:30" json:"id"`                     // 客户唯一标识符（主键）
	Name             string        `gorm:"size:100;not null" json:"name" binding:"required"` // 客户全称（必填）
	ShortName        string        `gorm:"size:50" json:"shortName"`                         // 客户简称
	Code             string        `gorm:"uniqueIndex;size:50" json:"code"`                  // 客户编码（唯一索引）
	Industry         string        `gorm:"size:50" json:"industry"`                          // 所属行业
	Scale            string        `gorm:"size:50" json:"scale"`                             // 企业规模
	Address          string        `gorm:"size:255" json:"address"`                          // 详细地址
	Website          string        `gorm:"size:255" json:"website"`                          // 官网地址

	// ========== 客户分级与跟进 ==========
	CustomerLevel    CustomerLevel `gorm:"size:20;default:LEAD" json:"customerLevel"` // 客户等级（默认为线索）
	SourceChannel    string        `gorm:"size:50" json:"sourceChannel"`               // 获客渠道
	FollowUserID     string        `gorm:"size:30" json:"followUserId"`                // 跟进人ID
	Tags             string        `gorm:"type:text" json:"tags"`                      // 标签（JSON 数组格式）
	Remark           string        `gorm:"type:text" json:"remark"`                    // 备注信息
	Status           int           `gorm:"default:1" json:"status"`                    // 状态（1:正常 0:禁用）

	// ========== 时间节点 ==========
	FirstContactDate *time.Time    `json:"firstContactDate"` // 首次接触日期
	ContractDate     *time.Time    `json:"contractDate"`     // 签约日期
	CreatedAt        time.Time     `json:"createdAt"`        // 创建时间
	UpdatedAt        time.Time     `json:"updatedAt"`        // 更新时间

	// ========== 关联关系 ==========
	// 跟进人信息（多对一关联）
	FollowUser *User           `gorm:"foreignKey:FollowUserID" json:"followUser,omitempty"`
	// 该客户的所有联系人（一对多关联）
	Contacts   []Contact       `gorm:"foreignKey:CustomerID" json:"contacts,omitempty"`
	// 该客户的所有合同（一对多关联）
	Contracts  []Contract      `gorm:"foreignKey:CustomerID" json:"contracts,omitempty"`
}

// TableName 指定客户模型对应的数据库表名
// 返回值：
//   - string: 数据库表名 "customers"
func (Customer) TableName() string {
	return "customers"
}
