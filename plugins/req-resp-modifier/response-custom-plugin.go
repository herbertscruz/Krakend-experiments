package main

import (
	"errors"
	"fmt"

	"github.com/herbertscruz/krakend-experiments/shared"
)

type ResponseCustomPlugin struct {
	ctx PluginContext
}

func NewResponseCustomPlugin(ctx PluginContext) (*ResponseCustomPlugin, error) {
	p := ResponseCustomPlugin{}
	p.ctx = ctx

	_, okGeneral := ctx.extra[ctx.pluginName].(map[string]interface{})
	_, ok := ctx.extra[ctx.pluginName+"-response"].(map[string]interface{})
	if !okGeneral && !ok {
		return nil, errors.New(fmt.Sprintf("general or response configuration of the %s plugin not found", ctx.pluginName))
	}

	return &p, nil
}

func (p *ResponseCustomPlugin) Bootstrap(resp *shared.ResponseWrapper) (*shared.ResponseWrapper, error) {
	fmt.Println("data:", resp.Data())
	fmt.Println("is complete:", resp.IsComplete())
	fmt.Println("headers:", resp.Headers())
	fmt.Println("status code:", resp.StatusCode())

	return resp, nil
}
