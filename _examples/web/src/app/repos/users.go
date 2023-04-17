package repos

import (
	"context"

	"github.com/68696c6c/goat"
	"github.com/68696c6c/goat/query"
	"github.com/68696c6c/goat/repo"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/68696c6c/web/app/enums"
	"github.com/68696c6c/web/app/models"
)

type UsersRepo interface {
	repo.CRUD[*models.User, models.UserRequest]
	GetByEmail(cx context.Context, email string) (*models.User, error)
}

type usersRepo struct {
	db *gorm.DB
}

func NewUsersRepo(db *gorm.DB) UsersRepo {
	return usersRepo{
		db: db,
	}
}

func (r usersRepo) Create(cx context.Context, u models.UserRequest) (*models.User, error) {
	m := models.NewUser()
	var errs []error
	if u.OrganizationId == nil {
		errs = append(errs, goat.NewValidationError("organizationId", "required"))
	} else {
		m.OrganizationID = *u.OrganizationId
	}
	if u.Level == nil {
		errs = append(errs, goat.NewValidationError("level", "required"))
	} else {
		m.Level = *u.Level
	}
	if u.Name == nil {
		errs = append(errs, goat.NewValidationError("name", "required"))
	} else {
		m.Name = *u.Name
	}
	if u.Email == nil {
		errs = append(errs, goat.NewValidationError("email", "required"))
	} else {
		_, err := repo.First[models.User](r.db.WithContext(cx), "email = ?", u.Email)
		if !goat.RecordNotFound(err) {
			errs = append(errs, goat.NewValidationError("email", "unique"))
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
	m, err := r.GetByID(cx, id)
	if err != nil {
		return nil, err
	}
	m.OrganizationID = goat.ValueOrDefault[goat.ID](u.OrganizationId, m.OrganizationID)
	m.Level = goat.ValueOrDefault[enums.UserLevel](u.Level, m.Level)
	m.Name = goat.ValueOrDefault[string](u.Name, m.Name)
	return m, nil
}

func (r usersRepo) Filter(cx context.Context, q query.Builder) ([]*models.User, query.Builder, error) {
	base := r.db.WithContext(cx).Model(&models.User{})
	return repo.Filter[models.User](base, q)
}

func (r usersRepo) GetByID(cx context.Context, id goat.ID, loadRelations ...bool) (*models.User, error) {
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

func (r usersRepo) GetByEmail(cx context.Context, email string) (*models.User, error) {
	q := r.db.WithContext(cx)
	return repo.First[models.User](q, "email = ?", email)
}
