package model

import (
	"time"
)

type Request struct {
	ID          uint      `gorm:"primaryKey;autoIncrement"`
	WorkflowID  uint      `gorm:"not null"`
	CurrentStep uint      `gorm:"not null"`
	Status      string    `gorm:"not null"` // "pending", "approved", "rejected"
	Amount      float64   `gorm:"not null"`
	CreatedAt   time.Time `gorm:"autoCreateTime:milli"`
}
