package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"log"

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
		h.list(w, r)
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
		h.get(w, r, id)
	case http.MethodPut:
		h.update(w, r, id)
	case http.MethodDelete:
		h.delete(w, r, id)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *ProductHandler) list(w http.ResponseWriter, r *http.Request) {
	prods, err := h.svc.ListProducts()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(prods)
}

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

func (h *ProductHandler) get(w http.ResponseWriter, r *http.Request, id int64) {
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

func (h *ProductHandler) delete(w http.ResponseWriter, r *http.Request, id int64) {
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
