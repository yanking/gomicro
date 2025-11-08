// Package repository provides data access functionality for orders.
package repository

import (
	"errors"
	"sync"

	"github.com/yanking/gomicro/examples/order-service/internal/model"
)

// OrderRepository defines the interface for order data access.
type OrderRepository interface {
	// Save saves an order.
	Save(order *model.Order) error

	// FindByID finds an order by ID.
	FindByID(id string) (*model.Order, error)

	// FindByUserID finds orders by user ID.
	FindByUserID(userID string) ([]*model.Order, error)

	// Update updates an order.
	Update(order *model.Order) error

	// Delete deletes an order by ID.
	Delete(id string) error
}

// InMemoryOrderRepository implements OrderRepository using in-memory storage.
type InMemoryOrderRepository struct {
	orders map[string]*model.Order
	mutex  sync.RWMutex
}

// NewInMemoryOrderRepository creates a new InMemoryOrderRepository.
func NewInMemoryOrderRepository() *InMemoryOrderRepository {
	return &InMemoryOrderRepository{
		orders: make(map[string]*model.Order),
	}
}

// Save saves an order.
func (r *InMemoryOrderRepository) Save(order *model.Order) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.orders[order.ID] = order
	return nil
}

// FindByID finds an order by ID.
func (r *InMemoryOrderRepository) FindByID(id string) (*model.Order, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	order, exists := r.orders[id]
	if !exists {
		return nil, errors.New("order not found")
	}

	return order, nil
}

// FindByUserID finds orders by user ID.
func (r *InMemoryOrderRepository) FindByUserID(userID string) ([]*model.Order, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var orders []*model.Order
	for _, order := range r.orders {
		if order.UserID == userID {
			orders = append(orders, order)
		}
	}

	return orders, nil
}

// Update updates an order.
func (r *InMemoryOrderRepository) Update(order *model.Order) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	_, exists := r.orders[order.ID]
	if !exists {
		return errors.New("order not found")
	}

	r.orders[order.ID] = order
	return nil
}

// Delete deletes an order by ID.
func (r *InMemoryOrderRepository) Delete(id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	_, exists := r.orders[id]
	if !exists {
		return errors.New("order not found")
	}

	delete(r.orders, id)
	return nil
}
