package repo

import (
	"errors"
	"sync"

	"github.com/edupooter/golang-api-ecommerce/internal/model"
	"github.com/edupooter/golang-api-ecommerce/internal/ports"
)

var ErrNotFound = errors.New("product not found")

type InMemoryRepo struct {
	mu     sync.RWMutex
	data   map[int64]*model.Product
	nextID int64
}

func NewInMemoryRepo() *InMemoryRepo {
	return &InMemoryRepo{
		data:   make(map[int64]*model.Product),
		nextID: 1,
	}
}

func (r *InMemoryRepo) Create(p *model.Product) (*model.Product, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	p.ID = r.nextID
	r.nextID++
	cp := *p
	r.data[p.ID] = &cp
	return &cp, nil
}

func (r *InMemoryRepo) GetAll() ([]*model.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]*model.Product, 0, len(r.data))
	for _, v := range r.data {
		cp := *v
		res = append(res, &cp)
	}
	return res, nil
}

func (r *InMemoryRepo) GetByID(id int64) (*model.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.data[id]
	if !ok {
		return nil, ErrNotFound
	}
	cp := *p
	return &cp, nil
}

func (r *InMemoryRepo) Update(p *model.Product) (*model.Product, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, ok := r.data[p.ID]
	if !ok {
		return nil, ErrNotFound
	}
	cp := *p
	r.data[p.ID] = &cp
	return &cp, nil
}

func (r *InMemoryRepo) Delete(id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, ok := r.data[id]
	if !ok {
		return ErrNotFound
	}
	delete(r.data, id)
	return nil
}

// DecrementStock decrements stock if enough quantity exists.
func (r *InMemoryRepo) DecrementStock(id int64, qty int) error {
	if qty <= 0 {
		return nil
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	p, ok := r.data[id]
	if !ok {
		return ErrNotFound
	}
	if p.Stock < qty {
		return ports.ErrInsufficientStock
	}
	p.Stock -= qty
	return nil
}

// IncrementStock increases stock (used for compensation)
func (r *InMemoryRepo) IncrementStock(id int64, qty int) error {
	if qty <= 0 {
		return nil
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	p, ok := r.data[id]
	if !ok {
		return ErrNotFound
	}
	p.Stock += qty
	return nil
}
