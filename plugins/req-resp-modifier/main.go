// SPDX-License-Identifier: Apache-2.0

package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/herbertscruz/krakend-experiments/shared"
)

var pluginName = "krakend-debugger"

var ModifierRegisterer = registerer(pluginName)

type registerer string

func (r registerer) RegisterModifiers(f func(
	name string,
	factoryFunc func(map[string]interface{}) func(interface{}) (interface{}, error),
	appliesToRequest bool,
	appliesToResponse bool,
)) {
	f(string(r)+"-request", r.requestDump, true, false)
	f(string(r)+"-response", r.responseDump, false, true)
}

func (r registerer) requestDump(
	extra map[string]interface{},
) func(interface{}) (interface{}, error) {
	return func(input interface{}) (interface{}, error) {
		req, ok := input.(shared.RequestWrapperInterface)
		if !ok {
			err := shared.ErrorToHTTPResponseError(unkownTypeErr, http.StatusInternalServerError)
			logger.Error(err)
			return nil, err
		}

		request := shared.RequestWrapper{}
		request.SetValues(
			req.Params(),
			req.Headers(),
			req.Body(),
			req.Method(),
			req.URL(),
			req.Query(),
			req.Path(),
		)

		pluginContext := PluginContext{r, pluginName, extra}
		customPlugin, err := NewRequestCustomPlugin(pluginContext)
		if err != nil {
			err := shared.ErrorToHTTPResponseError(err, http.StatusInternalServerError)
			logger.Error(err)
			return nil, err
		}

		return customPlugin.Bootstrap(&request)
	}
}

func (r registerer) responseDump(
	extra map[string]interface{},
) func(interface{}) (interface{}, error) {
	return func(input interface{}) (interface{}, error) {
		resp, ok := input.(shared.ResponseWrapperInterface)
		var err error
		if !ok {
			err = shared.ErrorToHTTPResponseError(unkownTypeErr, http.StatusInternalServerError)
			logger.Error(err)
		}

		response := shared.ResponseWrapper{}
		response.SetValues(
			resp.Data(),
			resp.Io(),
			resp.IsComplete(),
			resp.StatusCode(),
			resp.Headers(),
		)

		if err != nil {
			return shared.WriteErrorToResponseWrapper(err, &response), nil
		}

		pluginContext := PluginContext{r, pluginName, extra}
		customPlugin, err := NewResponseCustomPlugin(pluginContext)
		if err != nil {
			err := shared.ErrorToHTTPResponseError(err, http.StatusInternalServerError)
			logger.Error(err)
			return shared.WriteErrorToResponseWrapper(err, &response), nil
		}

		wrapper, err := customPlugin.Bootstrap(&response)
		if err != nil {
			logger.Error(err)
			return shared.WriteErrorToResponseWrapper(err, &response), nil
		}

		return wrapper, nil
	}
}

func main() {}

var unkownTypeErr = errors.New("unknow request type")

// This logger is replaced by the RegisterLogger method to load the one from KrakenD
var logger Logger = noopLogger{}

func (registerer) RegisterLogger(v interface{}) {
	l, ok := v.(Logger)
	if !ok {
		return
	}
	logger = l
	logger.Debug(fmt.Sprintf("[PLUGIN: %s] Logger loaded", pluginName))
}

type Logger interface {
	Debug(v ...interface{})
	Info(v ...interface{})
	Warning(v ...interface{})
	Error(v ...interface{})
	Critical(v ...interface{})
	Fatal(v ...interface{})
}

// Empty logger implementation
type noopLogger struct{}

func (n noopLogger) Debug(_ ...interface{})    {}
func (n noopLogger) Info(_ ...interface{})     {}
func (n noopLogger) Warning(_ ...interface{})  {}
func (n noopLogger) Error(_ ...interface{})    {}
func (n noopLogger) Critical(_ ...interface{}) {}
func (n noopLogger) Fatal(_ ...interface{})    {}

type PluginContext struct {
	r          registerer
	pluginName string
	extra      map[string]interface{}
}
