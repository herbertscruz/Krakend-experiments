package main

import (
	"errors"
	"fmt"

	"github.com/herbertscruz/krakend-experiments/shared"
)

type RequestCustomPlugin struct {
	ctx PluginContext
}

func NewRequestCustomPlugin(ctx PluginContext) (*RequestCustomPlugin, error) {
	p := RequestCustomPlugin{}
	p.ctx = ctx

	_, okGeneral := ctx.extra[ctx.pluginName].(map[string]interface{})
	_, ok := ctx.extra[ctx.pluginName+"-request"].(map[string]interface{})
	if !okGeneral && !ok {
		return nil, errors.New(fmt.Sprintf("general or request configuration of the %s plugin not found", ctx.pluginName))
	}

	return &p, nil
}

func (p *RequestCustomPlugin) Bootstrap(req *shared.RequestWrapper) (*shared.RequestWrapper, error) {
	fmt.Println("params:", req.Params())
	fmt.Println("headers:", req.Headers())
	fmt.Println("method:", req.Method())
	fmt.Println("url:", req.URL())
	fmt.Println("query:", req.Query())
	fmt.Println("path:", req.Path())

	// TODO Keep only when testing JSON type error response.
	if 0 == 0 {
		logger.Debug("testing JSON type error response")
		return nil, shared.StringToHTTPResponseError("my intentional error", 500)
	}

	return req, nil
}
