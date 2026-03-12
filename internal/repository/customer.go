// Package repository 定义了数据访问层（仓储层）
// 负责与数据库交互，提供 CRUD 操作接口
package repository

import (
	"context"

	"github.com/qzt/backend/internal/domain"
	"gorm.io/gorm"
)

// CustomerRepository 客户仓储接口
// 定义了客户数据访问层的标准方法
type CustomerRepository interface {
	// Create 创建客户记录
	Create(ctx context.Context, customer *domain.Customer) error
	// GetByID 根据 ID 获取客户详情
	GetByID(ctx context.Context, id string) (*domain.Customer, error)
	// List 获取客户列表，支持分页和过滤条件
	List(ctx context.Context, params ListParams) ([]domain.Customer, int64, error)
	// Update 更新客户信息
	Update(ctx context.Context, customer *domain.Customer) error
	// Delete 删除客户记录
	Delete(ctx context.Context, id string) error
}

// ListParams 客户列表查询参数
type ListParams struct {
	Page     int    // 页码
	PageSize int    // 每页数量
	Keyword  string // 搜索关键词
	Level    string // 客户等级
	Status   *int   // 状态（指针类型，支持可选）
}

// customerRepository 客户仓储实现
type customerRepository struct {
	db *gorm.DB // GORM 数据库实例
}

// NewCustomerRepository 创建客户仓储实例
// 参数：
//   - db: GORM 数据库实例
// 返回值：
//   - CustomerRepository: 客户仓储接口
func NewCustomerRepository(db *gorm.DB) CustomerRepository {
	return &customerRepository{db: db}
}

// Create 创建新的客户记录
// 参数：
//   - ctx: 上下文
//   - customer: 客户对象
// 返回值：
//   - error: 错误信息
func (r *customerRepository) Create(ctx context.Context, customer *domain.Customer) error {
	return r.db.WithContext(ctx).Create(customer).Error
}

// GetByID 根据 ID 获取客户详情
// 参数：
//   - ctx: 上下文
//   - id: 客户 ID
// 返回值：
//   - *domain.Customer: 客户对象
//   - error: 错误信息
func (r *customerRepository) GetByID(ctx context.Context, id string) (*domain.Customer, error) {
	var customer domain.Customer
	// 查询客户并预加载关联的跟进人和联系人
	err := r.db.WithContext(ctx).
		Preload("FollowUser").
		Preload("Contacts").
		First(&customer, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

// List 获取客户列表
// 参数：
//   - ctx: 上下文
//   - params: 查询参数（分页、关键词、等级、状态）
// 返回值：
//   - []domain.Customer: 客户列表
//   - int64: 总记录数
//   - error: 错误信息
func (r *customerRepository) List(ctx context.Context, params ListParams) ([]domain.Customer, int64, error) {
	var customers []domain.Customer
	var total int64

	// 构建查询
	query := r.db.WithContext(ctx).Model(&domain.Customer{})

	// 添加关键词过滤条件（搜索客户名称或简称）
	if params.Keyword != "" {
		query = query.Where("name LIKE ? OR short_name LIKE ?",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}

	// 添加客户等级过滤条件
	if params.Level != "" {
		query = query.Where("customer_level = ?", params.Level)
	}

	// 添加状态过滤条件
	if params.Status != nil {
		query = query.Where("status = ?", *params.Status)
	}

	// 统计总记录数
	query.Count(&total)

	// 计算分页偏移量
	offset := (params.Page - 1) * params.PageSize
	// 执行查询，预加载跟进人并按创建时间倒序排列
	err := query.Preload("FollowUser").
		Offset(offset).Limit(params.PageSize).
		Order("created_at DESC").
		Find(&customers).Error

	return customers, total, err
}

// Update 更新客户信息
// 参数：
//   - ctx: 上下文
//   - customer: 客户对象
// 返回值：
//   - error: 错误信息
func (r *customerRepository) Update(ctx context.Context, customer *domain.Customer) error {
	return r.db.WithContext(ctx).Save(customer).Error
}

// Delete 删除客户记录
// 参数：
//   - ctx: 上下文
//   - id: 客户 ID
// 返回值：
//   - error: 错误信息
func (r *customerRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&domain.Customer{}, "id = ?", id).Error
}
