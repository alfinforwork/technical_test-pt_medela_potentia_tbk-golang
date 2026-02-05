package model

import (
	"time"
)

type Request struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`     // id
	WorkflowID  uint      `gorm:"not null" json:"workflow_id"`            // workflow_id
	CurrentStep uint      `gorm:"not null" json:"current_step"`           // current_step
	Status      string    `gorm:"not null" json:"status"`                 // status: "pending", "approved", "rejected"
	Amount      float64   `gorm:"not null" json:"amount"`                 // amount
	CreatedAt   time.Time `gorm:"autoCreateTime:milli" json:"created_at"` // created_at
}
