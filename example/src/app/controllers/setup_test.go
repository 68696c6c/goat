package controllers

import (
	"context"
	"os"
	"testing"

	"github.com/68696c6c/goat"
	"github.com/68696c6c/goat/query2"
	"github.com/68696c6c/goat/resource"

	"github.com/68696c6c/example/app/models"
	"github.com/68696c6c/example/app/repos"
)

var tc *testContainer

func TestMain(m *testing.M) {
	goat.Init()

	tc = newTestContainer()

	os.Exit(m.Run())
}

type testContainer struct {
	usersRepo repos.UsersRepo
}

func newTestContainer() *testContainer {
	if tc == nil {
		tc = &testContainer{
			usersRepo: newDummyUsersRepo(),
		}
	}
	return tc
}

type dummyUsersRepo struct{}

func newDummyUsersRepo() repos.UsersRepo {
	return dummyUsersRepo{}
}

func (r dummyUsersRepo) Make() *models.User {
	return models.MakeUser()
}

func (r dummyUsersRepo) Create(cx context.Context, u models.UserRequest) (*models.User, error) {
	return nil, nil
}

func (r dummyUsersRepo) Update(cx context.Context, id goat.ID, u models.UserRequest) (*models.User, error) {
	return nil, nil
}

func (r dummyUsersRepo) Filter(cx context.Context, q query2.Builder, p resource.Pagination) ([]*models.User, resource.Pagination, error) {
	return []*models.User{}, resource.Pagination{}, nil
}

func (r dummyUsersRepo) ApplyFilterForUser(q query2.Builder, user *models.User) error {
	return nil
}

func (r dummyUsersRepo) GetById(cx context.Context, id goat.ID, loadRelations ...bool) (*models.User, error) {
	return nil, nil
}

func (r dummyUsersRepo) Save(cx context.Context, m *models.User) error {
	return nil
}

func (r dummyUsersRepo) Delete(cx context.Context, m *models.User) error {
	return nil
}

func (r dummyUsersRepo) FilterStrings(q query2.Builder, fields map[string][]string) error {
	return nil
}
