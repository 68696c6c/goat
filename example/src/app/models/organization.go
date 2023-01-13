package models

import (
	"encoding/json"

	"github.com/68696c6c/goat"
	"github.com/68696c6c/goat/model"
	"github.com/68696c6c/goat/resource"
)

const OrganizationLinkKey = "organizations"

type Organization struct {
	*model.Model

	// ParentId goat.ID `json:"parentId"`
	Name    string `json:"name" binding:"required"`
	Website string `json:"website" binding:"required"`

	*model.Timestamps
	*model.SoftDelete
	*resource.Embedded[*OrganizationEmbeds]
}

type OrganizationRequest struct {
	ParentId *goat.ID `json:"parentId,omitempty"`
	Name     *string  `json:"name,omitempty"`
	Website  *string  `json:"website,omitempty"`
}

type OrganizationEmbeds struct {
	Users []*User `json:"users,omitempty" gorm:"foreignKey:OrganizationId"`
}

func MakeOrganization() *Organization {
	return &Organization{
		Model: model.NewModel(),
		// ParentId: goat.NilID(),
		Name:    "",
		Website: "",
	}
}

// getEmbedded returns nil if the embedded Users array is empty to avoid rending JSON values like `"_embedded": {}`
func (m *Organization) getEmbedded() *resource.Embedded[*OrganizationEmbeds] {
	if m.Embedded == nil || len(m.Embedded.Embeds.Users) == 0 {
		return nil
	}
	return m.Embedded
}

func (m *Organization) MarshalJSON() ([]byte, error) {
	type Alias Organization
	return json.Marshal(&struct {
		*Alias
		*resource.Embedded[*OrganizationEmbeds]
		*resource.Links
	}{
		Alias:    (*Alias)(m),
		Embedded: m.getEmbedded(),
		Links:    goat.MakeResourceLinks(OrganizationLinkKey, m.ID.String()),
	})
}
