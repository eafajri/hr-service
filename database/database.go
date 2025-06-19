package database

import (
	"fmt"
	"log"
	"sync"

	"github.com/eafajri/hr-service.git/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	dbInstance *gorm.DB
	once       sync.Once
)

// GetDB returns a singleton GORM DB instance
func GetDB() *gorm.DB {
	once.Do(func() {
		cfg := config.GetConfig()

		dsn := fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
		)

		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}

		dbInstance = db
	})

	return dbInstance
}
