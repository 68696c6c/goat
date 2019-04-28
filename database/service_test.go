package database

import (
	"strings"
	"testing"

	"github.com/icrowley/fake"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewServiceGORM_GetMainDBName_Default(t *testing.T) {
	s := NewServiceGORM("")
	sn := s.GetMainDBName()

	assert.Equal(t, mainDBNameDefault, sn, "unexpected default main connection name")
}

func TestNewServiceGORM_GetMainDBName_Custom(t *testing.T) {
	n := fake.Word()
	s := NewServiceGORM(n)
	sn := s.GetMainDBName()

	assert.Equal(t, n, sn, "unexpected custom main connection name")
}

func TestNewServiceGORM_GetMainDB(t *testing.T) {
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	s := NewServiceGORM("db")
	d, err := s.GetMainDB()

	require.Nil(t, err, "unexpected error returned")
	assert.NotNil(t, d, "nil database connection returned")
}

func TestNewServiceGORM_GetCustomDB(t *testing.T) {
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	s := NewServiceGORM("db")
	d, err := s.GetCustomDB("db")

	require.Nil(t, err, "unexpected error returned")
	assert.NotNil(t, d, "nil database connection returned")
}
