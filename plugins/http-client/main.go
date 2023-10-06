// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/herbertscruz/krakend-experiments/shared"
)

var pluginName = "krakend-client-example"

var ClientRegisterer = registerer(pluginName)

type registerer string

func (r registerer) RegisterClients(f func(
	name string,
	handler func(context.Context, map[string]interface{}) (http.Handler, error),
)) {
	f(string(r), r.registerClients)
}

func (r registerer) registerClients(_ context.Context, extra map[string]interface{}) (http.Handler, error) {
	pluginContext := PluginContext{r, pluginName, extra}
	customPlugin, err := NewCustomPlugin(pluginContext)
	if err != nil {
		return nil, err
	}

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		resp, err := customPlugin.Bootstrap(w, req)
		if err != nil {
			shared.WriteErrorToHttpResponseWriter(err, resp, w)
			return
		}

		if resp != nil {
			shared.WriteHttpResponseToHttpResponseWriter(resp, w)
			return
		}
	}), nil
}

func main() {}

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
