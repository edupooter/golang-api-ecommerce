package server

import (
	"net/http"
	"os"

	_ "github.com/edupooter/golang-api-ecommerce/docs"
	"github.com/edupooter/golang-api-ecommerce/internal/handler"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(ph *handler.ProductHandler) http.Handler {
	return newRouter(ph, nil)
}

func NewRouterWithCheckout(ph *handler.ProductHandler, ch *handler.CheckoutHandler) http.Handler {
	return newRouter(ph, ch)
}

func newRouter(ph *handler.ProductHandler, ch *handler.CheckoutHandler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/products", ph.HandleProducts)
	mux.HandleFunc("/products/", ph.HandleProductByID)
	if ch != nil {
		mux.HandleFunc("/checkout", ch.HandleCheckout)
	}
	// allow overriding the swagger doc URL in production via SWAGGER_URL
	swaggerURL := os.Getenv("SWAGGER_URL")
	if swaggerURL == "" {
		swaggerURL = "http://localhost:8080/swagger/doc.json"
	}
	mux.Handle("/swagger/", httpSwagger.Handler(
		httpSwagger.URL(swaggerURL),
	))
	return mux
}
