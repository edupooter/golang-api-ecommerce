package repo

import (
	"sync"

	"github.com/edupooter/golang-api-ecommerce/internal/model"
)

type InMemoryOrderRepo struct {
	mu     sync.Mutex
	data   map[int64]*model.Order
	nextID int64
}

func NewInMemoryOrderRepo() *InMemoryOrderRepo {
	return &InMemoryOrderRepo{data: make(map[int64]*model.Order), nextID: 1}
}

func (r *InMemoryOrderRepo) Create(o *model.Order) (*model.Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	o.ID = r.nextID
	r.nextID++
	cp := *o
	r.data[o.ID] = &cp
	return &cp, nil
}
