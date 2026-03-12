package service

import (
	"context"
	"time"

	"github.com/qzt/backend/internal/domain"
	"github.com/qzt/backend/internal/repository"
)

// InvoiceService 发票服务
type InvoiceService interface {
	Create(ctx context.Context, req *CreateInvoiceReq) (*domain.Invoice, error)
	GetByID(ctx context.Context, id string) (*domain.Invoice, error)
	List(ctx context.Context, params ListInvoiceParams) (*ListResult, error)
	Update(ctx context.Context, id string, req *UpdateInvoiceReq) (*domain.Invoice, error)
	Delete(ctx context.Context, id string) error
	// 发票操作
	Submit(ctx context.Context, id string) (*domain.Invoice, error)
	Approve(ctx context.Context, id string) (*domain.Invoice, error)
	Reject(ctx context.Context, id string, reason string) (*domain.Invoice, error)
	Invoice(ctx context.Context, id string) (*domain.Invoice, error)
	Receive(ctx context.Context, id string) (*domain.Invoice, error)
	// 统计
	GetStats(ctx context.Context, params InvoiceStatsParams) (*InvoiceStatsReport, error)
}

type invoiceService struct {
	repo       repository.InvoiceRepository
	contract   repository.ContractRepository
	payment    repository.PaymentRepository
}

func NewInvoiceService(repo repository.InvoiceRepository, contract repository.ContractRepository, payment repository.PaymentRepository) InvoiceService {
	return &invoiceService{
		repo:     repo,
		contract: contract,
		payment:  payment,
	}
}

type CreateInvoiceReq struct {
	ContractID    string  `json:"contractId" binding:"required"`
	CustomerID    string  `json:"customerId" binding:"required"`
	InvoiceNumber string  `json:"invoiceNumber" binding:"required"`
	InvoiceType   string  `json:"invoiceType" binding:"required"` // 普票/专票
	Title         string  `json:"title" binding:"required"`
	TaxNumber     string  `json:"taxNumber"`
	Amount        float64 `json:"amount" binding:"required"`
	TaxRate       float64 `json:"taxRate"`
	TaxAmount     float64 `json:"taxAmount"`
	BankName      string  `json:"bankName"`
	BankAccount   string  `json:"bankAccount"`
	Address       string  `json:"address"`
	Phone         string  `json:"phone"`
	Remark        string  `json:"remark"`
}

type UpdateInvoiceReq struct {
	InvoiceNumber string  `json:"invoiceNumber"`
	InvoiceType   string  `json:"invoiceType"`
	Title         string  `json:"title"`
	TaxNumber     string  `json:"taxNumber"`
	Amount        float64 `json:"amount"`
	TaxRate       float64 `json:"taxRate"`
	TaxAmount     float64 `json:"taxAmount"`
	BankName      string  `json:"bankName"`
	BankAccount   string  `json:"bankAccount"`
	Address       string  `json:"address"`
	Phone         string  `json:"phone"`
	Remark        string  `json:"remark"`
}

type ListInvoiceParams struct {
	Page         int
	PageSize     int
	ContractID   string
	CustomerID   string
	Status       string
	InvoiceType  string
	StartDate    *time.Time
	EndDate      *time.Time
}

type InvoiceStatsParams struct {
	UserID    string
	StartDate time.Time
	EndDate   time.Time
}

type InvoiceStatsReport struct {
	TotalAmount     float64             `json:"totalAmount"`
	DraftCount      int64               `json:"draftCount"`
	SubmittedCount  int64               `json:"submittedCount"`
	ApprovedCount   int64               `json:"approvedCount"`
	InvoicedCount   int64               `json:"invoicedCount"`
	ReceivedCount   int64               `json:"receivedCount"`
	MonthlyData     []MonthlyInvoice    `json:"monthlyData"`
}

type MonthlyInvoice struct {
	Month  string  `json:"month"`
	Count  int64   `json:"count"`
	Amount float64 `json:"amount"`
}

func (s *invoiceService) Create(ctx context.Context, req *CreateInvoiceReq) (*domain.Invoice, error) {
	invoice := &domain.Invoice{
		ContractID:    req.ContractID,
		CustomerID:    req.CustomerID,
		InvoiceNumber: req.InvoiceNumber,
		InvoiceType:   req.InvoiceType,
		Title:         req.Title,
		TaxNumber:     req.TaxNumber,
		Amount:        req.Amount,
		TaxRate:       req.TaxRate,
		TaxAmount:     req.TaxAmount,
		BankName:      req.BankName,
		BankAccount:   req.BankAccount,
		Address:       req.Address,
		Phone:         req.Phone,
		Remark:        req.Remark,
		Status:        "DRAFT",
	}

	if err := s.repo.Create(ctx, invoice); err != nil {
		return nil, err
	}

	return invoice, nil
}

func (s *invoiceService) GetByID(ctx context.Context, id string) (*domain.Invoice, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *invoiceService) List(ctx context.Context, params ListInvoiceParams) (*ListResult, error) {
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

func (s *invoiceService) Update(ctx context.Context, id string, req *UpdateInvoiceReq) (*domain.Invoice, error) {
	invoice, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 只有草稿状态可以修改
	if invoice.Status != "DRAFT" {
		return nil, ErrInvalidStatus
	}

	if req.InvoiceNumber != "" {
		invoice.InvoiceNumber = req.InvoiceNumber
	}
	if req.InvoiceType != "" {
		invoice.InvoiceType = req.InvoiceType
	}
	if req.Title != "" {
		invoice.Title = req.Title
	}
	if req.Amount > 0 {
		invoice.Amount = req.Amount
	}

	if err := s.repo.Update(ctx, invoice); err != nil {
		return nil, err
	}

	return invoice, nil
}

func (s *invoiceService) Delete(ctx context.Context, id string) error {
	invoice, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// 只有草稿状态可以删除
	if invoice.Status != "DRAFT" {
		return ErrInvalidStatus
	}

	return s.repo.Delete(ctx, id)
}

// Submit 提交审批
func (s *invoiceService) Submit(ctx context.Context, id string) (*domain.Invoice, error) {
	invoice, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if invoice.Status != "DRAFT" {
		return nil, ErrInvalidStatus
	}

	invoice.Status = "SUBMITTED"
	if err := s.repo.Update(ctx, invoice); err != nil {
		return nil, err
	}

	return invoice, nil
}

// Approve 审批通过
func (s *invoiceService) Approve(ctx context.Context, id string) (*domain.Invoice, error) {
	invoice, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if invoice.Status != "SUBMITTED" {
		return nil, ErrInvalidStatus
	}

	invoice.Status = "APPROVED"
	now := time.Now()
	invoice.ApprovedAt = &now

	if err := s.repo.Update(ctx, invoice); err != nil {
		return nil, err
	}

	return invoice, nil
}

// Reject 审批拒绝
func (s *invoiceService) Reject(ctx context.Context, id string, reason string) (*domain.Invoice, error) {
	invoice, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if invoice.Status != "SUBMITTED" {
		return nil, ErrInvalidStatus
	}

	invoice.Status = "REJECTED"
	invoice.Remark = reason

	if err := s.repo.Update(ctx, invoice); err != nil {
		return nil, err
	}

	return invoice, nil
}

// Invoice 已开票
func (s *invoiceService) Invoice(ctx context.Context, id string) (*domain.Invoice, error) {
	invoice, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if invoice.Status != "APPROVED" {
		return nil, ErrInvalidStatus
	}

	invoice.Status = "INVOICED"
	now := time.Now()
	invoice.InvoiceDate = &now

	if err := s.repo.Update(ctx, invoice); err != nil {
		return nil, err
	}

	return invoice, nil
}

// Receive 已收到
func (s *invoiceService) Receive(ctx context.Context, id string) (*domain.Invoice, error) {
	invoice, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if invoice.Status != "INVOICED" {
		return nil, ErrInvalidStatus
	}

	invoice.Status = "RECEIVED"
	now := time.Now()
	invoice.ReceiveDate = &now

	if err := s.repo.Update(ctx, invoice); err != nil {
		return nil, err
	}

	return invoice, nil
}

// GetStats 发票统计
func (s *invoiceService) GetStats(ctx context.Context, params InvoiceStatsParams) (*InvoiceStatsReport, error) {
	stats, err := s.repo.GetStats(ctx, params)
	if err != nil {
		return nil, err
	}

	return stats, nil
}

// 错误定义
var (
	ErrInvalidStatus = New(10020, "状态不允许此操作")
)
