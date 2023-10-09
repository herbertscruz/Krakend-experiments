package main

import (
	"errors"
	"fmt"
)

type ResponseModifierPlugin struct {
	ctx PluginContext
}

func NewResponseModifierPlugin(ctx PluginContext) (*ResponseModifierPlugin, error) {
	p := ResponseModifierPlugin{}
	p.ctx = ctx

	_, okGeneral := ctx.extra[ctx.pluginName].(map[string]interface{})
	_, ok := ctx.extra[ctx.pluginName+"-response"].(map[string]interface{})
	if !okGeneral && !ok {
		return nil, errors.New(fmt.Sprintf("general or response configuration of the %s plugin not found", ctx.pluginName))
	}

	return &p, nil
}

func (p *ResponseModifierPlugin) Bootstrap(resp *ResponseWrapper) (*ResponseWrapper, error) {
	fmt.Println("data:", resp.Data())
	fmt.Println("is complete:", resp.IsComplete())
	fmt.Println("headers:", resp.Headers())
	fmt.Println("status code:", resp.StatusCode())

	return resp, nil
}
