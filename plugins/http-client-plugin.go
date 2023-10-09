package main

import (
	"errors"
	"fmt"
	"net/http"
)

type HttpClientPlugin struct {
	ctx PluginContext
}

func NewHttpClientPlugin(ctx PluginContext) (*HttpClientPlugin, error) {
	p := HttpClientPlugin{}
	p.ctx = ctx

	_, ok := ctx.extra[ctx.pluginName].(map[string]interface{})
	if !ok {
		return nil, errors.New(fmt.Sprintf("configuration of the %s plugin not found", ctx.pluginName))
	}

	return &p, nil
}

func (p *HttpClientPlugin) Bootstrap(w http.ResponseWriter, req *http.Request) (*http.Response, error) {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, ErrorToHTTPResponseError(err, http.StatusInternalServerError)
	}

	return resp, nil
}
