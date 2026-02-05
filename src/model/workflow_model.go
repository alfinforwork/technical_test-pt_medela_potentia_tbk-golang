package model

import (
	"time"
)

type Workflow struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`     // id
	Name      string    `gorm:"not null;unique" json:"name"`            // name
	CreatedAt time.Time `gorm:"autoCreateTime:milli" json:"created_at"` // created_at
}
