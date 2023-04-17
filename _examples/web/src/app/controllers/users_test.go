package controllers

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

// TODO: add tests for the remaining endpoints.

func Test_UsersController_List(t *testing.T) {
	w, tcx := requireTestResponseAndContext(t)

	c := NewUsersController(tc.UsersRepo, tc.Auth)
	c.List(tcx)
	require.Equal(t, http.StatusOK, w.Code)
}
