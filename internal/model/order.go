package model

type OrderItem struct {
    ProductID int64 `json:"product_id"`
    Quantity  int   `json:"quantity"`
    Price     float64 `json:"price"`
}

type Order struct {
    ID         int64       `json:"id"`
    CustomerID int64       `json:"customer_id"`
    Items      []OrderItem `json:"items"`
    Total      float64     `json:"total"`
}

func (o *Order) CalculateTotal() {
    var t float64
    for _, it := range o.Items {
        t += float64(it.Quantity) * it.Price
    }
    o.Total = t
}
