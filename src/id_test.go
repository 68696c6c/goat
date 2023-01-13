package goat

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_ID_JsonMarshal(t *testing.T) {
	uid := uuid.New()
	i := ID(uid)

	uidStr, err := json.Marshal(uid)
	assert.Nil(t, err)
	idStr, err := json.Marshal(i)
	assert.Nil(t, err)
	assert.Equal(t, uidStr, idStr)
}

func Test_ID_JsonUnmarshal(t *testing.T) {
	str := fmt.Sprintf("\"%s\"", uuid.New().String())

	uid := uuid.UUID{}
	id := ID{}

	err := json.Unmarshal([]byte(str), &uid)
	assert.Nil(t, err)
	err = json.Unmarshal([]byte(str), &id)
	assert.Nil(t, err)

	assert.NotEqual(t, uuid.Nil, uid)
	assert.NotEqual(t, ID(uuid.Nil), id)
	assert.Equal(t, ([16]byte)(uid), ([16]byte)(id))
}

func Test_ID_Scan(t *testing.T) {
	uid := uuid.New()
	idStr := uid.String()

	id := ID{}
	err := id.Scan(idStr)
	assert.Nil(t, err)
	assert.Equal(t, uid.String(), id.String())

	id = ID{}
	err = id.Scan([]byte(idStr))
	assert.Nil(t, err)
	assert.Equal(t, uid.String(), id.String())

	id = ID{}
	err = id.Scan(1)
	assert.NotNil(t, err)
}

func Test_ID_Value(t *testing.T) {
	id := ID(uuid.New())
	v, err := id.Value()
	b, ok := v.([]byte)

	assert.Nil(t, err)
	assert.True(t, ok)
	assert.Len(t, b, 16)
}

func Test_ID_Nil(t *testing.T) {
	uid := uuid.UUID{}
	id := ID{}

	assert.False(t, id.Valid())
	assert.True(t, ID(uid) == id)
	assert.True(t, uid == uuid.UUID(id))
}
