package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"create-orders-service/models"
)

var PostgresDB *gorm.DB

func ConnectPostgres() {
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Error al conectar a PostgreSQL: %v", err)
	}

	// Migración automática de la tabla 'orders'
	if err := db.AutoMigrate(&models.Order{}); err != nil {
		log.Fatalf("❌ Error en AutoMigrate para Order: %v", err)
	}

	PostgresDB = db
	log.Println("✅ Conectado a PostgreSQL")
}

func GetPostgresDB() *gorm.DB {
	return PostgresDB
}
