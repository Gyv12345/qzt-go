// Package service 定义了业务逻辑层
// 实现了核心业务规则和流程编排
package service

import (
	"context"
	"encoding/json"

	"github.com/qzt/backend/internal/domain"
	"github.com/qzt/backend/internal/repository"
)

// CustomerService 客户服务接口
// 定义了客户管理相关的业务方法
type CustomerService interface {
	// Create 创建客户
	Create(ctx context.Context, req *CreateCustomerReq) (*domain.Customer, error)
	// GetByID 获取客户详情
	GetByID(ctx context.Context, id string) (*domain.Customer, error)
	// List 获取客户列表
	List(ctx context.Context, params repository.ListParams) (*ListResult, error)
	// Update 更新客户信息
	Update(ctx context.Context, id string, req *UpdateCustomerReq) (*domain.Customer, error)
	// Delete 删除客户
	Delete(ctx context.Context, id string) error
}

// customerService 客户服务实现
type customerService struct {
	repo repository.CustomerRepository // 客户仓储
}

// NewCustomerService 创建客户服务实例
// 参数：
//   - repo: 客户仓储接口
// 返回值：
//   - CustomerService: 客户服务接口
func NewCustomerService(repo repository.CustomerRepository) CustomerService {
	return &customerService{repo: repo}
}

// CreateCustomerReq 创建客户请求
type CreateCustomerReq struct {
	Name          string   `json:"name" binding:"required"`          // 客户名称（必填）
	ShortName     string   `json:"shortName"`                        // 客户简称
	Industry      string   `json:"industry"`                         // 所属行业
	Scale         string   `json:"scale"`                            // 企业规模
	Address       string   `json:"address"`                          // 详细地址
	Website       string   `json:"website"`                          // 官网地址
	SourceChannel string   `json:"sourceChannel"`                    // 获客渠道
	Tags          []string `json:"tags"`                             // 标签列表
	Remark        string   `json:"remark"`                           // 备注信息
	FollowUserID  string   `json:"followUserId"`                     // 跟进人 ID
}

// Create 创建客户
// 创建新的客户记录，默认等级为线索（LEAD）
// 参数：
//   - ctx: 上下文
//   - req: 创建客户请求
// 返回值：
//   - *domain.Customer: 创建的客户对象
//   - error: 错误信息
func (s *customerService) Create(ctx context.Context, req *CreateCustomerReq) (*domain.Customer, error) {
	// 将标签数组序列化为 JSON 字符串
	tagsJSON, _ := json.Marshal(req.Tags)

	// 构建客户对象
	customer := &domain.Customer{
		Name:          req.Name,
		ShortName:     req.ShortName,
		Industry:      req.Industry,
		Scale:         req.Scale,
		Address:       req.Address,
		Website:       req.Website,
		SourceChannel: req.SourceChannel,
		Tags:          string(tagsJSON),
		Remark:        req.Remark,
		FollowUserID:  req.FollowUserID,
		CustomerLevel: domain.CustomerLevelLead, // 默认等级为线索
		Status:        1,                       // 默认状态为正常
	}

	// 调用仓储层创建客户
	if err := s.repo.Create(ctx, customer); err != nil {
		return nil, err
	}

	return customer, nil
}

// ListResult 列表查询结果
// 通用列表结果结构体
type ListResult struct {
	Items    []domain.Customer `json:"items"`    // 数据列表
	Total    int64             `json:"total"`    // 总记录数
	Page     int               `json:"page"`     // 当前页码
	PageSize int               `json:"pageSize"` // 每页数量
}

// List 获取客户列表
// 支持分页、关键词搜索、等级过滤和状态过滤
// 参数：
//   - ctx: 上下文
//   - params: 查询参数
// 返回值：
//   - *ListResult: 列表结果（包含数据和分页信息）
//   - error: 错误信息
func (s *customerService) List(ctx context.Context, params repository.ListParams) (*ListResult, error) {
	// 调用仓储层查询客户列表
	items, total, err := s.repo.List(ctx, params)
	if err != nil {
		return nil, err
	}

	return &ListResult{
		Items:    items,
		Total:    total,
		Page:     params.Page,
		PageSize: params.PageSize,
	}, nil
}

// GetByID 获取客户详情
// 根据客户 ID 获取详细信息
// 参数：
//   - ctx: 上下文
//   - id: 客户 ID
// 返回值：
//   - *domain.Customer: 客户对象
//   - error: 错误信息
func (s *customerService) GetByID(ctx context.Context, id string) (*domain.Customer, error) {
	return s.repo.GetByID(ctx, id)
}

// UpdateCustomerReq 更新客户请求
// 所有字段都是可选的（omitempty）
type UpdateCustomerReq struct {
	Name           string   `json:"name,omitempty"`           // 客户名称
	ShortName      string   `json:"shortName,omitempty"`      // 客户简称
	Industry       string   `json:"industry,omitempty"`       // 所属行业
	Scale          string   `json:"scale,omitempty"`          // 企业规模
	Address        string   `json:"address,omitempty"`        // 详细地址
	Website        string   `json:"website,omitempty"`        // 官网地址
	CustomerLevel  string   `json:"customerLevel,omitempty"`  // 客户等级
	SourceChannel  string   `json:"sourceChannel,omitempty"`  // 获客渠道
	Tags           []string `json:"tags,omitempty"`           // 标签列表
	Remark         string   `json:"remark,omitempty"`         // 备注信息
	FollowUserID   string   `json:"followUserId,omitempty"`   // 跟进人 ID
	Status         *int     `json:"status,omitempty"`         // 状态（使用指针支持可选）
}

// Update 更新客户信息
// 根据请求参数更新客户的部分或全部信息
// 参数：
//   - ctx: 上下文
//   - id: 客户 ID
//   - req: 更新客户请求
// 返回值：
//   - *domain.Customer: 更新后的客户对象
//   - error: 错误信息
func (s *customerService) Update(ctx context.Context, id string, req *UpdateCustomerReq) (*domain.Customer, error) {
	// 获取现有客户信息
	customer, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 根据请求参数更新客户信息（仅更新非空字段）
	if req.Name != "" {
		customer.Name = req.Name
	}
	if req.ShortName != "" {
		customer.ShortName = req.ShortName
	}
	if req.Industry != "" {
		customer.Industry = req.Industry
	}
	if req.Scale != "" {
		customer.Scale = req.Scale
	}
	if req.Address != "" {
		customer.Address = req.Address
	}
	if req.Website != "" {
		customer.Website = req.Website
	}
	if req.CustomerLevel != "" {
		customer.CustomerLevel = domain.CustomerLevel(req.CustomerLevel)
	}
	if req.SourceChannel != "" {
		customer.SourceChannel = req.SourceChannel
	}
	if req.Tags != nil {
		// 序列化标签数组
		tagsJSON, _ := json.Marshal(req.Tags)
		customer.Tags = string(tagsJSON)
	}
	if req.Remark != "" {
		customer.Remark = req.Remark
	}
	if req.FollowUserID != "" {
		customer.FollowUserID = req.FollowUserID
	}
	if req.Status != nil {
		customer.Status = *req.Status
	}

	// 调用仓储层更新客户
	if err := s.repo.Update(ctx, customer); err != nil {
		return nil, err
	}

	return customer, nil
}

// Delete 删除客户
// 根据客户 ID 删除客户记录
// 参数：
//   - ctx: 上下文
//   - id: 客户 ID
// 返回值：
//   - error: 错误信息
func (s *customerService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
