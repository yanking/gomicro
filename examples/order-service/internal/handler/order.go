// Package handler provides HTTP handlers for orders.
package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/yanking/gomicro/examples/order-service/internal/model"
	"github.com/yanking/gomicro/examples/order-service/internal/service"
)

// OrderHandler handles HTTP requests for orders.
type OrderHandler struct {
	service *service.OrderService
}

// NewOrderHandler creates a new OrderHandler.
func NewOrderHandler(service *service.OrderService) *OrderHandler {
	return &OrderHandler{
		service: service,
	}
}

// CreateOrderRequest represents the request body for creating an order.
type CreateOrderRequest struct {
	ID     string            `json:"id"`
	UserID string            `json:"user_id"`
	Items  []model.OrderItem `json:"items"`
}

// OrderResponse represents the response body for order operations.
type OrderResponse struct {
	ID        string            `json:"id"`
	UserID    string            `json:"user_id"`
	Items     []model.OrderItem `json:"items"`
	Status    model.OrderStatus `json:"status"`
	Total     float64           `json:"total"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

// CreateOrder handles creating a new order.
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	order, err := h.service.CreateOrder(req.ID, req.UserID, req.Items)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := &OrderResponse{
		ID:        order.ID,
		UserID:    order.UserID,
		Items:     order.Items,
		Status:    order.Status,
		Total:     order.Total,
		CreatedAt: order.CreatedAt,
		UpdatedAt: order.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(resp)
}

// GetOrder handles getting an order by ID.
func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "missing id parameter", http.StatusBadRequest)
		return
	}

	order, err := h.service.GetOrder(id)
	if err != nil {
		http.Error(w, "order not found", http.StatusNotFound)
		return
	}

	resp := &OrderResponse{
		ID:        order.ID,
		UserID:    order.UserID,
		Items:     order.Items,
		Status:    order.Status,
		Total:     order.Total,
		CreatedAt: order.CreatedAt,
		UpdatedAt: order.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

// PayOrder handles paying an order.
func (h *OrderHandler) PayOrder(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "missing id parameter", http.StatusBadRequest)
		return
	}

	if err := h.service.PayOrder(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("order paid successfully"))
}

// ShipOrder handles shipping an order.
func (h *OrderHandler) ShipOrder(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "missing id parameter", http.StatusBadRequest)
		return
	}

	if err := h.service.ShipOrder(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("order shipped successfully"))
}

// DeliverOrder handles delivering an order.
func (h *OrderHandler) DeliverOrder(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "missing id parameter", http.StatusBadRequest)
		return
	}

	if err := h.service.DeliverOrder(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("order delivered successfully"))
}

// CancelOrder handles cancelling an order.
func (h *OrderHandler) CancelOrder(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "missing id parameter", http.StatusBadRequest)
		return
	}

	if err := h.service.CancelOrder(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("order cancelled successfully"))
}
