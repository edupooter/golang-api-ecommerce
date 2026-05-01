package model

type OrderItem struct {
	ProductID int64   `json:"product_id" example:"1"`
	Quantity  int     `json:"quantity" example:"2"`
	Price     float64 `json:"price" example:"49.9"`
}

type Order struct {
	ID         int64       `json:"id" example:"1"`
	CustomerID int64       `json:"customer_id" example:"1"`
	Items      []OrderItem `json:"items"`
	Total      float64     `json:"total" example:"119.7"`
}

func (o *Order) CalculateTotal() {
	var t float64
	for _, it := range o.Items {
		t += float64(it.Quantity) * it.Price
	}
	o.Total = t
}
