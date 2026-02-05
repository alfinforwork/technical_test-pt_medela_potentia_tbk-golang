package model

import (
	"time"

	"gorm.io/datatypes"
)

type Step struct {
	ID         uint           `gorm:"primaryKey;autoIncrement" json:"id"`     // id
	WorkflowID uint           `gorm:"not null" json:"workflow_id"`            // workflow_id
	Level      uint           `gorm:"not null" json:"level"`                  // level
	Actor      string         `gorm:"not null" json:"actor"`                  // actor
	Conditions datatypes.JSON `gorm:"type:json" json:"conditions"`            // conditions
	CreatedAt  time.Time      `gorm:"autoCreateTime:milli" json:"created_at"` // created_at
}
