package repos

import (
	"context"
	"testing"

	"github.com/68696c6c/goat"
	"github.com/68696c6c/goat/query"
	"github.com/68696c6c/goat/resource"
	"github.com/icrowley/fake"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/68696c6c/example/app/models"
)

// Filter tests must come first in order to have a reliable expected count since other tests add records.

func Test_UsersRepo_Filter(t *testing.T) {
	q := query.NewQuery()
	p := resource.NewPagination()
	result, _, err := tc.usersRepo.Filter(context.Background(), q, p)
	require.Nil(t, err, "unexpected error")
	assert.Len(t, result, len(f.users), "unexpected number of rows returned")
}

func Test_UsersRepo_GetById(t *testing.T) {
	id := f.users[0].ID
	result, err := tc.usersRepo.GetById(context.Background(), id)
	require.Nil(t, err, "unexpected error")
	assert.Equal(t, id, result.ID, "unexpected record returned")
	assert.Nil(t, result.Embedded, "should not have loaded relations")
}

func Test_UsersRepo_GetById_Preloading(t *testing.T) {
	id := f.users[0].ID
	result, err := tc.usersRepo.GetById(context.Background(), id, true)
	require.Nil(t, err, "unexpected error")
	assert.Equal(t, id, result.ID, "unexpected record returned")
	assert.NotNil(t, result.Embedded, "should have loaded relations")
}

func Test_UsersRepo_GetById_NotFound(t *testing.T) {
	result, err := tc.usersRepo.GetById(context.Background(), goat.NewID())
	require.NotNil(t, err, "expected an error")
	require.True(t, goat.RecordNotFound(err), "expected a not found error")
	assert.Nil(t, result, "unexpected record returned")
}

func Test_UsersRepo_GetById_NilId(t *testing.T) {
	result, err := tc.usersRepo.GetById(context.Background(), goat.NilID())
	require.NotNil(t, err, "expected an error")
	require.True(t, goat.RecordNotFound(err), "expected a not found error")
	assert.Nil(t, result, "unexpected record returned")
}

func Test_UsersRepo_Save_Create(t *testing.T) {
	input := fakeUser(f.organizations[0].ID)

	// Create the record.
	err := tc.usersRepo.Save(context.Background(), input)
	require.Empty(t, err, "save returned errors")
	assert.NotEqual(t, goat.NilID(), input.ID, "saving a new record didn't set the id")
	assert.NotNil(t, input.CreatedAt, "saving a new record didn't set created_at")
	assert.Nil(t, input.UpdatedAt, "saving a new record set updated_at")
}

func Test_UsersRepo_Save_Update_Success(t *testing.T) {
	input := f.users[0]
	expected := fake.Word()
	input.Name = expected

	// Do the update.
	err := tc.usersRepo.Save(context.Background(), &input)
	require.Nil(t, err, "save returned an error")

	// Assert that the value was updated.
	result, err := tc.usersRepo.GetById(context.Background(), input.ID)
	require.Nil(t, err, "unexpected error")
	assert.NotNil(t, input.UpdatedAt, "updating a record didn't set updated_at")
	assert.Equal(t, expected, result.Name, "updating a record failed to change the value")
}

func Test_UsersRepo_Delete(t *testing.T) {

	// Persist a new user so that we can delete it without affecting other tests.
	input := fakeUser(f.organizations[0].ID)
	err := tc.db.Save(input).Error
	require.Nil(t, err, "failed to save temporary record")

	err = tc.usersRepo.Delete(context.Background(), input)
	require.Empty(t, err, "unexpected error")
	assertRecordDeleted[*models.User](t, tc.db, input, "failed to delete user")
}
