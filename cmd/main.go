package main

import (
	"create-orders-service/controllers"
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	router := gin.Default()

	// 🔓 CORS completamente habilitado incluyendo Authorization
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Puedes cambiar esto a ["http://54.175.97.19"] si quieres restringir
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 📦 Configurar conexión a PostgreSQL
	db, err := sql.Open("postgres", "host="+os.Getenv("DB_HOST")+" port="+os.Getenv("DB_PORT")+" user="+os.Getenv("DB_USER")+" password="+os.Getenv("DB_PASSWORD")+" dbname="+os.Getenv("DB_NAME")+" sslmode=disable")
	if err != nil {
		log.Fatal("❌ Error al conectar a PostgreSQL:", err)
	}

	// 📦 Conexión a MongoDB
	mongoURI := os.Getenv("PRODUCT_MS_URL") // Puedes cambiar si tienes un URI real de MongoDB
	mongoClient, err := mongo.Connect(nil, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal("❌ Error al conectar a MongoDB:", err)
	}

	// 📂 Rutas
	api := router.Group("/orders")
	{
		api.POST("", controllers.CreateOrder(db, mongoClient))
		api.GET("/:id", controllers.GetOrderByID(db))
		api.PUT("/:id", controllers.UpdateOrder(db))
		api.DELETE("/:id", controllers.DeleteOrder(db))
	}

	log.Println("🚀 Microservicio de órdenes corriendo en puerto 5001")
	router.Run(":5001")
}
