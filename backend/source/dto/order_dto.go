package dto

type OrderItemRequest struct {
	SKU string `json:"sku" binding:"required"`
	Qty int64  `json:"qty" binding:"required,gt=0"`
}

type OrderCreateRequest struct {
	Customer string             `json:"customer" binding:"required"`
	Items    []OrderItemRequest `json:"items" binding:"required,min=1,dive"`
}
