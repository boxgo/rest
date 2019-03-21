package rest

import (
	"strings"

	"github.com/boxgo/kit/logger"
	"github.com/gin-gonic/gin"
)

const (
	// extensionPrefix prefix
	extensionPrefix = "x-box-"
)

type (
	// Plugin rest api plugin
	Plugin interface {
		Name() string
		Extensions() []string
		Middleware(api API, ext interface{}) gin.HandlerFunc
	}

	// PluginRegister plugin register
	PluginRegister interface {
		Mount(plugins ...Plugin)
		Middlewarify(api API) gin.HandlersChain
	}

	pluginRegister struct {
		plugins []Plugin
	}
)

// Mount mount plugins to register
func (register *pluginRegister) Mount(plugins ...Plugin) {
	for _, plugin := range plugins {
		logger.Default.Debugw("Plugin Mount", "name", plugin.Name(), "extension", plugin.Extensions())
	}

	register.plugins = append(register.plugins, plugins...)
}

// Middlewarify convert plugin to middleware
func (register *pluginRegister) Middlewarify(api API) gin.HandlersChain {
	handlers := gin.HandlersChain{}

	for _, plugin := range register.plugins {
		for _, extension := range plugin.Extensions() {
			if strings.Index(extension, extensionPrefix) != 0 {
				continue
			}

			if api.Operation == nil {
				continue
			}

			if ext, ok := api.Operation.Extensions[extension]; ok {
				logger.Default.Debugw("Plugin middlewarify", "name", plugin.Name(), "extension", plugin.Extensions(), "extData", ext)
				handlers = append(handlers, plugin.Middleware(api, ext))
			}
		}
	}

	return handlers
}
