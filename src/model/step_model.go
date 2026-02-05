package model

import (
	"time"

	"gorm.io/datatypes"
)

type Step struct {
	ID         uint           `gorm:"primaryKey;autoIncrement"`
	WorkflowID uint           `gorm:"not null"`
	Level      uint           `gorm:"not null"`
	Actor      string         `gorm:"not null"`
	Conditions datatypes.JSON `gorm:"type:json"`
	CreatedAt  time.Time      `gorm:"autoCreateTime:milli"`
}
