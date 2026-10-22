package dto

type ProductRequest struct {
	Name  string `json:"name" binding:"required"`
	SKU   string `json:"sku" binding:"required"`
	Price *int64 `json:"price" binding:"required,gte=0"`
}

type ProductQueryRequest struct {
	Search string `form:"search"`
	Page   int    `form:"page" binding:"omitempty,min=1"`
	Limit  int    `form:"limit" binding:"omitempty,min=1"`
}

type ProductUpdateRequest struct {
	Name     *string `json:"name"`
	SKU      *string `json:"sku"`
	Price    *int64  `json:"price" binding:"omitempty,gte=0"`
	IsActive *bool   `json:"is_active"`
}
