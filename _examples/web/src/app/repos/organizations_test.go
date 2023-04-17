package repos

import (
	"context"
	"testing"

	"github.com/68696c6c/goat"
	"github.com/68696c6c/goat/query"
	"github.com/icrowley/fake"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/68696c6c/web/app/models"
	"github.com/68696c6c/web/test"
)

func Test_OrganizationsRepo_Create(t *testing.T) {
	inputName := fake.Brand()
	inputWebsite := fake.DomainName()
	result, err := tc.organizationsRepo.Create(context.Background(), models.OrganizationRequest{
		Name:    &inputName,
		Website: &inputWebsite,
	})
	require.Nil(t, err, "unexpected error")
	assert.Equal(t, inputName, result.Name, "unexpected name")
	assert.Equal(t, inputWebsite, result.Website, "unexpected website")
}

func Test_OrganizationsRepo_Create_NilFields(t *testing.T) {
	result, err := tc.organizationsRepo.Create(context.Background(), models.OrganizationRequest{
		Name:    nil,
		Website: nil,
	})
	require.NotNil(t, err, "expected an error")
	require.Nil(t, result, "unexpected record returned")
}

func Test_OrganizationsRepo_Update(t *testing.T) {
	id := f.Organizations[0].ID
	inputName := fake.Brand()
	inputWebsite := fake.DomainName()
	result, err := tc.organizationsRepo.Update(context.Background(), id, models.OrganizationRequest{
		Name:    &inputName,
		Website: &inputWebsite,
	})
	require.Nil(t, err, "unexpected error")
	assert.Equal(t, id, result.ID, "unexpected record returned")
	assert.Equal(t, inputName, result.Name, "unexpected name")
	assert.Equal(t, inputWebsite, result.Website, "unexpected website")
}

func Test_OrganizationsRepo_Update_NotFound(t *testing.T) {
	result, err := tc.organizationsRepo.Update(context.Background(), goat.NewID(), models.OrganizationRequest{
		Name:    nil,
		Website: nil,
	})
	require.NotNil(t, err, "expected an error")
	require.True(t, goat.RecordNotFound(err), "expected a not found error")
	require.Nil(t, result, "unexpected record returned")
}

func Test_OrganizationsRepo_Update_NilFields(t *testing.T) {
	input := f.Organizations[0]
	result, err := tc.organizationsRepo.Update(context.Background(), input.ID, models.OrganizationRequest{
		Name:    nil,
		Website: nil,
	})
	require.Nil(t, err, "unexpected error")
	assert.Equal(t, input.Name, result.Name, "unexpected name")
	assert.Equal(t, input.Website, result.Website, "unexpected website")
}

// TODO: add tests for filtering and pagination.

func Test_OrganizationsRepo_Filter(t *testing.T) {
	q := query.NewQuery()
	result, _, err := tc.organizationsRepo.Filter(context.Background(), q)
	require.Nil(t, err, "unexpected error")
	assert.Len(t, result, len(f.Organizations), "unexpected number of rows returned")
}

func Test_OrganizationsRepo_GetByID(t *testing.T) {
	id := f.Organizations[0].ID
	result, err := tc.organizationsRepo.GetByID(context.Background(), id)
	require.Nil(t, err, "unexpected error")
	assert.Equal(t, id, result.ID, "unexpected record returned")
	assert.Nil(t, result.ResourceEmbeds, "should not have loaded relations")
}

func Test_OrganizationsRepo_GetByID_Preload(t *testing.T) {
	id := f.Organizations[0].ID
	result, err := tc.organizationsRepo.GetByID(context.Background(), id, true)
	require.Nil(t, err, "unexpected error")
	assert.Equal(t, id, result.ID, "unexpected record returned")
	assert.NotNil(t, result.ResourceEmbeds, "should have loaded relations")
}

func Test_OrganizationsRepo_GetByID_NotFound(t *testing.T) {
	result, err := tc.organizationsRepo.GetByID(context.Background(), goat.NewID())
	require.NotNil(t, err, "expected an error")
	require.True(t, goat.RecordNotFound(err), "expected a not found error")
	assert.Nil(t, result, "unexpected record returned")
}

func Test_OrganizationsRepo_GetByID_NilId(t *testing.T) {
	result, err := tc.organizationsRepo.GetByID(context.Background(), goat.NilID())
	require.NotNil(t, err, "expected an error")
	require.True(t, goat.RecordNotFound(err), "expected a not found error")
	assert.Nil(t, result, "unexpected record returned")
}

func Test_OrganizationsRepo_Save_Create(t *testing.T) {
	input := test.FakeOrganization()

	err := tc.organizationsRepo.Save(context.Background(), input)
	require.Empty(t, err, "save returned an error")

	assert.NotEqual(t, goat.NilID(), input.ID, "saving a new record didn't set the id")
	assert.NotNil(t, input.CreatedAt, "saving a new record didn't set created timestamp")
	assert.Nil(t, input.UpdatedAt, "saving a new record set updated timestamp")
}

func Test_OrganizationsRepo_Save_Update_Success(t *testing.T) {
	input := f.Organizations[0]
	expected := fake.Word()
	input.Name = expected

	err := tc.organizationsRepo.Save(context.Background(), &input)
	require.Nil(t, err, "save returned an error")

	result, err := tc.organizationsRepo.GetByID(context.Background(), input.ID)
	require.Nil(t, err, "unexpected error")
	assert.NotNil(t, input.UpdatedAt, "updating a record didn't set updated timestamp")
	assert.Equal(t, expected, result.Name, "updating a record failed to change the value")
}

func Test_OrganizationsRepo_Delete(t *testing.T) {

	// Persist a new record so that we can delete it without affecting other tests.
	input := test.FakeOrganization()
	err := tc.db.Save(input).Error
	require.Nil(t, err, "failed to save temporary record")

	err = tc.organizationsRepo.Delete(context.Background(), input)
	require.Empty(t, err, "unexpected error")
	goat.AssertRecordDeleted[*models.Organization](t, tc.db, input, "failed to delete organization")
}
