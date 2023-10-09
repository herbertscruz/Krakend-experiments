// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"customErrors"
	"errors"
	"fmt"
	"net/http"
	"utils"
)

var pluginName = "custom"

var ModifierRegisterer = registerer(pluginName + "-modifier")
var HandlerRegisterer = registerer(pluginName + "-http-server")
var ClientRegisterer = registerer(pluginName + "-http-client")

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
		logger.Info(fmt.Sprintf("--- Started %s-request ---", string(ModifierRegisterer)))

		defer func() {
			logger.Info(fmt.Sprintf("--- Completed %s-request ---", string(ModifierRegisterer)))
		}()

		req, ok := input.(utils.RequestWrapperInterface)
		if !ok {
			err := customErrors.ErrorToHTTPResponseError(unkownTypeErr, http.StatusInternalServerError)
			logger.Error(err)
			return nil, err
		}

		request := utils.RequestWrapper{}
		request.SetValues(
			req.Params(),
			req.Headers(),
			req.Body(),
			req.Method(),
			req.URL(),
			req.Query(),
			req.Path(),
		)

		pluginContext := PluginContext{r, string(ModifierRegisterer), extra, nil}
		plugin, err := NewRequestModifierPlugin(pluginContext)
		if err != nil {
			err := customErrors.ErrorToHTTPResponseError(err, http.StatusInternalServerError)
			logger.Error(err)
			return nil, err
		}

		return plugin.Bootstrap(&request)
	}
}

func (r registerer) responseDump(
	extra map[string]interface{},
) func(interface{}) (interface{}, error) {
	return func(input interface{}) (interface{}, error) {
		logger.Info(fmt.Sprintf("--- Started %s-response ---", string(ModifierRegisterer)))

		defer func() {
			logger.Info(fmt.Sprintf("--- Completed %s-response ---", string(ModifierRegisterer)))
		}()

		resp, ok := input.(utils.ResponseWrapperInterface)
		var err error
		if !ok {
			err = customErrors.ErrorToHTTPResponseError(unkownTypeErr, http.StatusInternalServerError)
			logger.Error(err)
		}

		response := utils.ResponseWrapper{}
		response.SetValues(
			resp.Data(),
			resp.Io(),
			resp.IsComplete(),
			resp.StatusCode(),
			resp.Headers(),
		)

		if err != nil {
			return utils.WriteErrorToResponseWrapper(err, &response), nil
		}

		pluginContext := PluginContext{r, string(ModifierRegisterer), extra, nil}
		plugin, err := NewResponseModifierPlugin(pluginContext)
		if err != nil {
			err := customErrors.ErrorToHTTPResponseError(err, http.StatusInternalServerError)
			logger.Error(err)
			return utils.WriteErrorToResponseWrapper(err, &response), nil
		}

		wrapper, err := plugin.Bootstrap(&response)
		if err != nil {
			logger.Error(err)
			return utils.WriteErrorToResponseWrapper(err, &response), nil
		}

		return wrapper, nil
	}
}

func (r registerer) RegisterHandlers(f func(
	name string,
	handler func(context.Context, map[string]interface{}, http.Handler) (http.Handler, error),
)) {
	f(string(r), r.registerHandlers)
}

func (r registerer) registerHandlers(_ context.Context, extra map[string]interface{}, h http.Handler) (http.Handler, error) {
	pluginContext := PluginContext{r, string(HandlerRegisterer), extra, h}
	plugin, err := NewHttpServerPlugin(pluginContext)
	if err != nil {
		return h, err
	}

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		logger.Info(fmt.Sprintf("--- Started %s ---", string(HandlerRegisterer)))

		defer func() {
			logger.Info(fmt.Sprintf("--- Completed %s ---", string(HandlerRegisterer)))
		}()

		resp, err := plugin.Bootstrap(w, req)
		if err != nil {
			utils.WriteErrorToHttpResponseWriter(err, resp, w)
			return
		}

		if resp != nil {
			utils.WriteHttpResponseToHttpResponseWriter(resp, w)
			return
		}
	}), nil
}

func (r registerer) RegisterClients(f func(
	name string,
	handler func(context.Context, map[string]interface{}) (http.Handler, error),
)) {
	f(string(r), r.registerClients)
}

func (r registerer) registerClients(_ context.Context, extra map[string]interface{}) (http.Handler, error) {
	pluginContext := PluginContext{r, string(ClientRegisterer), extra, nil}
	plugin, err := NewHttpClientPlugin(pluginContext)
	if err != nil {
		return nil, err
	}

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		logger.Info(fmt.Sprintf("--- Started %s ---", string(ClientRegisterer)))

		defer func() {
			logger.Info(fmt.Sprintf("--- Completed %s ---", string(ClientRegisterer)))
		}()

		resp, err := plugin.Bootstrap(w, req)
		if err != nil {
			utils.WriteErrorToHttpResponseWriter(err, resp, w)
			return
		}

		if resp != nil {
			utils.WriteHttpResponseToHttpResponseWriter(resp, w)
			return
		}
	}), nil
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
	h          http.Handler
}
