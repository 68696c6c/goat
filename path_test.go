package goat

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"goat/types"
)

const (
	pathTestMsgInitPath = "initPath() returned an unexpected value"
	pathTestMsgSetRoot  = "failed to set root path"
)

var (
	pathTestContainer *Container
)

// A path test analog of goat.Init().
func pathTestInit() {
	u := types.NewGoatUtils()
	p, _ := initPath(u)

	pathTestContainer = newContainer(u, p, false)
	errs := GetErrors()
	if len(errs) == 0 {
		pathTestContainer.Utils.SetInitialized(true)
		return
	}
	errString := ErrorsToString(errs)
	panic("failed to initialize path test: " + errString)
}

func TestExecutableClean(t *testing.T) {
	path, err := executableClean()
	assert.Nil(t, err, "executableClean returned an error")
	assert.NotEmpty(t, path, "executableClean failed to return a path")
}

func TestInitPath(t *testing.T) {
	pathTestInit()
	assert.NotNil(t, pathTestContainer.Path, pathTestMsgInitPath)
}

func TestSetRoot(t *testing.T) {
	p := "/some/path"
	SetRoot(p)
	pathTestInit()
	assert.Equal(t, p, pathTestContainer.Path.Root(), pathTestMsgSetRoot)
}
