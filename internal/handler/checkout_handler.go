package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/edupooter/golang-api-ecommerce/internal/model"
	"github.com/edupooter/golang-api-ecommerce/internal/service"
)

type CheckoutHandler struct {
	svc *service.OrderService
}

func NewCheckoutHandler(svc *service.OrderService) *CheckoutHandler {
	return &CheckoutHandler{svc: svc}
}

type checkoutRequest struct {
	CustomerID int64       `json:"customer_id"`
	Cart       *model.Cart `json:"cart"`
}

// Checkout godoc
// @Summary Checkout cart and create order
// @Description Decrements stock and creates an order for the given customer and cart
// @Tags checkout
// @Accept json
// @Produce json
// @Param payload body checkoutRequest true "checkout payload"
// @Success 201 {object} model.Order
// @Header 201 {string} Location "Location of created order"
// @Header 201 {string} X-Order-ID "Created order id"
// @Failure 400 {object} model.ErrorResponse
// @Failure 409 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /checkout [post]
func (h *CheckoutHandler) HandleCheckout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req checkoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	ord, err := h.svc.Checkout(req.CustomerID, req.Cart)
	if err != nil {
		// simple error mapping (map ports.ErrInsufficientStock -> 409 elsewhere if needed)
		w.WriteHeader(http.StatusConflict)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if ord != nil && ord.ID != 0 {
		loc := "/orders/" + strconv.FormatInt(ord.ID, 10)
		w.Header().Set("Location", loc)
		w.Header().Set("X-Order-ID", strconv.FormatInt(ord.ID, 10))
	}
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(ord)
}
