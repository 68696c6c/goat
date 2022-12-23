package goat

import (
	"time"

	"gorm.io/gorm"
)

type Model struct {
	ID        ID         `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

func (m Model) BeforeCreate(tx *gorm.DB) error {
	if m.ID == NilID() {
		m.ID = NewID()
	}
	return nil
}

type ModelHardDelete struct {
	ID        ID         `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

func (m ModelHardDelete) BeforeCreate(tx *gorm.DB) error {
	if m.ID == NilID() {
		m.ID = NewID()
	}
	return nil
}
