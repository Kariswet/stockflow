package service

// import (
// 	"stockflow/config"
// 	"stockflow/model"

// 	"gorm.io/gorm"
// )

// type StockService struct {
// 	db *gorm.DB
// }

// func NewStockService() *StockService {
// 	pgInstance := config.GetPosgresInstance()
// 	db := pgInstance.GetDB()

// 	db.AutoMigrate(&model.StockMovement{})
// 	return &StockService{
// 		db: db,
// 	}
// }
