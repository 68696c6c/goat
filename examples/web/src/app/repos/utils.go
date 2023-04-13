package repos

import (
	"fmt"

	"github.com/68696c6c/goat/query"

	"github.com/68696c6c/web/app/enums"
	"github.com/68696c6c/web/app/lib/auth"
	"github.com/68696c6c/web/app/models"
)

func FilterUsersQuery(q query.Builder, fields map[string][]string) {
	for fieldName, values := range fields {
		switch fieldName {
		case "email":
			fallthrough
		case "name":
			if len(values) > 0 {
				q.AndLike(fieldName, fmt.Sprintf("%%%s%%", values[0]))
			}
		}
	}
}

func FilterUsersForUser(q query.Builder, user *models.User, claims *auth.Claims) {
	switch claims.Users[auth.Read] {
	case auth.All:
		break
	case auth.Org:
		q.AndEq("organization_id", user.OrganizationID)
	case auth.Self:
		q.AndEq("id", user.ID)
	default:
		// No access; apply an impossible condition to return no records.
		q.AndNotEq("1", "1")
	}
}

func UserHasUserAccess(targetUser, currentUser *models.User, claims *auth.Claims, action auth.Action) bool {
	switch claims.Users[action] {
	case auth.All:
		return true
	case auth.Org:
		return targetUser.OrganizationID == currentUser.OrganizationID
	case auth.Self:
		return targetUser.ID == currentUser.ID
	default:
		return false
	}
}

func CanUserWriteUserLevel(level enums.UserLevel, currentUser *models.User) bool {
	switch currentUser.Level {
	case enums.UserLevelSuper:
		return true
	case enums.UserLevelAdmin:
		return level != enums.UserLevelSuper
	default:
		return false
	}
}

func FilterOrganizationsQuery(q query.Builder, fields map[string][]string) {
	for fieldName, values := range fields {
		switch fieldName {
		case "website":
			fallthrough
		case "name":
			if len(values) > 0 {
				q.AndLike(fieldName, fmt.Sprintf("%%%s%%", values[0]))
			}
		}
	}
}

func FilterOrganizationsForUser(q query.Builder, user *models.User, claims *auth.Claims) {
	switch claims.Organizations[auth.Read] {
	case auth.All:
		break
	case auth.Org:
		q.AndEq("id", user.OrganizationID)
	default:
		// No access; apply an impossible condition to return no records.
		q.AndNotEq("1", "1")
	}
}

func UserHasOrganizationAccess(targetOrg *models.Organization, currentUser *models.User, claims *auth.Claims, action auth.Action) bool {
	switch claims.Organizations[action] {
	case auth.All:
		return true
	case auth.Org:
		return targetOrg.ID == currentUser.OrganizationID
	default:
		return false
	}
}
