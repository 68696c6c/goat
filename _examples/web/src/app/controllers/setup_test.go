package controllers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/68696c6c/goat"
	"github.com/68696c6c/goat/query"
	"github.com/gin-gonic/gin"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/68696c6c/web/app"
	"github.com/68696c6c/web/app/lib/auth"
	"github.com/68696c6c/web/app/models"
	"github.com/68696c6c/web/app/repos"
	"github.com/68696c6c/web/database"
	"github.com/68696c6c/web/test"
)

// import (
//
//	"context"
//	"net/url"
//	"os"
//	"testing"
//
//	"github.com/68696c6c/goat"
//	"github.com/68696c6c/goat/query"
//	"github.com/pkg/errors"
//	"github.com/pressly/goose/v3"
//	"gorm.io/gorm"
//
//	"github.com/68696c6c/web/app"
//	"github.com/68696c6c/web/app/models"
//	"github.com/68696c6c/web/app/repos"
//	_ "github.com/68696c6c/web/db/migrations"
//	"github.com/68696c6c/web/test"
//
// )
var (
	tc app.App
	f  test.Fixtures
	// router goat.Router

)

func TestMain(m *testing.M) {
	// err := os.Setenv("LOG_LEVEL", "debug")
	// if err != nil {
	// 	panic(err)
	// }
	// err = os.Setenv("LOG_STACKTRACE", "1")
	// if err != nil {
	// 	panic(err)
	// }
	// err = os.Setenv("BASE_URL", "https://test.com")
	// if err != nil {
	// 	panic(err)
	// }
	test.MustSetEnv("LOG_LEVEL", "debug")
	test.MustSetEnv("LOG_STACKTRACE", "1")
	test.MustSetEnv("BASE_URL", "https://test.com")

	goat.MustInit()

	var err error
	tc, err = app.InitApp(func() (*gorm.DB, error) {
		tdb, err := goat.GetDB(goat.DatabaseConfig{
			Debug:    goat.EnvBool("TEST_DB_DEBUG", true),
			Host:     goat.EnvString("TEST_DB_HOST", "db"),
			Port:     goat.EnvInt("TEST_DB_PORT", 3306),
			Database: goat.EnvString("TEST_DB_DATABASE", "web_test_controllers"),
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

	tc.Auth = newDummyAuthService()

	os.Exit(m.Run())
}

func requireTestResponseAndContext(t *testing.T) (*httptest.ResponseRecorder, *gin.Context) {
	w := httptest.NewRecorder()
	tcx, _ := gin.CreateTestContext(w)

	// The request method and url don't matter since controllers aren't
	// concerned with routing so the tests call the handlers directly.
	req, err := http.NewRequest(http.MethodGet, "https://test.com", nil)
	require.Nil(t, err, "failed to create ")
	tcx.Request = req

	return w, tcx
}

type dummyAuthService struct{}

func newDummyAuthService() auth.Service {
	return dummyAuthService{}
}

func (a dummyAuthService) GenerateClientCredentials() (auth.Client, error) {
	return auth.Client{}, nil
}

func (a dummyAuthService) Authorize(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (a dummyAuthService) GetToken(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (a dummyAuthService) ValidateToken(r *http.Request) (oauth2.TokenInfo, error) {
	return nil, nil
}

func (a dummyAuthService) GetTokenClaims(token string) (*auth.Claims, error) {
	return nil, nil
}

func (a dummyAuthService) GetCurrentUserAndClaims(cx *gin.Context) (*models.User, *auth.Claims, error) {
	return &models.User{}, &auth.Claims{
		CustomClaims: auth.CustomClaims{
			Users: auth.ActionPermissions{
				auth.Create: auth.All,
				auth.Read:   auth.All,
				auth.Update: auth.All,
				auth.Delete: auth.All,
			},
			Organizations: auth.ActionPermissions{
				auth.Create: auth.All,
				auth.Read:   auth.All,
				auth.Update: auth.All,
				auth.Delete: auth.All,
			},
		},
		StandardClaims: jwt.StandardClaims{},
	}, nil
}

type dummyOrganizationsRepo struct{}

func newDummyOrganizationsRepo() repos.OrganizationsRepo {
	return dummyOrganizationsRepo{}
}

func (r dummyOrganizationsRepo) Create(cx context.Context, u models.OrganizationRequest) (*models.Organization, error) {
	return nil, nil
}

func (r dummyOrganizationsRepo) Update(cx context.Context, id goat.ID, u models.OrganizationRequest) (*models.Organization, error) {
	return nil, nil
}

func (r dummyOrganizationsRepo) Filter(cx context.Context, q query.Builder) ([]*models.Organization, query.Builder, error) {
	return []*models.Organization{}, nil, nil
}

func (r dummyOrganizationsRepo) GetByID(cx context.Context, id goat.ID, loadRelations ...bool) (*models.Organization, error) {
	return nil, nil
}

func (r dummyOrganizationsRepo) Save(cx context.Context, m *models.Organization) error {
	return nil
}

func (r dummyOrganizationsRepo) Delete(cx context.Context, m *models.Organization) error {
	return nil
}
