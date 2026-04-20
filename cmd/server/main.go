package main

import (
	"log"
	"net/http"
	"os"

	"github.com/edupooter/golang-api-ecommerce/internal/handler"
	"github.com/edupooter/golang-api-ecommerce/internal/repo"
	"github.com/edupooter/golang-api-ecommerce/internal/server"
	"github.com/edupooter/golang-api-ecommerce/internal/service"
)

func main() {
	var rRepo repo.ProductRepository
	// if SQLITE_PATH is set, use SQLite-backed repository, otherwise in-memory
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
		rRepo = srepo
	} else {
		rRepo = repo.NewInMemoryRepo()
	}

	s := service.NewProductService(rRepo)
	h := handler.NewProductHandler(s)
	router := server.NewRouter(h)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	log.Printf("server listening on %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}
