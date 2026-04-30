package handler

import (
    "encoding/json"
    "net/http"

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
    CustomerID int64      `json:"customer_id"`
    Cart       *model.Cart `json:"cart"`
}

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
        // simple error mapping
        w.WriteHeader(http.StatusConflict)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    _ = json.NewEncoder(w).Encode(ord)
}
