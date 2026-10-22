package model

import "time"

type SalesOrder struct {
	ID          int64      `json:"id" db:"id"`
	Customer    string     `json:"customer" db:"customer"`
	Status      string     `json:"status" db:"status"`
	Total       int64      `json:"total" db:"total"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	CancelledAt *time.Time `json:"cancelled_at" db:"cancelled_at"`

	Items []SalesOrderItem `json:"items,omitempty" gorm:"foreignKey:OrderID"`
}

type SalesOrderItem struct {
	ID        int64  `json:"id" db:"id"`
	OrderID   int64  `json:"order_id" db:"order_id"`
	SKU       string `json:"sku" db:"sku"`
	Qty       int64  `json:"qty" db:"qty"`
	UnitPrice int64  `json:"unit_price" db:"unit_price"`

	Order SalesOrder `json:"-" gorm:"foreignKey:OrderID"`
}
