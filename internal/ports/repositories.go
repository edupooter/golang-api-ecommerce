package ports

import (
	"errors"

	"github.com/edupooter/golang-api-ecommerce/internal/model"
)

// Errors exported for coordination between services and adapters
var ErrInsufficientStock = errors.New("insufficient stock")

// StockRepository provides atomic-ish stock operations used by checkout
type StockRepository interface {
	DecrementStock(id int64, qty int) error
	IncrementStock(id int64, qty int) error
}

// OrderRepository persists orders
type OrderRepository interface {
	Create(o *model.Order) (*model.Order, error)
}

// CustomerRepository provides minimal customer access
type CustomerRepository interface {
	GetByID(id int64) (*model.Customer, error)
	Create(c *model.Customer) (*model.Customer, error)
}

// CartRepository stores/retrieves carts
type CartRepository interface {
	Get(id string) (*model.Cart, error)
	Save(c *model.Cart) error
}
