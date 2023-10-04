// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/herbertscruz/krakend-experiments/shared"
)

var pluginName = "krakend-server-example"

var HandlerRegisterer = registerer(pluginName)

type registerer string

func (r registerer) RegisterHandlers(f func(
	name string,
	handler func(context.Context, map[string]interface{}, http.Handler) (http.Handler, error),
)) {
	f(string(r), r.registerHandlers)
}

func (r registerer) registerHandlers(_ context.Context, extra map[string]interface{}, h http.Handler) (http.Handler, error) {
	pluginContext := PluginContext{r, pluginName, extra, logger, h}
	customPlugin, err := NewCustomPlugin(pluginContext)
	if err != nil {
		return h, err
	}

	return http.HandlerFunc(customPlugin.Bootstrap), nil
}

func main() {}

func init() {
	fmt.Printf("\n--- %s loaded ---\n", string(HandlerRegisterer))
}

var logger shared.Logger = shared.CustomLogger{}

func (registerer) RegisterLogger(v interface{}) {
	shared.RegisterCustomLogger(logger, pluginName, v)
}

type PluginContext struct {
	r          registerer
	pluginName string
	extra      map[string]interface{}
	logger     shared.Logger
	h          http.Handler
}
