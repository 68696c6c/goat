package goat

import (
	"testing"
	"github.com/icrowley/fake"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"runtime"
)

func setupPathTypeTest(root string) *path {
	exePath := fake.Word()
	return newPath(exePath, nil, root, runtime.Caller)
}

func TestPath_Root_Success(t *testing.T) {
	root := fake.Word()
	p := setupPathTypeTest(root)
	assert.Equal(t, root, p.Root(), "failed to set root")
}

func TestPath_RootDefault_Success(t *testing.T) {
	p := setupPathTypeTest("")
	assert.NotEmpty(t, p.Root(), "failed to set root")
}

func TestPath_Root_Panic(t *testing.T) {
	p := path{}
	defer func() {
		r := recover()
		assert.NotNil(t, r, "Root() failed to panic when the root path wasn't set")
	}()
	p.Root()
}

func TestPath_RootPath(t *testing.T) {
	root := fake.Word()
	p := setupPathTypeTest(root)

	file := fake.Word()
	expected := root + "/" + file

	assert.Equal(t, expected, p.RootPath(file), "RootPath() returned an unexpected value")
}

func TestPath_ExePath(t *testing.T) {
	root := fake.Word()
	exePath := fake.Word()
	p := newPath(exePath, nil, root, runtime.Caller)

	assert.Equal(t, exePath, p.ExePath(), "ExePath() returned an unexpected value")
}

func TestPath_ExeDir(t *testing.T) {
	root := fake.Word()
	exePath := fake.Word()
	exeDir := filepath.Dir(exePath)
	assert.NotEmpty(t, exeDir, "exeDir is empty")

	p := newPath(exePath, nil, root, runtime.Caller)
	assert.Equal(t, exeDir, p.ExeDir(), "ExeDir() returned an unexpected value")
}

func TestPath_CWD(t *testing.T) {
	root := fake.Word()
	p := setupPathTypeTest(root)
	assert.NotEmpty(t, p.CWD(), "CWD() returned an unexpected value")
}

func TestPath_CWD_Fail(t *testing.T) {
	root := fake.Word()
	exePath := fake.Word()

	f := func(int) (uintptr, string, int, bool) {
		return 1, "", 0, false
	}

	p := newPath(exePath, nil, root, f)

	assert.Empty(t, p.CWD(), "CWD() returned an unexpected value")
}
