// Package repository 定义了数据访问层（仓储层）
// 负责与数据库交互，提供 CRUD 操作接口
package repository

import (
	"context"

	"github.com/qzt/backend/internal/domain"
	"gorm.io/gorm"
)

// APIKeyRepository API Key 仓储接口
// 定义了 API Key 数据访问层的标准方法
type APIKeyRepository interface {
	// Create 创建 API Key 记录
	Create(ctx context.Context, apiKey *domain.APIKey) error
	// GetByHash 根据密钥哈希获取 API Key 信息
	GetByHash(ctx context.Context, keyHash string) (*domain.APIKey, error)
	// GetByID 根据 ID 获取 API Key 信息
	GetByID(ctx context.Context, id string) (*domain.APIKey, error)
	// List 获取指定用户的所有 API Key 列表
	List(ctx context.Context, userID string) ([]domain.APIKey, error)
	// Update 更新 API Key 信息
	Update(ctx context.Context, apiKey *domain.APIKey) error
	// Delete 删除 API Key 记录
	Delete(ctx context.Context, id string) error
	// UpdateLastUsed 更新 API Key 的最后使用时间
	UpdateLastUsed(id string) error
}

// apiKeyRepository API Key 仓储实现
type apiKeyRepository struct {
	db *gorm.DB // GORM 数据库实例
}

// NewAPIKeyRepository 创建 API Key 仓储实例
// 参数：
//   - db: GORM 数据库实例
// 返回值：
//   - APIKeyRepository: API Key 仓储接口
func NewAPIKeyRepository(db *gorm.DB) APIKeyRepository {
	return &apiKeyRepository{db: db}
}

// Create 创建新的 API Key 记录
// 参数：
//   - ctx: 上下文
//   - apiKey: API Key 对象
// 返回值：
//   - error: 错误信息
func (r *apiKeyRepository) Create(ctx context.Context, apiKey *domain.APIKey) error {
	return r.db.WithContext(ctx).Create(apiKey).Error
}

// GetByHash 根据密钥哈希获取 API Key 信息
// 用于 API 认证时验证 Key 的有效性
// 参数：
//   - ctx: 上下文
//   - keyHash: 密钥的 SHA256 哈希值
// 返回值：
//   - *domain.APIKey: API Key 对象
//   - error: 错误信息
func (r *apiKeyRepository) GetByHash(ctx context.Context, keyHash string) (*domain.APIKey, error) {
	var apiKey domain.APIKey
	// 查询 API Key 并预加载关联的用户信息
	err := r.db.WithContext(ctx).
		Preload("User").
		First(&apiKey, "key_hash = ?", keyHash).Error
	if err != nil {
		return nil, err
	}
	return &apiKey, nil
}

// GetByID 根据 ID 获取 API Key 信息
// 参数：
//   - ctx: 上下文
//   - id: API Key ID
// 返回值：
//   - *domain.APIKey: API Key 对象
//   - error: 错误信息
func (r *apiKeyRepository) GetByID(ctx context.Context, id string) (*domain.APIKey, error) {
	var apiKey domain.APIKey
	// 查询 API Key 并预加载关联的用户信息
	err := r.db.WithContext(ctx).
		Preload("User").
		First(&apiKey, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &apiKey, nil
}

// List 获取指定用户的所有 API Key 列表
// 参数：
//   - ctx: 上下文
//   - userID: 用户 ID
// 返回值：
//   - []domain.APIKey: API Key 列表
//   - error: 错误信息
func (r *apiKeyRepository) List(ctx context.Context, userID string) ([]domain.APIKey, error) {
	var apiKeys []domain.APIKey
	// 查询指定用户的所有 API Key，按创建时间倒序排列
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&apiKeys).Error
	return apiKeys, err
}

// Update 更新 API Key 信息
// 参数：
//   - ctx: 上下文
//   - apiKey: API Key 对象
// 返回值：
//   - error: 错误信息
func (r *apiKeyRepository) Update(ctx context.Context, apiKey *domain.APIKey) error {
	return r.db.WithContext(ctx).Save(apiKey).Error
}

// Delete 删除 API Key 记录
// 参数：
//   - ctx: 上下文
//   - id: API Key ID
// 返回值：
//   - error: 错误信息
func (r *apiKeyRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&domain.APIKey{}, "id = ?", id).Error
}

// UpdateLastUsed 更新 API Key 的最后使用时间
// 在每次成功使用 API Key 后调用，记录最后使用时间
// 参数：
//   - id: API Key ID
// 返回值：
//   - error: 错误信息
func (r *apiKeyRepository) UpdateLastUsed(id string) error {
	// 使用数据库的 NOW() 函数更新时间戳
	return r.db.Model(&domain.APIKey{}).
		Where("id = ?", id).
		Update("last_used_at", gorm.Expr("NOW()")).Error
}
