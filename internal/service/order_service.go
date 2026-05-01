package service

import (
	"errors"

	"github.com/edupooter/golang-api-ecommerce/internal/model"
	"github.com/edupooter/golang-api-ecommerce/internal/ports"
)

var (
	ErrInvalidCart = errors.New("invalid cart")
)

type OrderService struct {
	stockRepo ports.StockRepository
	orderRepo ports.OrderRepository
	custRepo  ports.CustomerRepository
}

func NewOrderService(s ports.StockRepository, o ports.OrderRepository, c ports.CustomerRepository) *OrderService {
	return &OrderService{stockRepo: s, orderRepo: o, custRepo: c}
}

// Checkout attempts to decrement stock for each cart item and create an order.
// It uses a simple compensation strategy: if any decrement fails, already decremented
// items are rolled back by incrementing stock back.
func (s *OrderService) Checkout(customerID int64, cart *model.Cart) (*model.Order, error) {
	if cart == nil || len(cart.Items) == 0 {
		return nil, ErrInvalidCart
	}
	// validate customer exists (if repository provided)
	if s.custRepo != nil {
		if _, err := s.custRepo.GetByID(customerID); err != nil {
			return nil, err
		}
	}

	// track succeeded decrements to allow compensation
	type dec struct {
		id  int64
		qty int
	}
	var succeeded []dec

	for _, it := range cart.Items {
		if it.Quantity <= 0 {
			// rollback and return
			for _, d := range succeeded {
				_ = s.stockRepo.IncrementStock(d.id, d.qty)
			}
			return nil, errors.New("invalid item quantity")
		}
		if err := s.stockRepo.DecrementStock(it.ProductID, it.Quantity); err != nil {
			// compensate
			for _, d := range succeeded {
				_ = s.stockRepo.IncrementStock(d.id, d.qty)
			}
			return nil, err
		}
		succeeded = append(succeeded, dec{id: it.ProductID, qty: it.Quantity})
	}

	// build order
	ord := &model.Order{CustomerID: customerID}
	for _, it := range cart.Items {
		// price resolution could be done via product repo; keep simple
		ord.Items = append(ord.Items, model.OrderItem{ProductID: it.ProductID, Quantity: it.Quantity, Price: 0})
	}
	ord.CalculateTotal()

	if s.orderRepo != nil {
		o, err := s.orderRepo.Create(ord)
		if err != nil {
			// rollback stock
			for _, d := range succeeded {
				_ = s.stockRepo.IncrementStock(d.id, d.qty)
			}
			return nil, err
		}
		return o, nil
	}

	return ord, nil
}
