package goat

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"

	"github.com/68696c6c/goat/query"
	"github.com/68696c6c/goat/sys/utils"
)

const (
	contextKeyRequest = "goat_request"
	contextKeyQuery   = "goat_query"
)

func getContextValue(cx *gin.Context, key string) any {
	r, exists := cx.Get(key)
	if !exists {
		return nil
	}
	return r
}

// BindMiddleware validates an incoming request and binds the request body to the provided
// struct if the validation passes.
//
// Returns a 400 error with validation errors if binding fails.
//
// Sets the bound request as an any in the Gin registry if binding
// succeeds.  You can retrieve in your handlers it like this:
//
// req, err := goat.GetRequest[*MyRequestType](cx)
//
// This middleware allows you to annotate your request struct fields with
// `binding:"required"` to make required fields.
//
// TODO: update the above comment
func BindMiddleware[T any]() gin.HandlerFunc {
	return func(cx *gin.Context) {
		var obj T
		if err := cx.ShouldBindWith(&obj, binding.JSON); err != nil {
			respondRequestBindingError(cx, err)
			return
		}
		cx.Set(contextKeyRequest, obj)
		return
	}
}

func BindRequest[T any](cx *gin.Context) (T, error) {
	var result T
	err := cx.ShouldBindWith(&result, binding.JSON)
	if err != nil {
		return result, err
	}
	return result, nil
}

// GetRequest returns the bound request struct from the provided Gin context or nil if a goat request has not been bound.
// After binding a request using BindMiddleware, call this function to retrieve it in your handler:
//
//	req, err := goat.GetRequest[*MyRequestType](cx)
func GetRequest[T any](cx *gin.Context) (T, error) {
	result, ok := getContextValue(cx, contextKeyRequest).(T)
	if !ok {
		return result, errors.New("failed to get request")
	}
	return result, nil
}

type ParamParser[T any] func(string) (T, error)

func ParseParam[T any](cx *gin.Context, key string, parser ParamParser[T]) (T, error) {
	param := cx.Param(key)
	result, err := parser(param)
	if err != nil {
		return result, errors.Wrapf(err, "failed to parse request param '%s' from value '%s'", key, param)
	}
	return result, nil
}

func GetIdParam(cx *gin.Context) (ID, error) {
	return ParseParam[ID](cx, "id", ParseID)
}

func respondRequestBindingError(c *gin.Context, err error) {
	// Cast the error to a ValidationErrors struct so that we can access detailed
	// information about the error.  If the error cannot be cast, a generic
	// Bad Request response will be sent instead of the detailed message.
	ve, ok := err.(validator.ValidationErrors)
	if !ok {
		RespondBadRequest(c, err)
		return
	}

	var errs []string
	for _, e := range ve {
		errs = append(errs, utils.MakeValidationError(e.Field(), e.Tag()).Error())
	}
	RespondValidationError(c, errors.New(strings.Join(errs, "; ")))
	return
}

func FilterMiddleware() gin.HandlerFunc {
	return func(cx *gin.Context) {
		q := query.NewQueryBuilder(cx)
		cx.Set(contextKeyQuery, q)
		return
	}
}

// TODO: since we need to cast the query here, we probably won't be able to use the query.Builder interface afterall...
func GetFilter(cx *gin.Context) (query.Builder, error) {
	result, ok := getContextValue(cx, contextKeyQuery).(query.Builder)
	if !ok {
		return result, errors.New("failed to get filter")
	}
	return result, nil
	// r := s.getByKey(cx, contextKeyQuery)
	// return r.(*query.Query)
}

// // Goals:
// // - bind requests to structs with as much type-safety as possible
// // - request validation:
// // 	- required/ignore for create
// //  - required/ignore for update
// //  - all go-playground validators
//
// func bind[T any](cx *gin.Context, target T) error {
// 	err := cx.ShouldBindWith(target, binding.JSON)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
