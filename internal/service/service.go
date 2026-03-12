// Package service 定义了业务逻辑层
// 实现了核心业务规则和流程编排
package service

import (
	"github.com/qzt/backend/internal/config"
	"github.com/qzt/backend/internal/repository"
)

// Service 服务聚合
// 聚合了所有业务服务接口，便于统一管理和依赖注入
type Service struct {
	// Auth 认证服务（登录、注册、令牌刷新）
	Auth         AuthService
	// Customer 客户服务（客户管理）
	Customer     CustomerService
	// Contact 联系人服务（联系人管理）
	Contact      ContactService
	// Contract 合同服务（合同管理）
	Contract     ContractService
	// APIKey API Key 服务（第三方接口认证）
	APIKey       APIKeyService
	// RBAC 权限服务（角色权限管理）
	RBAC         RBACService
	// FollowRecord 跟进记录服务
	FollowRecord FollowRecordService
	// Opportunity 商机服务
	Opportunity  OpportunityService
	// Payment 回款服务
	Payment      PaymentService
	// Report 报表服务
	Report       ReportService
}

// NewService 创建服务聚合实例
// 初始化所有业务服务，并进行依赖注入
// 参数：
//   - repos: 仓储聚合对象（数据访问层）
//   - cfg: 应用配置对象
// 返回值：
//   - *Service: 服务聚合对象
func NewService(repos *repository.Repository, cfg *config.Config) *Service {
	return &Service{
		// 创建认证服务，注入用户仓储和配置
		Auth:         NewAuthService(repos.User, cfg),
		// 创建客户服务，注入客户仓储
		Customer:     NewCustomerService(repos.Customer),
		// 创建联系人服务，注入联系人仓储
		Contact:      NewContactService(repos.Contact),
		// 创建合同服务，注入合同仓储
		Contract:     NewContractService(repos.Contract),
		// 创建 API Key 服务，注入 API Key 仓储和配置
		APIKey:       NewAPIKeyService(repos.APIKey, cfg),
		// 创建 RBAC 服务，注入角色、权限、用户仓储
		RBAC:         NewRBACService(repos.Role, repos.Permission, repos.User),
		// 创建跟进记录服务
		FollowRecord: NewFollowRecordService(repos.FollowRecord),
		// 创建商机服务
		Opportunity:  NewOpportunityService(repos.Opportunity, repos.Customer, repos.Contract),
		// 创建回款服务
		Payment:      NewPaymentService(repos.Payment),
		// 创建报表服务
		Report:       NewReportService(repos.Customer, repos.Opportunity, repos.Contract, repos.Payment),
	}
}
