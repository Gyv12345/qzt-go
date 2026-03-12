package service

import (
	"context"

	"github.com/qzt/backend/internal/domain"
	"github.com/qzt/backend/internal/repository"
)

// RBACService RBAC 服务接口
type RBACService interface {
	// 角色管理
	CreateRole(ctx context.Context, role *domain.Role) error
	GetRoleByID(ctx context.Context, id string) (*domain.Role, error)
	ListRoles(ctx context.Context) ([]domain.Role, error)
	UpdateRole(ctx context.Context, role *domain.Role) error
	DeleteRole(ctx context.Context, id string) error
	AssignPermissionsToRole(ctx context.Context, roleID string, permissionIDs []string) error

	// 权限管理
	CreatePermission(ctx context.Context, permission *domain.Permission) error
	GetPermissionByCode(ctx context.Context, code string) (*domain.Permission, error)
	ListPermissions(ctx context.Context) ([]domain.Permission, error)

	// 用户角色管理
	AssignRolesToUser(ctx context.Context, userID string, roleIDs []string) error
	GetUserRoles(ctx context.Context, userID string) ([]domain.Role, error)
	GetUserPermissions(ctx context.Context, userID string) ([]domain.Permission, error)

	// 权限检查
	CheckPermission(ctx context.Context, userID, permissionCode string) (bool, error)
}

type rbacService struct {
	roleRepo       repository.RoleRepository
	permissionRepo repository.PermissionRepository
	userRepo       repository.UserRepository
}

func NewRBACService(roleRepo repository.RoleRepository, permissionRepo repository.PermissionRepository, userRepo repository.UserRepository) RBACService {
	return &rbacService{
		roleRepo:       roleRepo,
		permissionRepo: permissionRepo,
		userRepo:       userRepo,
	}
}

func (s *rbacService) CreateRole(ctx context.Context, role *domain.Role) error {
	return s.roleRepo.Create(ctx, role)
}

func (s *rbacService) GetRoleByID(ctx context.Context, id string) (*domain.Role, error) {
	return s.roleRepo.GetByID(ctx, id)
}

func (s *rbacService) ListRoles(ctx context.Context) ([]domain.Role, error) {
	return s.roleRepo.List(ctx)
}

func (s *rbacService) UpdateRole(ctx context.Context, role *domain.Role) error {
	return s.roleRepo.Update(ctx, role)
}

func (s *rbacService) DeleteRole(ctx context.Context, id string) error {
	return s.roleRepo.Delete(ctx, id)
}

func (s *rbacService) AssignPermissionsToRole(ctx context.Context, roleID string, permissionIDs []string) error {
	return s.roleRepo.AssignPermissions(ctx, roleID, permissionIDs)
}

func (s *rbacService) CreatePermission(ctx context.Context, permission *domain.Permission) error {
	return s.permissionRepo.Create(ctx, permission)
}

func (s *rbacService) GetPermissionByCode(ctx context.Context, code string) (*domain.Permission, error) {
	return s.permissionRepo.GetByCode(ctx, code)
}

func (s *rbacService) ListPermissions(ctx context.Context) ([]domain.Permission, error) {
	return s.permissionRepo.List(ctx)
}

func (s *rbacService) AssignRolesToUser(ctx context.Context, userID string, roleIDs []string) error {
	return s.userRepo.AssignRoles(ctx, userID, roleIDs)
}

func (s *rbacService) GetUserRoles(ctx context.Context, userID string) ([]domain.Role, error) {
	return s.userRepo.GetRoles(ctx, userID)
}

func (s *rbacService) GetUserPermissions(ctx context.Context, userID string) ([]domain.Permission, error) {
	return s.userRepo.GetPermissions(ctx, userID)
}

func (s *rbacService) CheckPermission(ctx context.Context, userID, permissionCode string) (bool, error) {
	permissions, err := s.userRepo.GetPermissions(ctx, userID)
	if err != nil {
		return false, err
	}

	for _, perm := range permissions {
		if perm.Code == permissionCode {
			return true, nil
		}
	}
	return false, nil
}

// FollowRecordService 跟进记录服务
type FollowRecordService interface {
	Create(ctx context.Context, req *CreateFollowRecordReq) (*domain.FollowRecord, error)
	GetByID(ctx context.Context, id string) (*domain.FollowRecord, error)
	List(ctx context.Context, customerID string, page, pageSize int) (*ListResult, error)
	Update(ctx context.Context, id string, req *UpdateFollowRecordReq) (*domain.FollowRecord, error)
	Delete(ctx context.Context, id string) error
}

type followRecordService struct {
	repo repository.FollowRecordRepository
}

func NewFollowRecordService(repo repository.FollowRecordRepository) FollowRecordService {
	return &followRecordService{repo: repo}
}

type CreateFollowRecordReq struct {
	CustomerID   string  `json:"customerId" binding:"required"`
	Type         string  `json:"type" binding:"required"`
	Content      string  `json:"content" binding:"required"`
	NextFollowAt *string `json:"nextFollowAt"`
	Status       string  `json:"status"`
}

type UpdateFollowRecordReq struct {
	Type         string  `json:"type"`
	Content      string  `json:"content"`
	NextFollowAt *string `json:"nextFollowAt"`
	Status       string  `json:"status"`
}

func (s *followRecordService) Create(ctx context.Context, req *CreateFollowRecordReq) (*domain.FollowRecord, error) {
	// 从上下文获取用户ID
	userID, _ := ctx.Value("userId").(string)

	record := &domain.FollowRecord{
		CustomerID: req.CustomerID,
		UserID:     userID,
		Type:       req.Type,
		Content:    req.Content,
		Status:     req.Status,
	}

	if err := s.repo.Create(ctx, record); err != nil {
		return nil, err
	}

	return record, nil
}

func (s *followRecordService) GetByID(ctx context.Context, id string) (*domain.FollowRecord, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *followRecordService) List(ctx context.Context, customerID string, page, pageSize int) (*ListResult, error) {
	items, total, err := s.repo.List(ctx, customerID, page, pageSize)
	if err != nil {
		return nil, err
	}

	return &ListResult{
		Items:    items,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func (s *followRecordService) Update(ctx context.Context, id string, req *UpdateFollowRecordReq) (*domain.FollowRecord, error) {
	record, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Type != "" {
		record.Type = req.Type
	}
	if req.Content != "" {
		record.Content = req.Content
	}
	if req.Status != "" {
		record.Status = req.Status
	}

	if err := s.repo.Update(ctx, record); err != nil {
		return nil, err
	}

	return record, nil
}

func (s *followRecordService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
