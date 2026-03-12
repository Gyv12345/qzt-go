// Package repository 定义了数据访问层（仓储层）
// 负责与数据库交互，提供 CRUD 操作接口
package repository

import (
	"context"

	"github.com/qzt/backend/internal/domain"
	"gorm.io/gorm"
)

// UserRepository 用户仓储接口
// 定义了用户数据访问层的标准方法
type UserRepository interface {
	// Create 创建用户记录
	Create(ctx context.Context, user *domain.User) error
	// GetByID 根据 ID 获取用户详情
	GetByID(ctx context.Context, id string) (*domain.User, error)
	// GetByUsername 根据用户名获取用户信息
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	// Update 更新用户信息
	Update(ctx context.Context, user *domain.User) error
}

// userRepository 用户仓储实现
type userRepository struct {
	db *gorm.DB // GORM 数据库实例
}

// NewUserRepository 创建用户仓储实例
// 参数：
//   - db: GORM 数据库实例
// 返回值：
//   - UserRepository: 用户仓储接口
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// Create 创建新的用户记录
// 参数：
//   - ctx: 上下文
//   - user: 用户对象
// 返回值：
//   - error: 错误信息
func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// GetByID 根据 ID 获取用户详情
// 参数：
//   - ctx: 上下文
//   - id: 用户 ID
// 返回值：
//   - *domain.User: 用户对象
//   - error: 错误信息
func (r *userRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	var user domain.User
	// 查询用户并预加载关联的角色和部门信息
	err := r.db.WithContext(ctx).
		Preload("Roles").
		Preload("Department").
		First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByUsername 根据用户名获取用户信息
// 用于登录认证
// 参数：
//   - ctx: 上下文
//   - username: 用户名
// 返回值：
//   - *domain.User: 用户对象
//   - error: 错误信息
func (r *userRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	var user domain.User
	// 查询用户并预加载关联的角色和部门信息
	err := r.db.WithContext(ctx).
		Preload("Roles").
		Preload("Department").
		First(&user, "username = ?", username).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update 更新用户信息
// 参数：
//   - ctx: 上下文
//   - user: 用户对象
// 返回值：
//   - error: 错误信息
func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// Repository 仓储聚合
// 聚合了所有业务实体的仓储接口，便于依赖注入和管理
type Repository struct {
	Customer CustomerRepository // 客户仓储
	Contact  ContactRepository  // 联系人仓储
	Contract ContractRepository // 合同仓储
	APIKey   APIKeyRepository   // API Key 仓储
	User     UserRepository     // 用户仓储
}

// NewRepository 创建仓储聚合实例
// 参数：
//   - db: GORM 数据库实例
// 返回值：
//   - *Repository: 仓储聚合对象
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		Customer: NewCustomerRepository(db),
		Contact:  NewContactRepository(db),
		Contract: NewContractRepository(db),
		APIKey:   NewAPIKeyRepository(db),
		User:     NewUserRepository(db),
	}
}
