package request

// import (
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"
//
// 	"github.com/gin-gonic/gin"
// 	"github.com/stretchr/testify/assert"
// )
//
// type example struct {
// 	Name  string `json:"name" binding:"required"`
// 	Email string `json:"email" binding:"required,email"`
// }
//
// // var isEmail validator.Func = func(fl validator.FieldLevel) bool {
// // 	date, ok := fl.Field().Interface().(time.Time)
// // 	if ok {
// // 		today := time.Now()
// // 		if today.After(date) {
// // 			return false
// // 		}
// // 	}
// // 	return true
// // }
//
// type message struct {
// 	Message string `json:"message"`
// }
//
// func setupRouter() *gin.Engine {
// 	r := gin.Default()
// 	r.GET("/ping", func(cx *gin.Context) {
// 		req := example{}
// 		err := bind[example](cx, req)
// 		if err != nil {
// 			cx.AbortWithStatusJSON(http.StatusBadRequest, message{err.Error()})
// 			return
// 		}
// 		cx.AbortWithStatusJSON(http.StatusOK, req)
// 		return
// 	})
// 	return r
// }
//
// // Assert that fields can be marked as required.
// func Test_Required(t *testing.T) {
// 	router := setupRouter()
//
// 	w := httptest.NewRecorder()
// 	req, _ := http.NewRequest("GET", "/ping")
// 	router.ServeHTTP(w, req)
//
// 	assert.Equal(t, http.StatusOK, w.Code)
// 	assert.Equal(t, "pong", w.Body.String())
// }
//
// // // Assert that fields can be marked as required for creation.
// // func Test_Required_Create(t *testing.T) {
// //
// // }
// //
// // // Assert that fields can be marked as required for updating.
// // func Test_Required_Update(t *testing.T) {
// //
// // }
