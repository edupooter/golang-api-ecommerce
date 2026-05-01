// @title Golang API E-commerce
// @version 1.0
// @description API simples de e-commerce
// @host localhost:8080
// @BasePath /
// @schemes http
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/edupooter/golang-api-ecommerce/internal/handler"
	"github.com/edupooter/golang-api-ecommerce/internal/ports"
	"github.com/edupooter/golang-api-ecommerce/internal/repo"
	"github.com/edupooter/golang-api-ecommerce/internal/server"
	"github.com/edupooter/golang-api-ecommerce/internal/service"
)

func main() {
	var prodRepo repo.ProductRepository
	var orderRepo interface{}
	var custRepo interface{}

	// if SQLITE_PATH is set, use SQLite-backed repositories, otherwise in-memory
	if path := os.Getenv("SQLITE_PATH"); path != "" {
		srepo, err := repo.NewSQLiteRepo(path)
		if err != nil {
			log.Fatalf("failed to open sqlite db (%s): %v", path, err)
		}
		defer func() {
			if err := srepo.Close(); err != nil {
				log.Printf("error closing sqlite db: %v", err)
			}
		}()
		prodRepo = srepo

		or, err := repo.NewSQLiteOrderRepo(path)
		if err != nil {
			log.Fatalf("failed to open sqlite order repo (%s): %v", path, err)
		}
		defer or.Close()
		orderRepo = or

		cr, err := repo.NewSQLiteCustomerRepo(path)
		if err != nil {
			log.Fatalf("failed to open sqlite customer repo (%s): %v", path, err)
		}
		defer cr.Close()
		custRepo = cr
	} else {
		prod := repo.NewInMemoryRepo()
		prodRepo = prod
		orderRepo = repo.NewInMemoryOrderRepo()
		custRepo = repo.NewInMemoryCustomerRepo()
	}

	psvc := service.NewProductService(prodRepo)
	ph := handler.NewProductHandler(psvc)

	// build OrderService using concrete adapters (type assertions)
	var stock ports.StockRepository
	var or ports.OrderRepository
	var cr ports.CustomerRepository

	// stock repo is implemented by product repo implementations
	switch v := prodRepo.(type) {
	case *repo.InMemoryRepo:
		stock = v
	case *repo.SQLiteRepo:
		stock = v
	}
	switch v := orderRepo.(type) {
	case *repo.InMemoryOrderRepo:
		or = v
	case *repo.SQLiteOrderRepo:
		or = v
	}
	switch v := custRepo.(type) {
	case *repo.InMemoryCustomerRepo:
		cr = v
	case *repo.SQLiteCustomerRepo:
		cr = v
	}

	osvc := service.NewOrderService(stock, or, cr)
	ch := handler.NewCheckoutHandler(osvc)

	router := server.NewRouterWithCheckout(ph, ch)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	log.Printf("server listening on %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}
