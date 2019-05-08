package rest

import (
	"github.com/boxgo/logger"
	"github.com/gin-gonic/gin"
)

type (
	// Plugin rest api plugin
	Plugin interface {
		Name() string
		Extensions() []string
		Middleware(api *API) gin.HandlerFunc
	}

	// PluginRegister plugin register
	PluginRegister struct {
		plugins []Plugin
	}
)

// Mount mount plugins to register
func (register *PluginRegister) Mount(plugins ...Plugin) {
	for _, plugin := range plugins {
		logger.Default.Debugw("Plugin Mount", "name", plugin.Name(), "extension", plugin.Extensions())
	}

	register.plugins = append(register.plugins, plugins...)
}

// Middlewarify convert plugin to middleware
func (register *PluginRegister) Middlewarify(api API) gin.HandlersChain {
	handlers := gin.HandlersChain{}

	for _, plugin := range register.plugins {
		for _, extension := range plugin.Extensions() {
			if extension == "parameters" || extension == "responses" {
				logger.Default.Debugw("Plugin middlewarify", "name", plugin.Name(), "extension", plugin.Extensions())
				handlers = append(handlers, plugin.Middleware(&api))
			} else if ext, ok := api.Extension[extension]; ok {
				logger.Default.Debugw("Plugin middlewarify", "name", plugin.Name(), "extension", plugin.Extensions(), "ext", ext)
				handlers = append(handlers, plugin.Middleware(&api))
			}
		}
	}

	return handlers
}

// NewPluginRegister new a plugin register
func NewPluginRegister() *PluginRegister {
	return &PluginRegister{}
}
