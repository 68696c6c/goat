package goat

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ErrorResponder func(cx *gin.Context, e error)

type ErrorHandler interface {
	HandleError(cx *gin.Context, err error, responder ErrorResponder)
}

type ErrorHandlerGin struct {
	logger *logrus.Logger
}

func NewErrorHandler(l *logrus.Logger) ErrorHandler {
	return ErrorHandlerGin{
		logger: l,
	}
}

func (h ErrorHandlerGin) HandleError(cx *gin.Context, err error, responder ErrorResponder) {
	h.logger.Error(fmt.Sprintf("%s | %s", cx.HandlerName(), err))
	responder(cx, err)
	return
}

// TODO: build this into goat instead of a global???
var _errorHandler ErrorHandler

func InitErrorHandler() ErrorHandler {
	if _errorHandler == nil {
		_errorHandler = NewErrorHandler(GetLogger())
	}
	return _errorHandler
}

func HandleErrorResponse(cx *gin.Context, err error, responder ErrorResponder) {
	e := InitErrorHandler()
	e.HandleError(cx, err, responder)
}
