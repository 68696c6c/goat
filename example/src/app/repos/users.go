package repos

import (
	"context"
	"fmt"

	"github.com/68696c6c/goat"
	"github.com/68696c6c/goat/query"
	"github.com/68696c6c/goat/repo"
	"github.com/68696c6c/goat/resource"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/68696c6c/example/app/enums"
	"github.com/68696c6c/example/app/models"
)

type UsersRepo interface {
	repo.CRUD[*models.User, models.UserRequest]
	ApplyFilterForUser(q query.Builder, user *models.User) error
	FilterStrings(query query.Builder, fields map[string][]string) error
}

type usersRepo struct {
	db *gorm.DB
}

func NewUsersRepo(db *gorm.DB) UsersRepo {
	return usersRepo{
		db: db,
	}
}

func (r usersRepo) Make() *models.User {
	return models.MakeUser()
}

func (r usersRepo) Create(cx context.Context, u models.UserRequest) (*models.User, error) {
	m := r.Make()
	var errs []error
	if u.OrganizationId == nil {
		errs = append(errs, goat.MakeValidationError("organizationId", "required"))
	} else {
		m.OrganizationId = *u.OrganizationId
	}
	if u.Level == nil {
		errs = append(errs, goat.MakeValidationError("level", "required"))
	} else {
		m.Level = *u.Level
	}
	if u.Name == nil {
		errs = append(errs, goat.MakeValidationError("name", "required"))
	} else {
		m.Name = *u.Name
	}
	if u.Email == nil {
		errs = append(errs, goat.MakeValidationError("email", "required"))
	} else {
		_, err := repo.First[models.User](r.db.WithContext(cx), "email = ?", u.Email)
		if !goat.RecordNotFound(err) {
			errs = append(errs, goat.MakeValidationError("email", "unique"))
		} else {
			m.Email = *u.Email
		}
	}
	if len(errs) > 0 {
		return nil, goat.ErrorsToError(errs)
	}
	return m, nil
}

func (r usersRepo) Update(cx context.Context, id goat.ID, u models.UserRequest) (*models.User, error) {
	m, err := r.GetById(cx, id)
	if err != nil {
		return nil, err
	}
	m.OrganizationId = goat.ValueOrDefault[goat.ID](u.OrganizationId, m.OrganizationId)
	m.Level = goat.ValueOrDefault[enums.UserLevel](u.Level, m.Level)
	m.Name = goat.ValueOrDefault[string](u.Name, m.Name)
	return m, nil
}

// TODO: context???
func (r usersRepo) getBaseQuery() *gorm.DB {
	return r.db.Model(&models.User{})
}

func (r usersRepo) Filter(cx context.Context, q query.Builder, p resource.Pagination) ([]*models.User, resource.Pagination, error) {
	base := r.db.WithContext(cx).Model(&models.User{})
	return repo.Filter[models.User](base, q, p)
}

func (r usersRepo) ApplyFilterForUser(q query.Builder, user *models.User) error {
	switch user.Level {
	case enums.UserLevelSuper:
		// Super-users can see everything.
		break
	case enums.UserLevelCustomer:
		// Customers can only see themselves.
		q.WhereEq("id", user.ID)
		// q.Where("id", query2.Equals, user.ID)
	default:
		// Admin and normal users can only see users in their own org.
		q.WhereEq("organization_id", user.OrganizationId)
		// q.Where("organization_id", query2.Equals, user.OrganizationId)
	}
	return nil
}

func (r usersRepo) GetById(cx context.Context, id goat.ID, loadRelations ...bool) (*models.User, error) {
	q := r.db.WithContext(cx)
	if len(loadRelations) > 0 {
		q = q.Preload(clause.Associations)
	}
	return repo.First[models.User](q, "id = ?", id)
}

func (r usersRepo) Save(cx context.Context, m *models.User) error {
	if m.Model.ID.Valid() {
		return repo.Update[*models.User](r.db.WithContext(cx), m)
	} else {
		return repo.Create[*models.User](r.db.WithContext(cx), m)
	}
}

func (r usersRepo) Delete(cx context.Context, m *models.User) error {
	return repo.Delete[*models.User](r.db.WithContext(cx), m)
}

func (r usersRepo) FilterStrings(q query.Builder, fields map[string][]string) error {
	for fieldName, values := range fields {
		switch fieldName {
		case "email":
			fallthrough
		case "name":
			if len(values) > 0 {
				q.WhereLike(fieldName, fmt.Sprintf("%%%s%%", values[0]))
				// q.Where(fieldName, query2.Like, fmt.Sprintf("%%%s%%", values[0]))
			}
		}
	}
	return nil
}
