package customErrors

import "fmt"

type responseError interface {
	error
	StatusCode() int
}

type encodedResponseError interface {
	responseError
	Encoding() string
}

type HTTPResponseError struct {
	Code         int    `json:"http_status_code"`
	Msg          string `json:"http_body,omitempty"`
	HTTPEncoding string `json:"http_encoding"`
}

func StringToHTTPResponseError(message string, statusCode int) HTTPResponseError {
	jsonError := fmt.Sprintf(`{"status": false,"code":%d,"message":"%s"}`, statusCode, message)
	return HTTPResponseError{
		Code:         statusCode,
		Msg:          jsonError,
		HTTPEncoding: "application/json",
	}
}

func ErrorToHTTPResponseError(err error, statusCode int) HTTPResponseError {
	return StringToHTTPResponseError(err.Error(), statusCode)
}

func (r HTTPResponseError) Error() string {
	return r.Msg
}

func (r HTTPResponseError) StatusCode() int {
	return r.Code
}

func (r HTTPResponseError) Encoding() string {
	return r.HTTPEncoding
}
