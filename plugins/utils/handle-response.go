package utils

import (
	"bytes"
	"customErrors"
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
	// Copy headers
	for k, hs := range resp.Header {
		for _, h := range hs {
			w.Header().Add(k, h)
		}
	}

	// Copy status code
	w.WriteHeader(resp.StatusCode)

	// Copy body
	if resp.Body == nil {
		return
	}
	io.Copy(w, resp.Body)
	resp.Body.Close()
}

func WriteErrorToHttpResponseWriter(err error, resp *http.Response, w http.ResponseWriter) {
	var oHTTPResponseError customErrors.HTTPResponseError
	switch e := err.(type) {
	case customErrors.HTTPResponseError:
		oHTTPResponseError = e
	default:
		oHTTPResponseError = customErrors.ErrorToHTTPResponseError(e, http.StatusInternalServerError)
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

func WriteErrorToResponseWrapper(err error, resp *ResponseWrapper) *ResponseWrapper {
	var oHTTPResponseError customErrors.HTTPResponseError
	switch e := err.(type) {
	case customErrors.HTTPResponseError:
		oHTTPResponseError = e
	default:
		oHTTPResponseError = customErrors.ErrorToHTTPResponseError(e, http.StatusInternalServerError)
	}

	if len(resp.Data()) > 0 {
		// Handle body (output_encoding different than no-op)
		var data map[string]interface{}
		json.Unmarshal([]byte(oHTTPResponseError.Error()), &data)
		resp.SetData(data)

		return resp
	}

	payload := []byte(oHTTPResponseError.Error())
	contentLength := int64(len(payload))

	// Handle headers
	resp.Headers()["Content-Length"] = []string{fmt.Sprint(contentLength)}
	resp.Headers()["Content-Type"] = []string{oHTTPResponseError.HTTPEncoding}

	// Handle status code
	resp.SetStatusCode(oHTTPResponseError.StatusCode())

	// Handle body (output_encoding equals no-op)
	resp.SetIo(bytes.NewReader(payload))

	return resp
}
