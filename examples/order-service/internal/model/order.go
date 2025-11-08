// Package model contains the domain models for the order service.
package model

import (
	"errors"
	"time"
)

// OrderStatus represents the status of an order.
type OrderStatus string

const (
	// OrderStatusPending represents an order that has been created but not yet paid.
	OrderStatusPending OrderStatus = "pending"

	// OrderStatusPaid represents an order that has been paid.
	OrderStatusPaid OrderStatus = "paid"

	// OrderStatusShipped represents an order that has been shipped.
	OrderStatusShipped OrderStatus = "shipped"

	// OrderStatusDelivered represents an order that has been delivered.
	OrderStatusDelivered OrderStatus = "delivered"

	// OrderStatusCancelled represents an order that has been cancelled.
	OrderStatusCancelled OrderStatus = "cancelled"
)

// Order represents an order in the system.
type Order struct {
	ID        string      `json:"id" db:"id"`
	UserID    string      `json:"user_id" db:"user_id"`
	Items     []OrderItem `json:"items" db:"items"`
	Status    OrderStatus `json:"status" db:"status"`
	Total     float64     `json:"total" db:"total"`
	CreatedAt time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt time.Time   `json:"updated_at" db:"updated_at"`
}

// OrderItem represents an item in an order.
type OrderItem struct {
	ProductID string  `json:"product_id" db:"product_id"`
	Name      string  `json:"name" db:"name"`
	Price     float64 `json:"price" db:"price"`
	Quantity  int     `json:"quantity" db:"quantity"`
}

// NewOrder creates a new order.
func NewOrder(id, userID string, items []OrderItem) (*Order, error) {
	if id == "" {
		return nil, errors.New("order id cannot be empty")
	}

	if userID == "" {
		return nil, errors.New("user id cannot be empty")
	}

	if len(items) == 0 {
		return nil, errors.New("order must have at least one item")
	}

	total := 0.0
	for _, item := range items {
		if item.Quantity <= 0 {
			return nil, errors.New("item quantity must be greater than zero")
		}
		total += item.Price * float64(item.Quantity)
	}

	now := time.Now()
	return &Order{
		ID:        id,
		UserID:    userID,
		Items:     items,
		Status:    OrderStatusPending,
		Total:     total,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// Pay marks the order as paid.
func (o *Order) Pay() error {
	if o.Status != OrderStatusPending {
		return errors.New("only pending orders can be paid")
	}

	o.Status = OrderStatusPaid
	o.UpdatedAt = time.Now()
	return nil
}

// Ship marks the order as shipped.
func (o *Order) Ship() error {
	if o.Status != OrderStatusPaid {
		return errors.New("only paid orders can be shipped")
	}

	o.Status = OrderStatusShipped
	o.UpdatedAt = time.Now()
	return nil
}

// Deliver marks the order as delivered.
func (o *Order) Deliver() error {
	if o.Status != OrderStatusShipped {
		return errors.New("only shipped orders can be delivered")
	}

	o.Status = OrderStatusDelivered
	o.UpdatedAt = time.Now()
	return nil
}

// Cancel cancels the order.
func (o *Order) Cancel() error {
	if o.Status == OrderStatusDelivered {
		return errors.New("delivered orders cannot be cancelled")
	}

	o.Status = OrderStatusCancelled
	o.UpdatedAt = time.Now()
	return nil
}
