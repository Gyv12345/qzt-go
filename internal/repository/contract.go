// Package repository 定义了数据访问层（仓储层）
// 负责与数据库交互，提供 CRUD 操作接口
package repository

import (
	"context"

	"github.com/qzt/backend/internal/domain"
	"gorm.io/gorm"
)

// ContractRepository 合同仓储接口
// 定义了合同数据访问层的标准方法
type ContractRepository interface {
	// Create 创建合同记录
	Create(ctx context.Context, contract *domain.Contract) error
	// GetByID 根据 ID 获取合同详情
	GetByID(ctx context.Context, id string) (*domain.Contract, error)
	// List 获取合同列表，支持分页和过滤条件
	List(ctx context.Context, params ContractListParams) ([]domain.Contract, int64, error)
	// Update 更新合同信息
	Update(ctx context.Context, contract *domain.Contract) error
	// Delete 删除合同记录
	Delete(ctx context.Context, id string) error
}

// ContractListParams 合同列表查询参数
type ContractListParams struct {
	Page       int    // 页码
	PageSize   int    // 每页数量
	CustomerID string // 客户 ID（过滤指定客户的合同）
	Status     string // 合同状态
	Keyword    string // 搜索关键词
}

// contractRepository 合同仓储实现
type contractRepository struct {
	db *gorm.DB // GORM 数据库实例
}

// NewContractRepository 创建合同仓储实例
// 参数：
//   - db: GORM 数据库实例
// 返回值：
//   - ContractRepository: 合同仓储接口
func NewContractRepository(db *gorm.DB) ContractRepository {
	return &contractRepository{db: db}
}

// Create 创建新的合同记录
// 参数：
//   - ctx: 上下文
//   - contract: 合同对象
// 返回值：
//   - error: 错误信息
func (r *contractRepository) Create(ctx context.Context, contract *domain.Contract) error {
	return r.db.WithContext(ctx).Create(contract).Error
}

// GetByID 根据 ID 获取合同详情
// 参数：
//   - ctx: 上下文
//   - id: 合同 ID
// 返回值：
//   - *domain.Contract: 合同对象
//   - error: 错误信息
func (r *contractRepository) GetByID(ctx context.Context, id string) (*domain.Contract, error) {
	var contract domain.Contract
	// 查询合同并预加载关联的客户和负责人信息
	err := r.db.WithContext(ctx).
		Preload("Customer").
		Preload("Owner").
		First(&contract, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &contract, nil
}

// List 获取合同列表
// 参数：
//   - ctx: 上下文
//   - params: 查询参数（分页、客户ID、状态、关键词）
// 返回值：
//   - []domain.Contract: 合同列表
//   - int64: 总记录数
//   - error: 错误信息
func (r *contractRepository) List(ctx context.Context, params ContractListParams) ([]domain.Contract, int64, error) {
	var contracts []domain.Contract
	var total int64

	// 构建查询
	query := r.db.WithContext(ctx).Model(&domain.Contract{})

	// 添加客户 ID 过滤条件（查询指定客户的合同）
	if params.CustomerID != "" {
		query = query.Where("customer_id = ?", params.CustomerID)
	}

	// 添加合同状态过滤条件
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	// 添加关键词过滤条件（搜索合同标题或合同编号）
	if params.Keyword != "" {
		query = query.Where("title LIKE ? OR contract_no LIKE ?",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}

	// 统计总记录数
	query.Count(&total)

	// 计算分页偏移量
	offset := (params.Page - 1) * params.PageSize
	// 执行查询，预加载客户和负责人信息，并按创建时间倒序排列
	err := query.Preload("Customer").Preload("Owner").
		Offset(offset).Limit(params.PageSize).
		Order("created_at DESC").
		Find(&contracts).Error

	return contracts, total, err
}

// Update 更新合同信息
// 参数：
//   - ctx: 上下文
//   - contract: 合同对象
// 返回值：
//   - error: 错误信息
func (r *contractRepository) Update(ctx context.Context, contract *domain.Contract) error {
	return r.db.WithContext(ctx).Save(contract).Error
}

// Delete 删除合同记录
// 参数：
//   - ctx: 上下文
//   - id: 合同 ID
// 返回值：
//   - error: 错误信息
func (r *contractRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&domain.Contract{}, "id = ?", id).Error
}
