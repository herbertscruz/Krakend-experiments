package main

import (
	"fmt"
	"utils"
)

type ResponseModifierPlugin struct {
	ctx PluginContext
}

func NewResponseModifierPlugin(ctx PluginContext) (*ResponseModifierPlugin, error) {
	p := ResponseModifierPlugin{}
	p.ctx = ctx

	configGeneral, okGeneral := ctx.extra[ctx.pluginName].(map[string]interface{})

	logger.Debug(loggerFormatter(ctx.pluginName, fmt.Sprintf("configGeneral: %v", configGeneral)))

	config, ok := ctx.extra[ctx.pluginName+"-response"].(map[string]interface{})
	if !okGeneral && !ok {
		info := fmt.Sprintf("general or response configuration of the %s plugin not found", ctx.pluginName)
		logger.Info(loggerFormatter(ctx.pluginName, info))
	}

	logger.Debug(loggerFormatter(ctx.pluginName, fmt.Sprintf("config: %v", config)))

	return &p, nil
}

func (p *ResponseModifierPlugin) Bootstrap(resp *utils.ResponseWrapper) (*utils.ResponseWrapper, error) {
	logger.Debug(loggerFormatter(p.ctx.pluginName, fmt.Sprintf("data: %v", resp.Data())))
	logger.Debug(loggerFormatter(p.ctx.pluginName, fmt.Sprintf("is complete: %v", resp.IsComplete())))
	logger.Debug(loggerFormatter(p.ctx.pluginName, fmt.Sprintf("headers: %v", resp.Headers())))
	logger.Debug(loggerFormatter(p.ctx.pluginName, fmt.Sprintf("status code: %v", resp.StatusCode())))

	return resp, nil
}
