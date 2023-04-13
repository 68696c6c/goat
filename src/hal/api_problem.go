package hal

import "net/http"

// ApiProblem describes an API error according to the API-Problem proposed standard.
// https://phlyrestfully.readthedocs.io/en/latest/problems.html
// https://datatracker.ietf.org/doc/html/draft-nottingham-http-problem-02
type ApiProblem struct {
	DescribedBy string `json:"describedBy"`          // a URL to a document describing the error condition (required)
	Title       string `json:"title"`                // a brief title for the error condition (required)
	HttpStatus  int    `json:"httpStatus,omitempty"` // the HTTP status code for the current request (optional)
	Details     string `json:"details,omitempty"`    // error details specific to this request (optional)
	SupportID   string `json:"supportId,omitempty"`  // a URL to the specific problem occurrence (e.g., to a log message) (optional)
}

func NewApiProblem(statusCode int, err error) ApiProblem {
	var details string
	if err != nil {
		details = err.Error()
	}
	if rfcLink, ok := httpStatusRfcLinkMap[statusCode]; ok {
		return ApiProblem{
			DescribedBy: rfcLink,
			Title:       http.StatusText(statusCode),
			HttpStatus:  statusCode,
			Details:     details,
		}
	}
	code := http.StatusInternalServerError
	return ApiProblem{
		DescribedBy: httpStatusRfcLinkMap[code],
		Title:       http.StatusText(code),
		HttpStatus:  code,
		Details:     details,
	}
}

var httpStatusRfcLinkMap = map[int]string{
	// 100
	http.StatusContinue:           "https://www.rfc-editor.org/rfc/rfc9110.html#name-100-continue",
	http.StatusSwitchingProtocols: "https://www.rfc-editor.org/rfc/rfc9110.html#name-101-switching-protocols",
	http.StatusProcessing:         "https://www.rfc-editor.org/rfc/rfc2518.html#section-10.1",
	http.StatusEarlyHints:         "https://www.rfc-editor.org/rfc/rfc8297.html#section-2",

	// 200
	http.StatusOK:                   "https://www.rfc-editor.org/rfc/rfc9110.html#name-200-ok",
	http.StatusCreated:              "https://www.rfc-editor.org/rfc/rfc9110.html#name-201-created",
	http.StatusAccepted:             "https://www.rfc-editor.org/rfc/rfc9110.html#name-202-accepted",
	http.StatusNonAuthoritativeInfo: "https://www.rfc-editor.org/rfc/rfc9110.html#name-203-non-authoritative-infor",
	http.StatusNoContent:            "https://www.rfc-editor.org/rfc/rfc9110.html#name-204-no-content",
	http.StatusResetContent:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-205-reset-content",
	http.StatusPartialContent:       "https://www.rfc-editor.org/rfc/rfc9110.html#name-206-partial-content",
	http.StatusMultiStatus:          "https://www.rfc-editor.org/rfc/rfc4918#section-11.1",
	http.StatusAlreadyReported:      "https://www.rfc-editor.org/rfc/rfc5842#section-7.1",
	http.StatusIMUsed:               "https://www.rfc-editor.org/rfc/rfc3229#section-10.4.1",

	// 300
	http.StatusMultipleChoices:  "https://www.rfc-editor.org/rfc/rfc9110.html#name-300-multiple-choices",
	http.StatusMovedPermanently: "https://www.rfc-editor.org/rfc/rfc9110.html#name-301-moved-permanently",
	http.StatusFound:            "https://www.rfc-editor.org/rfc/rfc9110.html#name-302-found",
	http.StatusSeeOther:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-303-see-other",
	http.StatusNotModified:      "https://www.rfc-editor.org/rfc/rfc9110.html#name-304-not-modified",
	http.StatusUseProxy:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-305-use-proxy",
	// http._: "https://www.rfc-editor.org/rfc/rfc9110.html#name-306-unused", // Status code 306 is unused.
	http.StatusTemporaryRedirect: "https://www.rfc-editor.org/rfc/rfc9110.html#name-307-temporary-redirect",
	http.StatusPermanentRedirect: "https://www.rfc-editor.org/rfc/rfc9110.html#name-308-permanent-redirect",

	// 400
	http.StatusBadRequest:                   "https://www.rfc-editor.org/rfc/rfc9110.html#name-400-bad-request",
	http.StatusUnauthorized:                 "https://www.rfc-editor.org/rfc/rfc9110.html#name-401-unauthorized",
	http.StatusPaymentRequired:              "https://www.rfc-editor.org/rfc/rfc9110.html#name-402-payment-required",
	http.StatusForbidden:                    "https://www.rfc-editor.org/rfc/rfc9110.html#name-403-forbidden",
	http.StatusNotFound:                     "https://www.rfc-editor.org/rfc/rfc9110.html#name-404-not-found",
	http.StatusMethodNotAllowed:             "https://www.rfc-editor.org/rfc/rfc9110.html#name-405-method-not-allowed",
	http.StatusNotAcceptable:                "https://www.rfc-editor.org/rfc/rfc9110.html#name-406-not-acceptable",
	http.StatusProxyAuthRequired:            "https://www.rfc-editor.org/rfc/rfc9110.html#name-407-proxy-authentication-re",
	http.StatusRequestTimeout:               "https://www.rfc-editor.org/rfc/rfc9110.html#name-408-request-timeout",
	http.StatusConflict:                     "https://www.rfc-editor.org/rfc/rfc9110.html#name-409-conflict",
	http.StatusGone:                         "https://www.rfc-editor.org/rfc/rfc9110.html#name-410-gone",
	http.StatusLengthRequired:               "https://www.rfc-editor.org/rfc/rfc9110.html#name-411-length-required",
	http.StatusPreconditionFailed:           "https://www.rfc-editor.org/rfc/rfc9110.html#name-412-precondition-failed",
	http.StatusRequestEntityTooLarge:        "https://www.rfc-editor.org/rfc/rfc9110.html#name-413-content-too-large",
	http.StatusRequestURITooLong:            "https://www.rfc-editor.org/rfc/rfc9110.html#name-414-uri-too-long",
	http.StatusUnsupportedMediaType:         "https://www.rfc-editor.org/rfc/rfc9110.html#name-415-unsupported-media-type",
	http.StatusRequestedRangeNotSatisfiable: "https://www.rfc-editor.org/rfc/rfc9110.html#name-416-range-not-satisfiable",
	http.StatusExpectationFailed:            "https://www.rfc-editor.org/rfc/rfc9110.html#name-417-expectation-failed",
	http.StatusTeapot:                       "https://www.rfc-editor.org/rfc/rfc9110.html#name-418-unused",
	http.StatusMisdirectedRequest:           "https://www.rfc-editor.org/rfc/rfc9110.html#name-421-misdirected-request",
	http.StatusUnprocessableEntity:          "https://www.rfc-editor.org/rfc/rfc9110.html#name-422-unprocessable-content",
	http.StatusLocked:                       "https://www.rfc-editor.org/rfc/rfc4918#section-11.3",
	http.StatusFailedDependency:             "https://www.rfc-editor.org/rfc/rfc4918#section-11.4",
	http.StatusTooEarly:                     "https://www.rfc-editor.org/rfc/rfc8470#section-5.2",
	http.StatusUpgradeRequired:              "https://www.rfc-editor.org/rfc/rfc9110.html#name-426-upgrade-required",
	http.StatusPreconditionRequired:         "https://www.rfc-editor.org/rfc/rfc6585.html#section-3",
	http.StatusTooManyRequests:              "https://www.rfc-editor.org/rfc/rfc6585.html#section-4",
	http.StatusRequestHeaderFieldsTooLarge:  "https://www.rfc-editor.org/rfc/rfc6585.html#section-5",
	http.StatusUnavailableForLegalReasons:   "https://www.rfc-editor.org/rfc/rfc7725#section-3",

	// 500
	http.StatusInternalServerError:           "https://www.rfc-editor.org/rfc/rfc9110.html#section-15.6.1",
	http.StatusNotImplemented:                "https://www.rfc-editor.org/rfc/rfc9110.html#section-15.6.2",
	http.StatusBadGateway:                    "https://www.rfc-editor.org/rfc/rfc9110.html#name-502-bad-gateway",
	http.StatusServiceUnavailable:            "https://www.rfc-editor.org/rfc/rfc9110.html#name-503-service-unavailable",
	http.StatusGatewayTimeout:                "https://www.rfc-editor.org/rfc/rfc9110.html#name-504-gateway-timeout",
	http.StatusHTTPVersionNotSupported:       "https://www.rfc-editor.org/rfc/rfc9110.html#name-505-http-version-not-suppor",
	http.StatusVariantAlsoNegotiates:         "https://www.rfc-editor.org/rfc/rfc2295#section-8.1",
	http.StatusInsufficientStorage:           "https://www.rfc-editor.org/rfc/rfc4918#section-11.5",
	http.StatusLoopDetected:                  "https://www.rfc-editor.org/rfc/rfc5842#section-7.2",
	http.StatusNotExtended:                   "https://www.rfc-editor.org/rfc/rfc2774#section-7",
	http.StatusNetworkAuthenticationRequired: "https://www.rfc-editor.org/rfc/rfc6585.html#section-6",
}
