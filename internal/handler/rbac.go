package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/qzt/backend/internal/domain"
	"github.com/qzt/backend/internal/service"
	"github.com/qzt/backend/pkg/response"
)

// RBACHandler RBAC 处理器
type RBACHandler struct {
	rbacSvc service.RBACService
}

func NewRBACHandler(rbacSvc service.RBACService) *RBACHandler {
	return &RBACHandler{rbacSvc: rbacSvc}
}

// CreateRole 创建角色
func (h *RBACHandler) CreateRole(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Code        string `json:"code" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParams, err.Error())
		return
	}

	role := &domain.Role{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Status:      1,
	}

	if err := h.rbacSvc.CreateRole(c.Request.Context(), role); err != nil {
		response.Error(c, response.CodeInternalError)
		return
	}

	response.Created(c, role)
}

// ListRoles 获取角色列表
func (h *RBACHandler) ListRoles(c *gin.Context) {
	roles, err := h.rbacSvc.ListRoles(c.Request.Context())
	if err != nil {
		response.Error(c, response.CodeInternalError)
		return
	}

	response.Success(c, roles)
}

// AssignPermissions 分配权限给角色
func (h *RBACHandler) AssignPermissions(c *gin.Context) {
	roleID := c.Param("id")

	var req struct {
		PermissionIDs []string `json:"permissionIds" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParams, err.Error())
		return
	}

	if err := h.rbacSvc.AssignPermissionsToRole(c.Request.Context(), roleID, req.PermissionIDs); err != nil {
		response.Error(c, response.CodeInternalError)
		return
	}

	response.Success(c, nil)
}

// CreatePermission 创建权限
func (h *RBACHandler) CreatePermission(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Code        string `json:"code" binding:"required"`
		Type        string `json:"type"`
		Resource    string `json:"resource"`
		Method      string `json:"method"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParams, err.Error())
		return
	}

	permission := &domain.Permission{
		Name:        req.Name,
		Code:        req.Code,
		Type:        req.Type,
		Resource:    req.Resource,
		Method:      req.Method,
		Description: req.Description,
		Status:      1,
	}

	if err := h.rbacSvc.CreatePermission(c.Request.Context(), permission); err != nil {
		response.Error(c, response.CodeInternalError)
		return
	}

	response.Created(c, permission)
}

// ListPermissions 获取权限列表
func (h *RBACHandler) ListPermissions(c *gin.Context) {
	permissions, err := h.rbacSvc.ListPermissions(c.Request.Context())
	if err != nil {
		response.Error(c, response.CodeInternalError)
		return
	}

	response.Success(c, permissions)
}

// AssignRolesToUser 给用户分配角色
func (h *RBACHandler) AssignRolesToUser(c *gin.Context) {
	userID := c.Param("id")

	var req struct {
		RoleIDs []string `json:"roleIds" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParams, err.Error())
		return
	}

	if err := h.rbacSvc.AssignRolesToUser(c.Request.Context(), userID, req.RoleIDs); err != nil {
		response.Error(c, response.CodeInternalError)
		return
	}

	response.Success(c, nil)
}

// GetUserPermissions 获取用户权限
func (h *RBACHandler) GetUserPermissions(c *gin.Context) {
	userID := c.Param("id")

	permissions, err := h.rbacSvc.GetUserPermissions(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, response.CodeInternalError)
		return
	}

	response.Success(c, permissions)
}

// FollowRecordHandler 跟进记录处理器
type FollowRecordHandler struct {
	svc service.FollowRecordService
}

func NewFollowRecordHandler(svc service.FollowRecordService) *FollowRecordHandler {
	return &FollowRecordHandler{svc: svc}
}

// CreateFollowRecord 创建跟进记录
func (h *FollowRecordHandler) Create(c *gin.Context) {
	var req service.CreateFollowRecordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParams, err.Error())
		return
	}

	record, err := h.svc.Create(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, response.CodeInternalError)
		return
	}

	response.Created(c, record)
}

// ListFollowRecords 获取跟进记录列表
func (h *FollowRecordHandler) List(c *gin.Context) {
	customerID := c.Query("customerId")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	result, err := h.svc.List(c.Request.Context(), customerID, page, pageSize)
	if err != nil {
		response.Error(c, response.CodeInternalError)
		return
	}

	response.SuccessPage(c, result.Items, result.Total, result.Page, result.PageSize)
}

// GetFollowRecord 获取跟进记录详情
func (h *FollowRecordHandler) Get(c *gin.Context) {
	id := c.Param("id")

	record, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, response.CodeNotFound)
		return
	}

	response.Success(c, record)
}

// UpdateFollowRecord 更新跟进记录
func (h *FollowRecordHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req service.UpdateFollowRecordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParams, err.Error())
		return
	}

	record, err := h.svc.Update(c.Request.Context(), id, &req)
	if err != nil {
		response.Error(c, response.CodeInternalError)
		return
	}

	response.Success(c, record)
}

// DeleteFollowRecord 删除跟进记录
func (h *FollowRecordHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		response.Error(c, response.CodeInternalError)
		return
	}

	response.NoContent(c)
}
