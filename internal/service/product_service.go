package service

import (
	"errors"

	"github.com/edupooter/golang-api-ecommerce/internal/model"
	"github.com/edupooter/golang-api-ecommerce/internal/repo"
)

var ErrInvalidProduct = errors.New("invalid product")

type ProductService struct {
	repo repo.ProductRepository
}

func NewProductService(r repo.ProductRepository) *ProductService {
	return &ProductService{repo: r}
}

func (s *ProductService) CreateProduct(p *model.Product) (*model.Product, error) {
	if p.Name == "" || p.Price < 0 || p.Stock < 0 {
		return nil, ErrInvalidProduct
	}
	return s.repo.Create(p)
}

func (s *ProductService) ListProducts() ([]*model.Product, error) {
	return s.repo.GetAll()
}

func (s *ProductService) GetProduct(id int64) (*model.Product, error) {
	return s.repo.GetByID(id)
}

func (s *ProductService) UpdateProduct(p *model.Product) (*model.Product, error) {
	if p.Name == "" || p.Price < 0 || p.Stock < 0 {
		return nil, ErrInvalidProduct
	}
	return s.repo.Update(p)
}

func (s *ProductService) DeleteProduct(id int64) error {
	return s.repo.Delete(id)
}
