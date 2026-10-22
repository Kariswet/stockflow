package dto

type StockMovementRequest struct {
	Type string  `json:"type" binding:"required,oneof=in out"`
	Qty  int64   `json:"qty" binding:"required,gt=0"`
	Note *string `json:"note"`
}
