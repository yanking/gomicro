// Package mysql provides MySQL implementation of the order repository.
package mysql

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/yanking/gomicro/examples/order-service/internal/model"
	"github.com/yanking/gomicro/examples/order-service/internal/repository"
)

// MySQLRepository implements OrderRepository interface for MySQL database.
type MySQLRepository struct {
	db *sql.DB
}

// NewMySQLRepository creates a new MySQL repository instance.
func NewMySQLRepository(db *sql.DB) repository.OrderRepository {
	return &MySQLRepository{
		db: db,
	}
}

// Save saves an order to the database.
func (r *MySQLRepository) Save(order *model.Order) error {
	// Convert items to JSON
	itemsJSON, err := json.Marshal(order.Items)
	if err != nil {
		return fmt.Errorf("failed to marshal order items: %w", err)
	}

	// Insert or update the order
	query := `
		INSERT INTO orders (id, user_id, items, status, total, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			user_id = VALUES(user_id),
			items = VALUES(items),
			status = VALUES(status),
			total = VALUES(total),
			updated_at = VALUES(updated_at)
	`

	_, err = r.db.Exec(query, order.ID, order.UserID, itemsJSON, order.Status, order.Total, order.CreatedAt, order.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to save order: %w", err)
	}

	return nil
}

// FindByID finds an order by ID.
func (r *MySQLRepository) FindByID(id string) (*model.Order, error) {
	query := `
		SELECT id, user_id, items, status, total, created_at, updated_at
		FROM orders
		WHERE id = ?
	`

	var order model.Order
	var itemsJSON []byte

	err := r.db.QueryRow(query, id).Scan(
		&order.ID,
		&order.UserID,
		&itemsJSON,
		&order.Status,
		&order.Total,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("order not found")
		}
		return nil, fmt.Errorf("failed to find order: %w", err)
	}

	// Unmarshal items from JSON
	if err := json.Unmarshal(itemsJSON, &order.Items); err != nil {
		return nil, fmt.Errorf("failed to unmarshal order items: %w", err)
	}

	return &order, nil
}

// FindByUserID finds orders by user ID.
func (r *MySQLRepository) FindByUserID(userID string) ([]*model.Order, error) {
	query := `
		SELECT id, user_id, items, status, total, created_at, updated_at
		FROM orders
		WHERE user_id = ?
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query orders: %w", err)
	}
	defer rows.Close()

	var orders []*model.Order
	for rows.Next() {
		var order model.Order
		var itemsJSON []byte

		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&itemsJSON,
			&order.Status,
			&order.Total,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}

		// Unmarshal items from JSON
		if err := json.Unmarshal(itemsJSON, &order.Items); err != nil {
			return nil, fmt.Errorf("failed to unmarshal order items: %w", err)
		}

		orders = append(orders, &order)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return orders, nil
}

// Update updates an order.
func (r *MySQLRepository) Update(order *model.Order) error {
	// Convert items to JSON
	itemsJSON, err := json.Marshal(order.Items)
	if err != nil {
		return fmt.Errorf("failed to marshal order items: %w", err)
	}

	query := `
		UPDATE orders
		SET user_id = ?, items = ?, status = ?, total = ?, updated_at = ?
		WHERE id = ?
	`

	result, err := r.db.Exec(query, order.UserID, itemsJSON, order.Status, order.Total, order.UpdatedAt, order.ID)
	if err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("order not found")
	}

	return nil
}

// Delete deletes an order by ID.
func (r *MySQLRepository) Delete(id string) error {
	query := `DELETE FROM orders WHERE id = ?`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete order: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("order not found")
	}

	return nil
}
