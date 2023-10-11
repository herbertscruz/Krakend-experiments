package main

import (
	"customErrors"
	"fmt"
	"net/http"
)

type HttpClientPlugin struct {
	ctx PluginContext
}

func NewHttpClientPlugin(ctx PluginContext) (*HttpClientPlugin, error) {
	p := HttpClientPlugin{}
	p.ctx = ctx

	config, ok := ctx.extra[ctx.pluginName].(map[string]interface{})
	if !ok {
		info := fmt.Sprintf("configuration of the %s plugin not found", ctx.pluginName)
		logger.Info(loggerFormatter(ctx.pluginName, info))
	}

	logger.Debug(loggerFormatter(ctx.pluginName, fmt.Sprintf("config: %v", config)))

	return &p, nil
}

func (p *HttpClientPlugin) Bootstrap(w http.ResponseWriter, req *http.Request) (*http.Response, error) {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, customErrors.ErrorToHTTPResponseError(err, http.StatusInternalServerError)
	}

	return resp, nil
}
