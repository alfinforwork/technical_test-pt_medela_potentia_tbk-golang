package model

import (
	"time"
)

type Workflow struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	Name      string    `gorm:"not null;unique"`
	CreatedAt time.Time `gorm:"autoCreateTime:milli"`
}
