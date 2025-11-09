// Package mongo provides MongoDB implementation of the order repository.
package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/yanking/gomicro/examples/order-service/internal/model"
	"github.com/yanking/gomicro/examples/order-service/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// MongoRepository implements OrderRepository interface for MongoDB.
type MongoRepository struct {
	collection *mongo.Collection
	timeout    time.Duration
}

// NewMongoRepository creates a new MongoDB repository instance.
func NewMongoRepository(collection *mongo.Collection) repository.OrderRepository {
	return &MongoRepository{
		collection: collection,
		timeout:    10 * time.Second,
	}
}

// Save saves an order to the database.
func (r *MongoRepository) Save(order *model.Order) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	// Use MongoDB's upsert functionality
	filter := bson.M{"id": order.ID}
	update := bson.M{
		"$set": bson.M{
			"id":         order.ID,
			"user_id":    order.UserID,
			"items":      order.Items,
			"status":     order.Status,
			"total":      order.Total,
			"created_at": order.CreatedAt,
			"updated_at": order.UpdatedAt,
		},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update, &mongo.UpdateOptions{
		Upsert: &[]bool{true}[0],
	})
	if err != nil {
		return fmt.Errorf("failed to save order: %w", err)
	}

	return nil
}

// FindByID finds an order by ID.
func (r *MongoRepository) FindByID(id string) (*model.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	var order model.Order
	filter := bson.M{"id": id}

	err := r.collection.FindOne(ctx, filter).Decode(&order)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("order not found")
		}
		return nil, fmt.Errorf("failed to find order: %w", err)
	}

	return &order, nil
}

// FindByUserID finds orders by user ID.
func (r *MongoRepository) FindByUserID(userID string) ([]*model.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	filter := bson.M{"user_id": userID}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to query orders: %w", err)
	}
	defer cursor.Close(ctx)

	var orders []*model.Order
	for cursor.Next(ctx) {
		var order model.Order
		if err := cursor.Decode(&order); err != nil {
			return nil, fmt.Errorf("failed to decode order: %w", err)
		}
		orders = append(orders, &order)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return orders, nil
}

// Update updates an order.
func (r *MongoRepository) Update(order *model.Order) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	filter := bson.M{"id": order.ID}
	update := bson.M{
		"$set": bson.M{
			"user_id":    order.UserID,
			"items":      order.Items,
			"status":     order.Status,
			"total":      order.Total,
			"updated_at": order.UpdatedAt,
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("order not found")
	}

	return nil
}

// Delete deletes an order by ID.
func (r *MongoRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	filter := bson.M{"id": id}
	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete order: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("order not found")
	}

	return nil
}
