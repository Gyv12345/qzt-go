package service

import (
	"context"
	"time"

	"github.com/qzt/backend/internal/domain"
	"github.com/qzt/backend/internal/repository"
)

// ReportService 报表服务
type ReportService interface {
	// 销售漏斗
	GetSalesFunnel(ctx context.Context, userID string) (*SalesFunnelReport, error)
	// 业绩统计
	GetPerformance(ctx context.Context, params PerformanceParams) (*PerformanceReport, error)
	// 回款统计
	GetPaymentStats(ctx context.Context, params PaymentStatsParams) (*PaymentStatsReport, error)
	// 客户统计
	GetCustomerStats(ctx context.Context, userID string) (*CustomerStatsReport, error)
	// 综合报表
	GetDashboard(ctx context.Context, userID string) (*DashboardReport, error)
}

type reportService struct {
	customerRepo    repository.CustomerRepository
	opportunityRepo repository.OpportunityRepository
	contractRepo    repository.ContractRepository
	paymentRepo     repository.PaymentRepository
}

func NewReportService(
	customerRepo repository.CustomerRepository,
	opportunityRepo repository.OpportunityRepository,
	contractRepo repository.ContractRepository,
	paymentRepo repository.PaymentRepository,
) ReportService {
	return &reportService{
		customerRepo:    customerRepo,
		opportunityRepo: opportunityRepo,
		contractRepo:    contractRepo,
		paymentRepo:     paymentRepo,
	}
}

// SalesFunnelReport 销售漏斗报表
type SalesFunnelReport struct {
	Stages []FunnelStage `json:"stages"`
	Total  int64         `json:"total"`
}

type FunnelStage struct {
	Stage       string  `json:"stage"`
	Count       int64   `json:"count"`
	Amount      float64 `json:"amount"`
	Probability int     `json:"probability"`
	Weighted    float64 `json:"weighted"` // 加权金额 = amount * probability / 100
}

func (s *reportService) GetSalesFunnel(ctx context.Context, userID string) (*SalesFunnelReport, error) {
	stats, err := s.opportunityRepo.GetStageStats(ctx, userID)
	if err != nil {
		return nil, err
	}

	stages := []FunnelStage{}
	stageOrder := []string{
		"LEAD", "QUALIFY", "PROPOSAL", "NEGOTIATE", "CONTRACT",
	}

	stageProbability := map[string]int{
		"LEAD":       10,
		"QUALIFY":    20,
		"PROPOSAL":   40,
		"NEGOTIATE":  60,
		"CONTRACT":   80,
		"CLOSED_WON": 100,
	}

	var total int64
	for _, stage := range stageOrder {
		count := stats[stage]
		total += count

		stages = append(stages, FunnelStage{
			Stage:       stage,
			Count:       count,
			Probability: stageProbability[stage],
		})
	}

	return &SalesFunnelReport{
		Stages: stages,
		Total:  total,
	}, nil
}

// PerformanceParams 业绩统计参数
type PerformanceParams struct {
	UserID    string
	StartDate time.Time
	EndDate   time.Time
}

// PerformanceReport 业绩报表
type PerformanceReport struct {
	TotalOpportunities int64     `json:"totalOpportunities"`
	TotalAmount        float64   `json:"totalAmount"`
	WonAmount          float64   `json:"wonAmount"`
	LostAmount         float64   `json:"lostAmount"`
	WinRate            float64   `json:"winRate"` // 赢单率
	TotalContracts     int64     `json:"totalContracts"`
	ContractAmount     float64   `json:"contractAmount"`
	TotalPayments      float64   `json:"totalPayments"`
	MonthlyData        []MonthlyPerformance `json:"monthlyData"`
}

type MonthlyPerformance struct {
	Month      string  `json:"month"`
	Opportunities int64 `json:"opportunities"`
	Amount     float64 `json:"amount"`
	Contracts  int64   `json:"contracts"`
	Payments   float64 `json:"payments"`
}

func (s *reportService) GetPerformance(ctx context.Context, params PerformanceParams) (*PerformanceReport, error) {
	// TODO: 实现详细的业绩统计
	return &PerformanceReport{
		TotalOpportunities: 0,
		TotalAmount:        0,
		WonAmount:          0,
		LostAmount:         0,
		WinRate:            0,
		TotalContracts:     0,
		ContractAmount:     0,
		TotalPayments:      0,
		MonthlyData:        []MonthlyPerformance{},
	}, nil
}

// PaymentStatsParams 回款统计参数
type PaymentStatsParams struct {
	UserID    string
	StartDate time.Time
	EndDate   time.Time
}

// PaymentStatsReport 回款统计报表
type PaymentStatsReport struct {
	TotalAmount     float64            `json:"totalAmount"`
	ConfirmedAmount float64            `json:"confirmedAmount"`
	PendingAmount   float64            `json:"pendingAmount"`
	MonthlyData     []MonthlyPayment   `json:"monthlyData"`
	PaymentTypes    []PaymentTypeStats `json:"paymentTypes"`
}

type MonthlyPayment struct {
	Month  string  `json:"month"`
	Amount float64 `json:"amount"`
}

type PaymentTypeStats struct {
	Type   string  `json:"type"`
	Amount float64 `json:"amount"`
	Count  int64   `json:"count"`
}

func (s *reportService) GetPaymentStats(ctx context.Context, params PaymentStatsParams) (*PaymentStatsReport, error) {
	// TODO: 实现详细的回款统计
	return &PaymentStatsReport{
		TotalAmount:     0,
		ConfirmedAmount: 0,
		PendingAmount:   0,
		MonthlyData:     []MonthlyPayment{},
		PaymentTypes:    []PaymentTypeStats{},
	}, nil
}

// CustomerStatsReport 客户统计报表
type CustomerStatsReport struct {
	TotalCustomers  int64                   `json:"totalCustomers"`
	LeadCount       int64                   `json:"leadCount"`
	ProspectCount   int64                   `json:"prospectCount"`
	CustomerCount   int64                   `json:"customerCount"`
	VIPCount        int64                   `json:"vipCount"`
	IndustryStats   []IndustryStats         `json:"industryStats"`
	LevelStats      []CustomerLevelStats    `json:"levelStats"`
}

type IndustryStats struct {
	Industry string `json:"industry"`
	Count    int64  `json:"count"`
}

type CustomerLevelStats struct {
	Level string `json:"level"`
	Count int64  `json:"count"`
}

func (s *reportService) GetCustomerStats(ctx context.Context, userID string) (*CustomerStatsReport, error) {
	// TODO: 实现详细的客户统计
	return &CustomerStatsReport{
		TotalCustomers: 0,
		LeadCount:      0,
		ProspectCount:  0,
		CustomerCount:  0,
		VIPCount:       0,
		IndustryStats:  []IndustryStats{},
		LevelStats:     []CustomerLevelStats{},
	}, nil
}

// DashboardReport 综合报表（仪表盘）
type DashboardReport struct {
	Summary      DashboardSummary      `json:"summary"`
	SalesFunnel  *SalesFunnelReport    `json:"salesFunnel"`
	RecentItems  RecentItems           `json:"recentItems"`
	Trends       []TrendData           `json:"trends"`
}

type DashboardSummary struct {
	TotalCustomers    int64   `json:"totalCustomers"`
	TotalOpportunities int64  `json:"totalOpportunities"`
	TotalContracts    int64   `json:"totalContracts"`
	TotalPayments     float64 `json:"totalPayments"`
	WinRate           float64 `json:"winRate"`
}

type RecentItems struct {
	RecentCustomers    []domain.Customer    `json:"recentCustomers"`
	RecentOpportunities []domain.Opportunity `json:"recentOpportunities"`
	RecentPayments     []domain.Payment     `json:"recentPayments"`
}

type TrendData struct {
	Date        string  `json:"date"`
	Customers   int64   `json:"customers"`
	Opportunities int64 `json:"opportunities"`
	Amount      float64 `json:"amount"`
}

func (s *reportService) GetDashboard(ctx context.Context, userID string) (*DashboardReport, error) {
	// 获取销售漏斗
	funnel, err := s.GetSalesFunnel(ctx, userID)
	if err != nil {
		return nil, err
	}

	// TODO: 实现完整的仪表盘数据
	return &DashboardReport{
		Summary: DashboardSummary{
			TotalCustomers:     0,
			TotalOpportunities: funnel.Total,
			TotalContracts:     0,
			TotalPayments:      0,
			WinRate:            0,
		},
		SalesFunnel: funnel,
		RecentItems: RecentItems{
			RecentCustomers:     []domain.Customer{},
			RecentOpportunities: []domain.Opportunity{},
			RecentPayments:      []domain.Payment{},
		},
		Trends: []TrendData{},
	}, nil
}
