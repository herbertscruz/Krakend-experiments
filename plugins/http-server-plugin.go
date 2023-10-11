package main

import (
	"fmt"
	"net/http"
)

type HttpServerPlugin struct {
	ctx PluginContext
}

func NewHttpServerPlugin(ctx PluginContext) (*HttpServerPlugin, error) {
	p := HttpServerPlugin{}
	p.ctx = ctx

	config, ok := ctx.extra[ctx.pluginName].(map[string]interface{})
	if !ok {
		info := fmt.Sprintf("configuration of the %s plugin not found", ctx.pluginName)
		logger.Info(loggerFormatter(ctx.pluginName, info))
	}

	logger.Debug(loggerFormatter(ctx.pluginName, fmt.Sprintf("config: %v", config)))

	return &p, nil
}

func (p *HttpServerPlugin) Bootstrap(w http.ResponseWriter, req *http.Request) (*http.Response, error) {
	p.ctx.h.ServeHTTP(w, req)
	return nil, nil
}
