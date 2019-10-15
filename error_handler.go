package goat

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type ContextResponder func(c *gin.Context)
type ErrorResponder func(c *gin.Context, e error)

type ErrorHandler interface {
	HandleMessage(c *gin.Context, m string, responder ErrorResponder)
	HandleError(c *gin.Context, err error, responder ErrorResponder)
	HandleErrorM(c *gin.Context, err error, m string, responder ErrorResponder)
	HandleErrors(c *gin.Context, errs []error, responder ErrorResponder)
	HandleErrorsM(c *gin.Context, errs []error, m string, responder ErrorResponder)
}

type ErrorHandlerGin struct {
	logger *logrus.Logger
}

func NewErrorHandler(l *logrus.Logger) ErrorHandler {
	return ErrorHandlerGin{
		logger: l,
	}
}

func (h ErrorHandlerGin) HandleContext(c *gin.Context, m string, responder ContextResponder) {
	err := errors.New(m)
	h.logger.Error(fmt.Sprintf("%s | %s", c.HandlerName(), err))
	responder(c)
	return
}

func (h ErrorHandlerGin) HandleMessage(c *gin.Context, m string, responder ErrorResponder) {
	err := errors.New(m)
	h.logger.Error(fmt.Sprintf("%s | %s", c.HandlerName(), err))
	responder(c, err)
	return
}

func (h ErrorHandlerGin) HandleError(c *gin.Context, err error, responder ErrorResponder) {
	h.logger.Error(fmt.Sprintf("%s | %s", c.HandlerName(), err))
	responder(c, err)
	return
}

func (h ErrorHandlerGin) HandleErrorM(c *gin.Context, err error, m string, responder ErrorResponder) {
	e := errors.Wrap(err, m)
	h.logger.Error(fmt.Sprintf("%s | %s", c.HandlerName(), e))
	responder(c, e)
	return
}

func (h ErrorHandlerGin) HandleErrors(c *gin.Context, errs []error, responder ErrorResponder) {
	e := ErrorsToError(errs)
	h.logger.Error(fmt.Sprintf("%s | %s", c.HandlerName(), e))
	responder(c, e)
	return
}

func (h ErrorHandlerGin) HandleErrorsM(c *gin.Context, errs []error, m string, responder ErrorResponder) {
	e := errors.Wrap(ErrorsToError(errs), m)
	h.logger.Error(fmt.Sprintf("%s | %s", c.HandlerName(), e))
	responder(c, e)
	return
}
