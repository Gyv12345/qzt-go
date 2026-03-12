package repository

import (
	"context"
	"time"

	"github.com/qzt/backend/internal/domain"
	"github.com/qzt/backend/internal/service"
	"gorm.io/gorm"
)

// InvoiceRepository 发票仓储
type InvoiceRepository interface {
	Create(ctx context.Context, invoice *domain.Invoice) error
	GetByID(ctx context.Context, id string) (*domain.Invoice, error)
	List(ctx context.Context, params service.ListInvoiceParams) ([]domain.Invoice, int64, error)
	Update(ctx context.Context, invoice *domain.Invoice) error
	Delete(ctx context.Context, id string) error
	GetStats(ctx context.Context, params service.InvoiceStatsParams) (*service.InvoiceStatsReport, error)
}

type invoiceRepository struct {
	db *gorm.DB
}

func NewInvoiceRepository(db *gorm.DB) InvoiceRepository {
	return &invoiceRepository{db: db}
}

func (r *invoiceRepository) Create(ctx context.Context, invoice *domain.Invoice) error {
	return r.db.WithContext(ctx).Create(invoice).Error
}

func (r *invoiceRepository) GetByID(ctx context.Context, id string) (*domain.Invoice, error) {
	var invoice domain.Invoice
	err := r.db.WithContext(ctx).
		Preload("Contract").
		Preload("Customer").
		First(&invoice, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &invoice, nil
}

func (r *invoiceRepository) List(ctx context.Context, params service.ListInvoiceParams) ([]domain.Invoice, int64, error) {
	var invoices []domain.Invoice
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.Invoice{})

	if params.ContractID != "" {
		query = query.Where("contract_id = ?", params.ContractID)
	}
	if params.CustomerID != "" {
		query = query.Where("customer_id = ?", params.CustomerID)
	}
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}
	if params.InvoiceType != "" {
		query = query.Where("invoice_type = ?", params.InvoiceType)
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
		Preload("Contract").
		Preload("Customer").
		Offset(offset).
		Limit(params.PageSize).
		Order("created_at DESC").
		Find(&invoices).Error

	return invoices, total, err
}

func (r *invoiceRepository) Update(ctx context.Context, invoice *domain.Invoice) error {
	return r.db.WithContext(ctx).Save(invoice).Error
}

func (r *invoiceRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&domain.Invoice{}, "id = ?", id).Error
}

func (r *invoiceRepository) GetStats(ctx context.Context, params service.InvoiceStatsParams) (*service.InvoiceStatsReport, error) {
	report := &service.InvoiceStatsReport{}

	query := r.db.WithContext(ctx).Model(&domain.Invoice{})

	// 统计各状态数量
	type StatusCount struct {
		Status string
		Count  int64
	}

	var statusCounts []StatusCount
	err := query.Select("status, count(*) as count").
		Group("status").
		Scan(&statusCounts).Error
	if err != nil {
		return nil, err
	}

	for _, sc := range statusCounts {
		switch sc.Status {
		case "DRAFT":
			report.DraftCount = sc.Count
		case "SUBMITTED":
			report.SubmittedCount = sc.Count
		case "APPROVED":
			report.ApprovedCount = sc.Count
		case "INVOICED":
			report.InvoicedCount = sc.Count
		case "RECEIVED":
			report.ReceivedCount = sc.Count
		}
	}

	// 统计总金额
	var totalAmount float64
	err = r.db.WithContext(ctx).
		Model(&domain.Invoice{}).
		Where("status IN ?", []string{"APPROVED", "INVOICED", "RECEIVED"}).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalAmount).Error
	if err != nil {
		return nil, err
	}
	report.TotalAmount = totalAmount

	// 按月统计
	type MonthlyResult struct {
		Month  string
		Count  int64
		Amount float64
	}

	var monthlyResults []MonthlyResult
	err = r.db.WithContext(ctx).
		Model(&domain.Invoice{}).
		Select("TO_CHAR(created_at, 'YYYY-MM') as month, count(*) as count, COALESCE(SUM(amount), 0) as amount").
		Where("created_at >= ? AND created_at <= ?", params.StartDate, params.EndDate).
		Group("TO_CHAR(created_at, 'YYYY-MM')").
		Order("month").
		Scan(&monthlyResults).Error
	if err != nil {
		return nil, err
	}

	report.MonthlyData = make([]service.MonthlyInvoice, len(monthlyResults))
	for i, mr := range monthlyResults {
		report.MonthlyData[i] = service.MonthlyInvoice{
			Month:  mr.Month,
			Count:  mr.Count,
			Amount: mr.Amount,
		}
	}

	return report, nil
}
