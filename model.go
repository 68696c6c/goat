package goat

import (
	"time"

	"gorm.io/gorm"
)

type Model struct {
	ID        ID         `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" gorm:"-"`
	DeletedAt *time.Time `json:"deleted_at"`
}

func (m *Model) BeforeCreate(tx *gorm.DB) error {
	return setId(m.ID, tx)
}

func (m *Model) BeforeUpdate(tx *gorm.DB) error {
	return setUpdatedAt(tx)
}

type ModelHardDelete struct {
	ID        ID         `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" gorm:"-"`
}

func (m *ModelHardDelete) BeforeCreate(tx *gorm.DB) error {
	return setId(m.ID, tx)
}

func (m *ModelHardDelete) BeforeUpdate(tx *gorm.DB) error {
	return setUpdatedAt(tx)
}

func setId(id ID, tx *gorm.DB) error {
	if id == NilID() {
		tx.Statement.SetColumn("ID", NewID())
	}
	return nil
}

func setUpdatedAt(tx *gorm.DB) error {
	tx.Statement.SetColumn("UpdatedAt", time.Now())
	return nil
}
