package service

import (
	"errors"
	"log"
	"stockflow/config"
	"stockflow/dto"
	"stockflow/model"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProductService struct {
	db *gorm.DB
}

func NewProductService() *ProductService {
	pgInstance := config.GetPosgresInstance()
	db := pgInstance.GetDB()

	db.AutoMigrate(&model.Product{})
	db.AutoMigrate(&model.StockMovement{})
	return &ProductService{
		db: db,
	}
}

func (p *ProductService) Create(param dto.ProductRequest) model.Response {
	var product model.Product

	product.Name = param.Name
	product.SKU = param.SKU
	product.Price = param.Price
	product.IsActive = true
	product.Stock = 0

	err := p.db.Create(&product).Error
	if err != nil {
		log.Println(err)

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return model.Response{
				Metadata: model.MetadataResponse{
					Message: "SKU already exist",
				},
				Data: map[string]interface{}{
					"error": map[string]interface{}{
						"code": "DUPLICATE_SKU",
					},
				},
			}
		}

		return model.Response{
			Metadata: model.MetadataResponse{
				Message: "failed to create product",
			},
			Data: nil,
		}
	}

	var resp model.Response
	resp.Data = product

	return resp
}

func (p *ProductService) Search(param dto.ProductQueryRequest) model.Response {
	page := param.Page
	if page <= 0 {
		page = 1
	}
	limit := param.Limit
	if limit <= 0 {
		limit = 10
	}

	var product model.Product

	query := p.db.Model(&product)

	if param.Search != "" {
		searchTerm := "%" + param.Search + "%"
		query = query.Where("sku ILIKE ? OR name ILIKE ?", searchTerm, searchTerm)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return model.Response{
			Metadata: model.MetadataResponse{
				Message: "failed to count product",
			},
		}
	}

	var products []model.Product
	offset := (page - 1) * limit

	subQuery := "(SELECT created_at FROM stock_movements WHERE product_id = products.id ORDER BY created_at DESC LIMIT 1) as last_movement_at"

	err := query.Select("products.*, " + subQuery).Order("products.created_at DESC").Offset(offset).Limit(limit).Find(&products).Error
	if err != nil {
		return model.Response{
			Metadata: model.MetadataResponse{
				Message: "failed to fetch product",
			},
		}
	}

	return model.Response{
		Data: model.PaginatedResponse{
			Data:  products,
			Page:  page,
			Limit: limit,
			Total: total,
		},
	}
}

func (p *ProductService) Delete(id int64) model.Response {
	var product model.Product
	err := p.db.Delete(&product, id).Error
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return model.Response{
				Metadata: model.MetadataResponse{
					Message: "product has movement hoistory",
				},
				Data: map[string]interface{}{
					"code": map[string]interface{}{
						"error": "PRODUCT_HAS_HISTORY",
					},
				},
			}
		}

		return model.Response{
			Metadata: model.MetadataResponse{
				Message: "failed to delte product",
			},
		}
	}

	return model.Response{
		Metadata: model.MetadataResponse{
			Message: "OK",
		},
	}
}

func (p *ProductService) Update(id int64, param dto.ProductUpdateRequest) model.Response {
	var product model.Product

	err := p.db.First(&product, id).Error
	if err != nil {
		log.Println(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Response{
				Metadata: model.MetadataResponse{
					Message: "product not found",
				},
			}
		}
		return model.Response{
			Metadata: model.MetadataResponse{
				Message: "error fetching product",
			},
		}
	}

	updates := make(map[string]interface{})

	if param.Name != nil {
		updates["name"] = *param.Name
	}
	if param.SKU != nil {
		updates["sku"] = *param.SKU
	}
	if param.Price != nil {
		updates["price"] = *param.Price
	}
	if param.IsActive != nil {
		updates["is_active"] = *param.IsActive
	}

	if len(updates) == 0 {
		return model.Response{
			Metadata: model.MetadataResponse{
				Message: "OK",
			},
			Data: product,
		}
	}

	err = p.db.Model(&product).Updates(updates).Error
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return model.Response{
				Metadata: model.MetadataResponse{
					Message: "SKU already exists",
				},
				Data: map[string]interface{}{
					"code": map[string]interface{}{
						"error": "DUPLICATE_SKU",
					},
				},
			}
		}
		return model.Response{
			Metadata: model.MetadataResponse{
				Message: "product: failed to update",
			},
		}
	}
	return model.Response{
		Metadata: model.MetadataResponse{
			Message: "OK",
		},
		Data: product,
	}
}

func (p *ProductService) AdjustStock(id int64, param dto.StockMovementRequest) model.Response {
	var resp model.Response

	err := p.db.Transaction(func(tx *gorm.DB) error {
		var product model.Product

		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&product, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				resp.Metadata.Message = "product not found"
				return err
			}
			return err
		}
		if param.Type == "out" {
			if product.Stock < param.Qty {
				resp.Metadata.Message = "insufficient stock"
				resp.Data = map[string]interface{}{
					"error": map[string]interface{}{
						"code":      "INSUFFICIENT_STOCK",
						"available": product.Stock,
					},
				}
				return errors.New("insufficient stock")
			}
			product.Stock -= param.Qty
		} else {
			product.Stock += param.Qty
		}

		if err := tx.Model(&product).Update("stock", product.Stock).Error; err != nil {
			return err
		}

		movement := model.StockMovement{
			ProductID: id,
			Type:      param.Type,
			Qty:       param.Qty,
			Note:      param.Note,
		}
		if err := tx.Create(&movement).Error; err != nil {
			return err
		}
		resp.Metadata.Message = "OK"
		resp.Data = movement

		return nil
	})

	if err != nil && resp.Metadata.Message == "" {
		resp.Metadata.Message = "failed to adjust stock"
	}
	return resp
}

func (p *ProductService) Movement(id int64) model.Response {
	var movements []model.StockMovement

	err := p.db.Where("product_id = ?", id).Order("created_at DESC").Find(&movements).Error
	if err != nil {
		log.Println(err)
		return model.Response{
			Metadata: model.MetadataResponse{
				Message: "failed to fetch movement history",
			},
		}
	}

	if len(movements) == 0 {
		return model.Response{
			Metadata: model.MetadataResponse{
				Message: "no record found",
			},
			Data: movements,
		}
	}

	return model.Response{
		Metadata: model.MetadataResponse{
			Message: "OK",
		},
		Data: movements,
	}
}
