package repository

import (
	"context"

	"github.com/qzt/backend/internal/domain"
	"gorm.io/gorm"
)

// OpportunityRepository 商机仓储
type OpportunityRepository interface {
	Create(ctx context.Context, opp *domain.Opportunity) error
	GetByID(ctx context.Context, id string) (*domain.Opportunity, error)
	List(ctx context.Context, params ListOpportunityParams) ([]domain.Opportunity, int64, error)
	Update(ctx context.Context, opp *domain.Opportunity) error
	Delete(ctx context.Context, id string) error
	GetStageStats(ctx context.Context, userID string) (map[string]int64, error)
}

type opportunityRepository struct {
	db *gorm.DB
}

func NewOpportunityRepository(db *gorm.DB) OpportunityRepository {
	return &opportunityRepository{db: db}
}

func (r *opportunityRepository) Create(ctx context.Context, opp *domain.Opportunity) error {
	return r.db.WithContext(ctx).Create(opp).Error
}

func (r *opportunityRepository) GetByID(ctx context.Context, id string) (*domain.Opportunity, error) {
	var opp domain.Opportunity
	err := r.db.WithContext(ctx).
		Preload("Customer").
		Preload("User").
		First(&opp, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &opp, nil
}

type ListOpportunityParams struct {
	Page     int
	PageSize int
	UserID   string
	Stage    string
	Status   string
}

func (r *opportunityRepository) List(ctx context.Context, params ListOpportunityParams) ([]domain.Opportunity, int64, error) {
	var opps []domain.Opportunity
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.Opportunity{})

	if params.UserID != "" {
		query = query.Where("user_id = ?", params.UserID)
	}
	if params.Stage != "" {
		query = query.Where("stage = ?", params.Stage)
	}
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	query.Count(&total)

	offset := (params.Page - 1) * params.PageSize
	err := query.
		Preload("Customer").
		Preload("User").
		Offset(offset).
		Limit(params.PageSize).
		Order("created_at DESC").
		Find(&opps).Error

	return opps, total, err
}

func (r *opportunityRepository) Update(ctx context.Context, opp *domain.Opportunity) error {
	return r.db.WithContext(ctx).Save(opp).Error
}

func (r *opportunityRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&domain.Opportunity{}, "id = ?", id).Error
}

func (r *opportunityRepository) GetStageStats(ctx context.Context, userID string) (map[string]int64, error) {
	type Result struct {
		Stage string
		Count int64
	}

	var results []Result
	query := r.db.WithContext(ctx).
		Model(&domain.Opportunity{}).
		Select("stage, count(*) as count").
		Where("status = ?", "ACTIVE")

	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	err := query.Group("stage").Scan(&results).Error
	if err != nil {
		return nil, err
	}

	stats := make(map[string]int64)
	for _, r := range results {
		stats[r.Stage] = r.Count
	}

	return stats, nil
}

// PaymentRepository 回款仓储
type PaymentRepository interface {
	Create(ctx context.Context, payment *domain.Payment) error
	GetByID(ctx context.Context, id string) (*domain.Payment, error)
	List(ctx context.Context, params ListPaymentParams) ([]domain.Payment, int64, error)
	Update(ctx context.Context, payment *domain.Payment) error
	Delete(ctx context.Context, id string) error
	GetTotalByContract(ctx context.Context, contractID string) (float64, error)
}

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) Create(ctx context.Context, payment *domain.Payment) error {
	return r.db.WithContext(ctx).Create(payment).Error
}

func (r *paymentRepository) GetByID(ctx context.Context, id string) (*domain.Payment, error) {
	var payment domain.Payment
	err := r.db.WithContext(ctx).
		Preload("Contract").
		Preload("Customer").
		First(&payment, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

type ListPaymentParams struct {
	Page       int
	PageSize   int
	ContractID string
	CustomerID string
	Status     string
}

func (r *paymentRepository) List(ctx context.Context, params ListPaymentParams) ([]domain.Payment, int64, error) {
	var payments []domain.Payment
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.Payment{})

	if params.ContractID != "" {
		query = query.Where("contract_id = ?", params.ContractID)
	}
	if params.CustomerID != "" {
		query = query.Where("customer_id = ?", params.CustomerID)
	}
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	query.Count(&total)

	offset := (params.Page - 1) * params.PageSize
	err := query.
		Preload("Contract").
		Preload("Customer").
		Offset(offset).
		Limit(params.PageSize).
		Order("payment_date DESC").
		Find(&payments).Error

	return payments, total, err
}

func (r *paymentRepository) Update(ctx context.Context, payment *domain.Payment) error {
	return r.db.WithContext(ctx).Save(payment).Error
}

func (r *paymentRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&domain.Payment{}, "id = ?", id).Error
}

func (r *paymentRepository) GetTotalByContract(ctx context.Context, contractID string) (float64, error) {
	var total float64
	err := r.db.WithContext(ctx).
		Model(&domain.Payment{}).
		Where("contract_id = ? AND status = ?", contractID, "CONFIRMED").
		Select("COALESCE(SUM(amount), 0)").
		Scan(&total).Error
	return total, err
}
