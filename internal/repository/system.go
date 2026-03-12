package repository

import (
	"context"
	"time"

	"github.com/qzt/backend/internal/domain"
	"gorm.io/gorm"
)

// NotificationRepository 通知仓储
type NotificationRepository interface {
	Create(ctx context.Context, notification *domain.Notification) error
	List(ctx context.Context, userID string, params ListNotificationParams) ([]domain.Notification, int64, error)
	GetUnreadCount(ctx context.Context, userID string) (int64, error)
	MarkAsRead(ctx context.Context, id string, readAt *time.Time) error
	MarkAllAsRead(ctx context.Context, userID string, readAt *time.Time) error
	Delete(ctx context.Context, id string) error
}

type notificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) Create(ctx context.Context, notification *domain.Notification) error {
	return r.db.WithContext(ctx).Create(notification).Error
}

func (r *notificationRepository) List(ctx context.Context, userID string, params ListNotificationParams) ([]domain.Notification, int64, error) {
	var notifications []domain.Notification
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.Notification{}).Where("user_id = ?", userID)

	if params.Type != "" {
		query = query.Where("type = ?", params.Type)
	}
	if params.Category != "" {
		query = query.Where("category = ?", params.Category)
	}
	if params.IsRead != nil {
		query = query.Where("is_read = ?", *params.IsRead)
	}

	query.Count(&total)

	offset := (params.Page - 1) * params.PageSize
	err := query.
		Offset(offset).
		Limit(params.PageSize).
		Order("created_at DESC").
		Find(&notifications).Error

	return notifications, total, err
}

func (r *notificationRepository) GetUnreadCount(ctx context.Context, userID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&domain.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Count(&count).Error
	return count, err
}

func (r *notificationRepository) MarkAsRead(ctx context.Context, id string, readAt *time.Time) error {
	return r.db.WithContext(ctx).
		Model(&domain.Notification{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": readAt,
		}).Error
}

func (r *notificationRepository) MarkAllAsRead(ctx context.Context, userID string, readAt *time.Time) error {
	return r.db.WithContext(ctx).
		Model(&domain.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": readAt,
		}).Error
}

func (r *notificationRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&domain.Notification{}, "id = ?", id).Error
}

// OperationLogRepository 操作日志仓储
type OperationLogRepository interface {
	Create(ctx context.Context, log *domain.OperationLog) error
	List(ctx context.Context, params ListOperationLogParams) ([]domain.OperationLog, int64, error)
}

type operationLogRepository struct {
	db *gorm.DB
}

func NewOperationLogRepository(db *gorm.DB) OperationLogRepository {
	return &operationLogRepository{db: db}
}

func (r *operationLogRepository) Create(ctx context.Context, log *domain.OperationLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *operationLogRepository) List(ctx context.Context, params ListOperationLogParams) ([]domain.OperationLog, int64, error) {
	var logs []domain.OperationLog
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.OperationLog{})

	if params.UserID != "" {
		query = query.Where("user_id = ?", params.UserID)
	}
	if params.Method != "" {
		query = query.Where("method = ?", params.Method)
	}
	if params.Path != "" {
		query = query.Where("path LIKE ?", "%"+params.Path+"%")
	}
	if params.StartDate != nil {
		query = query.Where("created_at >= ?", params.StartDate)
	}
	if params.EndDate != nil {
		query = query.Where("created_at <= ?", params.EndDate)
	}

	query.Count(&total)

	offset := (params.Page - 1) * params.PageSize
	err := query.
		Offset(offset).
		Limit(params.PageSize).
		Order("created_at DESC").
		Find(&logs).Error

	return logs, total, err
}
