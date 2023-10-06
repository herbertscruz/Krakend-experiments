package main

import (
	"errors"
	"fmt"
	"html"
	"net/http"

	"github.com/herbertscruz/krakend-experiments/shared"
)

type CustomPlugin struct {
	ctx  PluginContext
	path string
	name string
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

	// check the passed configuration and initialize the plugin
	name, ok := config["name"].(string)
	if !ok {
		return nil, errors.New("wrong config")
	}
	if name != string(p.ctx.r) {
		return nil, fmt.Errorf("unknown register %s", name)
	}

	p.path = path
	p.name = name

	return &p, nil
}

func (p *CustomPlugin) Bootstrap(w http.ResponseWriter, req *http.Request) (*http.Response, error) {
	// The path matches, it has to be hijacked and no call to the backend happens.
	// The path is the the call to the backend, not the original request by the user.
	if req.URL.Path == p.path {
		// Return a custom JSON object:
		res := map[string]string{"message": html.EscapeString(req.URL.Path)}
		logger.Debug("request:", html.EscapeString(req.URL.Path))

		return shared.WriteHttpResponseFromMap(res, http.StatusOK), nil
	}

	// If the requested path is not what we defined, continue.
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, shared.ErrorToHTTPResponseError(err, http.StatusInternalServerError)
	}

	return resp, nil
}
