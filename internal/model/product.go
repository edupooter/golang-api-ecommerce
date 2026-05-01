package model

type Product struct {
	ID    int64   `json:"id" example:"1"`
	Name  string  `json:"name" example:"Camiseta Golang"`
	Price float64 `json:"price" example:"49.9"`
	Stock int     `json:"stock" example:"10"`
}
