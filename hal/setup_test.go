package hal

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type user struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	*ResourceEmbeds[userEmbeds]
}

func (m *user) getEmbedded() *ResourceEmbeds[userEmbeds] {
	if m.ResourceEmbeds == nil || len(m.ResourceEmbeds.Embeds.Phones) == 0 {
		return nil
	}
	return m.ResourceEmbeds
}

func (m *user) MarshalJSON() ([]byte, error) {
	type Alias user
	return json.Marshal(&struct {
		*Alias
		*ResourceEmbeds[userEmbeds]
		*ResourceLinks
	}{
		Alias:          (*Alias)(m),
		ResourceEmbeds: m.getEmbedded(),
		ResourceLinks:  NewLinks().AddLink("self", NewLink(fmt.Sprintf("%s/%s", exampleUsersUrl, strconv.Itoa(m.ID)))),
	})
}

type userEmbeds struct {
	Phones []phone `json:"phones,omitempty"`
}

func makeUserEmbeds(phoneIDs ...int) *ResourceEmbeds[userEmbeds] {
	var phones []phone
	for _, id := range phoneIDs {
		phones = append(phones, makePhone(id))
	}
	return &ResourceEmbeds[userEmbeds]{
		Embeds: userEmbeds{
			Phones: phones,
		},
	}
}

type phone struct {
	ID    int    `json:"id"`
	Phone string `json:"phone"`
}

func makePhone(id int) phone {
	digit := strconv.Itoa(id)
	var digits []string
	for i := 0; i < 10; i++ {
		digits = append(digits, digit)
	}
	return phone{
		ID:    id,
		Phone: strings.Join(digits, ""),
	}
}

const exampleUsersUrl = "https://example.com/users"
