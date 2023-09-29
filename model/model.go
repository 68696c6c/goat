package model

import (
	"time"

	"gorm.io/gorm"

	"github.com/68696c6c/goat"
)

type SoftDelete struct {
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

type Timestamps struct {
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" gorm:"<-:update"`
}

func (m *Timestamps) BeforeUpdate(tx *gorm.DB) error {
	tx.Statement.SetColumn("UpdatedAt", time.Now())
	return nil
}

type Model struct {
	ID goat.ID `json:"id" gorm:"primaryKey"`
}

func NewModel() *Model {
	return &Model{}
}

func (m *Model) BeforeCreate(tx *gorm.DB) error {
	if m.ID == goat.NilID() {
		tx.Statement.SetColumn("ID", goat.NewID())
	}
	return nil
}
