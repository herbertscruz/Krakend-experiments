package main

import (
	"fmt"
	"utils"
)

type RequestModifierPlugin struct {
	ctx PluginContext
}

func NewRequestModifierPlugin(ctx PluginContext) (*RequestModifierPlugin, error) {
	p := RequestModifierPlugin{}
	p.ctx = ctx

	configGeneral, okGeneral := ctx.extra[ctx.pluginName].(map[string]interface{})

	logger.Debug(loggerFormatter(ctx.pluginName, fmt.Sprintf("configGeneral: %v", configGeneral)))

	config, ok := ctx.extra[ctx.pluginName+"-request"].(map[string]interface{})
	if !okGeneral && !ok {
		info := fmt.Sprintf("general or request configuration of the %s plugin not found", ctx.pluginName)
		logger.Info(loggerFormatter(ctx.pluginName, info))
	}

	logger.Debug(loggerFormatter(ctx.pluginName, fmt.Sprintf("config: %v", config)))

	return &p, nil
}

func (p *RequestModifierPlugin) Bootstrap(req *utils.RequestWrapper) (*utils.RequestWrapper, error) {
	logger.Debug(loggerFormatter(p.ctx.pluginName, fmt.Sprintf("params: %v", req.Params())))
	logger.Debug(loggerFormatter(p.ctx.pluginName, fmt.Sprintf("headers: %v", req.Headers())))
	logger.Debug(loggerFormatter(p.ctx.pluginName, fmt.Sprintf("method: %v", req.Method())))
	logger.Debug(loggerFormatter(p.ctx.pluginName, fmt.Sprintf("url: %v", req.URL())))
	logger.Debug(loggerFormatter(p.ctx.pluginName, fmt.Sprintf("query: %v", req.Query())))
	logger.Debug(loggerFormatter(p.ctx.pluginName, fmt.Sprintf("path: %v", req.Path())))

	return req, nil
}
