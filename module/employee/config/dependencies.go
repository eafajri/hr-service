package transport

import (
	"log"

	"github.com/eafajri/hr-service.git/database"
	"github.com/eafajri/hr-service.git/internal/cache"

	"gorm.io/gorm"
)

type ModuleDependencies struct {
	MemoryCache cache.MemoryCache
	Database    gorm.DB
}

func NewModuleDependencies() *ModuleDependencies {
	db := database.GetDB()

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get generic DB: ", err)
	}

	err = sqlDB.Ping()
	if err != nil {
		log.Fatal("Failed to ping DB: ", err)
	}

	return &ModuleDependencies{
		MemoryCache: cache.NewMemoryCache(),
		Database:    *db,
	}
}
