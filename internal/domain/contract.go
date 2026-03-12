// Package domain 定义了系统的核心领域模型
// 这些模型代表了业务实体和它们之间的关系
package domain

import "time"

// ContractStatus 合同状态类型定义
// 用于表示合同在不同生命周期阶段的状态
type ContractStatus string

// 合同状态常量定义
const (
	// ContractStatusDraft 草稿状态 - 合同正在编辑中，尚未提交
	ContractStatusDraft ContractStatus = "DRAFT"
	// ContractStatusPending 待审批状态 - 合同已提交，等待审批
	ContractStatusPending ContractStatus = "PENDING"
	// ContractStatusActive 执行中状态 - 合同已生效，正在履行
	ContractStatusActive ContractStatus = "ACTIVE"
	// ContractStatusCompleted 已完成状态 - 合同已履行完毕
	ContractStatusCompleted ContractStatus = "COMPLETED"
	// ContractStatusCancelled 已取消状态 - 合同已取消或作废
	ContractStatusCancelled ContractStatus = "CANCELLED"
)

// Contract 合同模型
// 用于存储企业与客户之间签订的合同信息
// 包括合同基本信息、金额、状态和时间节点等
type Contract struct {
	// ========== 基本信息 ==========
	ID           string         `gorm:"primaryKey;size:30" json:"id"`                    // 合同唯一标识符（主键）
	CustomerID   string         `gorm:"size:30;index" json:"customerId"`                 // 客户ID（外键，带索引）
	ContractNo   string         `gorm:"uniqueIndex;size:50" json:"contractNo"`           // 合同编号（唯一索引）
	Title        string         `gorm:"size:200;not null" json:"title" binding:"required"` // 合同标题（必填）
	Type         string         `gorm:"size:50" json:"type"`                             // 合同类型

	// ========== 金额信息 ==========
	Amount       float64        `gorm:"type:decimal(12,2)" json:"amount"`               // 合同总金额
	PaidAmount   float64        `gorm:"type:decimal(12,2);default:0" json:"paidAmount"` // 已支付金额

	// ========== 状态与时间 ==========
	Status       ContractStatus `gorm:"size:20;default:DRAFT" json:"status"` // 合同状态（默认为草稿）
	StartDate    *time.Time     `json:"startDate"`                           // 合同开始日期
	EndDate      *time.Time     `json:"endDate"`                             // 合同结束日期
	SignDate     *time.Time     `json:"signDate"`                            // 合同签订日期

	// ========== 管理信息 ==========
	OwnerUserID  string         `gorm:"size:30" json:"ownerUserId"`       // 合同负责人ID
	Tags         string         `gorm:"type:text" json:"tags"`            // 标签（JSON 数组格式）
	Remark       string         `gorm:"type:text" json:"remark"`          // 备注信息

	// ========== 时间戳 ==========
	CreatedAt    time.Time      `json:"createdAt"`  // 创建时间
	UpdatedAt    time.Time      `json:"updatedAt"`  // 更新时间

	// ========== 关联关系 ==========
	// 所属客户信息（多对一关联）
	Customer *Customer `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	// 合同负责人信息（多对一关联）
	Owner    *User     `gorm:"foreignKey:OwnerUserID" json:"owner,omitempty"`
}

// TableName 指定合同模型对应的数据库表名
// 返回值：
//   - string: 数据库表名 "contracts"
func (Contract) TableName() string {
	return "contracts"
}
