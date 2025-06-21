package entity

import (
	"time"

	"gorm.io/datatypes"
)

type AuditLog struct {
	ID        int64          `gorm:"id" json:"id"`
	RequestID string         `gorm:"request_id" json:"request_id"`
	IPAddress string         `gorm:"ip_address" json:"ip_address"`
	TableName string         `gorm:"table_name" json:"table_name"`
	Action    string         `gorm:"action" json:"action"`
	Target    string         `gorm:"target" json:"target"`
	Payload   datatypes.JSON `gorm:"type:jsonb" json:"payload"`
	CreatedBy string         `gorm:"created_by" json:"created_by"`
	CreatedAt time.Time      `gorm:"created_at" json:"created_at"`
}
