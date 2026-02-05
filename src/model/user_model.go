package model

import "time"

type User struct {
	ID           uint      `gorm:"primaryKey;autoIncrement"`
	Name         string    `gorm:"not null"`
	Email        string    `gorm:"not null;unique"`
	PasswordHash string    `gorm:"not null"`
	CreatedAt    time.Time `gorm:"autoCreateTime:milli"`
}
