package goat

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/68696c6c/goat/hal"
)

type ApiProblem hal.ApiProblem

type Resource hal.Resource

type Collection[T any] hal.Collection[T]

func debugError(err error) error {
	if g.HTTP.DebugEnabled() {
		return err
	}
	return nil
}

func logHandlerError(cx *gin.Context, err error) {
	GetLogger().With("error", err, "handler", cx.HandlerName).Errorf("%s | %s", cx.HandlerName(), err)
}

func logHandlerWarn(cx *gin.Context, err error) {
	GetLogger().With("error", err, "handler", cx.HandlerName).Warnf("%s | %s", cx.HandlerName(), err)
}

func respondRequestError(cx *gin.Context, err error, status int) {
	logHandlerWarn(cx, err)
	cx.AbortWithStatusJSON(status, hal.NewApiProblem(status, debugError(err)))
}

func respondServerError(cx *gin.Context, err error, status int) {
	logHandlerError(cx, err)
	cx.AbortWithStatusJSON(status, hal.NewApiProblem(status, debugError(err)))
}

// 1XX responses

func RespondContinue(cx *gin.Context) {
	cx.AbortWithStatus(http.StatusContinue)
}

func RespondSwitchingProtocols(cx *gin.Context) {
	cx.AbortWithStatus(http.StatusSwitchingProtocols)
}

func RespondProcessing(cx *gin.Context) {
	cx.AbortWithStatus(http.StatusProcessing)
}

func RespondEarlyHints(cx *gin.Context) {
	cx.AbortWithStatus(http.StatusEarlyHints)
}

// 2XX responses

func RespondOk(cx *gin.Context, data any) {
	cx.AbortWithStatusJSON(http.StatusOK, data)
}

func RespondCreated(cx *gin.Context, data any) {
	cx.AbortWithStatusJSON(http.StatusCreated, data)
}

func RespondAccepted(cx *gin.Context, data any) {
	cx.AbortWithStatusJSON(http.StatusAccepted, data)
}

func RespondNonAuthoritativeInfo(cx *gin.Context, data any) {
	cx.AbortWithStatusJSON(http.StatusNonAuthoritativeInfo, data)
}

func RespondNoContent(cx *gin.Context) {
	cx.AbortWithStatus(http.StatusNoContent)
}

func RespondResetContent(cx *gin.Context, data any) {
	cx.AbortWithStatusJSON(http.StatusResetContent, data)
}

func RespondPartialContent(cx *gin.Context, data any) {
	cx.AbortWithStatusJSON(http.StatusPartialContent, data)
}

func RespondMultiStatus(cx *gin.Context, data any) {
	cx.AbortWithStatusJSON(http.StatusMultiStatus, data)
}

func RespondAlreadyReported(cx *gin.Context, data any) {
	cx.AbortWithStatusJSON(http.StatusAlreadyReported, data)
}

func RespondUsed(cx *gin.Context, data any) {
	cx.AbortWithStatusJSON(http.StatusIMUsed, data)
}

// 3XX responses

func RespondMultipleChoices(cx *gin.Context) {
	cx.AbortWithStatus(http.StatusMultipleChoices)
}

func RespondMovedPermanently(cx *gin.Context) {
	cx.AbortWithStatus(http.StatusMovedPermanently)
}

func RespondFound(cx *gin.Context) {
	cx.AbortWithStatus(http.StatusFound)
}

func RespondSeeOther(cx *gin.Context) {
	cx.AbortWithStatus(http.StatusSeeOther)
}

func RespondNotModified(cx *gin.Context) {
	cx.AbortWithStatus(http.StatusNotModified)
}

func RespondUseProxy(cx *gin.Context) {
	cx.AbortWithStatus(http.StatusUseProxy)
}

func RespondTemporaryRedirect(cx *gin.Context) {
	cx.AbortWithStatus(http.StatusTemporaryRedirect)
}

func RespondPermanentRedirect(cx *gin.Context) {
	cx.AbortWithStatus(http.StatusPermanentRedirect)
}

// 4XX responses

func RespondBadRequest(cx *gin.Context, err error) {
	respondRequestError(cx, err, http.StatusBadRequest)
}

func RespondValidationError(cx *gin.Context, err error) {
	respondRequestError(cx, err, http.StatusBadRequest)
}

func RespondUnauthorized(cx *gin.Context, err error) {
	respondRequestError(cx, err, http.StatusUnauthorized)
}

func RespondForbidden(cx *gin.Context, err error) {
	respondRequestError(cx, err, http.StatusForbidden)
}

func RespondNotFound(cx *gin.Context, err error) {
	respondRequestError(cx, err, http.StatusNotFound)
}

func RespondMethodNotAllowed(cx *gin.Context, err error) {
	respondRequestError(cx, err, http.StatusMethodNotAllowed)
}

func RespondNotAcceptable(cx *gin.Context, err error) {
	respondRequestError(cx, err, http.StatusNotAcceptable)
}

func RespondProxyAuthRequired(cx *gin.Context, err error) {
	respondRequestError(cx, err, http.StatusProxyAuthRequired)
}

func RespondRequestTimeout(cx *gin.Context, err error) {
	respondRequestError(cx, err, http.StatusRequestTimeout)
}

func RespondConflict(cx *gin.Context, err error) {
	respondRequestError(cx, err, http.StatusConflict)
}

func RespondGone(cx *gin.Context, err error) {
	respondRequestError(cx, err, http.StatusGone)
}

func RespondLengthRequired(cx *gin.Context, err error) {
	respondRequestError(cx, err, http.StatusLengthRequired)
}

func RespondPreconditionFailed(cx *gin.Context, err error) {
	respondRequestError(cx, err, http.StatusPreconditionFailed)
}

func RespondEntityTooLarge(cx *gin.Context, err error) {
	respondRequestError(cx, err, http.StatusRequestEntityTooLarge)
}

func RespondURITooLong(cx *gin.Context, err error) {
	respondRequestError(cx, err, http.StatusRequestURITooLong)
}

func RespondUnsupportedMediaType(cx *gin.Context, err error) {
	respondRequestError(cx, err, http.StatusUnsupportedMediaType)
}

func RespondRangeNotSatisfiable(cx *gin.Context, err error) {
	respondRequestError(cx, err, http.StatusRequestedRangeNotSatisfiable)
}

func RespondExpectationFailed(cx *gin.Context, err error) {
	respondRequestError(cx, err, http.StatusExpectationFailed)
}

func RespondTeapot(cx *gin.Context, err error) {
	respondRequestError(cx, err, http.StatusTeapot)
}

func RespondMisdirected(cx *gin.Context, err error) {
	respondRequestError(cx, err, http.StatusMisdirectedRequest)
}

func RespondUnprocessableEntity(cx *gin.Context, err error) {
	respondRequestError(cx, err, http.StatusUnprocessableEntity)
}

func RespondLocked(cx *gin.Context, err error) {
	respondRequestError(cx, err, http.StatusLocked)
}

func RespondFailedDependency(cx *gin.Context, err error) {
	respondRequestError(cx, err, http.StatusFailedDependency)
}

func RespondTooEarly(cx *gin.Context, err error) {
	respondRequestError(cx, err, http.StatusTooEarly)
}

func RespondUpgradeRequired(cx *gin.Context, err error) {
	respondRequestError(cx, err, http.StatusUpgradeRequired)
}

func RespondPreconditionRequired(cx *gin.Context, err error) {
	respondRequestError(cx, err, http.StatusPreconditionRequired)
}

func RespondTooManyRequests(cx *gin.Context, err error) {
	respondRequestError(cx, err, http.StatusTooManyRequests)
}

func RespondHeaderFieldsTooLarge(cx *gin.Context, err error) {
	respondRequestError(cx, err, http.StatusRequestHeaderFieldsTooLarge)
}

func RespondUnavailableForLegalReasons(cx *gin.Context, err error) {
	respondRequestError(cx, err, http.StatusUnavailableForLegalReasons)
}

// 5XX responses

func RespondServerError(cx *gin.Context, err error) {
	respondServerError(cx, err, http.StatusInternalServerError)
}

func RespondNotImplemented(cx *gin.Context, err error) {
	respondServerError(cx, err, http.StatusNotImplemented)
}

func RespondBadGateway(cx *gin.Context, err error) {
	respondServerError(cx, err, http.StatusBadGateway)
}

func RespondServiceUnavailable(cx *gin.Context, err error) {
	respondServerError(cx, err, http.StatusServiceUnavailable)
}

func RespondGatewayTimeout(cx *gin.Context, err error) {
	respondServerError(cx, err, http.StatusGatewayTimeout)
}

func RespondVersionNotSupported(cx *gin.Context, err error) {
	respondServerError(cx, err, http.StatusHTTPVersionNotSupported)
}

func RespondVariantAlsoNegotiates(cx *gin.Context, err error) {
	respondServerError(cx, err, http.StatusVariantAlsoNegotiates)
}

func RespondInsufficientStorage(cx *gin.Context, err error) {
	respondServerError(cx, err, http.StatusInsufficientStorage)
}

func RespondLoopDetected(cx *gin.Context, err error) {
	respondServerError(cx, err, http.StatusLoopDetected)
}

func RespondNotExtended(cx *gin.Context, err error) {
	respondServerError(cx, err, http.StatusNotExtended)
}

func RespondNetworkAuthenticationRequired(cx *gin.Context, err error) {
	respondServerError(cx, err, http.StatusNetworkAuthenticationRequired)
}
