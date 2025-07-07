package models

import (
	"time"

	"gorm.io/gorm"
)

type OrderItem struct {
	ProductID string  `json:"productId"` // No persistido
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

type Order struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Email     string         `gorm:"column:user_email" json:"email"` // <- ESTE CAMBIO ES CLAVE
	Total     float64        `json:"total"`
	Status    string         `json:"status"`
	CreatedAt time.Time      `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Items []OrderItem `gorm:"-" json:"items"`
}
