package service

import (
	"context"
	"time"

	"github.com/qzt/backend/internal/domain"
)

// NotificationService 通知服务
type NotificationService interface {
	// 发送通知
	Send(ctx context.Context, req *SendNotificationReq) error
	SendToUsers(ctx context.Context, userIDs []string, req *SendNotificationReq) error
	Broadcast(ctx context.Context, req *SendNotificationReq) error

	// 查询通知
	List(ctx context.Context, userID string, params ListNotificationParams) (*ListResult, error)
	GetUnreadCount(ctx context.Context, userID string) (int64, error)

	// 操作通知
	MarkAsRead(ctx context.Context, id string) error
	MarkAllAsRead(ctx context.Context, userID string) error
	Delete(ctx context.Context, id string) error
}

type notificationService struct {
	repo NotificationRepository
}

func NewNotificationService(repo NotificationRepository) NotificationService {
	return &notificationService{repo: repo}
}

type SendNotificationReq struct {
	UserID      string `json:"userId"`
	Title       string `json:"title" binding:"required"`
	Content     string `json:"content"`
	Type        string `json:"type"`        // SYSTEM/BUSINESS/REMINDER
	Category    string `json:"category"`    // 客户/商机/合同/回款
	RelatedID   string `json:"relatedId"`
	RelatedType string `json:"relatedType"`
	Priority    int    `json:"priority"`
}

type ListNotificationParams struct {
	Page     int
	PageSize int
	Type     string
	Category string
	IsRead   *bool
}

func (s *notificationService) Send(ctx context.Context, req *SendNotificationReq) error {
	notification := &domain.Notification{
		UserID:      req.UserID,
		Title:       req.Title,
		Content:     req.Content,
		Type:        req.Type,
		Category:    req.Category,
		RelatedID:   req.RelatedID,
		RelatedType: req.RelatedType,
		Priority:    req.Priority,
		IsRead:      false,
	}

	return s.repo.Create(ctx, notification)
}

func (s *notificationService) SendToUsers(ctx context.Context, userIDs []string, req *SendNotificationReq) error {
	for _, userID := range userIDs {
		notification := &domain.Notification{
			UserID:      userID,
			Title:       req.Title,
			Content:     req.Content,
			Type:        req.Type,
			Category:    req.Category,
			RelatedID:   req.RelatedID,
			RelatedType: req.RelatedType,
			Priority:    req.Priority,
			IsRead:      false,
		}

		if err := s.repo.Create(ctx, notification); err != nil {
			return err
		}
	}
	return nil
}

func (s *notificationService) Broadcast(ctx context.Context, req *SendNotificationReq) error {
	// TODO: 实现广播通知
	return nil
}

func (s *notificationService) List(ctx context.Context, userID string, params ListNotificationParams) (*ListResult, error) {
	items, total, err := s.repo.List(ctx, userID, params)
	if err != nil {
		return nil, err
	}

	return &ListResult{
		Items:    items,
		Total:    total,
		Page:     params.Page,
		PageSize: params.PageSize,
	}, nil
}

func (s *notificationService) GetUnreadCount(ctx context.Context, userID string) (int64, error) {
	return s.repo.GetUnreadCount(ctx, userID)
}

func (s *notificationService) MarkAsRead(ctx context.Context, id string) error {
	now := time.Now()
	return s.repo.MarkAsRead(ctx, id, &now)
}

func (s *notificationService) MarkAllAsRead(ctx context.Context, userID string) error {
	now := time.Now()
	return s.repo.MarkAllAsRead(ctx, userID, &now)
}

func (s *notificationService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

// OperationLogService 操作日志服务
type OperationLogService interface {
	Log(ctx context.Context, log *domain.OperationLog) error
	List(ctx context.Context, params ListOperationLogParams) (*ListResult, error)
}

type operationLogService struct {
	repo OperationLogRepository
}

func NewOperationLogService(repo OperationLogRepository) OperationLogService {
	return &operationLogService{repo: repo}
}

type ListOperationLogParams struct {
	Page     int
	PageSize int
	UserID   string
	Method   string
	Path     string
	StartDate *time.Time
	EndDate   *time.Time
}

func (s *operationLogService) Log(ctx context.Context, log *domain.OperationLog) error {
	return s.repo.Create(ctx, log)
}

func (s *operationLogService) List(ctx context.Context, params ListOperationLogParams) (*ListResult, error) {
	items, total, err := s.repo.List(ctx, params)
	if err != nil {
		return nil, err
	}

	return &ListResult{
		Items:    items,
		Total:    total,
		Page:     params.Page,
		PageSize: params.PageSize,
	}, nil
}
