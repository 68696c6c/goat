package resource

import "net/http"

// ApiProblem describes an API error according to the API-Problem proposed standard.
// https://phlyrestfully.readthedocs.io/en/latest/problems.html
type ApiProblem struct {
	DescribedBy string `json:"describedBy"`          // a URL to a document describing the error condition (required)
	Title       string `json:"title"`                // a brief title for the error condition (required)
	HttpStatus  int    `json:"httpStatus,omitempty"` // the HTTP status code for the current request (optional)
	Details     string `json:"details,omitempty"`    // error details specific to this request (optional)
	SupportId   string `json:"supportId,omitempty"`  // a URL to the specific problem occurrence (e.g., to a log message) (optional)
}

func NewApiProblem(status int, title string, err error) *ApiProblem {
	var details string
	if err != nil {
		details = err.Error()
	}
	return &ApiProblem{
		DescribedBy: GetStatusDescription(status),
		Title:       title,
		HttpStatus:  status,
		Details:     details,
	}
}

// SetDescribedBy attempts to set the DescribedBy field using the provided text or a link to the RFC describing the HttpStatus.
func (a *ApiProblem) SetDescribedBy(text string) *ApiProblem {
	a.DescribedBy = text
	return a
}

func GetStatusDescription(status int) string {
	if rfcLink, ok := httpStatusRfcLinkMap[status]; ok {
		return rfcLink
	}
	return ""
}

// TODO: finish filling this out
var httpStatusRfcLinkMap = map[int]string{
	// 100
	http.StatusContinue:           "https://www.rfc-editor.org/rfc/rfc9110.html#name-100-continue",
	http.StatusSwitchingProtocols: "https://www.rfc-editor.org/rfc/rfc9110.html#name-101-switching-protocols",
	http.StatusProcessing:         "https://www.rfc-editor.org/rfc/rfc2518.html#section-10.1",
	http.StatusEarlyHints:         "https://www.rfc-editor.org/rfc/rfc8297.html#section-2",

	// 200
	http.StatusOK:                   "",
	http.StatusCreated:              "",
	http.StatusAccepted:             "",
	http.StatusNonAuthoritativeInfo: "",
	http.StatusNoContent:            "",
	http.StatusResetContent:         "",
	http.StatusPartialContent:       "",
	http.StatusMultiStatus:          "",
	http.StatusAlreadyReported:      "",
	http.StatusIMUsed:               "",

	// 300
	http.StatusMultipleChoices:  "",
	http.StatusMovedPermanently: "",
	http.StatusFound:            "",
	http.StatusSeeOther:         "",
	http.StatusNotModified:      "",
	http.StatusUseProxy:         "",
	// http._: "",
	http.StatusTemporaryRedirect: "",
	http.StatusPermanentRedirect: "",

	// 400
	http.StatusBadRequest:                   "",
	http.StatusUnauthorized:                 "",
	http.StatusPaymentRequired:              "",
	http.StatusForbidden:                    "",
	http.StatusNotFound:                     "",
	http.StatusMethodNotAllowed:             "",
	http.StatusNotAcceptable:                "",
	http.StatusProxyAuthRequired:            "",
	http.StatusRequestTimeout:               "",
	http.StatusConflict:                     "",
	http.StatusGone:                         "",
	http.StatusLengthRequired:               "",
	http.StatusPreconditionFailed:           "",
	http.StatusRequestEntityTooLarge:        "",
	http.StatusRequestURITooLong:            "",
	http.StatusUnsupportedMediaType:         "",
	http.StatusRequestedRangeNotSatisfiable: "",
	http.StatusExpectationFailed:            "",
	http.StatusTeapot:                       "",
	http.StatusMisdirectedRequest:           "",
	http.StatusUnprocessableEntity:          "",
	http.StatusLocked:                       "",
	http.StatusFailedDependency:             "",
	http.StatusTooEarly:                     "",
	http.StatusUpgradeRequired:              "",
	http.StatusPreconditionRequired:         "",
	http.StatusTooManyRequests:              "",
	http.StatusRequestHeaderFieldsTooLarge:  "",
	http.StatusUnavailableForLegalReasons:   "",

	// 500
	http.StatusInternalServerError:           "",
	http.StatusNotImplemented:                "",
	http.StatusBadGateway:                    "",
	http.StatusServiceUnavailable:            "",
	http.StatusGatewayTimeout:                "",
	http.StatusHTTPVersionNotSupported:       "",
	http.StatusVariantAlsoNegotiates:         "",
	http.StatusInsufficientStorage:           "",
	http.StatusLoopDetected:                  "",
	http.StatusNotExtended:                   "",
	http.StatusNetworkAuthenticationRequired: "",
}
