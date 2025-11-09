// Package handler provides HTTP handlers for orders.
package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
	CreatedAt string            `json:"created_at"`
	UpdatedAt string            `json:"updated_at"`
}

// CreateOrder handles creating a new order.
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	order, err := h.service.CreateOrder(req.ID, req.UserID, req.Items)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp := &OrderResponse{
		ID:        order.ID,
		UserID:    order.UserID,
		Items:     order.Items,
		Status:    order.Status,
		Total:     order.Total,
		CreatedAt: order.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: order.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	c.JSON(http.StatusCreated, resp)
}

// GetOrder handles getting an order by ID.
func (h *OrderHandler) GetOrder(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing id parameter"})
		return
	}

	order, err := h.service.GetOrder(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	resp := &OrderResponse{
		ID:        order.ID,
		UserID:    order.UserID,
		Items:     order.Items,
		Status:    order.Status,
		Total:     order.Total,
		CreatedAt: order.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: order.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	c.JSON(http.StatusOK, resp)
}

// PayOrder handles paying an order.
func (h *OrderHandler) PayOrder(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing id parameter"})
		return
	}

	if err := h.service.PayOrder(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "order paid successfully"})
}

// ShipOrder handles shipping an order.
func (h *OrderHandler) ShipOrder(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing id parameter"})
		return
	}

	if err := h.service.ShipOrder(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "order shipped successfully"})
}

// DeliverOrder handles delivering an order.
func (h *OrderHandler) DeliverOrder(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing id parameter"})
		return
	}

	if err := h.service.DeliverOrder(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "order delivered successfully"})
}

// CancelOrder handles cancelling an order.
func (h *OrderHandler) CancelOrder(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing id parameter"})
		return
	}

	if err := h.service.CancelOrder(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "order cancelled successfully"})
}
