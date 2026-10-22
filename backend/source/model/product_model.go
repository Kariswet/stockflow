package model

import "time"

type Product struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	SKU       string    `json:"sku" gorm:"uniqueIndex;not null"`
	Name      string    `json:"name" gorm:"not null"`
	Price     *int64    `json:"price" gorm:"not null"`
	Stock     int64     `json:"stock" gorm:"not null;default:0"`
	IsActive  bool      `json:"is_active" gorm:"not null;default:true"`
	CreatedAt time.Time `json:"created_at" gorm:"not null;default:current_timestamp"`

	LastMovementAt *time.Time `json:"last_movement_at" gorm:"-"`
}
