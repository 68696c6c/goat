package request

import (
	"github.com/gin-gonic/gin"
)

// Goals:
// - bind requests to structs with as much type-safety as possible
// - request validation:
// 	- required/ignore for create
//  - required/ignore for update
//  - all go-playground validators

func bind[T any](cx *gin.Context, target T) error {
	err := cx.Bind(target)
	if err != nil {
		return err
	}
	return nil
}
