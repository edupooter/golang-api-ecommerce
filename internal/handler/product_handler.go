package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/edupooter/golang-api-ecommerce/internal/model"
	"github.com/edupooter/golang-api-ecommerce/internal/repo"
	"github.com/edupooter/golang-api-ecommerce/internal/service"
)

type ProductHandler struct {
	svc *service.ProductService
}

func NewProductHandler(svc *service.ProductService) *ProductHandler {
	return &ProductHandler{svc: svc}
}

func (h *ProductHandler) HandleProducts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.list(w)
	case http.MethodPost:
		h.create(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *ProductHandler) HandleProductByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/products/")
	if idStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	switch r.Method {
	case http.MethodGet:
		h.get(w, id)
	case http.MethodPut:
		h.update(w, r, id)
	case http.MethodDelete:
		h.delete(w, id)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// ListProducts godoc
// @Summary List products
// @Tags products
// @Produce json
// @Success 200 {array} model.Product
// @Failure 500 {object} map[string]string
// @Router /products [get]
func (h *ProductHandler) list(w http.ResponseWriter) {
	prods, err := h.svc.ListProducts()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(prods)
}

// CreateProduct godoc
// @Summary Create a product
// @Tags products
// @Accept json
// @Produce json
// @Param product body model.Product true "product"
// @Success 201 {object} model.Product
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /products [post]
func (h *ProductHandler) create(w http.ResponseWriter, r *http.Request) {
	var p model.Product
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(body, &p); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	created, err := h.svc.CreateProduct(&p)
	if err != nil {
		if err == service.ErrInvalidProduct {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	// expose created id in Location and a custom header to help REST clients extract it
	loc := "/products/" + strconv.FormatInt(created.ID, 10)
	idstr := strconv.FormatInt(created.ID, 10)
	w.Header().Set("Location", loc)
	w.Header().Set("X-Resource-ID", idstr)
	w.WriteHeader(http.StatusCreated)
	// log headers for debugging
	log.Printf("created product id=%s, Location=%s, X-Resource-ID=%s", idstr, loc, idstr)
	_ = json.NewEncoder(w).Encode(created)
}

// GetProduct godoc
// @Summary Get a product by ID
// @Tags products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} model.Product
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /products/{id} [get]
func (h *ProductHandler) get(w http.ResponseWriter, id int64) {
	p, err := h.svc.GetProduct(id)
	if err != nil {
		if err == repo.ErrNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(p)
}

// UpdateProduct godoc
// @Summary Update a product
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param product body model.Product true "product"
// @Success 200 {object} model.Product
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /products/{id} [put]
func (h *ProductHandler) update(w http.ResponseWriter, r *http.Request, id int64) {
	var p model.Product
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(body, &p); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	p.ID = id
	updated, err := h.svc.UpdateProduct(&p)
	if err != nil {
		if err == service.ErrInvalidProduct {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err == repo.ErrNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(updated)
}

// DeleteProduct godoc
// @Summary Delete a product
// @Tags products
// @Param id path int true "Product ID"
// @Success 204 {string} string "No Content"
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /products/{id} [delete]
func (h *ProductHandler) delete(w http.ResponseWriter, id int64) {
	if err := h.svc.DeleteProduct(id); err != nil {
		if err == repo.ErrNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
