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

func Test_HttpUsers_List(t *testing.T) {
	h := goat.NewHandlerTest(router)

	response := requireTestResponse(t, h, httpTestRequest{
		accessToken: requireAccessToken(t, h),
		method:      http.MethodGet,
		path:        "/api/users",
	})
	require.Equal(t, http.StatusOK, response.Code, "unexpected status code")

	result := make(map[string]any)
	err := json.Unmarshal(response.Body.Bytes(), &result)
	require.Nil(t, err, "failed to parse response body")
	assert.Len(t, result["_embedded"], len(f.Users))
}

func Test_HttpUsers_View(t *testing.T) {
	h := goat.NewHandlerTest(router)
	m := f.Users[0]

	response := requireTestResponse(t, h, httpTestRequest{
		accessToken: requireAccessToken(t, h),
		method:      http.MethodGet,
		path:        "/api/users/" + m.ID.String(),
	})
	require.Equal(t, http.StatusOK, response.Code, "unexpected status code")

	result := models.User{}
	err := json.Unmarshal(response.Body.Bytes(), &result)
	require.Nil(t, err, "failed to parse response body")
	assert.Equal(t, result.Name, m.Name)
	assert.Equal(t, result.Email, m.Email)
	assert.NotNil(t, result.Embeds.Organization)
}

func Test_HttpUsers_Create(t *testing.T) {
	h := goat.NewHandlerTest(router)
	requestData := map[string]any{
		"organizationId": f.Organizations[0].ID.String(),
		"level":          "user",
		"name":           "Example User",
		"email":          "user@test.com",
	}

	response := requireTestResponse(t, h, httpTestRequest{
		accessToken: requireAccessToken(t, h),
		method:      http.MethodPost,
		path:        "/api/users",
		jsonBody:    requestData,
	})
	require.Equal(t, http.StatusCreated, response.Code, "unexpected status code")

	result := models.User{}
	err := json.Unmarshal(response.Body.Bytes(), &result)
	require.Nil(t, err, "failed to parse response body")
	assert.Equal(t, result.OrganizationID.String(), requestData["organizationId"])
	assert.Equal(t, result.Name, requestData["name"])
	assert.Equal(t, result.Email, requestData["email"])
}

func Test_HttpUsers_Update(t *testing.T) {
	h := goat.NewHandlerTest(router)
	m := f.Users[0]
	requestData := map[string]any{
		"name":  "Updated Name",
		"level": "admin",
	}

	response := requireTestResponse(t, h, httpTestRequest{
		accessToken: requireAccessToken(t, h),
		method:      http.MethodPut,
		path:        "/api/users/" + m.ID.String(),
		jsonBody:    requestData,
	})
	require.Equal(t, http.StatusOK, response.Code, "unexpected status code")

	result := models.User{}
	err := json.Unmarshal(response.Body.Bytes(), &result)
	require.Nil(t, err, "failed to parse response body")
	assert.Equal(t, result.Name, requestData["name"])
	assert.Equal(t, result.Level.String(), requestData["level"])
}

func Test_HttpUsers_Delete(t *testing.T) {
	h := goat.NewHandlerTest(router)

	// Create a new record to delete to avoid deleting one of the fixtures.
	requestData := map[string]any{
		"organizationId": f.Organizations[0].ID,
		"level":          "user",
		"name":           "Delete Me",
		"email":          "delete@test.com",
	}
	createResponse := requireTestResponse(t, h, httpTestRequest{
		accessToken: requireAccessToken(t, h),
		method:      http.MethodPost,
		path:        "/api/users",
		jsonBody:    requestData,
	})
	require.Equal(t, http.StatusCreated, createResponse.Code, "unexpected status code")

	// Get the id of the created record.
	subject := models.User{}
	err := json.Unmarshal(createResponse.Body.Bytes(), &subject)
	require.Nil(t, err, "failed to parse response body")

	// Delete the record.
	response := requireTestResponse(t, h, httpTestRequest{
		accessToken: requireAccessToken(t, h),
		method:      http.MethodDelete,
		path:        "/api/users/" + subject.ID.String(),
	})
	assert.Equal(t, http.StatusNoContent, response.Code, "unexpected status code")
}
