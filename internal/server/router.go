package server

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/edupooter/golang-api-ecommerce/docs"
	"github.com/edupooter/golang-api-ecommerce/internal/handler"
)

func NewRouter(ph *handler.ProductHandler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/products", ph.HandleProducts)
	mux.HandleFunc("/products/", ph.HandleProductByID)
	mux.Handle("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))
	return mux
}
