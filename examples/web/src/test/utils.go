package test

import (
	"os"

	"github.com/68696c6c/goat"
	"github.com/icrowley/fake"
	"gorm.io/gorm"

	"github.com/68696c6c/web/app/enums"
	"github.com/68696c6c/web/app/models"
)

type Fixtures struct {
	Organizations []models.Organization
	Users         []models.User
	SuperUser     models.User
}

func MustGetPersistedFixtures(db *gorm.DB) Fixtures {
	var organizations []models.Organization
	var users []models.User

	// Organization 1
	org := FakeOrganization()
	mustPersistFixture[*models.Organization](db, org)
	organizations = append(organizations, *org)

	superUser := FakeUser(org.ID)
	superUser.Level = enums.UserLevelSuper
	mustPersistFixture[*models.User](db, superUser)

	// Organization 2
	{
		org := FakeOrganization()
		mustPersistFixture[*models.Organization](db, org)
		organizations = append(organizations, *org)
		{
			user := FakeUser(org.ID)
			user.Level = enums.UserLevelAdmin
			mustPersistFixture[*models.User](db, user)
			users = append(users, *user)
		}
		{
			user := FakeUser(org.ID)
			mustPersistFixture[*models.User](db, user)
			users = append(users, *user)
		}
		{
			user := FakeUser(org.ID)
			mustPersistFixture[*models.User](db, user)
			users = append(users, *user)
		}
	}

	// Organization 3
	{
		org := FakeOrganization()
		mustPersistFixture[*models.Organization](db, org)
		organizations = append(organizations, *org)
		{
			user := FakeUser(org.ID)
			user.Level = enums.UserLevelAdmin
			mustPersistFixture[*models.User](db, user)
			users = append(users, *user)
		}
		{
			user := FakeUser(org.ID)
			mustPersistFixture[*models.User](db, user)
			users = append(users, *user)
		}
		{
			user := FakeUser(org.ID)
			mustPersistFixture[*models.User](db, user)
			users = append(users, *user)
		}
	}

	// Tests will often use the superUser for authentication and f.Users[0] as the test subject.
	// To avoid tests accidentally deleting or modifying the superUser, we will add it to the end of the users array.
	users = append(users, *superUser)

	return Fixtures{
		Organizations: organizations,
		Users:         users,
		SuperUser:     *superUser,
	}
}

func MustSetEnv(key, value string) {
	err := os.Setenv(key, value)
	if err != nil {
		panic(err)
	}
}

func FakeOrganization() *models.Organization {
	result := models.NewOrganization()
	result.Name = fake.Brand()
	result.Website = fake.DomainName()
	return result
}

func FakeUser(organizationId goat.ID) *models.User {
	result := models.NewUser()
	result.OrganizationID = organizationId
	result.Level = enums.UserLevelUser
	result.Name = fake.FullName()
	result.Email = fake.EmailAddress()
	return result
}

func mustPersistFixture[T any](db *gorm.DB, m T) {
	if err := db.Create(m).Error; err != nil {
		panic(err)
	}
}
