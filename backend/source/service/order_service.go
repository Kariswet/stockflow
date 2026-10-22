package service

import (
	"errors"
	"log"
	"sort"
	"stockflow/config"
	"stockflow/dto"
	"stockflow/model"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OrderService struct {
	db *gorm.DB
}

func NewOrderService() *OrderService {
	pgInstance := config.GetPosgresInstance()
	db := pgInstance.GetDB()

	db.AutoMigrate(&model.SalesOrder{})
	db.AutoMigrate(&model.SalesOrderItem{})
	return &OrderService{
		db: db,
	}
}

func (o *OrderService) Create(param dto.OrderCreateRequest) model.Response {
	skuMap := make(map[string]bool)
	var skus []string

	for _, item := range param.Items {
		if skuMap[item.SKU] {
			return model.Response{
				Metadata: model.MetadataResponse{
					Message: "duplicate sku in request",
				},
			}
		}
		skuMap[item.SKU] = true
		skus = append(skus, item.SKU)
	}

	sort.Strings(skus)
	var resp model.Response

	err := o.db.Transaction(func(tx *gorm.DB) error {
		var products []model.Product
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("sku IN ?", skus).
			Order("sku ASC").Find(&products).Error; err != nil {
			return err
		}
		if len(products) != len(skus) {
			resp.Metadata.Message = "unknown sku"
			return errors.New("unknown sku")
		}

		dbProductMAp := make(map[string]*model.Product)
		for i := range products {
			dbProductMAp[products[i].SKU] = &products[i]
		}

		var totalOrder int64 = 0
		var orderItmes []model.SalesOrderItem
		var movements []model.StockMovement

		for _, paramItem := range param.Items {
			prod := dbProductMAp[paramItem.SKU]

			if prod.Stock < paramItem.Qty {
				resp.Metadata.Message = "insufficient stock"
				resp.Data = map[string]interface{}{
					"error": map[string]interface{}{
						"code":      "INSUFFICIENT_STOCK",
						"sku":       paramItem.SKU,
						"available": prod.Stock,
					},
				}
				return errors.New("insuffcicient stock")
			}
			totalOrder += *prod.Price * paramItem.Qty
			prod.Stock -= paramItem.Qty

			orderItmes = append(orderItmes, model.SalesOrderItem{
				SKU:       paramItem.SKU,
				Qty:       paramItem.Qty,
				UnitPrice: *prod.Price,
			})

			note := "sales order"
			movements = append(movements, model.StockMovement{
				ProductID: prod.ID,
				Type:      "out",
				Qty:       paramItem.Qty,
				Note:      &note,
			})
		}

		order := model.SalesOrder{
			Customer: param.Customer,
			Status:   "craeted",
			Total:    totalOrder,
		}
		if err := tx.Create(&order).Error; err != nil {
			return err
		}

		for i := range orderItmes {
			orderItmes[i].OrderID = order.ID
		}
		if err := tx.Create(&orderItmes).Error; err != nil {
			return err
		}
		if err := tx.Create(&movements).Error; err != nil {
			return err
		}

		for _, p := range products {
			if err := tx.Model(p).Update("stock", p.Stock).Error; err != nil {
				return err
			}
		}

		order.Items = orderItmes
		resp.Metadata.Message = "OK"
		resp.Data = order

		return nil
	})

	if err != nil && resp.Metadata.Message == "" {
		resp.Metadata.Message = "failed create order"
	}

	return resp
}

func (o *OrderService) Cancel(id int64) model.Response {
	var order model.SalesOrder

	err := o.db.Preload("Items").First(&order, id).Error
	if err != nil {
		log.Println(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Response{
				Metadata: model.MetadataResponse{
					Message: "order not found",
				},
			}
		}
		return model.Response{
			Metadata: model.MetadataResponse{
				Message: "failed fetch order data",
			},
		}
	}

	if order.Status == "CANCELLED" {
		return model.Response{
			Metadata: model.MetadataResponse{
				Message: "order already cancelled",
			},
			Data: map[string]interface{}{
				"error": map[string]interface{}{
					"code": "ORDER_ALREADY_CANCELLED",
				},
			},
		}
	}

	var skus []string
	for _, item := range order.Items {
		skus = append(skus, item.SKU)
	}
	sort.Strings(skus)

	err = o.db.Transaction(func(tx *gorm.DB) error {
		var products []model.Product
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("sku IN ?", skus).Order("sku ASC").
			Find(&products).Error; err != nil {
			return err
		}

		productMap := make(map[string]*model.Product)
		for i := range products {
			productMap[products[i].SKU] = &products[i]
		}

		var movements []model.StockMovement

		for _, item := range order.Items {
			prod := productMap[item.SKU]
			if prod != nil {
				prod.Stock += item.Qty
				note := "order cancelled"
				movements = append(movements, model.StockMovement{
					ProductID: prod.ID,
					Type:      "in",
					Qty:       item.Qty,
					Note:      &note,
				})
			}
		}

		for _, p := range products {
			if err := tx.Model(p).Update("stock", p.Stock).Error; err != nil {
				return err
			}
		}

		if err := tx.Create(&movements).Error; err != nil {
			return err
		}

		if err := tx.Model(&order).Updates(map[string]interface{}{
			"status":       "CANCELLED",
			"cancelled_at": time.Now(),
		}).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		log.Println("Transaction failed:", err)
		return model.Response{
			Metadata: model.MetadataResponse{
				Message: "failed to cancel 0order",
			},
		}
	}

	return model.Response{
		Metadata: model.MetadataResponse{
			Message: "OK",
		},
	}
}

func (o *OrderService) GetDetail(id int64) model.Response {
	var order model.SalesOrder

	err := o.db.Preload("Items").First(&order, id).Error
	if err != nil {
		log.Println(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Response{
				Metadata: model.MetadataResponse{
					Message: "order not found",
				},
			}
		}
		return model.Response{
			Metadata: model.MetadataResponse{
				Message: "failed fetch detail order",
			},
		}
	}

	return model.Response{
		Data: order,
	}
}
