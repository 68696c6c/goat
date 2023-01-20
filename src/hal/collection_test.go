package hal

import (
	"encoding/json"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/68696c6c/goat/query"
)

func Test_NewCollection(t *testing.T) {
	input, err := url.Parse(exampleUsersUrl + "?page=4&pageSize=3&total=18&totalPages=6")
	require.Nil(t, err, "failed to parse test url")
	result := NewCollection[string]([]string{"a", "b"}, query.NewQueryFromUrl(input.Query()), input)
	assert.Len(t, result.Embeds, 2)
	assert.Equal(t, 4, result.Page)
	assert.Equal(t, 3, result.PageSize)
	assert.Equal(t, 18, result.Total)
	assert.Equal(t, 6, result.TotalPages)
	assert.Equal(t, exampleUsersUrl+"?page=4&pageSize=3&total=18&totalPages=6", result.Links["self"].Href)
	assert.Equal(t, exampleUsersUrl+"?page=1&pageSize=3&total=18&totalPages=6", result.Links["first"].Href)
	assert.Equal(t, exampleUsersUrl+"?page=3&pageSize=3&total=18&totalPages=6", result.Links["previous"].Href)
	assert.Equal(t, exampleUsersUrl+"?page=5&pageSize=3&total=18&totalPages=6", result.Links["next"].Href)
	assert.Equal(t, exampleUsersUrl+"?page=6&pageSize=3&total=18&totalPages=6", result.Links["last"].Href)
}

func Test_Collection_Json(t *testing.T) {
	input, err := url.Parse(exampleUsersUrl + "?page=4&pageSize=3&total=18&totalPages=6")
	require.Nil(t, err, "failed to parse test url")
	resources := []user{
		{Id: 1, Name: "First", ResourceEmbeds: makeUserEmbeds(1, 2)},
		{Id: 2, Name: "Second", ResourceEmbeds: makeUserEmbeds(3, 4)},
		{Id: 3, Name: "Third", ResourceEmbeds: makeUserEmbeds(5, 6)},
	}
	collection := NewCollection[user](resources, query.NewQueryFromUrl(input.Query()), input)
	result, err := json.MarshalIndent(collection, "", "  ")
	require.Nil(t, err, "failed to marshal collection")
	assert.Equal(t, expectedCollectionJson, string(result))
}

const expectedCollectionJson = `{
  "page": 4,
  "pageSize": 3,
  "total": 18,
  "totalPages": 6,
  "_embedded": [
    {
      "id": 1,
      "name": "First",
      "_embedded": {
        "phones": [
          {
            "id": 1,
            "phone": "1111111111"
          },
          {
            "id": 2,
            "phone": "2222222222"
          }
        ]
      },
      "_links": {
        "self": {
          "href": "https://example.com/users/1"
        }
      }
    },
    {
      "id": 2,
      "name": "Second",
      "_embedded": {
        "phones": [
          {
            "id": 3,
            "phone": "3333333333"
          },
          {
            "id": 4,
            "phone": "4444444444"
          }
        ]
      },
      "_links": {
        "self": {
          "href": "https://example.com/users/2"
        }
      }
    },
    {
      "id": 3,
      "name": "Third",
      "_embedded": {
        "phones": [
          {
            "id": 5,
            "phone": "5555555555"
          },
          {
            "id": 6,
            "phone": "6666666666"
          }
        ]
      },
      "_links": {
        "self": {
          "href": "https://example.com/users/3"
        }
      }
    }
  ],
  "_links": {
    "first": {
      "href": "https://example.com/users?page=1\u0026pageSize=3\u0026total=18\u0026totalPages=6"
    },
    "last": {
      "href": "https://example.com/users?page=6\u0026pageSize=3\u0026total=18\u0026totalPages=6"
    },
    "next": {
      "href": "https://example.com/users?page=5\u0026pageSize=3\u0026total=18\u0026totalPages=6"
    },
    "previous": {
      "href": "https://example.com/users?page=3\u0026pageSize=3\u0026total=18\u0026totalPages=6"
    },
    "self": {
      "href": "https://example.com/users?page=4\u0026pageSize=3\u0026total=18\u0026totalPages=6"
    }
  }
}`
