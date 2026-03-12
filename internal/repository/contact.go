// Package repository 定义了数据访问层（仓储层）
// 负责与数据库交互，提供 CRUD 操作接口
package repository

import (
	"context"

	"github.com/qzt/backend/internal/domain"
	"gorm.io/gorm"
)

// ContactRepository 联系人仓储接口
// 定义了联系人数据访问层的标准方法
type ContactRepository interface {
	// Create 创建联系人记录
	Create(ctx context.Context, contact *domain.Contact) error
	// GetByID 根据 ID 获取联系人详情
	GetByID(ctx context.Context, id string) (*domain.Contact, error)
	// List 获取联系人列表，支持分页和过滤条件
	List(ctx context.Context, params ContactListParams) ([]domain.Contact, int64, error)
	// Update 更新联系人信息
	Update(ctx context.Context, contact *domain.Contact) error
	// Delete 删除联系人记录
	Delete(ctx context.Context, id string) error
}

// ContactListParams 联系人列表查询参数
type ContactListParams struct {
	Page       int    // 页码
	PageSize   int    // 每页数量
	CustomerID string // 客户 ID（过滤指定客户的联系人）
	Keyword    string // 搜索关键词
}

// contactRepository 联系人仓储实现
type contactRepository struct {
	db *gorm.DB // GORM 数据库实例
}

// NewContactRepository 创建联系人仓储实例
// 参数：
//   - db: GORM 数据库实例
// 返回值：
//   - ContactRepository: 联系人仓储接口
func NewContactRepository(db *gorm.DB) ContactRepository {
	return &contactRepository{db: db}
}

// Create 创建新的联系人记录
// 参数：
//   - ctx: 上下文
//   - contact: 联系人对象
// 返回值：
//   - error: 错误信息
func (r *contactRepository) Create(ctx context.Context, contact *domain.Contact) error {
	return r.db.WithContext(ctx).Create(contact).Error
}

// GetByID 根据 ID 获取联系人详情
// 参数：
//   - ctx: 上下文
//   - id: 联系人 ID
// 返回值：
//   - *domain.Contact: 联系人对象
//   - error: 错误信息
func (r *contactRepository) GetByID(ctx context.Context, id string) (*domain.Contact, error) {
	var contact domain.Contact
	// 查询联系人并预加载关联的客户信息
	err := r.db.WithContext(ctx).
		Preload("Customer").
		First(&contact, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &contact, nil
}

// List 获取联系人列表
// 参数：
//   - ctx: 上下文
//   - params: 查询参数（分页、客户ID、关键词）
// 返回值：
//   - []domain.Contact: 联系人列表
//   - int64: 总记录数
//   - error: 错误信息
func (r *contactRepository) List(ctx context.Context, params ContactListParams) ([]domain.Contact, int64, error) {
	var contacts []domain.Contact
	var total int64

	// 构建查询
	query := r.db.WithContext(ctx).Model(&domain.Contact{})

	// 添加客户 ID 过滤条件（查询指定客户的联系人）
	if params.CustomerID != "" {
		query = query.Where("customer_id = ?", params.CustomerID)
	}

	// 添加关键词过滤条件（搜索姓名、电话或手机号）
	if params.Keyword != "" {
		query = query.Where("name LIKE ? OR phone LIKE ? OR mobile LIKE ?",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}

	// 统计总记录数
	query.Count(&total)

	// 计算分页偏移量
	offset := (params.Page - 1) * params.PageSize
	// 执行查询，预加载客户信息并按创建时间倒序排列
	err := query.Preload("Customer").
		Offset(offset).Limit(params.PageSize).
		Order("created_at DESC").
		Find(&contacts).Error

	return contacts, total, err
}

// Update 更新联系人信息
// 参数：
//   - ctx: 上下文
//   - contact: 联系人对象
// 返回值：
//   - error: 错误信息
func (r *contactRepository) Update(ctx context.Context, contact *domain.Contact) error {
	return r.db.WithContext(ctx).Save(contact).Error
}

// Delete 删除联系人记录
// 参数：
//   - ctx: 上下文
//   - id: 联系人 ID
// 返回值：
//   - error: 错误信息
func (r *contactRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&domain.Contact{}, "id = ?", id).Error
}
