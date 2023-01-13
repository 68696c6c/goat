package goat

// import (
// 	"net/http"
// 	"testing"
//
// 	"github.com/gin-gonic/gin"
// 	"github.com/gin-gonic/gin/binding"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"
// )
//
// var initialized bool
// var _router Router
//
// type httpTestModel struct {
// 	ID   int    `json:"id"`
// 	Name string `json:"name"`
// }
//
// type responseItem struct {
// 	httpTestModel
// }
//
// type responseList struct {
// 	Data []httpTestModel `json:"data"`
// }
//
// func getTestRouter(setRoutes func(Router)) Router {
// 	if !initialized {
// 		Init()
// 		_router = InitRouter()
// 		setRoutes(_router)
// 		initialized = true
// 	}
// 	return _router
// }
//
// func setTestRoutes(r Router) {
// 	r.GET("/ping", func(cx *gin.Context) {
// 		RespondOk(cx, "pong")
// 	})
// 	r.GET("/headers", func(cx *gin.Context) {
// 		key := cx.Request.Header.Get("Authorization")
// 		m := httpTestModel{1, key}
// 		RespondOk(cx, responseItem{m})
// 	})
// 	r.POST("/post", func(cx *gin.Context) {
// 		m := httpTestModel{}
// 		if err := cx.ShouldBindWith(&m, binding.JSON); err != nil {
// 			RespondBadRequest(cx, err)
// 			return
// 		}
// 		RespondOk(cx, responseItem{m})
// 	})
// 	r.PUT("/put", func(cx *gin.Context) {
// 		m := httpTestModel{}
// 		if err := cx.ShouldBindWith(&m, binding.JSON); err != nil {
// 			RespondBadRequest(cx, err)
// 			return
// 		}
// 		RespondOk(cx, responseItem{m})
// 	})
// 	r.GET("/get/:id", func(cx *gin.Context) {
// 		m := httpTestModel{1, "model"}
// 		RespondOk(cx, responseItem{m})
// 	})
// 	r.GET("/list", func(cx *gin.Context) {
// 		data := []httpTestModel{
// 			{1, "model 1"},
// 			{2, "model 2"},
// 		}
// 		RespondOk(cx, responseList{data})
// 	})
// 	r.DELETE("/delete", func(cx *gin.Context) {
// 		RespondOk(cx, "ok")
// 	})
// }
//
// func TestHandlerTest_Simple(t *testing.T) {
// 	tr := getTestRouter(setTestRoutes)
//
// 	h := NewHandlerTest(tr)
// 	r := h.Get("/ping").Send()
//
// 	assert.Equal(t, http.StatusOK, r.Code, "unexpected status code")
// }
//
// func TestHandlerTest_Headers(t *testing.T) {
// 	tr := getTestRouter(setTestRoutes)
//
// 	c := NewHandlerTest(tr)
// 	h := map[string]string{
// 		"Authorization": "value",
// 	}
// 	r := c.Get("/headers").Headers(h).Send()
// 	response := &responseItem{}
// 	err := r.Map(response)
// 	require.Nil(t, err, "failed to unmarshal response body")
// 	assert.Equal(t, http.StatusOK, r.Code, "unexpected status code")
// 	assert.Equal(t, response.Name, "value", "unexpected result")
// }
//
// func TestHandlerTest_Post(t *testing.T) {
// 	tr := getTestRouter(setTestRoutes)
//
// 	c := NewHandlerTest(tr)
// 	data := &map[string]any{
// 		"id":   1,
// 		"name": "example",
// 	}
// 	r := c.Post("/post").Body(data).Send()
// 	require.Equal(t, http.StatusOK, r.Code, "unexpected status code, body: "+r.BodyString)
//
// 	response := &responseItem{}
// 	err := r.Map(response)
// 	assert.Nil(t, err, "failed to unmarshal response body")
// 	assert.Equal(t, response.Name, "example", "unexpected result")
// }
//
// func TestHandlerTest_Get(t *testing.T) {
// 	tr := getTestRouter(setTestRoutes)
//
// 	c := NewHandlerTest(tr)
// 	r := c.GetF("/get/%v", 1).Send()
// 	response := &responseItem{}
// 	err := r.Map(response)
// 	require.Nil(t, err, "failed to unmarshal response body")
// 	assert.Equal(t, http.StatusOK, r.Code, "unexpected status code")
// }
//
// func TestHandlerTest_List(t *testing.T) {
// 	tr := getTestRouter(setTestRoutes)
//
// 	c := NewHandlerTest(tr)
// 	r := c.Get("/list").Send()
// 	require.Equal(t, http.StatusOK, r.Code, "unexpected status code")
// 	response := &responseList{}
// 	err := r.Map(response)
// 	require.Nil(t, err, "failed to unmarshal response body")
// 	assert.Len(t, response.Data, 2, "unexpected response length")
// }
//
// func TestHandlerTest_Put(t *testing.T) {
// 	tr := getTestRouter(setTestRoutes)
//
// 	c := NewHandlerTest(tr)
// 	data := &map[string]any{
// 		"id":   1,
// 		"name": "example",
// 	}
// 	r := c.Put("/put").Body(data).Send()
// 	require.Equal(t, http.StatusOK, r.Code, "unexpected status code, body: "+r.BodyString)
//
// 	response := &responseItem{}
// 	err := r.Map(response)
// 	assert.Nil(t, err, "failed to unmarshal response body")
// 	assert.Equal(t, response.Name, "example", "unexpected result")
// }
//
// func TestHandlerTest_Delete(t *testing.T) {
// 	tr := getTestRouter(setTestRoutes)
//
// 	c := NewHandlerTest(tr)
// 	r := c.Delete("/delete").Send()
// 	assert.Equal(t, http.StatusOK, r.Code, "unexpected status code")
// }
