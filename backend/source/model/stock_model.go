package model

import "time"

type StockMovement struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	ProductID int64     `json:"product_id" gorm:"not null"`
	Type      string    `json:"type" gorm:"not null"`
	Qty       int64     `json:"qty" gorm:"not null"`
	Note      *string   `json:"note"`
	CreatedAt time.Time `json:"created_at" gorm:"not null;default:current_timestamp"`

	Product Product `json:"-" gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
