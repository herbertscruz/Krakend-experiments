package shared

import (
	"io"
	"net/url"
)

type RequestWrapperInterface interface {
	Params() map[string]string
	Headers() map[string][]string
	Body() io.ReadCloser
	Method() string
	URL() *url.URL
	Query() url.Values
	Path() string
}

type RequestWrapper struct {
	params  map[string]string
	headers map[string][]string
	body    io.ReadCloser
	method  string
	url     *url.URL
	query   url.Values
	path    string
}

func (r RequestWrapper) Params() map[string]string    { return r.params }
func (r RequestWrapper) Headers() map[string][]string { return r.headers }
func (r RequestWrapper) Body() io.ReadCloser          { return r.body }
func (r RequestWrapper) Method() string               { return r.method }
func (r RequestWrapper) URL() *url.URL                { return r.url }
func (r RequestWrapper) Query() url.Values            { return r.query }
func (r RequestWrapper) Path() string                 { return r.path }

func (r *RequestWrapper) SetBody(body io.ReadCloser) {
	r.body = body
}

type ResponseWrapperInterface interface {
	Data() map[string]interface{}
	Io() io.Reader
	IsComplete() bool
	StatusCode() int
	Headers() map[string][]string
}

type ResponseWrapper struct {
	data       map[string]interface{}
	io         io.Reader
	isComplete bool
	statusCode int
	headers    map[string][]string
}

func (r ResponseWrapper) Data() map[string]interface{} { return r.data }
func (r ResponseWrapper) Io() io.Reader                { return r.io }
func (r ResponseWrapper) IsComplete() bool             { return r.isComplete }
func (r ResponseWrapper) StatusCode() int              { return r.statusCode }
func (r ResponseWrapper) Headers() map[string][]string { return r.headers }

func (r *ResponseWrapper) SetIo(reader io.Reader) {
	r.io = reader
}

func (r *ResponseWrapper) SetStatusCode(statusCode int) {
	r.statusCode = statusCode
}
