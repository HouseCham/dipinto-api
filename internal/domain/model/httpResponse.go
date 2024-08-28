package model

// HTTPResponse is a generic response structure for HTTP responses
type HTTPResponse struct {
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	StatusCode int         `json:"status_code"`
}
