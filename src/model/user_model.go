package model

import "time"

type User struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`     // id
	Name         string    `gorm:"not null" json:"name"`                   // name
	Email        string    `gorm:"not null;unique" json:"email"`           // email
	PasswordHash string    `gorm:"not null" json:"-"`                      // password_hash
	CreatedAt    time.Time `gorm:"autoCreateTime:milli" json:"created_at"` // created_at
}
