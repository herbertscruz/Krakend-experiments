package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io"
	"net/http"
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

func (p *CustomPlugin) Bootstrap(w http.ResponseWriter, req *http.Request) {
	// The path matches, it has to be hijacked and no call to the backend happens.
	// The path is the the call to the backend, not the original request by the user.
	if req.URL.Path == p.path {
		w.Header().Add("Content-Type", "application/json")
		// Return a custom JSON object:
		res := map[string]string{"message": html.EscapeString(req.URL.Path)}
		b, _ := json.Marshal(res)
		w.Write(b)
		logger.Debug("request:", html.EscapeString(req.URL.Path))

		return
	}

	// If the requested path is not what we defined, continue.
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Copy headers, status codes, and body from the backend to the response writer
	for k, hs := range resp.Header {
		for _, h := range hs {
			w.Header().Add(k, h)
		}
	}
	w.WriteHeader(resp.StatusCode)
	if resp.Body == nil {
		return
	}
	io.Copy(w, resp.Body)
	resp.Body.Close()
}
