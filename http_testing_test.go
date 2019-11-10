package goat

import (
	"net/http"
	"os"
	"testing"

	ghttp "github.com/68696c6c/goat/src/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var initialized bool

type httpTestModel struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type responseItem struct {
	httpTestModel
}

type responseList struct {
	Data []httpTestModel `json:"data"`
}

func getTestRouter(t *testing.T) ghttp.Router {
	if !initialized {
		err := os.Setenv("ENV", "test")
		require.Nil(t, err, "failed to set env variable 'ENV'")

		Init()
		initialized = true
	}
	r := GetRouter()
	setTestRoutes(r)
	return r
}

func setTestRoutes(r ghttp.Router) {
	engine := r.GetEngine()
	engine.GET("/ping", func(c *gin.Context) {
		RespondMessage(c, "pong")
	})
	engine.GET("/headers", func(c *gin.Context) {
		key := c.Request.Header.Get("Authorization")
		m := httpTestModel{1, key}
		RespondData(c, responseItem{m})
	})
	engine.POST("/post", func(c *gin.Context) {
		m := httpTestModel{}
		if err := c.ShouldBindWith(&m, binding.JSON); err != nil {
			RespondBadRequestError(c, err)
			return
		}
		RespondData(c, responseItem{m})
	})
	engine.PUT("/put", func(c *gin.Context) {
		m := httpTestModel{}
		if err := c.ShouldBindWith(&m, binding.JSON); err != nil {
			RespondBadRequestError(c, err)
			return
		}
		RespondData(c, responseItem{m})
	})
	engine.GET("/get/:id", func(c *gin.Context) {
		m := httpTestModel{1, "model"}
		RespondData(c, responseItem{m})
	})
	engine.GET("/list", func(c *gin.Context) {
		data := []httpTestModel{
			{1, "model 1"},
			{2, "model 2"},
		}
		RespondData(c, responseList{data})
	})
	engine.DELETE("/delete", func(c *gin.Context) {
		RespondMessage(c, "ok")
	})
}

func TestHandlerTest_Simple(t *testing.T) {
	tr := getTestRouter(t)

	h := NewHandlerTest(tr)
	r := h.Get("/ping").Send()

	assert.Equal(t, http.StatusOK, r.Code, "unexpected status code")
}

func TestHandlerTest_Headers(t *testing.T) {
	tr := getTestRouter(t)

	c := NewHandlerTest(tr)
	h := map[string]string{
		"Authorization": "value",
	}
	r := c.Get("/headers").Headers(h).Send()
	response := &responseItem{}
	err := r.Map(response)
	require.Nil(t, err, "failed to unmarshal response body")
	assert.Equal(t, http.StatusOK, r.Code, "unexpected status code")
	assert.Equal(t, response.Name, "value", "unexpected result")
}

func TestHandlerTest_Post(t *testing.T) {
	tr := getTestRouter(t)

	c := NewHandlerTest(tr)
	data := &map[string]interface{}{
		"id":   1,
		"name": "example",
	}
	r := c.Post("/post").Body(data).Send()
	require.Equal(t, http.StatusOK, r.Code, "unexpected status code, body: "+r.BodyString)

	response := &responseItem{}
	err := r.Map(response)
	assert.Nil(t, err, "failed to unmarshal response body")
	assert.Equal(t, response.Name, "example", "unexpected result")
}

func TestHandlerTest_Get(t *testing.T) {
	tr := getTestRouter(t)

	c := NewHandlerTest(tr)
	r := c.GetF("/get/%v", 1).Send()
	response := &responseItem{}
	err := r.Map(response)
	require.Nil(t, err, "failed to unmarshal response body")
	assert.Equal(t, http.StatusOK, r.Code, "unexpected status code")
}

func TestHandlerTest_List(t *testing.T) {
	tr := getTestRouter(t)

	c := NewHandlerTest(tr)
	r := c.Get("/list").Send()
	require.Equal(t, http.StatusOK, r.Code, "unexpected status code")
	response := &responseList{}
	err := r.Map(response)
	require.Nil(t, err, "failed to unmarshal response body")
	assert.Len(t, response.Data, 2, "unexpected response length")
}

func TestHandlerTest_Put(t *testing.T) {
	tr := getTestRouter(t)

	c := NewHandlerTest(tr)
	data := &map[string]interface{}{
		"id":   1,
		"name": "example",
	}
	r := c.Put("/put").Body(data).Send()
	require.Equal(t, http.StatusOK, r.Code, "unexpected status code, body: "+r.BodyString)

	response := &responseItem{}
	err := r.Map(response)
	assert.Nil(t, err, "failed to unmarshal response body")
	assert.Equal(t, response.Name, "example", "unexpected result")
}

func TestHandlerTest_Delete(t *testing.T) {
	tr := getTestRouter(t)

	c := NewHandlerTest(tr)
	r := c.Delete("/delete").Send()
	assert.Equal(t, http.StatusOK, r.Code, "unexpected status code")
}
