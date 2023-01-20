package models

import (
	"encoding/json"

	"github.com/68696c6c/goat"
	"github.com/68696c6c/goat/hal"
	"github.com/68696c6c/goat/model"

	"github.com/68696c6c/example/app/enums"
)

const UserLinkKey = "users"

type User struct {
	*model.Model

	OrganizationId goat.ID         `json:"organizationId" binding:"required"`
	Level          enums.UserLevel `json:"level" binding:"required"`
	Name           string          `json:"name" binding:"required"`
	Email          string          `json:"email" binding:"required,email"`

	*model.Timestamps
	*model.SoftDelete
	*hal.ResourceEmbeds[*UserEmbeds]
}

type UserRequest struct {
	OrganizationId *goat.ID         `json:"organizationId,omitempty"`
	Level          *enums.UserLevel `json:"level,omitempty"`
	Name           *string          `json:"name,omitempty"`
	Email          *string          `json:"email,omitempty" binding:"omitempty,email"`
}

type UserEmbeds struct {
	Organization *Organization `json:"organization,omitempty" gorm:"foreignKey:OrganizationId"`
}

func MakeUser() *User {
	return &User{
		Model:          model.NewModel(),
		OrganizationId: goat.NilID(),
		Level:          enums.UserLevelUser,
		Name:           "",
		Email:          "",
	}
}

func (m *User) MarshalJSON() ([]byte, error) {
	type Alias User
	return json.Marshal(&struct {
		*Alias
		*hal.ResourceLinks
	}{
		Alias:         (*Alias)(m),
		ResourceLinks: goat.NewResourceLinks(UserLinkKey, m.ID.String()),
	})
}
