package main

import (
	"errors"
	"fmt"
	"utils"
)

type RequestModifierPlugin struct {
	ctx PluginContext
}

func NewRequestModifierPlugin(ctx PluginContext) (*RequestModifierPlugin, error) {
	p := RequestModifierPlugin{}
	p.ctx = ctx

	_, okGeneral := ctx.extra[ctx.pluginName].(map[string]interface{})
	_, ok := ctx.extra[ctx.pluginName+"-request"].(map[string]interface{})
	if !okGeneral && !ok {
		return nil, errors.New(fmt.Sprintf("general or request configuration of the %s plugin not found", ctx.pluginName))
	}

	return &p, nil
}

func (p *RequestModifierPlugin) Bootstrap(req *utils.RequestWrapper) (*utils.RequestWrapper, error) {
	fmt.Println("params:", req.Params())
	fmt.Println("headers:", req.Headers())
	fmt.Println("method:", req.Method())
	fmt.Println("url:", req.URL())
	fmt.Println("query:", req.Query())
	fmt.Println("path:", req.Path())

	return req, nil
}
