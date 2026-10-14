package domain

import (
	"time"

	"github.com/google/uuid"
)

type Todo struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID    uuid.UUID `json:"userId" gorm:"type:uuid;not null;index;column:user_id"`
	Title     string    `json:"title" gorm:"not null"`
	Done      bool      `json:"done" gorm:"not null;default:false"`
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}

func (Todo) TableName() string { return "todos" }
