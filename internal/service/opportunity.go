package service

import (
	"context"
	"time"

	"github.com/qzt/backend/internal/domain"
	"github.com/qzt/backend/internal/repository"
)

// OpportunityService 商机服务
type OpportunityService interface {
	Create(ctx context.Context, req *CreateOpportunityReq) (*domain.Opportunity, error)
	GetByID(ctx context.Context, id string) (*domain.Opportunity, error)
	List(ctx context.Context, params ListOpportunityParams) (*ListResult, error)
	Update(ctx context.Context, id string, req *UpdateOpportunityReq) (*domain.Opportunity, error)
	Delete(ctx context.Context, id string) error
	// 商机阶段流转
	MoveToNextStage(ctx context.Context, id string) (*domain.Opportunity, error)
	MoveToPrevStage(ctx context.Context, id string) (*domain.Opportunity, error)
	// 商机转化
	ConvertToContract(ctx context.Context, id string) (*domain.Contract, error)
	// 统计
	GetStageStats(ctx context.Context, userID string) (map[string]int64, error)
}

type opportunityService struct {
	repo           repository.OpportunityRepository
	customerRepo   repository.CustomerRepository
	contractRepo   repository.ContractRepository
}

func NewOpportunityService(repo repository.OpportunityRepository, customerRepo repository.CustomerRepository, contractRepo repository.ContractRepository) OpportunityService {
	return &opportunityService{
		repo:         repo,
		customerRepo: customerRepo,
		contractRepo: contractRepo,
	}
}

// 商机阶段定义
const (
	StageLead       = "LEAD"        // 线索
	StageQualify    = "QUALIFY"     // 资格审查
	StageProposal   = "PROPOSAL"    // 方案报价
	StageNegotiate  = "NEGOTIATE"   // 商务谈判
	StageContract   = "CONTRACT"    // 签订合同
	StageClosedWon  = "CLOSED_WON"  // 赢单
	StageClosedLost = "CLOSED_LOST" // 输单
)

// 阶段顺序
var stageOrder = []string{
	StageLead,
	StageQualify,
	StageProposal,
	StageNegotiate,
	StageContract,
	StageClosedWon,
}

// 阶段对应的成功概率
var stageProbability = map[string]int{
	StageLead:       10,
	StageQualify:    20,
	StageProposal:   40,
	StageNegotiate:  60,
	StageContract:   80,
	StageClosedWon:  100,
	StageClosedLost: 0,
}

type CreateOpportunityReq struct {
	Name            string  `json:"name" binding:"required"`
	CustomerID      string  `json:"customerId" binding:"required"`
	Amount          float64 `json:"amount" binding:"required"`
	ExpectedCloseAt string  `json:"expectedCloseAt"`
	Description     string  `json:"description"`
}

type UpdateOpportunityReq struct {
	Name            string  `json:"name"`
	Amount          float64 `json:"amount"`
	Stage           string  `json:"stage"`
	Probability     int     `json:"probability"`
	ExpectedCloseAt string  `json:"expectedCloseAt"`
	Description     string  `json:"description"`
	Status          string  `json:"status"`
}

type ListOpportunityParams struct {
	Page     int
	PageSize int
	UserID   string
	Stage    string
	Status   string
}

func (s *opportunityService) Create(ctx context.Context, req *CreateOpportunityReq) (*domain.Opportunity, error) {
	userID, _ := ctx.Value("userId").(string)

	var expectedCloseAt *time.Time
	if req.ExpectedCloseAt != "" {
		t, err := time.Parse("2006-01-02", req.ExpectedCloseAt)
		if err == nil {
			expectedCloseAt = &t
		}
	}

	opp := &domain.Opportunity{
		Name:            req.Name,
		CustomerID:      req.CustomerID,
		UserID:          userID,
		Amount:          req.Amount,
		Stage:           StageLead,
		Probability:     stageProbability[StageLead],
		ExpectedCloseAt: expectedCloseAt,
		Description:     req.Description,
		Status:          "ACTIVE",
	}

	if err := s.repo.Create(ctx, opp); err != nil {
		return nil, err
	}

	return opp, nil
}

func (s *opportunityService) GetByID(ctx context.Context, id string) (*domain.Opportunity, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *opportunityService) List(ctx context.Context, params ListOpportunityParams) (*ListResult, error) {
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

func (s *opportunityService) Update(ctx context.Context, id string, req *UpdateOpportunityReq) (*domain.Opportunity, error) {
	opp, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		opp.Name = req.Name
	}
	if req.Amount > 0 {
		opp.Amount = req.Amount
	}
	if req.Stage != "" {
		opp.Stage = req.Stage
		opp.Probability = stageProbability[req.Stage]
	}
	if req.Probability >= 0 {
		opp.Probability = req.Probability
	}
	if req.Description != "" {
		opp.Description = req.Description
	}
	if req.Status != "" {
		opp.Status = req.Status
	}
	if req.ExpectedCloseAt != "" {
		t, err := time.Parse("2006-01-02", req.ExpectedCloseAt)
		if err == nil {
			opp.ExpectedCloseAt = &t
		}
	}

	if err := s.repo.Update(ctx, opp); err != nil {
		return nil, err
	}

	return opp, nil
}

func (s *opportunityService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

// MoveToNextStage 推进到下一阶段
func (s *opportunityService) MoveToNextStage(ctx context.Context, id string) (*domain.Opportunity, error) {
	opp, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 找到当前阶段的索引
	currentIndex := -1
	for i, stage := range stageOrder {
		if stage == opp.Stage {
			currentIndex = i
			break
		}
	}

	// 如果已经是最后一个阶段或找不到，返回错误
	if currentIndex == -1 || currentIndex >= len(stageOrder)-1 {
		return nil, ErrInvalidStage
	}

	// 推进到下一阶段
	opp.Stage = stageOrder[currentIndex+1]
	opp.Probability = stageProbability[opp.Stage]

	if err := s.repo.Update(ctx, opp); err != nil {
		return nil, err
	}

	return opp, nil
}

// MoveToPrevStage 回退到上一阶段
func (s *opportunityService) MoveToPrevStage(ctx context.Context, id string) (*domain.Opportunity, error) {
	opp, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 找到当前阶段的索引
	currentIndex := -1
	for i, stage := range stageOrder {
		if stage == opp.Stage {
			currentIndex = i
			break
		}
	}

	// 如果是第一个阶段或找不到，返回错误
	if currentIndex <= 0 {
		return nil, ErrInvalidStage
	}

	// 回退到上一阶段
	opp.Stage = stageOrder[currentIndex-1]
	opp.Probability = stageProbability[opp.Stage]

	if err := s.repo.Update(ctx, opp); err != nil {
		return nil, err
	}

	return opp, nil
}

// ConvertToContract 商机转化为合同
func (s *opportunityService) ConvertToContract(ctx context.Context, id string) (*domain.Contract, error) {
	opp, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 检查商机状态
	if opp.Stage != StageContract {
		return nil, ErrInvalidStage
	}

	// 创建合同
	contract := &domain.Contract{
		Name:       opp.Name + " - 合同",
		CustomerID: opp.CustomerID,
		Amount:     opp.Amount,
		Status:     "DRAFT",
		StartDate:  time.Now(),
	}

	if err := s.contractRepo.Create(ctx, contract); err != nil {
		return nil, err
	}

	// 更新商机状态为赢单
	opp.Stage = StageClosedWon
	opp.Probability = 100
	if err := s.repo.Update(ctx, opp); err != nil {
		return nil, err
	}

	return contract, nil
}

// GetStageStats 获取商机阶段统计
func (s *opportunityService) GetStageStats(ctx context.Context, userID string) (map[string]int64, error) {
	return s.repo.GetStageStats(ctx, userID)
}
