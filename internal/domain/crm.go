package domain

import "time"

// FollowRecord 客户跟进记录
type FollowRecord struct {
	ID           string     `gorm:"primaryKey;size:30" json:"id"`
	CustomerID   string     `gorm:"size:30;not null;index" json:"customerId"` // 关联客户
	UserID       string     `gorm:"size:30;not null;index" json:"userId"`     // 跟进人
	Type         string     `gorm:"size:20" json:"type"`                      // 跟进方式（电话/邮件/拜访/微信）
	Content      string     `gorm:"type:text;not null" json:"content"`        // 跟进内容
	NextFollowAt *time.Time `json:"nextFollowAt"`                             // 下次跟进时间
	Status       string     `gorm:"size:20" json:"status"`                    // 跟进状态
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`

	// 关联关系
	Customer *Customer `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	User     *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (FollowRecord) TableName() string {
	return "follow_records"
}

// Opportunity 商机/销售线索
type Opportunity struct {
	ID              string     `gorm:"primaryKey;size:30" json:"id"`
	Name            string     `gorm:"size:100;not null" json:"name"`           // 商机名称
	CustomerID      string     `gorm:"size:30;not null;index" json:"customerId"` // 关联客户
	UserID          string     `gorm:"size:30;not null;index" json:"userId"`     // 负责人
	Amount          float64    `gorm:"type:decimal(12,2)" json:"amount"`        // 商机金额
	Stage           string     `gorm:"size:20;default:'LEAD'" json:"stage"`      // 商机阶段
	Probability     int        `gorm:"default:0" json:"probability"`            // 成功概率（0-100）
	ExpectedCloseAt *time.Time `json:"expectedCloseAt"`                         // 预计成交日期
	Description     string     `gorm:"type:text" json:"description"`            // 商机描述
	Status          string     `gorm:"size:20;default:'ACTIVE'" json:"status"`  // 状态
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       time.Time  `json:"updatedAt"`

	Customer *Customer `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	User     *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (Opportunity) TableName() string {
	return "opportunities"
}

// Payment 回款记录
type Payment struct {
	ID           string     `gorm:"primaryKey;size:30" json:"id"`
	ContractID   string     `gorm:"size:30;not null;index" json:"contractId"` // 关联合同
	CustomerID   string     `gorm:"size:30;not null;index" json:"customerId"` // 关联客户
	Amount       float64    `gorm:"type:decimal(12,2);not null" json:"amount"` // 回款金额
	PaymentDate  time.Time  `gorm:"not null" json:"paymentDate"`               // 回款日期
	PaymentType  string     `gorm:"size:20" json:"paymentType"`                // 回款方式（银行转账/现金/支票）
	Status       string     `gorm:"size:20;default:'PENDING'" json:"status"`  // 状态
	Remark       string     `gorm:"type:text" json:"remark"`                  // 备注
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`

	Contract *Contract `gorm:"foreignKey:ContractID" json:"contract,omitempty"`
	Customer *Customer `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
}

func (Payment) TableName() string {
	return "payments"
}

// Product 产品
type Product struct {
	ID          string    `gorm:"primaryKey;size:30" json:"id"`
	Name        string    `gorm:"size:100;not null" json:"name"`        // 产品名称
	Code        string    `gorm:"uniqueIndex;size:50" json:"code"`      // 产品编码
	Category    string    `gorm:"size:50" json:"category"`              // 产品分类
	Price       float64   `gorm:"type:decimal(10,2)" json:"price"`      // 价格
	Unit        string    `gorm:"size:20" json:"unit"`                  // 单位
	Description string    `gorm:"type:text" json:"description"`         // 描述
	Status      int       `gorm:"default:1" json:"status"`              // 状态
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func (Product) TableName() string {
	return "products"
}

// Invoice 发票
type Invoice struct {
	ID           string     `gorm:"primaryKey;size:30" json:"id"`
	ContractID   string     `gorm:"size:30;not null;index" json:"contractId"`
	CustomerID   string     `gorm:"size:30;not null;index" json:"customerId"`
	InvoiceNo    string     `gorm:"uniqueIndex;size:50" json:"invoiceNo"`     // 发票号码
	Amount       float64    `gorm:"type:decimal(12,2);not null" json:"amount"` // 开票金额
	Type         string     `gorm:"size:20" json:"type"`                       // 发票类型（增值税专用/普通）
	Status       string     `gorm:"size:20;default:'PENDING'" json:"status"`  // 状态
	InvoicedAt   *time.Time `json:"invoicedAt"`                                // 开票日期
	Remark       string     `gorm:"type:text" json:"remark"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`

	Contract *Contract `gorm:"foreignKey:ContractID" json:"contract,omitempty"`
	Customer *Customer `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
}

func (Invoice) TableName() string {
	return "invoices"
}
