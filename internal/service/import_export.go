package service

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/qzt/backend/internal/domain"
	"github.com/qzt/backend/internal/repository"
)

// ImportExportService 导入导出服务
type ImportExportService interface {
	// 客户导入导出
	ExportCustomers(ctx context.Context, params ExportCustomersParams) ([]byte, error)
	ImportCustomers(ctx context.Context, data []byte) (*ImportResult, error)
	
	// 商机导入导出
	ExportOpportunities(ctx context.Context, params ExportOpportunitiesParams) ([]byte, error)
	ImportOpportunities(ctx context.Context, data []byte) (*ImportResult, error)
	
	// 合同导入导出
	ExportContracts(ctx context.Context, params ExportContractsParams) ([]byte, error)
	ImportContracts(ctx context.Context, data []byte) (*ImportResult, error)
	
	// 回款导入导出
	ExportPayments(ctx context.Context, params ExportPaymentsParams) ([]byte, error)
	ImportPayments(ctx context.Context, data []byte) (*ImportResult, error)
}

type importExportService struct {
	customerRepo   repository.CustomerRepository
	opportunityRepo repository.OpportunityRepository
	contractRepo   repository.ContractRepository
	paymentRepo    repository.PaymentRepository
}

func NewImportExportService(
	customerRepo repository.CustomerRepository,
	opportunityRepo repository.OpportunityRepository,
	contractRepo repository.ContractRepository,
	paymentRepo repository.PaymentRepository,
) ImportExportService {
	return &importExportService{
		customerRepo:   customerRepo,
		opportunityRepo: opportunityRepo,
		contractRepo:   contractRepo,
		paymentRepo:    paymentRepo,
	}
}

// ExportCustomersParams 导出客户参数
type ExportCustomersParams struct {
	Format string // csv/json/xlsx
	UserID  string
}

// ExportOpportunitiesParams 导出商机参数
type ExportOpportunitiesParams struct {
	Format string
	UserID  string
	Stage   string
}

// ExportContractsParams 导出合同参数
type ExportContractsParams struct {
	Format string
	UserID  string
	Status  string
}

// ExportPaymentsParams 导出回款参数
type ExportPaymentsParams struct {
	Format string
	StartDate time.Time
	EndDate   time.Time
}

// ImportResult 导入结果
type ImportResult struct {
	SuccessCount int      `json:"successCount"`
	FailCount    int      `json:"failCount"`
	Errors       []string `json:"errors"`
}

// ========== 客户导入导出 ==========

func (s *importExportService) ExportCustomers(ctx context.Context, params ExportCustomersParams) ([]byte, error) {
	// 查询客户数据
	_, items, err := s.customerRepo.List(ctx, "", 1, 1000)
	if err != nil {
		return nil, err
	}

	switch params.Format {
	case "csv":
		return s.exportCustomersToCSV(items)
	case "json":
		return json.MarshalIndent(items, "", "  ")
	default:
		return json.Marshal(items)
	}
}

func (s *importExportService) exportCustomersToCSV(customers []domain.Customer) ([]byte, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	
	// 写入表头
	headers := []string{
		"ID", "Name", "Contact", "Phone", "Email", "Industry",
		"Address", "Level", "Status", "CreatedAt",
	}
	if err := writer.Write(headers); err != nil {
		return nil, err
	}
	
	// 写入数据
	for _, c := range customers {
		row := []string{
			c.ID,
			c.Name,
			c.Contact,
			c.Phone,
			c.Email,
			c.Industry,
			c.Address,
			c.Level,
			c.Status,
			c.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		if err := writer.Write(row); err != nil {
			return nil, err
		}
	}
	
	writer.Flush()
	return buf.Bytes(), nil
}

func (s *importExportService) ImportCustomers(ctx context.Context, data []byte) (*ImportResult, error) {
	result := &ImportResult{}
	
	// 解析 CSV
	reader := csv.NewReader(bytes.NewReader(data))
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	
	if len(records) == 0 {
		return result, nil
	}
	
	// 跳过表头
	dataRecords := records
	if len(records) > 0 && records[0][0] == "ID" {
		dataRecords = records[1:]
	}
	
	userID, _ := ctx.Value("userId").(string)
	
	// 批量导入
	for i, record := range dataRecords {
		if len(record) < 5 {
			result.FailCount++
			result.Errors = append(result.Errors, fmt.Sprintf("第%d行：数据不完整", i+2))
			continue
		}
		
		customer := &domain.Customer{
			Name:     record[1],
			Contact:  record[2],
			Phone:    record[3],
			Email:    record[4],
			Industry: getFieldValue(record, 5, ""),
			Address:  getFieldValue(record, 6, ""),
			Level:    getFieldValue(record, 7, "LEAD"),
			Status:   getFieldValue(record, 8, "ACTIVE"),
		}
		
		if err := s.customerRepo.Create(ctx, customer); err != nil {
			result.FailCount++
			result.Errors = append(result.Errors, fmt.Sprintf("第%d行：%v", i+2, err))
		} else {
			result.SuccessCount++
		}
	}
	
	return result, nil
}

// ========== 商机导入导出 ==========

func (s *importExportService) ExportOpportunities(ctx context.Context, params ExportOpportunitiesParams) ([]byte, error) {
	listParams := repository.ListOpportunityParams{
		Page:     1,
		PageSize: 1000,
		UserID:   params.UserID,
		Stage:    params.Stage,
	}
	
	_, items, err := s.opportunityRepo.List(ctx, listParams)
	if err != nil {
		return nil, err
	}

	switch params.Format {
	case "csv":
		return s.exportOpportunitiesToCSV(items)
	case "json":
		return json.MarshalIndent(items, "", "  ")
	default:
		return json.Marshal(items)
	}
}

func (s *importExportService) exportOpportunitiesToCSV(opps []domain.Opportunity) ([]byte, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	
	headers := []string{
		"ID", "Name", "CustomerID", "Amount", "Stage",
		"Probability", "Status", "CreatedAt",
	}
	if err := writer.Write(headers); err != nil {
		return nil, err
	}
	
	for _, o := range opps {
		row := []string{
			o.ID,
			o.Name,
			o.CustomerID,
			strconv.FormatFloat(o.Amount, 'f', 2, 64),
			o.Stage,
			strconv.Itoa(o.Probability),
			o.Status,
			o.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		if err := writer.Write(row); err != nil {
			return nil, err
		}
	}
	
	writer.Flush()
	return buf.Bytes(), nil
}

func (s *importExportService) ImportOpportunities(ctx context.Context, data []byte) (*ImportResult, error) {
	result := &ImportResult{}
	
	reader := csv.NewReader(bytes.NewReader(data))
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	
	if len(records) == 0 {
		return result, nil
	}
	
	dataRecords := records
	if len(records) > 0 && records[0][0] == "ID" {
		dataRecords = records[1:]
	}
	
	userID, _ := ctx.Value("userId").(string)
	
	for i, record := range dataRecords {
		if len(record) < 4 {
			result.FailCount++
			result.Errors = append(result.Errors, fmt.Sprintf("第%d行：数据不完整", i+2))
			continue
		}
		
		amount, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			result.FailCount++
			result.Errors = append(result.Errors, fmt.Sprintf("第%d行：金额格式错误", i+2))
			continue
		}
		
		probability, _ := strconv.Atoi(getFieldValue(record, 5, "10"))
		
		opp := &domain.Opportunity{
			Name:       record[1],
			CustomerID: record[2],
			Amount:     amount,
			UserID:     userID,
			Stage:      getFieldValue(record, 4, "LEAD"),
			Probability: probability,
			Status:     getFieldValue(record, 6, "ACTIVE"),
		}
		
		if err := s.opportunityRepo.Create(ctx, opp); err != nil {
			result.FailCount++
			result.Errors = append(result.Errors, fmt.Sprintf("第%d行：%v", i+2, err))
		} else {
			result.SuccessCount++
		}
	}
	
	return result, nil
}

// ========== 合同导入导出 ==========

func (s *importExportService) ExportContracts(ctx context.Context, params ExportContractsParams) ([]byte, error) {
	_, items, err := s.contractRepo.List(ctx, "", 1, 1000)
	if err != nil {
		return nil, err
	}

	switch params.Format {
	case "csv":
		return s.exportContractsToCSV(items)
	case "json":
		return json.MarshalIndent(items, "", "  ")
	default:
		return json.Marshal(items)
	}
}

func (s *importExportService) exportContractsToCSV(contracts []domain.Contract) ([]byte, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	
	headers := []string{
		"ID", "Name", "CustomerID", "Amount", "StartDate",
		"EndDate", "Status", "CreatedAt",
	}
	if err := writer.Write(headers); err != nil {
		return nil, err
	}
	
	for _, c := range contracts {
		row := []string{
			c.ID,
			c.Name,
			c.CustomerID,
			strconv.FormatFloat(c.Amount, 'f', 2, 64),
			c.StartDate.Format("2006-01-02"),
			c.EndDate.Format("2006-01-02"),
			c.Status,
			c.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		if err := writer.Write(row); err != nil {
			return nil, err
		}
	}
	
	writer.Flush()
	return buf.Bytes(), nil
}

func (s *importExportService) ImportContracts(ctx context.Context, data []byte) (*ImportResult, error) {
	// TODO: 实现合同导入
	return &ImportResult{SuccessCount: 0, FailCount: 0}, nil
}

// ========== 回款导入导出 ==========

func (s *importExportService) ExportPayments(ctx context.Context, params ExportPaymentsParams) ([]byte, error) {
	listParams := repository.ListPaymentParams{
		Page:     1,
		PageSize: 1000,
		StartDate: &params.StartDate,
		EndDate:   &params.EndDate,
	}
	
	_, items, err := s.paymentRepo.List(ctx, listParams)
	if err != nil {
		return nil, err
	}

	switch params.Format {
	case "csv":
		return s.exportPaymentsToCSV(items)
	case "json":
		return json.MarshalIndent(items, "", "  ")
	default:
		return json.Marshal(items)
	}
}

func (s *importExportService) exportPaymentsToCSV(payments []domain.Payment) ([]byte, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	
	headers := []string{
		"ID", "ContractID", "CustomerID", "Amount", "PaymentDate",
		"PaymentType", "Status", "CreatedAt",
	}
	if err := writer.Write(headers); err != nil {
		return nil, err
	}
	
	for _, p := range payments {
		row := []string{
			p.ID,
			p.ContractID,
			p.CustomerID,
			strconv.FormatFloat(p.Amount, 'f', 2, 64),
			p.PaymentDate.Format("2006-01-02"),
			p.PaymentType,
			p.Status,
			p.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		if err := writer.Write(row); err != nil {
			return nil, err
		}
	}
	
	writer.Flush()
	return buf.Bytes(), nil
}

func (s *importExportService) ImportPayments(ctx context.Context, data []byte) (*ImportResult, error) {
	// TODO: 实现回款导入
	return &ImportResult{SuccessCount: 0, FailCount: 0}, nil
}

// ========== 辅助函数 ==========

func getFieldValue(record []string, index int, defaultValue string) string {
	if index >= len(record) {
		return defaultValue
	}
	if record[index] == "" {
		return defaultValue
	}
	return record[index]
}
