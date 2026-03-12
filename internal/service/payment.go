package service

import (
	"context"
	"time"

	"github.com/qzt/backend/internal/domain"
	"github.com/qzt/backend/internal/repository"
)

// PaymentService 回款服务
type PaymentService interface {
	Create(ctx context.Context, req *CreatePaymentReq) (*domain.Payment, error)
	GetByID(ctx context.Context, id string) (*domain.Payment, error)
	List(ctx context.Context, params repository.ListPaymentParams) (*ListResult, error)
	Confirm(ctx context.Context, id string) (*domain.Payment, error)
	Delete(ctx context.Context, id string) error
}

type paymentService struct {
	repo repository.PaymentRepository
}

func NewPaymentService(repo repository.PaymentRepository) PaymentService {
	return &paymentService{repo: repo}
}

type CreatePaymentReq struct {
	ContractID  string
	CustomerID  string
	Amount      float64
	PaymentDate time.Time
	PaymentType string
	Remark      string
}

func (s *paymentService) Create(ctx context.Context, req *CreatePaymentReq) (*domain.Payment, error) {
	payment := &domain.Payment{
		ContractID:  req.ContractID,
		CustomerID:  req.CustomerID,
		Amount:      req.Amount,
		PaymentDate: req.PaymentDate,
		PaymentType: req.PaymentType,
		Status:      "PENDING",
		Remark:      req.Remark,
	}

	if err := s.repo.Create(ctx, payment); err != nil {
		return nil, err
	}

	return payment, nil
}

func (s *paymentService) GetByID(ctx context.Context, id string) (*domain.Payment, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *paymentService) List(ctx context.Context, params repository.ListPaymentParams) (*ListResult, error) {
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

func (s *paymentService) Confirm(ctx context.Context, id string) (*domain.Payment, error) {
	payment, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	payment.Status = "CONFIRMED"
	if err := s.repo.Update(ctx, payment); err != nil {
		return nil, err
	}

	return payment, nil
}

func (s *paymentService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

// 错误定义
var (
	ErrInvalidStage = New(10010, "无效的阶段")
)
