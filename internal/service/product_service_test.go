package service_test

import (
	"testing"

	"github.com/edupooter/golang-api-ecommerce/internal/model"
	"github.com/edupooter/golang-api-ecommerce/internal/repo"
	"github.com/edupooter/golang-api-ecommerce/internal/service"
)

func TestCreateAndGetProduct(t *testing.T) {
	r := repo.NewInMemoryRepo()
	s := service.NewProductService(r)

	p := &model.Product{Name: "T-shirt", Price: 19.99, Stock: 10}
	created, err := s.CreateProduct(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if created.ID == 0 {
		t.Fatalf("expected id assigned")
	}

	got, err := s.GetProduct(created.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Name != p.Name || got.Price != p.Price || got.Stock != p.Stock {
		t.Fatalf("retrieved product mismatch: got=%v want=%v", got, p)
	}
}

func TestUpdateAndDeleteProduct(t *testing.T) {
	r := repo.NewInMemoryRepo()
	s := service.NewProductService(r)

	p := &model.Product{Name: "Book", Price: 9.99, Stock: 5}
	created, _ := s.CreateProduct(p)
	created.Price = 12.50
	updated, err := s.UpdateProduct(created)
	if err != nil {
		t.Fatalf("update error: %v", err)
	}
	if updated.Price != 12.50 {
		t.Fatalf("price not updated")
	}

	if err := s.DeleteProduct(created.ID); err != nil {
		t.Fatalf("delete error: %v", err)
	}

	_, err = s.GetProduct(created.ID)
	if err == nil {
		t.Fatalf("expected error after delete")
	}
}
