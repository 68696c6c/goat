package repos

import (
	"context"
	"testing"

	"github.com/68696c6c/goat"
	"github.com/68696c6c/goat/query"
	"github.com/icrowley/fake"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/68696c6c/web/app/enums"
	"github.com/68696c6c/web/app/models"
	"github.com/68696c6c/web/test"
)

func Test_UsersRepo_Create(t *testing.T) {
	inputOrgID := f.Organizations[0].ID
	inputLevel := enums.UserLevelUser
	inputName := fake.FullName()
	inputEmail := fake.EmailAddress()
	result, err := tc.usersRepo.Create(context.Background(), models.UserRequest{
		OrganizationId: &inputOrgID,
		Level:          &inputLevel,
		Name:           &inputName,
		Email:          &inputEmail,
	})
	require.Nil(t, err, "unexpected error")
	assert.Equal(t, inputOrgID, result.OrganizationID, "unexpected organization id")
	assert.Equal(t, inputLevel, result.Level, "unexpected level")
	assert.Equal(t, inputName, result.Name, "unexpected name")
	assert.Equal(t, inputEmail, result.Email, "unexpected email")
}

func Test_UsersRepo_Create_NilFields(t *testing.T) {
	result, err := tc.usersRepo.Create(context.Background(), models.UserRequest{
		OrganizationId: nil,
		Name:           nil,
		Email:          nil,
	})
	require.NotNil(t, err, "expected an error")
	require.Nil(t, result, "unexpected record returned")
}

func Test_UsersRepo_Update(t *testing.T) {
	id := f.Users[0].ID
	inputOrgID := f.Organizations[0].ID
	inputName := fake.FullName()
	result, err := tc.usersRepo.Update(context.Background(), id, models.UserRequest{
		OrganizationId: &inputOrgID,
		Name:           &inputName,
	})
	require.Nil(t, err, "unexpected error")
	assert.Equal(t, id, result.ID, "unexpected record returned")
	assert.Equal(t, inputOrgID, result.OrganizationID, "unexpected organization id")
	assert.Equal(t, inputName, result.Name, "unexpected name")
}

func Test_UsersRepo_Update_NotFound(t *testing.T) {
	result, err := tc.usersRepo.Update(context.Background(), goat.NewID(), models.UserRequest{
		OrganizationId: nil,
		Name:           nil,
	})
	require.NotNil(t, err, "expected an error")
	require.True(t, goat.RecordNotFound(err), "expected a not found error")
	require.Nil(t, result, "unexpected record returned")
}

func Test_UsersRepo_Update_NilFields(t *testing.T) {
	input := f.Users[0]
	result, err := tc.usersRepo.Update(context.Background(), input.ID, models.UserRequest{
		OrganizationId: nil,
		Name:           nil,
	})
	require.Nil(t, err, "unexpected error")
	assert.Equal(t, input.OrganizationID, result.OrganizationID, "unexpected organization id")
	assert.Equal(t, input.Name, result.Name, "unexpected name")
}

// TODO: add tests for filtering and pagination.

func Test_UsersRepo_Filter(t *testing.T) {
	q := query.NewQuery()
	result, _, err := tc.usersRepo.Filter(context.Background(), q)
	require.Nil(t, err, "unexpected error")
	assert.Len(t, result, len(f.Users), "unexpected number of rows returned")
}

func Test_UsersRepo_GetByID(t *testing.T) {
	id := f.Users[0].ID
	result, err := tc.usersRepo.GetByID(context.Background(), id)
	require.Nil(t, err, "unexpected error")
	assert.Equal(t, id, result.ID, "unexpected record returned")
	assert.Nil(t, result.ResourceEmbeds, "should not have loaded relations")
}

func Test_UsersRepo_GetByID_Preloading(t *testing.T) {
	id := f.Users[0].ID
	result, err := tc.usersRepo.GetByID(context.Background(), id, true)
	require.Nil(t, err, "unexpected error")
	assert.Equal(t, id, result.ID, "unexpected record returned")
	assert.NotNil(t, result.ResourceEmbeds, "should have loaded relations")
}

func Test_UsersRepo_GetByID_NotFound(t *testing.T) {
	result, err := tc.usersRepo.GetByID(context.Background(), goat.NewID())
	require.NotNil(t, err, "expected an error")
	require.True(t, goat.RecordNotFound(err), "expected a not found error")
	assert.Nil(t, result, "unexpected record returned")
}

func Test_UsersRepo_GetByID_NilId(t *testing.T) {
	result, err := tc.usersRepo.GetByID(context.Background(), goat.NilID())
	require.NotNil(t, err, "expected an error")
	require.True(t, goat.RecordNotFound(err), "expected a not found error")
	assert.Nil(t, result, "unexpected record returned")
}

func Test_UsersRepo_Save_Create(t *testing.T) {
	input := test.FakeUser(f.Organizations[0].ID)

	// Create the record.
	err := tc.usersRepo.Save(context.Background(), input)
	require.Empty(t, err, "save returned errors")
	assert.NotEqual(t, goat.NilID(), input.ID, "saving a new record didn't set the id")
	assert.NotNil(t, input.CreatedAt, "saving a new record didn't set created_at")
	assert.Nil(t, input.UpdatedAt, "saving a new record set updated_at")
}

func Test_UsersRepo_Save_Update_Success(t *testing.T) {
	input := f.Users[0]
	expected := fake.Word()
	input.Name = expected

	// Do the update.
	err := tc.usersRepo.Save(context.Background(), &input)
	require.Nil(t, err, "save returned an error")

	// Assert that the value was updated.
	result, err := tc.usersRepo.GetByID(context.Background(), input.ID)
	require.Nil(t, err, "unexpected error")
	assert.NotNil(t, input.UpdatedAt, "updating a record didn't set updated_at")
	assert.Equal(t, expected, result.Name, "updating a record failed to change the value")
}

func Test_UsersRepo_Delete(t *testing.T) {

	// Persist a new user so that we can delete it without affecting other tests.
	input := test.FakeUser(f.Organizations[0].ID)
	err := tc.db.Save(input).Error
	require.Nil(t, err, "failed to save temporary record")

	err = tc.usersRepo.Delete(context.Background(), input)
	require.Empty(t, err, "unexpected error")
	goat.AssertRecordDeleted[*models.User](t, tc.db, input, "failed to delete user")
}

func Test_UsersRepo_GetByEmail(t *testing.T) {
	email := f.Users[0].Email
	result, err := tc.usersRepo.GetByEmail(context.Background(), email)
	require.Nil(t, err, "unexpected error")
	assert.Equal(t, email, result.Email, "unexpected record returned")
}
