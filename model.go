package goat

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Model struct {
	ID        ID         `json:"id" gorm:"primary_key"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

func (m Model) BeforeCreate(scope *gorm.Scope) error {
	id := NewID()
	scope.SetColumn("ID", id)
	return nil
}
