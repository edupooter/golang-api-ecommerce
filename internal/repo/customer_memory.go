package repo

import (
	"sync"

	"github.com/edupooter/golang-api-ecommerce/internal/model"
)

type InMemoryCustomerRepo struct {
	mu     sync.RWMutex
	data   map[int64]*model.Customer
	nextID int64
}

func NewInMemoryCustomerRepo() *InMemoryCustomerRepo {
	return &InMemoryCustomerRepo{data: make(map[int64]*model.Customer), nextID: 1}
}

func (r *InMemoryCustomerRepo) GetByID(id int64) (*model.Customer, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, ok := r.data[id]
	if !ok {
		return nil, ErrNotFound
	}
	cp := *c
	return &cp, nil
}

func (r *InMemoryCustomerRepo) Create(c *model.Customer) (*model.Customer, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	c.ID = r.nextID
	r.nextID++
	cp := *c
	r.data[c.ID] = &cp
	return &cp, nil
}
