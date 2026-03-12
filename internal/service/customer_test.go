// Package service_test 测试业务逻辑层
package service_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/qzt/backend/internal/domain"
	"github.com/qzt/backend/internal/repository"
	"github.com/qzt/backend/internal/service"
)

// MockCustomerRepository 模拟客户仓储
type MockCustomerRepository struct {
	customers map[string]*domain.Customer
}

func NewMockCustomerRepository() *MockCustomerRepository {
	return &MockCustomerRepository{
		customers: make(map[string]*domain.Customer),
	}
}

func (m *MockCustomerRepository) Create(ctx context.Context, customer *domain.Customer) error {
	m.customers[customer.ID] = customer
	return nil
}

func (m *MockCustomerRepository) GetByID(ctx context.Context, id string) (*domain.Customer, error) {
	if customer, ok := m.customers[id]; ok {
		return customer, nil
	}
	return nil, repository.ErrNotFound
}

func (m *MockCustomerRepository) List(ctx context.Context, params repository.ListParams) ([]domain.Customer, int64, error) {
	var result []domain.Customer
	for _, c := range m.customers {
		result = append(result, *c)
	}
	return result, int64(len(result)), nil
}

func (m *MockCustomerRepository) Update(ctx context.Context, customer *domain.Customer) error {
	m.customers[customer.ID] = customer
	return nil
}

func (m *MockCustomerRepository) Delete(ctx context.Context, id string) error {
	delete(m.customers, id)
	return nil
}

// TestCustomerService_Create 测试创建客户
func TestCustomerService_Create(t *testing.T) {
	repo := NewMockCustomerRepository()
	svc := service.NewCustomerService(repo)

	req := &service.CreateCustomerReq{
		Name:      "测试公司",
		ShortName: "测试",
		Industry:  "IT",
		Scale:     "100-500人",
		Address:   "北京市海淀区",
		Website:   "https://example.com",
		Tags:      []string{"VIP", "重要客户"},
	}

	customer, err := svc.Create(context.Background(), req)
	if err != nil {
		t.Fatalf("创建客户失败: %v", err)
	}

	if customer.Name != req.Name {
		t.Errorf("期望名称 %s, 实际 %s", req.Name, customer.Name)
	}

	if customer.CustomerLevel != domain.CustomerLevelLead {
		t.Errorf("期望等级 %s, 实际 %s", domain.CustomerLevelLead, customer.CustomerLevel)
	}

	// 验证标签序列化
	var tags []string
	if err := json.Unmarshal([]byte(customer.Tags), &tags); err != nil {
		t.Fatalf("解析标签失败: %v", err)
	}
	if len(tags) != 2 {
		t.Errorf("期望 2 个标签, 实际 %d", len(tags))
	}
}

// TestCustomerService_GetByID 测试获取客户
func TestCustomerService_GetByID(t *testing.T) {
	repo := NewMockCustomerRepository()
	svc := service.NewCustomerService(repo)

	// 先创建一个客户
	req := &service.CreateCustomerReq{
		Name: "测试公司",
	}
	created, _ := svc.Create(context.Background(), req)

	// 然后获取它
	customer, err := svc.GetByID(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("获取客户失败: %v", err)
	}

	if customer.Name != created.Name {
		t.Errorf("期望名称 %s, 实际 %s", created.Name, customer.Name)
	}
}

// TestCustomerService_Update 测试更新客户
func TestCustomerService_Update(t *testing.T) {
	repo := NewMockCustomerRepository()
	svc := service.NewCustomerService(repo)

	// 先创建一个客户
	createReq := &service.CreateCustomerReq{
		Name: "测试公司",
	}
	created, _ := svc.Create(context.Background(), createReq)

	// 更新客户
	newName := "新公司名称"
	status := 0
	updateReq := &service.UpdateCustomerReq{
		Name:   newName,
		Status: &status,
	}

	customer, err := svc.Update(context.Background(), created.ID, updateReq)
	if err != nil {
		t.Fatalf("更新客户失败: %v", err)
	}

	if customer.Name != newName {
		t.Errorf("期望名称 %s, 实际 %s", newName, customer.Name)
	}

	if customer.Status != 0 {
		t.Errorf("期望状态 0, 实际 %d", customer.Status)
	}
}

// TestCustomerService_Delete 测试删除客户
func TestCustomerService_Delete(t *testing.T) {
	repo := NewMockCustomerRepository()
	svc := service.NewCustomerService(repo)

	// 先创建一个客户
	req := &service.CreateCustomerReq{
		Name: "测试公司",
	}
	created, _ := svc.Create(context.Background(), req)

	// 删除客户
	err := svc.Delete(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("删除客户失败: %v", err)
	}

	// 验证已删除
	_, err = svc.GetByID(context.Background(), created.ID)
	if err == nil {
		t.Error("期望客户已删除，但仍能获取到")
	}
}

// TestCustomerService_List 测试获取客户列表
func TestCustomerService_List(t *testing.T) {
	repo := NewMockCustomerRepository()
	svc := service.NewCustomerService(repo)

	// 创建多个客户
	for i := 0; i < 5; i++ {
		req := &service.CreateCustomerReq{
			Name: "测试公司",
		}
		svc.Create(context.Background(), req)
	}

	// 获取列表
	params := repository.ListParams{
		Page:     1,
		PageSize: 10,
	}

	result, err := svc.List(context.Background(), params)
	if err != nil {
		t.Fatalf("获取列表失败: %v", err)
	}

	if result.Total != 5 {
		t.Errorf("期望总数 5, 实际 %d", result.Total)
	}
}
