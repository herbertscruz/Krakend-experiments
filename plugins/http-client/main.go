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
	pluginContext := PluginContext{r, pluginName, extra, logger}
	customPlugin, err := NewCustomPlugin(pluginContext)
	if err != nil {
		return nil, err
	}

	return http.HandlerFunc(customPlugin.Bootstrap), nil
}

func main() {}

func init() {
	fmt.Printf("\n--- %s loaded ---\n", string(ClientRegisterer))
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
}
