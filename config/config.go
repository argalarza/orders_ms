package config

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func LoadEnv() {
	_ = godotenv.Load()
}

func ConnectPostgres() *sql.DB {
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	dbname := os.Getenv("POSTGRES_DB")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("❌ Error al conectar a PostgreSQL:", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("❌ PostgreSQL no responde:", err)
	}

	log.Println("✅ Conectado a PostgreSQL")

	// Crear tabla orders si no existe
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS orders (
		id SERIAL PRIMARY KEY,
		user_email VARCHAR(255),
		total NUMERIC,
		status VARCHAR(50),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`)
	if err != nil {
		log.Fatal("❌ Error creando tabla orders:", err)
	}

	// Crear tabla order_items si no existe
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS order_items (
		id SERIAL PRIMARY KEY,
		order_id INT REFERENCES orders(id) ON DELETE CASCADE,
		product_id VARCHAR(100),
		quantity INT,
		price NUMERIC
	);`)
	if err != nil {
		log.Fatal("❌ Error creando tabla order_items:", err)
	}

	log.Println("✅ Tablas creadas o existentes")
	return db
}

func ConnectMongo() *mongo.Client {
	uri := os.Getenv("MONGO_URI")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("❌ Error al conectar a MongoDB:", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("❌ MongoDB no responde:", err)
	}

	log.Println("✅ Conectado a MongoDB")
	return client
}
