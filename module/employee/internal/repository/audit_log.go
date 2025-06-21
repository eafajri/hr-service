package repository

import (
	"encoding/json"

	"github.com/eafajri/hr-service.git/module/employee/internal/entity"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type AuditLogRepositoryImpl struct {
	DB *gorm.DB
}

func NewAuditLogRepository(db *gorm.DB) *AuditLogRepositoryImpl {
	return &AuditLogRepositoryImpl{
		DB: db,
	}
}

func (r *AuditLogRepositoryImpl) Create(log entity.AuditLog, payload any) error {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Payload = datatypes.JSON([]byte("{}"))
	} else {
		log.Payload = datatypes.JSON(payloadBytes)
	}

	return r.DB.Create(&log).Error
}
