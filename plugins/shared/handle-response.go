package shared

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func WriteHttpResponseFromByteArray(payload []byte, statusCode int) *http.Response {
	contentLength := int64(len(payload))

	header := http.Header{
		"Content-Type": {"application/json"},
	}

	return &http.Response{
		StatusCode:    statusCode,
		Header:        header,
		ContentLength: contentLength,
		Body:          io.NopCloser(bytes.NewBuffer(payload)),
	}
}

func WriteHttpResponseFromMap(payload map[string]string, statusCode int) *http.Response {
	b, _ := json.Marshal(payload)
	return WriteHttpResponseFromByteArray(b, statusCode)
}

func WriteHttpResponseToHttpResponseWriter(resp *http.Response, w http.ResponseWriter) {
	// Copy headers, status codes, and body
	for k, hs := range resp.Header {
		for _, h := range hs {
			w.Header().Add(k, h)
		}
	}

	// Copy status codes
	w.WriteHeader(resp.StatusCode)

	// Copy body
	if resp.Body == nil {
		return
	}
	io.Copy(w, resp.Body)
	resp.Body.Close()
}

func WriteErrorToHttpResponseWriter(err error, resp *http.Response, w http.ResponseWriter) {
	var oHTTPResponseError HTTPResponseError
	switch e := err.(type) {
	case HTTPResponseError:
		oHTTPResponseError = e
	default:
		oHTTPResponseError = ErrorToHTTPResponseError(e, http.StatusInternalServerError)
	}

	payload := []byte(err.Error())
	contentLength := int64(len(payload))

	var newresp *http.Response
	if resp == nil {
		newresp = WriteHttpResponseFromByteArray(payload, oHTTPResponseError.StatusCode())
	} else {
		header := resp.Header.Clone()
		header.Set("Content-Length", fmt.Sprint(contentLength))
		defer func() {
			resp.Body.Close()
		}()

		newresp = &http.Response{
			StatusCode:    resp.StatusCode,
			Header:        header,
			ContentLength: contentLength,
			Body:          io.NopCloser(bytes.NewBuffer(payload)),
		}
	}

	WriteHttpResponseToHttpResponseWriter(newresp, w)
}
