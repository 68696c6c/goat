package repos

import (
	"net/url"
	"os"
	"testing"

	"github.com/68696c6c/goat"
	"github.com/icrowley/fake"
	"github.com/pkg/errors"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/68696c6c/example/app/enums"
	"github.com/68696c6c/example/app/models"
	_ "github.com/68696c6c/example/db/migrations"
)

var (
	tc *testContainer
	f  testFixtures
)

func TestMain(m *testing.M) {
	goat.Init()

	tdb := mustInitTestDb(goat.DatabaseConnection{
		Debug:           goat.EnvBool("TEST_DB_DEBUG", true),
		Host:            goat.EnvString("TEST_DB_HOST", "db"),
		Port:            goat.EnvInt("TEST_DB_PORT", 3306),
		Database:        goat.EnvString("TEST_DB_DATABASE", "example_test_repos"),
		Username:        goat.EnvString("TEST_DB_USERNAME", "root"),
		Password:        url.QueryEscape(goat.EnvString("TEST_DB_PASSWORD", "secret")),
		MultiStatements: false,
	})
	tc = newTestContainer(tdb)
	f = mustGetPersistedFixtures(tc.db)

	os.Exit(m.Run())
}

type testContainer struct {
	db                *gorm.DB
	organizationsRepo OrganizationsRepo
	usersRepo         UsersRepo
}

func newTestContainer(db *gorm.DB) *testContainer {
	if tc != nil {
		return tc
	}

	tc = &testContainer{
		db:                db,
		organizationsRepo: NewOrganizationsRepo(db),
		usersRepo:         NewUsersRepo(db),
	}
	return tc
}

type testFixtures struct {
	organizations []models.Organization
	users         []models.User
}

func mustGetPersistedFixtures(db *gorm.DB) testFixtures {
	var organizations []models.Organization
	var users []models.User

	organization := fakeOrganization()
	mustPersistFixture[*models.Organization](db, organization)
	organizations = append(organizations, *organization)

	user := fakeUser(organization.ID)
	mustPersistFixture[*models.User](db, user)
	users = append(users, *user)

	return testFixtures{
		organizations: organizations,
		users:         users,
	}
}

func fakeOrganization() *models.Organization {
	result := models.MakeOrganization()
	result.Name = fake.FullName()
	result.Website = fake.DomainName()
	return result
}

func fakeUser(organizationId goat.ID) *models.User {
	result := models.MakeUser()
	result.OrganizationId = organizationId
	result.Level = enums.UserLevelUser
	result.Name = fake.FullName()
	result.Email = fake.EmailAddress()
	return result
}

func assertRecordDeleted[M models.Model](t *testing.T, db *gorm.DB, input M, msg string) {
	err := db.First(input).Error
	assert.NotNil(t, err)
	assert.True(t, goat.RecordNotFound(err), msg)
}

func mustPersistFixture[T any](db *gorm.DB, m T) {
	if err := db.Create(m).Error; err != nil {
		panic(err)
	}
}

func mustInitTestDb(connectionConfig goat.DatabaseConnection) *gorm.DB {
	tdb, err := goat.GetCustomDB(connectionConfig)
	if err != nil {
		panic(errors.Wrap(err, "failed to initialize test db connection"))
	}

	if err := goose.SetDialect("mysql"); err != nil {
		goat.ExitError(errors.Wrap(err, "error initializing goose"))
	}

	sqlDb, err := tdb.DB()
	if err != nil {
		panic(errors.Wrap(err, "failed to get sql db"))
	}

	err = goose.Run("down-to", sqlDb, ".", "0")
	if err != nil {
		panic(errors.Wrap(err, "failed to reset test db"))
	}

	err = goose.Run("up", sqlDb, ".")
	if err != nil {
		panic(errors.Wrap(err, "failed to migrate test db"))
	}

	return tdb
}
