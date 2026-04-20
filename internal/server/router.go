package server

import (
	"net/http"

	"github.com/edupooter/golang-api-ecommerce/internal/handler"
)

func NewRouter(ph *handler.ProductHandler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/products", ph.HandleProducts)
	mux.HandleFunc("/products/", ph.HandleProductByID)
	return mux
}
