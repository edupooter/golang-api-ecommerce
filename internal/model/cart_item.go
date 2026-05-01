package model

type CartItem struct {
	ProductID int64 `json:"product_id" example:"1"`
	Quantity  int   `json:"quantity" example:"2"`
}
