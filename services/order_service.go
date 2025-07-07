package services

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"create-orders-service/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Product struct {
	ID    primitive.ObjectID `bson:"_id"`
	Name  string             `bson:"name"`
	Price float64            `bson:"price"`
}

func CalculateTotal(items []models.OrderItem, mongo *mongo.Client) (float64, error) {
	var total float64
	collection := mongo.Database("products_db").Collection("products")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for i, item := range items {
		var product Product

		objectID, err := primitive.ObjectIDFromHex(item.ProductID)
		if err != nil {
			return 0, fmt.Errorf("ID de producto inválido: %s", item.ProductID)
		}

		err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&product)
		if err != nil {
			return 0, fmt.Errorf("producto no encontrado: %s", item.ProductID)
		}

		items[i].Price = product.Price
		total += product.Price * float64(item.Quantity)
	}

	return total, nil
}

func InsertOrder(db *sql.DB, order models.Order) (int, error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	var orderID int
	err = tx.QueryRow("INSERT INTO orders (user_email, total, status) VALUES ($1, $2, $3) RETURNING id",
		order.Email, order.Total, order.Status).Scan(&orderID)

	if err != nil {
		return 0, err
	}

	for _, item := range order.Items {
		_, err = tx.Exec("INSERT INTO order_items (order_id, product_id, quantity, price) VALUES ($1, $2, $3, $4)",
			orderID, item.ProductID, item.Quantity, item.Price)
		if err != nil {
			return 0, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return orderID, nil
}
