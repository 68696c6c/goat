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
	input, err := url.Parse(exampleUsersUrl + "?page=4&pages=6&size=3&total=18")
	require.Nil(t, err, "failed to parse test url")
	result := NewCollection[string]([]string{"a", "b"}, query.NewQueryFromUrl(input.Query()), input)
	assert.Len(t, result.Embeds, 2)
	assert.Equal(t, 4, result.Page)
	assert.Equal(t, 3, result.PageSize)
	assert.Equal(t, 18, result.Total)
	assert.Equal(t, 6, result.TotalPages)
	assert.Equal(t, exampleUsersUrl+"?page=4&pages=6&size=3&total=18", result.Links["self"].Href)
	assert.Equal(t, exampleUsersUrl+"?page=1&pages=6&size=3&total=18", result.Links["first"].Href)
	assert.Equal(t, exampleUsersUrl+"?page=3&pages=6&size=3&total=18", result.Links["previous"].Href)
	assert.Equal(t, exampleUsersUrl+"?page=5&pages=6&size=3&total=18", result.Links["next"].Href)
	assert.Equal(t, exampleUsersUrl+"?page=6&pages=6&size=3&total=18", result.Links["last"].Href)
}

func Test_Collection_Json(t *testing.T) {
	input, err := url.Parse(exampleUsersUrl + "?page=4&size=3&total=18&pages=6")
	require.Nil(t, err, "failed to parse test url")
	resources := []user{
		{ID: 1, Name: "First", ResourceEmbeds: makeUserEmbeds(1, 2)},
		{ID: 2, Name: "Second", ResourceEmbeds: makeUserEmbeds(3, 4)},
		{ID: 3, Name: "Third", ResourceEmbeds: makeUserEmbeds(5, 6)},
	}
	collection := NewCollection[user](resources, query.NewQueryFromUrl(input.Query()), input)
	result, err := json.MarshalIndent(collection, "", "  ")
	require.Nil(t, err, "failed to marshal collection")
	assert.Equal(t, expectedCollectionJson, string(result))
}

const expectedCollectionJson = `{
  "page": 4,
  "size": 3,
  "total": 18,
  "pages": 6,
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
      "href": "https://example.com/users?page=1\u0026pages=6\u0026size=3\u0026total=18"
    },
    "last": {
      "href": "https://example.com/users?page=6\u0026pages=6\u0026size=3\u0026total=18"
    },
    "next": {
      "href": "https://example.com/users?page=5\u0026pages=6\u0026size=3\u0026total=18"
    },
    "previous": {
      "href": "https://example.com/users?page=3\u0026pages=6\u0026size=3\u0026total=18"
    },
    "self": {
      "href": "https://example.com/users?page=4\u0026pages=6\u0026size=3\u0026total=18"
    }
  }
}`
