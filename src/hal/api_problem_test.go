package hal

import (
	"net/http"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func Test_NewApiProblem(t *testing.T) {
	result := NewApiProblem(http.StatusOK, errors.New("example"))
	assert.Equal(t, httpStatusRfcLinkMap[http.StatusOK], result.DescribedBy, "unexpected description")
	assert.Equal(t, http.StatusText(http.StatusOK), result.Title, "unexpected title")
	assert.Equal(t, http.StatusOK, result.HttpStatus, "unexpected status")
	assert.Equal(t, "example", result.Details, "unexpected details")
}

func Test_NewApiProblem_InvalidStatus(t *testing.T) {
	result := NewApiProblem(1, errors.New("example"))
	assert.Equal(t, httpStatusRfcLinkMap[http.StatusInternalServerError], result.DescribedBy, "unexpected description")
	assert.Equal(t, http.StatusText(http.StatusInternalServerError), result.Title, "unexpected title")
	assert.Equal(t, http.StatusInternalServerError, result.HttpStatus, "unexpected status")
	assert.Equal(t, "example", result.Details, "unexpected details")
}

func Test_NewApiProblem_NilError(t *testing.T) {
	result := NewApiProblem(http.StatusNotFound, nil)
	assert.Equal(t, httpStatusRfcLinkMap[http.StatusNotFound], result.DescribedBy, "unexpected description")
	assert.Equal(t, http.StatusText(http.StatusNotFound), result.Title, "unexpected title")
	assert.Equal(t, http.StatusNotFound, result.HttpStatus, "unexpected status")
	assert.Equal(t, "", result.Details, "unexpected details")
}
