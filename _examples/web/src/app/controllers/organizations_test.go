package controllers

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

// TODO: add tests for the remaining endpoints.

func Test_OrganizationsController_List(t *testing.T) {
	w, tcx := requireTestResponseAndContext(t)

	c := NewOrganizationsController(tc.OrganizationsRepo, tc.Auth)
	c.List(tcx)
	require.Equal(t, http.StatusOK, w.Code)
}
