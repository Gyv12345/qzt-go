package repository

import (
	"context"

	"github.com/qzt/backend/internal/domain"
	"gorm.io/gorm"
)

// PermissionRepository 权限仓储接口
type PermissionRepository interface {
	Create(ctx context.Context, permission *domain.Permission) error
	GetByID(ctx context.Context, id string) (*domain.Permission, error)
	GetByCode(ctx context.Context, code string) (*domain.Permission, error)
	List(ctx context.Context) ([]domain.Permission, error)
	Update(ctx context.Context, permission *domain.Permission) error
	Delete(ctx context.Context, id string) error
	// 获取角色的所有权限
	GetByRoleID(ctx context.Context, roleID string) ([]domain.Permission, error)
	// 获取用户的所有权限（通过角色）
	GetByUserID(ctx context.Context, userID string) ([]domain.Permission, error)
}

type permissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) PermissionRepository {
	return &permissionRepository{db: db}
}

func (r *permissionRepository) Create(ctx context.Context, permission *domain.Permission) error {
	return r.db.WithContext(ctx).Create(permission).Error
}

func (r *permissionRepository) GetByID(ctx context.Context, id string) (*domain.Permission, error) {
	var permission domain.Permission
	err := r.db.WithContext(ctx).First(&permission, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *permissionRepository) GetByCode(ctx context.Context, code string) (*domain.Permission, error) {
	var permission domain.Permission
	err := r.db.WithContext(ctx).First(&permission, "code = ?", code).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *permissionRepository) List(ctx context.Context) ([]domain.Permission, error) {
	var permissions []domain.Permission
	err := r.db.WithContext(ctx).Find(&permissions).Error
	return permissions, err
}

func (r *permissionRepository) Update(ctx context.Context, permission *domain.Permission) error {
	return r.db.WithContext(ctx).Save(permission).Error
}

func (r *permissionRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&domain.Permission{}, "id = ?", id).Error
}

func (r *permissionRepository) GetByRoleID(ctx context.Context, roleID string) ([]domain.Permission, error) {
	var permissions []domain.Permission
	err := r.db.WithContext(ctx).
		Table("permissions").
		Joins("JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Where("role_permissions.role_id = ?", roleID).
		Find(&permissions).Error
	return permissions, err
}

func (r *permissionRepository) GetByUserID(ctx context.Context, userID string) ([]domain.Permission, error) {
	var permissions []domain.Permission
	err := r.db.WithContext(ctx).
		Table("permissions").
		Joins("JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Joins("JOIN user_roles ON role_permissions.role_id = user_roles.role_id").
		Where("user_roles.user_id = ?", userID).
		Distinct().
		Find(&permissions).Error
	return permissions, err
}

// RoleRepository 角色仓储接口
type RoleRepository interface {
	Create(ctx context.Context, role *domain.Role) error
	GetByID(ctx context.Context, id string) (*domain.Role, error)
	GetByCode(ctx context.Context, code string) (*domain.Role, error)
	List(ctx context.Context) ([]domain.Role, error)
	Update(ctx context.Context, role *domain.Role) error
	Delete(ctx context.Context, id string) error
	// 角色权限管理
	AssignPermissions(ctx context.Context, roleID string, permissionIDs []string) error
	RemovePermissions(ctx context.Context, roleID string, permissionIDs []string) error
	// 获取角色的用户列表
	GetUsers(ctx context.Context, roleID string) ([]domain.User, error)
}

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) Create(ctx context.Context, role *domain.Role) error {
	return r.db.WithContext(ctx).Create(role).Error
}

func (r *roleRepository) GetByID(ctx context.Context, id string) (*domain.Role, error) {
	var role domain.Role
	err := r.db.WithContext(ctx).Preload("Menus").First(&role, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) GetByCode(ctx context.Context, code string) (*domain.Role, error) {
	var role domain.Role
	err := r.db.WithContext(ctx).First(&role, "code = ?", code).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) List(ctx context.Context) ([]domain.Role, error) {
	var roles []domain.Role
	err := r.db.WithContext(ctx).Preload("Menus").Find(&roles).Error
	return roles, err
}

func (r *roleRepository) Update(ctx context.Context, role *domain.Role) error {
	return r.db.WithContext(ctx).Save(role).Error
}

func (r *roleRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&domain.Role{}, "id = ?", id).Error
}

func (r *roleRepository) AssignPermissions(ctx context.Context, roleID string, permissionIDs []string) error {
	// 使用事务批量插入
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, permissionID := range permissionIDs {
			rp := domain.RolePermission{
				RoleID:       roleID,
				PermissionID: permissionID,
			}
			if err := tx.Create(&rp).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *roleRepository) RemovePermissions(ctx context.Context, roleID string, permissionIDs []string) error {
	return r.db.WithContext(ctx).
		Where("role_id = ? AND permission_id IN ?", roleID, permissionIDs).
		Delete(&domain.RolePermission{}).Error
}

func (r *roleRepository) GetUsers(ctx context.Context, roleID string) ([]domain.User, error) {
	var users []domain.User
	err := r.db.WithContext(ctx).
		Table("users").
		Joins("JOIN user_roles ON users.id = user_roles.user_id").
		Where("user_roles.role_id = ?", roleID).
		Find(&users).Error
	return users, err
}

// UserRepository 扩展方法
type UserRepository interface {
	// 原有方法
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id string) (*domain.User, error)
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id string) error
	// 新增 RBAC 方法
	AssignRoles(ctx context.Context, userID string, roleIDs []string) error
	RemoveRoles(ctx context.Context, userID string, roleIDs []string) error
	GetRoles(ctx context.Context, userID string) ([]domain.Role, error)
	GetPermissions(ctx context.Context, userID string) ([]domain.Permission, error)
}

// FollowRecordRepository 跟进记录仓储
type FollowRecordRepository interface {
	Create(ctx context.Context, record *domain.FollowRecord) error
	GetByID(ctx context.Context, id string) (*domain.FollowRecord, error)
	List(ctx context.Context, customerID string, page, pageSize int) ([]domain.FollowRecord, int64, error)
	Update(ctx context.Context, record *domain.FollowRecord) error
	Delete(ctx context.Context, id string) error
}

type followRecordRepository struct {
	db *gorm.DB
}

func NewFollowRecordRepository(db *gorm.DB) FollowRecordRepository {
	return &followRecordRepository{db: db}
}

func (r *followRecordRepository) Create(ctx context.Context, record *domain.FollowRecord) error {
	return r.db.WithContext(ctx).Create(record).Error
}

func (r *followRecordRepository) GetByID(ctx context.Context, id string) (*domain.FollowRecord, error) {
	var record domain.FollowRecord
	err := r.db.WithContext(ctx).
		Preload("Customer").
		Preload("User").
		First(&record, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (r *followRecordRepository) List(ctx context.Context, customerID string, page, pageSize int) ([]domain.FollowRecord, int64, error) {
	var records []domain.FollowRecord
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.FollowRecord{})
	if customerID != "" {
		query = query.Where("customer_id = ?", customerID)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.
		Preload("Customer").
		Preload("User").
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&records).Error

	return records, total, err
}

func (r *followRecordRepository) Update(ctx context.Context, record *domain.FollowRecord) error {
	return r.db.WithContext(ctx).Save(record).Error
}

func (r *followRecordRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&domain.FollowRecord{}, "id = ?", id).Error
}

// 其他仓储接口类似实现...
