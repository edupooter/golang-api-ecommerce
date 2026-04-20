package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/edupooter/golang-api-ecommerce/internal/handler"
	"github.com/edupooter/golang-api-ecommerce/internal/model"
	"github.com/edupooter/golang-api-ecommerce/internal/repo"
	"github.com/edupooter/golang-api-ecommerce/internal/server"
	"github.com/edupooter/golang-api-ecommerce/internal/service"
)

func setup() http.Handler {
	r := repo.NewInMemoryRepo()
	s := service.NewProductService(r)
	h := handler.NewProductHandler(s)
	return server.NewRouter(h)
}

func TestHTTPCreateAndGet(t *testing.T) {
	router := setup()

	prod := &model.Product{Name: "Sneakers", Price: 59.9, Stock: 3}
	b, _ := json.Marshal(prod)
	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected created, got %d", rr.Code)
	}

	var created model.Product
	if err := json.NewDecoder(rr.Body).Decode(&created); err != nil {
		t.Fatalf("decode created: %v", err)
	}

	req = httptest.NewRequest(http.MethodGet, "/products/"+strconv.FormatInt(created.ID, 10), nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected ok, got %d", rr.Code)
	}

	var got model.Product
	if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
		t.Fatalf("decode get: %v", err)
	}

	if got.ID != created.ID {
		t.Fatalf("id mismatch: got %d want %d", got.ID, created.ID)
	}
}
