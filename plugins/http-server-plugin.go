package main

import (
	"errors"
	"fmt"
	"net/http"
)

type HttpServerPlugin struct {
	ctx PluginContext
}

func NewHttpServerPlugin(ctx PluginContext) (*HttpServerPlugin, error) {
	p := HttpServerPlugin{}
	p.ctx = ctx

	_, ok := ctx.extra[ctx.pluginName].(map[string]interface{})
	if !ok {
		return nil, errors.New(fmt.Sprintf("configuration of the %s plugin not found", ctx.pluginName))
	}

	return &p, nil
}

func (p *HttpServerPlugin) Bootstrap(w http.ResponseWriter, req *http.Request) (*http.Response, error) {
	p.ctx.h.ServeHTTP(w, req)
	return nil, nil
}
