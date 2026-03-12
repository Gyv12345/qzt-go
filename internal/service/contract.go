// Package service 定义了业务逻辑层
// 实现了核心业务规则和流程编排
package service

import (
	"context"
	"encoding/json"

	"github.com/qzt/backend/internal/domain"
	"github.com/qzt/backend/internal/repository"
)

// ContractService 合同服务接口
// 定义了合同管理相关的业务方法
type ContractService interface {
	// Create 创建合同
	Create(ctx context.Context, req *CreateContractReq) (*domain.Contract, error)
	// GetByID 获取合同详情
	GetByID(ctx context.Context, id string) (*domain.Contract, error)
	// List 获取合同列表
	List(ctx context.Context, params repository.ContractListParams) (*ContractListResult, error)
	// Update 更新合同信息
	Update(ctx context.Context, id string, req *UpdateContractReq) (*domain.Contract, error)
	// Delete 删除合同
	Delete(ctx context.Context, id string) error
}

// contractService 合同服务实现
type contractService struct {
	repo repository.ContractRepository // 合同仓储
}

// NewContractService 创建合同服务实例
// 参数：
//   - repo: 合同仓储接口
// 返回值：
//   - ContractService: 合同服务接口
func NewContractService(repo repository.ContractRepository) ContractService {
	return &contractService{repo: repo}
}

// CreateContractReq 创建合同请求
type CreateContractReq struct {
	CustomerID  string    `json:"customerId" binding:"required"` // 客户 ID（必填）
	Title       string    `json:"title" binding:"required"`      // 合同标题（必填）
	Type        string    `json:"type"`                          // 合同类型
	Amount      float64   `json:"amount"`                        // 合同金额
	OwnerUserID string    `json:"ownerUserId"`                   // 合同负责人 ID
	Tags        []string  `json:"tags"`                          // 标签列表
	Remark      string    `json:"remark"`                        // 备注信息
}

// Create 创建合同
// 创建新的合同记录，默认状态为草稿（DRAFT）
// 参数：
//   - ctx: 上下文
//   - req: 创建合同请求
// 返回值：
//   - *domain.Contract: 创建的合同对象
//   - error: 错误信息
func (s *contractService) Create(ctx context.Context, req *CreateContractReq) (*domain.Contract, error) {
	// 将标签数组序列化为 JSON 字符串
	tagsJSON, _ := json.Marshal(req.Tags)

	// 构建合同对象
	contract := &domain.Contract{
		CustomerID:  req.CustomerID,
		Title:       req.Title,
		Type:        req.Type,
		Amount:      req.Amount,
		OwnerUserID: req.OwnerUserID,
		Tags:        string(tagsJSON),
		Remark:      req.Remark,
		Status:      domain.ContractStatusDraft, // 默认状态为草稿
	}

	// 调用仓储层创建合同
	if err := s.repo.Create(ctx, contract); err != nil {
		return nil, err
	}

	return contract, nil
}

// ContractListResult 合同列表查询结果
type ContractListResult struct {
	Items    []domain.Contract `json:"items"`    // 数据列表
	Total    int64             `json:"total"`    // 总记录数
	Page     int               `json:"page"`     // 当前页码
	PageSize int               `json:"pageSize"` // 每页数量
}

// List 获取合同列表
// 支持分页、按客户过滤、按状态过滤和关键词搜索
// 参数：
//   - ctx: 上下文
//   - params: 查询参数
// 返回值：
//   - *ContractListResult: 列表结果（包含数据和分页信息）
//   - error: 错误信息
func (s *contractService) List(ctx context.Context, params repository.ContractListParams) (*ContractListResult, error) {
	// 调用仓储层查询合同列表
	items, total, err := s.repo.List(ctx, params)
	if err != nil {
		return nil, err
	}

	return &ContractListResult{
		Items:    items,
		Total:    total,
		Page:     params.Page,
		PageSize: params.PageSize,
	}, nil
}

// GetByID 获取合同详情
// 根据合同 ID 获取详细信息
// 参数：
//   - ctx: 上下文
//   - id: 合同 ID
// 返回值：
//   - *domain.Contract: 合同对象
//   - error: 错误信息
func (s *contractService) GetByID(ctx context.Context, id string) (*domain.Contract, error) {
	return s.repo.GetByID(ctx, id)
}

// UpdateContractReq 更新合同请求
// 所有字段都是可选的（omitempty）
type UpdateContractReq struct {
	Title       string    `json:"title,omitempty"`       // 合同标题
	Type        string    `json:"type,omitempty"`        // 合同类型
	Amount      *float64  `json:"amount,omitempty"`      // 合同金额（使用指针支持可选）
	PaidAmount  *float64  `json:"paidAmount,omitempty"`  // 已支付金额（使用指针支持可选）
	Status      string    `json:"status,omitempty"`      // 合同状态
	OwnerUserID string    `json:"ownerUserId,omitempty"` // 合同负责人 ID
	Tags        []string  `json:"tags,omitempty"`        // 标签列表
	Remark      string    `json:"remark,omitempty"`      // 备注信息
}

// Update 更新合同信息
// 根据请求参数更新合同的部分或全部信息
// 参数：
//   - ctx: 上下文
//   - id: 合同 ID
//   - req: 更新合同请求
// 返回值：
//   - *domain.Contract: 更新后的合同对象
//   - error: 错误信息
func (s *contractService) Update(ctx context.Context, id string, req *UpdateContractReq) (*domain.Contract, error) {
	// 获取现有合同信息
	contract, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 根据请求参数更新合同信息（仅更新非空字段）
	if req.Title != "" {
		contract.Title = req.Title
	}
	if req.Type != "" {
		contract.Type = req.Type
	}
	if req.Amount != nil {
		contract.Amount = *req.Amount
	}
	if req.PaidAmount != nil {
		contract.PaidAmount = *req.PaidAmount
	}
	if req.Status != "" {
		contract.Status = domain.ContractStatus(req.Status)
	}
	if req.OwnerUserID != "" {
		contract.OwnerUserID = req.OwnerUserID
	}
	if req.Tags != nil {
		// 序列化标签数组
		tagsJSON, _ := json.Marshal(req.Tags)
		contract.Tags = string(tagsJSON)
	}
	if req.Remark != "" {
		contract.Remark = req.Remark
	}

	// 调用仓储层更新合同
	if err := s.repo.Update(ctx, contract); err != nil {
		return nil, err
	}

	return contract, nil
}

// Delete 删除合同
// 根据合同 ID 删除合同记录
// 参数：
//   - ctx: 上下文
//   - id: 合同 ID
// 返回值：
//   - error: 错误信息
func (s *contractService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
