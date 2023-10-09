package main

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

func (r *RequestWrapper) SetValues(
	params map[string]string,
	headers map[string][]string,
	body io.ReadCloser,
	method string,
	url *url.URL,
	query url.Values,
	path string,
) {
	r.params = params
	r.headers = headers
	r.body = body
	r.method = method
	r.url = url
	r.query = query
	r.path = path
}

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

func (r *ResponseWrapper) SetValues(
	data map[string]interface{},
	io io.Reader,
	isComplete bool,
	statusCode int,
	headers map[string][]string,
) {
	r.data = data
	r.io = io
	r.isComplete = isComplete
	r.statusCode = statusCode
	r.headers = headers
}

func (r *ResponseWrapper) SetData(data map[string]interface{}) {
	r.data = data
}

func (r *ResponseWrapper) SetIo(io io.Reader) {
	r.io = io
}

func (r *ResponseWrapper) SetStatusCode(statusCode int) {
	r.statusCode = statusCode
}

func (r *ResponseWrapper) SetHeaders(headers map[string][]string) {
	r.headers = headers
}