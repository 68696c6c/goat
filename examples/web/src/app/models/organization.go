package models

import (
	"encoding/json"

	"github.com/68696c6c/goat"
	"github.com/68696c6c/goat/hal"
	"github.com/68696c6c/goat/model"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const OrganizationLinkKey = "organizations"

type Organization struct {
	*model.Model

	Name    string `json:"name" binding:"required"`
	Website string `json:"website" binding:"required"`

	*model.Timestamps
	*model.SoftDelete
	*hal.ResourceEmbeds[*OrganizationEmbeds]
}

type OrganizationRequest struct {
	Name    *string `json:"name,omitempty"`
	Website *string `json:"website,omitempty"`
}

type OrganizationEmbeds struct {
	Users []*User `json:"users,omitempty" gorm:"foreignKey:OrganizationID"`
}

func NewOrganization() *Organization {
	return &Organization{
		Model:   model.NewModel(),
		Name:    "",
		Website: "",
	}
}

// getEmbedded returns nil if the embedded Users array is empty to avoid rending JSON values like `"_embedded": {}`
func (m *Organization) getEmbedded() *hal.ResourceEmbeds[*OrganizationEmbeds] {
	if m.ResourceEmbeds == nil || m.Embeds == nil || m.Embeds.Users == nil || len(m.Embeds.Users) == 0 {
		return nil
	}
	return m.ResourceEmbeds
}

func (m *Organization) MarshalJSON() ([]byte, error) {
	type Alias Organization
	return json.Marshal(&struct {
		*Alias
		*hal.ResourceEmbeds[*OrganizationEmbeds]
		*hal.ResourceLinks
	}{
		Alias:          (*Alias)(m),
		ResourceEmbeds: m.getEmbedded(),
		ResourceLinks:  goat.NewResourceLinks(OrganizationLinkKey, m.ID.String()),
	})
}

// BeforeDelete GORM hook deletes the record relations prior to deleting the record.
func (m *Organization) BeforeDelete(db *gorm.DB) error {
	err := db.Where("organization_id = ?", m.ID).Delete(&User{}).Error
	if err != nil {
		return errors.Wrap(err, "failed to delete related users")
	}
	return nil
}
