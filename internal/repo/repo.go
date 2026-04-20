package repo

import "github.com/edupooter/golang-api-ecommerce/internal/model"

type ProductRepository interface {
	Create(p *model.Product) (*model.Product, error)
	GetAll() ([]*model.Product, error)
	GetByID(id int64) (*model.Product, error)
	Update(p *model.Product) (*model.Product, error)
	Delete(id int64) error
}
