package cmd

import (
	"github.com/68696c6c/goat"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gorm.io/gorm"

	"github.com/68696c6c/web/app/enums"
	"github.com/68696c6c/web/app/models"
)

func init() {
	Root.AddCommand(&cobra.Command{
		Use:   "seed",
		Short: "Seeds the database with example data.",
		RunE: func(cmd *cobra.Command, args []string) error {
			goat.MustInit()

			db, err := goat.GetMainDB()
			if err != nil {
				return errors.Wrap(err, "failed to initialize migration connection")
			}

			var organizations []*models.Organization
			var users []*models.User

			// Organization 1
			{
				org := models.NewOrganization()
				org.ID = mustParseId("87c3cd85-2849-41b3-ae5e-afeee71c875c")
				org.Name = "Organization 1"
				org.Website = "https://example1.com"
				organizations = append(organizations, org)
				{
					user := models.NewUser()
					user.ID = mustParseId("d5eed2d6-d679-4f4c-a955-d15e05942caa")
					user.OrganizationID = org.ID
					user.Level = enums.UserLevelSuper
					user.Name = "Org1 Super1"
					user.Email = "super1@example1.com"
					users = append(users, user)
				}
			}

			// Organization 2
			{
				org := models.NewOrganization()
				org.ID = mustParseId("d5c9c8d6-8207-4703-9ca8-edd4ea498299")
				org.Name = "Organization 2"
				org.Website = "https://example2.com"
				organizations = append(organizations, org)
				{
					user := models.NewUser()
					user.ID = mustParseId("a873fba5-0ae8-4d27-bdbe-cce0fe434a92")
					user.OrganizationID = org.ID
					user.Level = enums.UserLevelAdmin
					user.Name = "Org2 Admin1"
					user.Email = "admin1@example2.com"
					users = append(users, user)
				}
				{
					user := models.NewUser()
					user.ID = mustParseId("11919df8-daed-40f6-805d-b5fa996659c2")
					user.OrganizationID = org.ID
					user.Level = enums.UserLevelUser
					user.Name = "Org2 User1"
					user.Email = "user1@example2.com"
					users = append(users, user)
				}
				{
					user := models.NewUser()
					user.ID = mustParseId("cfc834c7-f723-41da-b311-33155a0f452b")
					user.OrganizationID = org.ID
					user.Level = enums.UserLevelUser
					user.Name = "Org2 User2"
					user.Email = "user2@example2.com"
					users = append(users, user)
				}
			}

			// Organization 3
			{
				org := models.NewOrganization()
				org.ID = mustParseId("8ee0d831-e986-4a31-aa87-bd797c219201")
				org.Name = "Organization 3"
				org.Website = "https://example3.com"
				organizations = append(organizations, org)
				{
					user := models.NewUser()
					user.ID = mustParseId("d5e98ec6-ee8d-4abc-9405-20b53d4d3e76")
					user.OrganizationID = org.ID
					user.Level = enums.UserLevelAdmin
					user.Name = "Org3 Admin1"
					user.Email = "admin1@example3.com"
					users = append(users, user)
				}
				{
					user := models.NewUser()
					user.ID = mustParseId("959c4f0b-7763-4d52-9506-4b2a5714d29a")
					user.OrganizationID = org.ID
					user.Level = enums.UserLevelUser
					user.Name = "Org3 User1"
					user.Email = "user1@example3.com"
					users = append(users, user)
				}
				{
					user := models.NewUser()
					user.ID = mustParseId("e9aaa4d9-2a25-4739-8bc4-ef7851fbba86")
					user.OrganizationID = org.ID
					user.Level = enums.UserLevelUser
					user.Name = "Org3 User2"
					user.Email = "user2@example3.com"
					users = append(users, user)
				}
			}

			err = db.Transaction(func(tx *gorm.DB) error {
				for _, org := range organizations {
					if err := tx.Create(org).Error; err != nil {
						return errors.Wrap(err, "failed to persist organization")
					}
				}

				for _, user := range users {
					if err := tx.Create(user).Error; err != nil {
						return errors.Wrap(err, "failed to persist user")
					}
				}

				return nil
			})
			if err != nil {
				return errors.Wrap(err, "failed to execute transaction")
			}

			return nil
		},
	})
}

func mustParseId(id string) goat.ID {
	result, err := goat.ParseID(id)
	if err != nil {
		panic(errors.Wrapf(err, "failed to parse id: %s", id))
	}
	return result
}
