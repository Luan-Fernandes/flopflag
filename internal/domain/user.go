package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Email        string    `gorm:"not null;uniqueIndex"`
	PasswordHash string    `gorm:"column:password_hash;not null"`
	Role         string    `gorm:"not null"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
}

func (User) TableName() string { return "users" }
