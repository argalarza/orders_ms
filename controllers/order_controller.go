package controllers

import (
	"create-orders-service/models"
	"create-orders-service/services"
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// Crear nueva orden
func CreateOrder(db *sql.DB, mongo *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var order models.Order
		if err := c.ShouldBindJSON(&order); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Formato inválido"})
			return
		}

		email, exists := c.Get("email")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token sin email"})
			return
		}

		order.Email = email.(string)
		order.Status = "CREATED"

		log.Println("⏳ Calculando total...")
		total, err := services.CalculateTotal(order.Items, mongo)
		if err != nil {
			log.Println("❌ Error en CalculateTotal:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		order.Total = total
		log.Println("✅ Total calculado:", total)

		log.Println("💾 Insertando orden...")
		id, err := services.InsertOrder(db, order)
		if err != nil {
			log.Println("❌ Error al insertar orden:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo crear la orden"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"orderId": id})
	}
}

// Obtener orden por ID
func GetOrderByID(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var order models.Order
		row := db.QueryRow("SELECT id, user_email, total, status FROM orders WHERE id=$1", id)
		if err := row.Scan(&order.ID, &order.Email, &order.Total, &order.Status); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Orden no encontrada"})
			return
		}

		rows, err := db.Query("SELECT product_id, quantity, price FROM order_items WHERE order_id=$1", id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al consultar productos"})
			return
		}
		defer rows.Close()

		for rows.Next() {
			var item models.OrderItem
			if err := rows.Scan(&item.ProductID, &item.Quantity, &item.Price); err != nil {
				log.Println("❌ Error al escanear item:", err)
				continue
			}
			order.Items = append(order.Items, item)
		}

		c.JSON(http.StatusOK, order)
	}
}

// Actualizar estado de la orden
func UpdateOrder(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var body struct {
			Status string `json:"status"`
		}

		if err := c.ShouldBindJSON(&body); err != nil || body.Status == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Estado inválido"})
			return
		}

		_, err := db.Exec("UPDATE orders SET status=$1 WHERE id=$2", body.Status, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo actualizar"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Orden actualizada"})
	}
}

// Eliminar (cancelar) orden
func DeleteOrder(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		_, err := db.Exec("DELETE FROM orders WHERE id=$1", id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo eliminar"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Orden cancelada"})
	}
}
