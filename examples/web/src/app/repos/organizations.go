package repos

import (
	"context"

	"github.com/68696c6c/goat"
	"github.com/68696c6c/goat/query"
	"github.com/68696c6c/goat/repo"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/68696c6c/web/app/models"
)

type OrganizationsRepo interface {
	repo.CRUD[*models.Organization, models.OrganizationRequest]
}

type organizationsRepo struct {
	db *gorm.DB
}

func NewOrganizationsRepo(db *gorm.DB) OrganizationsRepo {
	return organizationsRepo{
		db: db,
	}
}

func (r organizationsRepo) Create(_ context.Context, u models.OrganizationRequest) (*models.Organization, error) {
	result := models.NewOrganization()
	var errs []error
	if u.Name == nil {
		errs = append(errs, goat.NewValidationError("name", "required"))
	} else {
		result.Name = *u.Name
	}
	if u.Website == nil {
		errs = append(errs, goat.NewValidationError("website", "required"))
	} else {
		result.Website = *u.Website
	}
	if len(errs) > 0 {
		return nil, goat.ErrorsToError(errs)
	}
	return result, nil
}

func (r organizationsRepo) Update(cx context.Context, id goat.ID, u models.OrganizationRequest) (*models.Organization, error) {
	m, err := r.GetByID(cx, id)
	if err != nil {
		return nil, err
	}
	m.Name = goat.ValueOrDefault[string](u.Name, m.Name)
	m.Website = goat.ValueOrDefault[string](u.Website, m.Website)
	return m, nil
}

func (r organizationsRepo) Filter(cx context.Context, q query.Builder) ([]*models.Organization, query.Builder, error) {
	base := r.db.WithContext(cx).Model(&models.Organization{})
	return repo.Filter[models.Organization](base, q)
}

func (r organizationsRepo) GetByID(cx context.Context, id goat.ID, loadRelations ...bool) (*models.Organization, error) {
	q := r.db.WithContext(cx)
	if len(loadRelations) > 0 {
		q = q.Preload(clause.Associations)
	}
	return repo.First[models.Organization](q, "id = ?", id)
}

func (r organizationsRepo) Save(cx context.Context, m *models.Organization) error {
	if m.Model.ID.Valid() {
		return repo.Update[*models.Organization](r.db.WithContext(cx), m)
	} else {
		return repo.Create[*models.Organization](r.db.WithContext(cx), m)
	}
}

func (r organizationsRepo) Delete(cx context.Context, m *models.Organization) error {
	return repo.Delete[*models.Organization](r.db.WithContext(cx), m)
}
