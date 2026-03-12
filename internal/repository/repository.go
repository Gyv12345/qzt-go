package repository

import (
	"context"

	"github.com/qzt/backend/internal/domain"
	"gorm.io/gorm"
)

// Repository 聚合所有仓储
type Repository struct {
	User         UserRepository
	Role         RoleRepository
	Permission   PermissionRepository
	Customer     CustomerRepository
	Contact      ContactRepository
	Contract     ContractRepository
	APIKey       APIKeyRepository
	FollowRecord FollowRecordRepository
	Opportunity  OpportunityRepository
	Payment      PaymentRepository
}

// NewRepository 创建仓储聚合实例
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		User:         NewUserRepository(db),
		Role:         NewRoleRepository(db),
		Permission:   NewPermissionRepository(db),
		Customer:     NewCustomerRepository(db),
		Contact:      NewContactRepository(db),
		Contract:     NewContractRepository(db),
		APIKey:       NewAPIKeyRepository(db),
		FollowRecord: NewFollowRecordRepository(db),
		Opportunity:  NewOpportunityRepository(db),
		Payment:      NewPaymentRepository(db),
	}
}

// UserRepository 用户仓储
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id string) (*domain.User, error)
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id string) error
	AssignRoles(ctx context.Context, userID string, roleIDs []string) error
	RemoveRoles(ctx context.Context, userID string, roleIDs []string) error
	GetRoles(ctx context.Context, userID string) ([]domain.Role, error)
	GetPermissions(ctx context.Context, userID string) ([]domain.Permission, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).
		Preload("Roles").
		Preload("Department").
		First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).
		Preload("Roles").
		First(&user, "username = ?", username).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&domain.User{}, "id = ?", id).Error
}

func (r *userRepository) AssignRoles(ctx context.Context, userID string, roleIDs []string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, roleID := range roleIDs {
			ur := domain.UserRole{
				UserID: userID,
				RoleID: roleID,
			}
			if err := tx.Create(&ur).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *userRepository) RemoveRoles(ctx context.Context, userID string, roleIDs []string) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND role_id IN ?", userID, roleIDs).
		Delete(&domain.UserRole{}).Error
}

func (r *userRepository) GetRoles(ctx context.Context, userID string) ([]domain.Role, error) {
	var roles []domain.Role
	err := r.db.WithContext(ctx).
		Table("roles").
		Joins("JOIN user_roles ON roles.id = user_roles.role_id").
		Where("user_roles.user_id = ?", userID).
		Find(&roles).Error
	return roles, err
}

func (r *userRepository) GetPermissions(ctx context.Context, userID string) ([]domain.Permission, error) {
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
