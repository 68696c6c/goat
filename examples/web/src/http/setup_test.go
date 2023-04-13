package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/68696c6c/goat"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/68696c6c/web/app"
	"github.com/68696c6c/web/database"
	"github.com/68696c6c/web/test"
)

var (
	tc     app.App
	f      test.Fixtures
	router goat.Router
)

const authClientID = "test_client_id"

func TestMain(m *testing.M) {
	test.MustSetEnv("LOG_LEVEL", "debug")
	test.MustSetEnv("LOG_STACKTRACE", "1")
	test.MustSetEnv("BASE_URL", "https://test.com")
	test.MustSetEnv("AUTH_CLIENT_ID", authClientID)
	test.MustSetEnv("AUTH_CLIENT_SECRET", "test_client_secret")
	test.MustSetEnv("AUTH_CLIENT_PUBLIC", "1")
	test.MustSetEnv("AUTH_SIGNATURE_KEY", "00000000")

	goat.MustInit()

	var err error
	tc, err = app.InitApp(func() (*gorm.DB, error) {
		tdb, err := goat.GetDB(goat.DatabaseConfig{
			Debug:    goat.EnvBool("TEST_DB_DEBUG", true),
			Host:     goat.EnvString("TEST_DB_HOST", "db"),
			Port:     goat.EnvInt("TEST_DB_PORT", 3306),
			Database: goat.EnvString("TEST_DB_DATABASE", "web_test_http"),
			Username: goat.EnvString("TEST_DB_USERNAME", "root"),
			Password: url.QueryEscape(goat.EnvString("TEST_DB_PASSWORD", "secret")),
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to initialize test db connection")
		}

		err = database.ResetDB(tdb)
		if err != nil {
			return nil, errors.Wrap(err, "failed to reset test db")
		}

		return tdb, nil
	})
	if err != nil {
		panic(err)
	}

	f = test.MustGetPersistedFixtures(tc.DB)

	r, err := InitRouter(tc)
	if err != nil {
		panic(err)
	}
	router = r

	os.Exit(m.Run())
}

type httpTestRequest struct {
	method      string
	path        string
	accessToken string
	jsonBody    any
}

func requireTestResponse(t *testing.T, h *goat.HandlerTest, r httpTestRequest) *httptest.ResponseRecorder {
	request, err := h.NewRequest(r.path)
	require.Nil(t, err, "failed to create request")

	request.SetMethod(r.method).SetHeader("Authorization", r.accessToken)
	if r.jsonBody != nil {
		request.SetBodyJSON(r.jsonBody)
	}
	response, err := h.SetRequest(request).Send()
	require.Nil(t, err, "failed to send request")

	return response
}

func requireAccessToken(t *testing.T, h *goat.HandlerTest) string {

	// Get an authorization code.
	authCodeRequest, err := h.NewRequest("/api/tokens/authorize")
	require.Nil(t, err, "failed to create authorization code request")

	authCodeRequest.SetMethod(http.MethodPost).
		SetBodyForm(goat.MakeRequestData().
			Set("grant_type", "authorization_code").
			Set("username", f.SuperUser.Email).
			Set("response_type", "code").
			Set("redirect_uri", "http://127.0.0.1/api/tokens/callback").
			Set("client_id", authClientID).
			Values())
	authCodeResponse, err := h.SetRequest(authCodeRequest).Send()
	require.Nil(t, err, "failed to send authorization code request")
	require.Equal(t, http.StatusFound, authCodeResponse.Code, "authorization code request returned an unexpected status code")

	// Follow the redirect to get the authorization code.
	redirectUrl, err := authCodeResponse.Result().Location()
	require.Nil(t, err, "failed to get authorization request redirect location")

	redirectRequest, err := goat.NewRequest(redirectUrl.String())
	require.Nil(t, err, "failed to create redirect request")

	redirectResponse, err := h.SetRequest(redirectRequest.SetMethod(http.MethodGet)).Send()
	require.Nil(t, err, "failed to send redirect request")
	require.Equal(t, http.StatusOK, redirectResponse.Code, "redirect request returned an unexpected status code")

	redirectResult := make(map[string]string)
	err = json.Unmarshal(redirectResponse.Body.Bytes(), &redirectResult)
	require.Nil(t, err, "failed to parse redirect response body")

	authorizationCode, ok := redirectResult["code"]
	require.True(t, ok, "redirect response did not include an authorization code")

	// Exchange the authorization code for a JWT.
	jwtRequest, err := h.NewRequest("/api/tokens/exchange")
	require.Nil(t, err, "failed to create jwt request")

	jwtRequest.SetMethod(http.MethodPost).
		SetBodyForm(goat.MakeRequestData().
			Set("grant_type", "authorization_code").
			Set("redirect_uri", "http://127.0.0.1/api/tokens/callback").
			Set("client_id", authClientID).
			Set("code", authorizationCode).
			Values())
	jwtResponse, err := h.SetRequest(jwtRequest).Send()
	require.Nil(t, err, "failed to send jwt request")
	require.Equal(t, http.StatusOK, jwtResponse.Code, "jwt request returned unexpected status code")

	jwtResult := make(map[string]any)
	err = json.Unmarshal(jwtResponse.Body.Bytes(), &jwtResult)
	require.Nil(t, err, "failed to parse jwt response body")

	tokenType, ok := jwtResult["token_type"]
	require.True(t, ok, "jwt response did not include a token_type")

	accessToken, ok := jwtResult["access_token"]
	require.True(t, ok, "jwt response did not include an access_token")

	return fmt.Sprintf("%v %v", tokenType, accessToken)
}
