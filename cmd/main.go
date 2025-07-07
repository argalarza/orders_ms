package main

import (
	"create-orders-service/config"
	"create-orders-service/controllers"
	"create-orders-service/middleware"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()
	db := config.ConnectPostgres()
	mongo := config.ConnectMongo()
	router := gin.Default()

	router.Use(cors.Default())
	router.Use(middleware.JWTMiddleware())

	router.POST("/orders", controllers.CreateOrder(db, mongo))
	router.GET("/orders/:id", controllers.GetOrderByID(db))
	router.PUT("/orders/:id", controllers.UpdateOrder(db))
	router.DELETE("/orders/:id", controllers.DeleteOrder(db))

	port := os.Getenv("PORT")
	if port == "" {
		port = "5001"
	}

	log.Fatal(router.Run(":" + port))
}
