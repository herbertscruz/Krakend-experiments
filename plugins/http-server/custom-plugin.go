package main

import (
	"errors"
	"fmt"
	"html"
	"net/http"
)

type CustomPlugin struct {
	ctx  PluginContext
	path string
}

func NewCustomPlugin(ctx PluginContext) (*CustomPlugin, error) {
	p := CustomPlugin{}
	p.ctx = ctx

	config, ok := ctx.extra[ctx.pluginName].(map[string]interface{})
	if !ok {
		return nil, errors.New(fmt.Sprintf("configuration of the %s plugin not found", ctx.pluginName))
	}

	// The plugin will look for this path:
	path, _ := config["path"].(string)
	logger.Debug(fmt.Sprintf("The plugin is now hijacking the path %s", path))

	p.path = path

	return &p, nil
}

func (p *CustomPlugin) Bootstrap(w http.ResponseWriter, req *http.Request) {
	// If the requested path is not what we defined, continue.
	if req.URL.Path != p.path {
		p.ctx.h.ServeHTTP(w, req)
		return
	}

	// The path has to be hijacked:
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(req.URL.Path))
	logger.Debug("request:", html.EscapeString(req.URL.Path))
}
