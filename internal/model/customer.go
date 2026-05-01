package model

type Customer struct {
	ID    int64  `json:"id" example:"1"`
	Name  string `json:"name" example:"João Exemplo"`
	Email string `json:"email" example:"joao@example.com"`
}
