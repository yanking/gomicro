// Package service provides business logic for orders.
package service

import (
	"github.com/yanking/gomicro/examples/order-service/internal/model"
	"github.com/yanking/gomicro/examples/order-service/internal/repository"
)

// OrderService provides order business logic.
type OrderService struct {
	repo repository.OrderRepository
}

// NewOrderService creates a new OrderService.
func NewOrderService(repo repository.OrderRepository) *OrderService {
	return &OrderService{
		repo: repo,
	}
}

// CreateOrder creates a new order.
func (s *OrderService) CreateOrder(id, userID string, items []model.OrderItem) (*model.Order, error) {
	order, err := model.NewOrder(id, userID, items)
	if err != nil {
		return nil, err
	}

	err = s.repo.Save(order)
	if err != nil {
		return nil, err
	}

	return order, nil
}

// GetOrder gets an order by ID.
func (s *OrderService) GetOrder(id string) (*model.Order, error) {
	return s.repo.FindByID(id)
}

// GetUserOrders gets orders by user ID.
func (s *OrderService) GetUserOrders(userID string) ([]*model.Order, error) {
	return s.repo.FindByUserID(userID)
}

// PayOrder pays an order.
func (s *OrderService) PayOrder(id string) error {
	order, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	err = order.Pay()
	if err != nil {
		return err
	}

	return s.repo.Update(order)
}

// ShipOrder ships an order.
func (s *OrderService) ShipOrder(id string) error {
	order, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	err = order.Ship()
	if err != nil {
		return err
	}

	return s.repo.Update(order)
}

// DeliverOrder delivers an order.
func (s *OrderService) DeliverOrder(id string) error {
	order, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	err = order.Deliver()
	if err != nil {
		return err
	}

	return s.repo.Update(order)
}

// CancelOrder cancels an order.
func (s *OrderService) CancelOrder(id string) error {
	order, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	err = order.Cancel()
	if err != nil {
		return err
	}

	return s.repo.Update(order)
}
