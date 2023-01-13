package models

import (
	"encoding/json"

	"github.com/68696c6c/goat"
	"github.com/68696c6c/goat/model"
	"github.com/68696c6c/goat/resource"

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
	*resource.Embedded[*UserEmbeds]
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
		*resource.Links
	}{
		Alias: (*Alias)(m),
		Links: goat.MakeResourceLinks(UserLinkKey, m.ID.String()),
	})
}

//
// func (m *User) ValidateCreate() {
// 	// Email is required, must be a valid email, and must be unique.
// }
//
// func ValidateUser(sl validator.StructLevel) {
// 	user := sl.Current().Interface().(User)
//
// 	// If createdAt is nil, Email is required and must be a valid email
//
// 	// If createdAt is not nil, Email must not be present.
//
// }
