// Package handler 提供 HTTP 请求处理功能
//
// 该包包含所有 RESTful API 的处理器，负责：
//   - 接收 HTTP 请求
//   - 验证请求参数
//   - 调用服务层处理业务逻辑
//   - 返回 HTTP 响应
//
// 处理器结构：
//   - Handler: 聚合所有处理器的容器
//   - AuthHandler: 处理认证相关请求（登录、注册、刷新令牌）
//   - CustomerHandler: 处理客户管理请求（CRUD）
//   - ContactHandler: 处理联系人管理请求（CRUD）
//   - ContractHandler: 处理合同管理请求（CRUD）
//   - APIKeyHandler: 处理 API Key 管理请求（创建、列表、撤销）
//
// 典型用法：
//   services := service.NewService(repos, cfg)
//   handlers := handler.NewHandler(services)
//   r.POST("/api/auth/login", handlers.Auth.Login)
package handler

import (
	"github.com/qzt/backend/internal/service"
)

// Handler 聚合所有 HTTP 处理器的容器
//
// 统一管理所有业务模块的处理器，提供依赖注入的入口点。
// 通过该结构体可以访问所有模块的处理器实例。
//
// 使用示例：
//   handlers := handler.NewHandler(services)
//   // 使用认证处理器
//   r.POST("/login", handlers.Auth.Login)
//   // 使用客户处理器
//   r.GET("/customers", handlers.Customer.List)
type Handler struct {
	Auth         *AuthHandler         // 认证处理器（登录、注册、刷新令牌）
	Customer     *CustomerHandler     // 客户处理器（客户 CRUD 操作）
	Contact      *ContactHandler      // 联系人处理器（联系人 CRUD 操作）
	Contract     *ContractHandler     // 合同处理器（合同 CRUD 操作）
	APIKey       *APIKeyHandler       // API Key 处理器（API Key 管理）
	RBAC         *RBACHandler         // RBAC 处理器（角色权限管理）
	FollowRecord *FollowRecordHandler // 跟进记录处理器
	Opportunity  *OpportunityHandler  // 商机处理器
	Payment      *PaymentHandler      // 回款处理器
	Invoice      *InvoiceHandler      // 发票处理器
	Report       *ReportHandler       // 报表处理器
	Notification *NotificationHandler // 通知处理器
	OperationLog *OperationLogHandler // 操作日志处理器
	ImportExport *ImportExportHandler // 导入导出处理器
}

// NewHandler 创建处理器容器实例
//
// 初始化所有业务模块的处理器，并将它们聚合到 Handler 结构体中。
// 这是依赖注入的入口点，需要传入服务层实例。
//
// 参数：
//   - svc: 服务层实例，包含所有业务服务
//
// 返回值：
//   - *Handler: 聚合了所有处理器的容器实例
//
// 使用示例：
//   services := service.NewService(repos, cfg)
//   handlers := handler.NewHandler(services)
//   
//   // 在路由中使用
//   r := gin.New()
//   r.POST("/api/auth/login", handlers.Auth.Login)
//   r.GET("/api/customers", handlers.Customer.List)
func NewHandler(svc *service.Service) *Handler {
	return &Handler{
		Auth:         NewAuthHandler(svc.Auth),                       // 创建认证处理器
		Customer:     NewCustomerHandler(svc.Customer),               // 创建客户处理器
		Contact:      NewContactHandler(svc.Contact),                 // 创建联系人处理器
		Contract:     NewContractHandler(svc.Contract),               // 创建合同处理器
		APIKey:       NewAPIKeyHandler(svc.APIKey),                   // 创建 API Key 处理器
		RBAC:         NewRBACHandler(svc.RBAC),                       // 创建 RBAC 处理器
		FollowRecord: NewFollowRecordHandler(svc.FollowRecord),       // 创建跟进记录处理器
		Opportunity:  NewOpportunityHandler(svc.Opportunity),         // 创建商机处理器
		Payment:      NewPaymentHandler(svc.Payment),                 // 创建回款处理器
		Invoice:      NewInvoiceHandler(svc.Invoice),                 // 创建发票处理器
		Report:       NewReportHandler(svc.Report),                   // 创建报表处理器
		Notification: NewNotificationHandler(svc.Notification),       // 创建通知处理器
		OperationLog: NewOperationLogHandler(svc.OperationLog),       // 创建操作日志处理器
		ImportExport: NewImportExportHandler(svc.ImportExport),         // 创建导入导出处理器
	}
}
