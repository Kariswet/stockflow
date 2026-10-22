package config

import (
	"log"
	"os"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresManager struct {
	db *gorm.DB
}

var (
	pgInstance *PostgresManager
	pgOnce     sync.Once
)

func GetPosgresInstance() *PostgresManager {
	pgOnce.Do(func() {
		pgInstance = &PostgresManager{}
		pgInstance.initPGConnection()
	})
	return pgInstance
}

func (p *PostgresManager) initPGConnection() {
	uri := os.Getenv("PG_URI")
	if uri == "" {
		log.Println("PG: missing env PG_URI")
		return
	}

	db, err := gorm.Open(postgres.Open(uri), &gorm.Config{})
	if err != nil {
		log.Printf("PG: conneciton failed: %v", err)
		return
	}

	p.db = db
	log.Println("PG: connected")
}

func (p *PostgresManager) GetDB() *gorm.DB {
	return p.db
}
