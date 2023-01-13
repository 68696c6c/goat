package repos

import (
	"context"
	"fmt"

	"github.com/68696c6c/goat"
	"github.com/68696c6c/goat/query2"
	"github.com/68696c6c/goat/repo"
	"github.com/68696c6c/goat/resource"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/68696c6c/example/app/enums"
	"github.com/68696c6c/example/app/models"
)

type UsersRepo interface {
	repo.CRUD[*models.User, models.UserRequest]
	// ValidateCreate(validator validator.StructLevel)
	ApplyFilterForUser(q query2.Builder, user *models.User) error
	FilterStrings(query query2.Builder, fields map[string][]string) error
	// UserLevelFilter(user models.User) repo.QueryFilter
	// SearchFilter(fields url.Values) repo.QueryFilter
	// FilterTemp(cx context.Context, p query.Pagination, scopes ...repo.QueryScope) ([]*models.User, query.Pagination, error)
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

// func (r usersRepo) ValidateCreate(validator validator.StructLevel) {
// 	req := validator.Current().Interface().(models.UserRequest)
//
// 	if req.Email == nil {
// 		validator.ReportError(req.Email, "email", "", "required", "")
// 	}
//
// 	_, err := repo.First[*models.User](r.db, &models.User{}, "email = ?", req.Email)
// 	if !goat.RecordNotFound(err) {
// 		validator.ReportError(req.Email, "email", "", "unique", "")
// 	}
// }
//
// func (r usersRepo) ValidateUpdate(validator validator.StructLevel) {
// 	req := validator.Current().Interface().(models.UserRequest)
//
// 	if req.Email == nil {
// 		validator.ReportError(req.Email, "email", "", "required", "")
// 	}
//
// 	_, err := repo.First[*models.User](r.db, &models.User{}, "email = ?", req.Email)
// 	if !goat.RecordNotFound(err) {
// 		validator.ReportError(req.Email, "email", "", "unique", "")
// 	}
// }

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

// func (r usersRepo) Filter(cx context.Context, q query.Builder) ([]*models.User, error) {
//
// 	return repo.FilterTemp[models.User](q, r.db.WithContext(cx).Model(&models.User{}))
// 	// return repo.Filter[models.User](q, r.getBaseQuery)
// }

// func (r usersRepo) Filter(cx context.Context, p resource.Pagination, filters ...repo.QueryFilter) ([]*models.User, resource.Pagination, error) {
func (r usersRepo) Filter(cx context.Context, q query2.Builder, p resource.Pagination) ([]*models.User, resource.Pagination, error) {
	base := r.db.WithContext(cx).Model(&models.User{})
	return repo.Filter[models.User](base, q, p)
	// // for _, scope := range scopes {
	// // 	base.Scopes(scope)
	// // }
	// result, pagination, err := repo.PaginatedFilter[models.User](base, p, filters...)
	// if err != nil {
	// 	return result, pagination, err
	// }
	// // pp := query.NewPaginationFromValues(p.Page, p.PageSize, count)
	// // return repo.Filter[models.User](q, r.getBaseQuery)
	// return result, pagination, nil
}

// func (r usersRepo) UserLevelFilter(user models.User) repo.QueryFilter {
// 	return func(db *gorm.DB) *gorm.DB {
// 		switch user.Level {
// 		case enums.UserLevelSuper:
// 			// Super-users can see everything.
// 			break
// 		case enums.UserLevelCustomer:
// 			// Customers can only see themselves.
// 			db.Where("id = ?", user.ID)
// 		default:
// 			// Admin and normal users can only see users in their own org.
// 			db.Where("organization_id = ?", user.OrganizationId)
// 		}
// 		return db
// 	}
// }

func (r usersRepo) ApplyFilterForUser(q query2.Builder, user *models.User) error {
	switch user.Level {
	case enums.UserLevelSuper:
		// Super-users can see everything.
		break
	case enums.UserLevelCustomer:
		// Customers can only see themselves.
		// q.WhereEq("id", user.ID)
		q.Where("id", query2.Equals, user.ID)
	default:
		// Admin and normal users can only see users in their own org.
		// q.WhereEq("organization_id", user.OrganizationId)
		q.Where("organization_id", query2.Equals, user.OrganizationId)
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

func (r usersRepo) FilterStrings(q query2.Builder, fields map[string][]string) error {
	for fieldName, values := range fields {
		switch fieldName {
		case "email":
			fallthrough
		case "name":
			if len(values) > 0 {
				// q.WhereLike(fieldName, fmt.Sprintf("%%%s%%", values[0]))
				q.Where(fieldName, query2.Like, fmt.Sprintf("%%%s%%", values[0]))
			}
		}
	}
	return nil
}

// func filterStrings(db *gorm.DB, fields url.Values) *gorm.DB {
// 	for fieldName, values := range fields {
// 		switch fieldName {
// 		case "email":
// 			fallthrough
// 		case "name":
// 			if len(values) > 0 {
// 				cond := fmt.Sprintf("%s LIKE ?", fieldName)
// 				db.Or(cond, fmt.Sprintf("%%%s%%", values[0]))
// 				// f.WhereField(fieldName, filter.OpLike, fmt.Sprintf("%%%s%%", values[0]))
// 			}
// 		}
// 	}
// 	return db
// }
//
// func (r usersRepo) SearchFilter(fields url.Values) repo.QueryFilter {
// 	return func(db *gorm.DB) *gorm.DB {
// 		if len(fields) == 0 {
// 			return db
// 		}
// 		return db.Where(filterStrings(db, fields))
// 	}
// }
//
// type tempFilter struct {
// 	condition string
// 	value     string
// }
//
// func (r usersRepo) filterStrings(db *gorm.DB, fields url.Values) *gorm.DB {
// 	// var result []tempFilter
// 	// q := db.Unscoped()
// 	var set bool
// 	for fieldName, values := range fields {
// 		switch fieldName {
// 		case "email":
// 			fallthrough
// 		case "name":
// 			if len(values) > 0 {
// 				// result = append(result, tempFilter{
// 				// 	condition: fmt.Sprintf("%s LIKE ?", fieldName),
// 				// 	value:     fmt.Sprintf("%%%s%%", values[0]),
// 				// })
// 				cond := fmt.Sprintf("%s LIKE ?", fieldName)
// 				value := fmt.Sprintf("%%%s%%", values[0])
// 				if set {
// 					db.Or(cond, value)
// 				} else {
// 					set = true
// 					db.Where(cond, value)
// 				}
// 				// // f.WhereField(fieldName, filter.OpLike, fmt.Sprintf("%%%s%%", values[0]))
// 			}
// 		}
// 	}
// 	return db
// }
//
// func (r usersRepo) filterStrings(fields url.Values) repo.QueryFilter {
// 	return func(db *gorm.DB) *gorm.DB {
// 		// db.Where("1 = 1")
// 		for fieldName, values := range fields {
// 			switch fieldName {
// 			case "email":
// 				fallthrough
// 			case "name":
// 				if len(values) > 0 {
// 					cond := fmt.Sprintf("%s LIKE ?", fieldName)
// 					db.Or(cond, fmt.Sprintf("%%%s%%", values[0]))
// 					// f.WhereField(fieldName, filter.OpLike, fmt.Sprintf("%%%s%%", values[0]))
// 				}
// 			}
// 		}
// 		return db
// 		// where, params, err := f.Apply()
// 		// if err != nil {
// 		// 	// return nil, errors.Wrap(err, "failed to apply filter")
// 		// }
// 		// var filters []filter
// 		// db.Where("1 = 1")
// 		// for fieldName, values := range fields {
// 		// 	switch fieldName {
// 		// 	case "email":
// 		// 		fallthrough
// 		// 	case "name":
// 		// 		if len(values) > 0 {
// 		// 			filters = append(filters, filter{
// 		// 				field: fmt.Sprintf("%s LIKE ?", fieldName),
// 		// 				value: fmt.Sprintf("%%%s%%", values[0]),
// 		// 			})
// 		// 			// cond := fmt.Sprintf("%s LIKE ?", fieldName)
// 		// 			// db.Where(cond, fmt.Sprintf("%%%s%%", values[0]))
// 		// 		}
// 		// 	}
// 		// }
// 		// if len(filters) > 0 {
// 		// 	for i, f := range filters {
// 		//
// 		// 	}
// 		// }
// 		// return db
// 	}
// }

// func temp(db *gorm.DB, filters []tempFilter) *gorm.DB {
// 	for _, filter := range filters {
// 		db.Or(filter.condition, filter.value)
// 	}
// }

// func (r usersRepo) SearchFilter(fields url.Values) repo.QueryFilter {
// 	return func(db *gorm.DB) *gorm.DB {
// 		if len(fields) == 0 {
// 			return db
// 		}
// 		// conditions := r.filterStrings(db, fields)
// 		// d := db.Unscoped()
// 		return db.Where(r.filterStrings(db, fields))
// 	}
// }
