package domain

import "time"

// Notification 通知消息
type Notification struct {
	ID          string    `json:"id" gorm:"primaryKey;size:30"`
	UserID      string    `json:"userId" gorm:"size:30;index;not null;comment:接收用户ID"`
	Title       string    `json:"title" gorm:"size:200;not null;comment:标题"`
	Content     string    `json:"content" gorm:"type:text;comment:内容"`
	Type        string    `json:"type" gorm:"size:20;not null;comment:类型(SYSTEM/BUSINESS/REMINDER)"`
	Category    string    `json:"category" gorm:"size:50;comment:分类"`
	RelatedID   string    `json:"relatedId" gorm:"size:30;comment:关联业务ID"`
	RelatedType string    `json:"relatedType" gorm:"size:30;comment:关联业务类型"`
	IsRead      bool      `json:"isRead" gorm:"default:false;comment:是否已读"`
	ReadAt      *time.Time `json:"readAt" gorm:"comment:阅读时间"`
	Priority    int       `json:"priority" gorm:"default:0;comment:优先级(0普通/1重要/2紧急)"`
	CreatedAt   time.Time `json:"createdAt" gorm:"autoCreateTime"`
}

// OperationLog 操作日志
type OperationLog struct {
	ID         string    `json:"id" gorm:"primaryKey;size:30"`
	UserID     string    `json:"userId" gorm:"size:30;index;comment:操作用户ID"`
	Username   string    `json:"username" gorm:"size:100;comment:操作用户名"`
	Method     string    `json:"method" gorm:"size:10;not null;comment:请求方法"`
	Path       string    `json:"path" gorm:"size:500;not null;comment:请求路径"`
	Query      string    `json:"query" gorm:"type:text;comment:查询参数"`
	Body       string    `json:"body" gorm:"type:text;comment:请求体"`
	IP         string    `json:"ip" gorm:"size:50;comment:IP地址"`
	UserAgent  string    `json:"userAgent" gorm:"size:500;comment:用户代理"`
	StatusCode int       `json:"statusCode" gorm:"comment:响应状态码"`
	Latency    int64     `json:"latency" gorm:"comment:响应时间(毫秒)"`
	Error      string    `json:"error" gorm:"type:text;comment:错误信息"`
	CreatedAt  time.Time `json:"createdAt" gorm:"autoCreateTime;index"`
}

// TableName 指定表名
func (Notification) TableName() string {
	return "notifications"
}

func (OperationLog) TableName() string {
	return "operation_logs"
}
