package auth

import (
	"github.com/68696c6c/goat"
	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"

	"github.com/68696c6c/web/app/enums"
)

type Permission string

const (
	All  Permission = "all"
	Org  Permission = "org"
	Self Permission = "self"
	None Permission = "none"
)

type Action string

const (
	Create Action = "create"
	Read   Action = "read"
	Update Action = "update"
	Delete Action = "delete"
)

type ActionPermissions map[Action]Permission

func (a ActionPermissions) Valid() error {
	var errs []error
	if p, ok := a[Create]; !ok || p == "" {
		errs = append(errs, errors.New("invalid create permission"))
	}
	if p, ok := a[Read]; !ok || p == "" {
		errs = append(errs, errors.New("invalid read permission"))
	}
	if p, ok := a[Update]; !ok || p == "" {
		errs = append(errs, errors.New("invalid update permission"))
	}
	if p, ok := a[Delete]; !ok || p == "" {
		errs = append(errs, errors.New("invalid delete permission"))
	}
	return goat.ErrorsToError(errs)
}

type CustomClaims struct {
	Users         ActionPermissions `json:"users"`
	Organizations ActionPermissions `json:"organizations"`
}

func GetCustomClaimsByUserLevel(level enums.UserLevel) CustomClaims {
	var users ActionPermissions
	var orgs ActionPermissions
	switch level {
	case enums.UserLevelSuper:
		// Super admins have full access to everything.
		users = ActionPermissions{
			Create: All,
			Read:   All,
			Update: All,
			Delete: All,
		}
		orgs = ActionPermissions{
			Create: All,
			Read:   All,
			Update: All,
			Delete: All,
		}
	case enums.UserLevelAdmin:
		// Admins can manage everything within their own organization but cannot create new organizations.
		users = ActionPermissions{
			Create: Org,
			Read:   Org,
			Update: Org,
			Delete: Org,
		}
		orgs = ActionPermissions{
			Create: None,
			Read:   Org,
			Update: Org,
			Delete: Org,
		}
	default:
		// Users can only see their own organization and only manage their own account.
		users = ActionPermissions{
			Create: None,
			Read:   Org,
			Update: Self,
			Delete: Self,
		}
		orgs = ActionPermissions{
			Create: None,
			Read:   Org,
			Update: None,
			Delete: None,
		}
	}
	return CustomClaims{
		Users:         users,
		Organizations: orgs,
	}
}

type Claims struct {
	CustomClaims
	jwt.StandardClaims
}

func (c *Claims) Valid() error {
	var errs []error
	err := c.StandardClaims.Valid()
	if err != nil {
		errs = append(errs, err)
	}
	err = c.Users.Valid()
	if err != nil {
		errs = append(errs, errors.Wrap(err, "invalid users permissions"))
	}
	err = c.Organizations.Valid()
	if err != nil {
		errs = append(errs, errors.Wrap(err, "invalid organizations permissions"))
	}
	return goat.ErrorsToError(errs)
}
