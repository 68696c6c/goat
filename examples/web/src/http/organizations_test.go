package http

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/68696c6c/goat"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/68696c6c/web/app/models"
)

// TODO: add more assertions for response embeds and links.
// TODO: add tests for different user types and permissions.

func Test_HttpOrganizations_List(t *testing.T) {
	h := goat.NewHandlerTest(router)

	response := requireTestResponse(t, h, httpTestRequest{
		accessToken: requireAccessToken(t, h),
		method:      http.MethodGet,
		path:        "/api/organizations",
	})
	require.Equal(t, http.StatusOK, response.Code, "unexpected status code")

	result := make(map[string]any)
	err := json.Unmarshal(response.Body.Bytes(), &result)
	require.Nil(t, err, "failed to parse response body")
	assert.Len(t, result["_embedded"], len(f.Organizations))
}

func Test_HttpOrganizations_View(t *testing.T) {
	h := goat.NewHandlerTest(router)
	m := f.Organizations[0]

	response := requireTestResponse(t, h, httpTestRequest{
		accessToken: requireAccessToken(t, h),
		method:      http.MethodGet,
		path:        "/api/organizations/" + m.ID.String(),
	})
	require.Equal(t, http.StatusOK, response.Code, "unexpected status code")

	result := models.Organization{}
	err := json.Unmarshal(response.Body.Bytes(), &result)
	require.Nil(t, err, "failed to parse response body")
	assert.Equal(t, result.Name, m.Name)
	assert.Equal(t, result.Website, m.Website)
	assert.NotNil(t, result.Embeds.Users)
}

func Test_HttpOrganizations_Create(t *testing.T) {
	h := goat.NewHandlerTest(router)
	requestData := map[string]any{
		"name":    "Example Org",
		"website": "https://test.com",
	}

	response := requireTestResponse(t, h, httpTestRequest{
		accessToken: requireAccessToken(t, h),
		method:      http.MethodPost,
		path:        "/api/organizations",
		jsonBody:    requestData,
	})
	require.Equal(t, http.StatusCreated, response.Code, "unexpected status code")

	result := models.Organization{}
	err := json.Unmarshal(response.Body.Bytes(), &result)
	require.Nil(t, err, "failed to parse response body")
	assert.Equal(t, result.Name, requestData["name"])
	assert.Equal(t, result.Website, requestData["website"])
}

func Test_HttpOrganizations_Update(t *testing.T) {
	h := goat.NewHandlerTest(router)
	m := f.Organizations[0]
	requestData := map[string]any{
		"name":    "Updated Name",
		"website": "https://test.com",
	}

	response := requireTestResponse(t, h, httpTestRequest{
		accessToken: requireAccessToken(t, h),
		method:      http.MethodPut,
		path:        "/api/organizations/" + m.ID.String(),
		jsonBody:    requestData,
	})
	require.Equal(t, http.StatusOK, response.Code, "unexpected status code")

	result := models.Organization{}
	err := json.Unmarshal(response.Body.Bytes(), &result)
	require.Nil(t, err, "failed to parse response body")
	assert.Equal(t, result.Name, requestData["name"])
	assert.Equal(t, result.Website, requestData["website"])
}

func Test_HttpOrganizations_Delete(t *testing.T) {
	h := goat.NewHandlerTest(router)

	// Create a new record to delete to avoid deleting one of the fixtures.
	requestData := map[string]any{
		"name":    "Delete Me",
		"website": "https://test.com",
	}
	createResponse := requireTestResponse(t, h, httpTestRequest{
		accessToken: requireAccessToken(t, h),
		method:      http.MethodPost,
		path:        "/api/organizations",
		jsonBody:    requestData,
	})
	require.Equal(t, http.StatusCreated, createResponse.Code, "unexpected status code")

	// Get the id of the created record.
	subject := models.Organization{}
	err := json.Unmarshal(createResponse.Body.Bytes(), &subject)
	require.Nil(t, err, "failed to parse response body")

	// Delete the record.
	response := requireTestResponse(t, h, httpTestRequest{
		accessToken: requireAccessToken(t, h),
		method:      http.MethodDelete,
		path:        "/api/organizations/" + subject.ID.String(),
	})
	assert.Equal(t, http.StatusNoContent, response.Code, "unexpected status code")
}
