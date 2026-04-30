package service

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/edupooter/golang-api-ecommerce/internal/model"
	"github.com/edupooter/golang-api-ecommerce/internal/ports"
	"github.com/edupooter/golang-api-ecommerce/internal/repo"
)

// mockStockRepo implements ports.StockRepository for unit tests
type mockStockRepo struct {
	mu          sync.Mutex
	decremented map[int64]int
	incremented map[int64]int
	fail        map[int64]error
}

func (m *mockStockRepo) DecrementStock(id int64, qty int) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.fail != nil {
		if err, ok := m.fail[id]; ok && err != nil {
			return err
		}
	}
	if m.decremented == nil {
		m.decremented = make(map[int64]int)
	}
	m.decremented[id] += qty
	return nil
}

func (m *mockStockRepo) IncrementStock(id int64, qty int) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.incremented == nil {
		m.incremented = make(map[int64]int)
	}
	m.incremented[id] += qty
	return nil
}

type mockOrderRepo struct {
	mu      sync.Mutex
	created []*model.Order
}

func (m *mockOrderRepo) Create(o *model.Order) (*model.Order, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	o.ID = int64(len(m.created) + 1)
	m.created = append(m.created, o)
	return o, nil
}

type mockCustomerRepo struct{ exists bool }

func (m *mockCustomerRepo) GetByID(id int64) (*model.Customer, error) {
	if m.exists {
		return &model.Customer{ID: id, Name: "Test", Email: "t@example.com"}, nil
	}
	return nil, errors.New("not found")
}

func (m *mockCustomerRepo) Create(c *model.Customer) (*model.Customer, error) {
	if c == nil {
		return nil, errors.New("nil customer")
	}
	c.ID = 1
	return c, nil
}

func TestOrderService_Checkout_Success(t *testing.T) {
	ms := &mockStockRepo{}
	mo := &mockOrderRepo{}
	mc := &mockCustomerRepo{exists: true}
	svc := NewOrderService(ms, mo, mc)

	cart := &model.Cart{Items: []model.CartItem{{ProductID: 1, Quantity: 2}, {ProductID: 2, Quantity: 1}}}
	ord, err := svc.Checkout(1, cart)
	if err != nil {
		t.Fatalf("expected success, got err: %v", err)
	}
	if ord == nil {
		t.Fatalf("expected order, got nil")
	}

	ms.mu.Lock()
	if ms.decremented[1] != 2 {
		t.Fatalf("expected decremented[1]=2 got %d", ms.decremented[1])
	}
	if ms.decremented[2] != 1 {
		t.Fatalf("expected decremented[2]=1 got %d", ms.decremented[2])
	}
	ms.mu.Unlock()

	mo.mu.Lock()
	if len(mo.created) != 1 {
		t.Fatalf("expected order repo called once, got %d", len(mo.created))
	}
	mo.mu.Unlock()
}

func TestOrderService_Checkout_RollbackOnInsufficientStock(t *testing.T) {
	ms := &mockStockRepo{fail: map[int64]error{2: ports.ErrInsufficientStock}}
	mo := &mockOrderRepo{}
	svc := NewOrderService(ms, mo, nil)

	cart := &model.Cart{Items: []model.CartItem{{ProductID: 1, Quantity: 1}, {ProductID: 2, Quantity: 1}}}
	_, err := svc.Checkout(1, cart)
	if err == nil {
		t.Fatalf("expected error due to insufficient stock")
	}

	ms.mu.Lock()
	if ms.incremented[1] != 1 {
		t.Fatalf("expected compensation increment for product 1, got %d", ms.incremented[1])
	}
	ms.mu.Unlock()
}

func TestOrderService_Concurrency_InMemoryRepo(t *testing.T) {
	r := repo.NewInMemoryRepo()
	p, _ := r.Create(&model.Product{Name: "P", Price: 10, Stock: 5})
	svc := NewOrderService(r, nil, nil)
	cart := &model.Cart{Items: []model.CartItem{{ProductID: p.ID, Quantity: 1}}}

	var wg sync.WaitGroup
	var success int32
	attempts := 10
	wg.Add(attempts)
	for i := 0; i < attempts; i++ {
		go func() {
			defer wg.Done()
			_, err := svc.Checkout(1, cart)
			if err == nil {
				atomic.AddInt32(&success, 1)
			}
		}()
	}
	wg.Wait()

	if atomic.LoadInt32(&success) != int32(5) {
		t.Fatalf("expected 5 successes, got %d", success)
	}
	prod, err := r.GetByID(p.ID)
	if err != nil {
		t.Fatalf("unexpected get error: %v", err)
	}
	if prod.Stock != 0 {
		t.Fatalf("expected final stock 0, got %d", prod.Stock)
	}
}
