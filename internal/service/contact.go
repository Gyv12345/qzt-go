// Package service 定义了业务逻辑层
// 实现了核心业务规则和流程编排
package service

import (
	"context"

	"github.com/qzt/backend/internal/domain"
	"github.com/qzt/backend/internal/repository"
)

// ContactService 联系人服务接口
// 定义了联系人管理相关的业务方法
type ContactService interface {
	// Create 创建联系人
	Create(ctx context.Context, req *CreateContactReq) (*domain.Contact, error)
	// GetByID 获取联系人详情
	GetByID(ctx context.Context, id string) (*domain.Contact, error)
	// List 获取联系人列表
	List(ctx context.Context, params repository.ContactListParams) (*ContactListResult, error)
	// Update 更新联系人信息
	Update(ctx context.Context, id string, req *UpdateContactReq) (*domain.Contact, error)
	// Delete 删除联系人
	Delete(ctx context.Context, id string) error
}

// contactService 联系人服务实现
type contactService struct {
	repo repository.ContactRepository // 联系人仓储
}

// NewContactService 创建联系人服务实例
// 参数：
//   - repo: 联系人仓储接口
// 返回值：
//   - ContactService: 联系人服务接口
func NewContactService(repo repository.ContactRepository) ContactService {
	return &contactService{repo: repo}
}

// CreateContactReq 创建联系人请求
type CreateContactReq struct {
	CustomerID string `json:"customerId" binding:"required"` // 客户 ID（必填）
	Name       string `json:"name" binding:"required"`       // 联系人姓名（必填）
	Position   string `json:"position"`                     // 职位
	Department string `json:"department"`                   // 部门
	Phone      string `json:"phone"`                        // 固定电话
	Mobile     string `json:"mobile"`                       // 手机号码
	Email      string `json:"email"`                        // 电子邮箱
	Wechat     string `json:"wechat"`                       // 微信号
	QQ         string `json:"qq"`                           // QQ 号码
	IsPrimary  bool   `json:"isPrimary"`                    // 是否为主要联系人
	IsDecision bool   `json:"isDecision"`                   // 是否为决策者
	Gender     string `json:"gender"`                       // 性别
	Address    string `json:"address"`                      // 联系地址
	Remark     string `json:"remark"`                       // 备注信息
}

// Create 创建联系人
// 创建新的联系人记录
// 参数：
//   - ctx: 上下文
//   - req: 创建联系人请求
// 返回值：
//   - *domain.Contact: 创建的联系人对象
//   - error: 错误信息
func (s *contactService) Create(ctx context.Context, req *CreateContactReq) (*domain.Contact, error) {
	// 构建联系人对象
	contact := &domain.Contact{
		CustomerID: req.CustomerID,
		Name:       req.Name,
		Position:   req.Position,
		Department: req.Department,
		Phone:      req.Phone,
		Mobile:     req.Mobile,
		Email:      req.Email,
		Wechat:     req.Wechat,
		QQ:         req.QQ,
		IsPrimary:  req.IsPrimary,
		IsDecision: req.IsDecision,
		Gender:     req.Gender,
		Address:    req.Address,
		Remark:     req.Remark,
		Status:     1, // 默认状态为正常
	}

	// 调用仓储层创建联系人
	if err := s.repo.Create(ctx, contact); err != nil {
		return nil, err
	}

	return contact, nil
}

// ContactListResult 联系人列表查询结果
type ContactListResult struct {
	Items    []domain.Contact `json:"items"`    // 数据列表
	Total    int64            `json:"total"`    // 总记录数
	Page     int              `json:"page"`     // 当前页码
	PageSize int              `json:"pageSize"` // 每页数量
}

// List 获取联系人列表
// 支持分页、按客户过滤和关键词搜索
// 参数：
//   - ctx: 上下文
//   - params: 查询参数
// 返回值：
//   - *ContactListResult: 列表结果（包含数据和分页信息）
//   - error: 错误信息
func (s *contactService) List(ctx context.Context, params repository.ContactListParams) (*ContactListResult, error) {
	// 调用仓储层查询联系人列表
	items, total, err := s.repo.List(ctx, params)
	if err != nil {
		return nil, err
	}

	return &ContactListResult{
		Items:    items,
		Total:    total,
		Page:     params.Page,
		PageSize: params.PageSize,
	}, nil
}

// GetByID 获取联系人详情
// 根据联系人 ID 获取详细信息
// 参数：
//   - ctx: 上下文
//   - id: 联系人 ID
// 返回值：
//   - *domain.Contact: 联系人对象
//   - error: 错误信息
func (s *contactService) GetByID(ctx context.Context, id string) (*domain.Contact, error) {
	return s.repo.GetByID(ctx, id)
}

// UpdateContactReq 更新联系人请求
// 所有字段都是可选的（omitempty）
type UpdateContactReq struct {
	Name       string `json:"name,omitempty"`       // 联系人姓名
	Position   string `json:"position,omitempty"`   // 职位
	Department string `json:"department,omitempty"` // 部门
	Phone      string `json:"phone,omitempty"`      // 固定电话
	Mobile     string `json:"mobile,omitempty"`     // 手机号码
	Email      string `json:"email,omitempty"`      // 电子邮箱
	Wechat     string `json:"wechat,omitempty"`     // 微信号
	QQ         string `json:"qq,omitempty"`         // QQ 号码
	IsPrimary  *bool  `json:"isPrimary,omitempty"`  // 是否为主要联系人（使用指针支持可选）
	IsDecision *bool  `json:"isDecision,omitempty"` // 是否为决策者（使用指针支持可选）
	Gender     string `json:"gender,omitempty"`     // 性别
	Address    string `json:"address,omitempty"`    // 联系地址
	Remark     string `json:"remark,omitempty"`     // 备注信息
	Status     *int   `json:"status,omitempty"`     // 状态（使用指针支持可选）
}

// Update 更新联系人信息
// 根据请求参数更新联系人的部分或全部信息
// 参数：
//   - ctx: 上下文
//   - id: 联系人 ID
//   - req: 更新联系人请求
// 返回值：
//   - *domain.Contact: 更新后的联系人对象
//   - error: 错误信息
func (s *contactService) Update(ctx context.Context, id string, req *UpdateContactReq) (*domain.Contact, error) {
	// 获取现有联系人信息
	contact, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 根据请求参数更新联系人信息（仅更新非空字段）
	if req.Name != "" {
		contact.Name = req.Name
	}
	if req.Position != "" {
		contact.Position = req.Position
	}
	if req.Department != "" {
		contact.Department = req.Department
	}
	if req.Phone != "" {
		contact.Phone = req.Phone
	}
	if req.Mobile != "" {
		contact.Mobile = req.Mobile
	}
	if req.Email != "" {
		contact.Email = req.Email
	}
	if req.Wechat != "" {
		contact.Wechat = req.Wechat
	}
	if req.QQ != "" {
		contact.QQ = req.QQ
	}
	if req.IsPrimary != nil {
		contact.IsPrimary = *req.IsPrimary
	}
	if req.IsDecision != nil {
		contact.IsDecision = *req.IsDecision
	}
	if req.Gender != "" {
		contact.Gender = req.Gender
	}
	if req.Address != "" {
		contact.Address = req.Address
	}
	if req.Remark != "" {
		contact.Remark = req.Remark
	}
	if req.Status != nil {
		contact.Status = *req.Status
	}

	// 调用仓储层更新联系人
	if err := s.repo.Update(ctx, contact); err != nil {
		return nil, err
	}

	return contact, nil
}

// Delete 删除联系人
// 根据联系人 ID 删除联系人记录
// 参数：
//   - ctx: 上下文
//   - id: 联系人 ID
// 返回值：
//   - error: 错误信息
func (s *contactService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
