package types

// A generic response.
// swagger:response Response
type Response struct {
	Message string                 `json:"message"`
	Errors  []error                `json:"errors,omitempty"`
	Data    map[string]interface{} `json:"data,omitempty"`
}

// A validation error response.
// swagger:response ValidationResponse
type ValidationResponse struct {
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors,omitempty"`
}

// A boolean response.
// swagger:response BoolResponse
type BoolResponse struct {
	Valid bool `json:"valid"`
}
