package goat

import (
	"time"

	"github.com/google/uuid"
)

type Model struct {
	ID        uuid.UUID  `json:"id" gorm:"primary_key"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
