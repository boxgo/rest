package main

import (
	"github.com/boxgo/kit/logger"
	"github.com/boxgo/rest"
	"github.com/gin-gonic/gin"
)

type (
	PluginX struct{}
	PluginY struct{}
)

func (p *PluginX) Name() string {
	return "PluginX"
}

func (p *PluginX) Extensions() []string {
	return []string{"x-box-x"}
}

func (p *PluginX) Middleware(api rest.API, ext interface{}) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		logger.Default.Infow("x-box-x", "api", api, "ext", ext)

		ctx.JSON(200, "x-box-x 留下买路钱")
		ctx.Abort()
		ctx.Next()
	}
}

func (p *PluginY) Name() string {
	return "PluginY"
}

func (p *PluginY) Extensions() []string {
	return []string{"x-box-y"}
}

func (p *PluginY) Middleware(api rest.API, ext interface{}) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		logger.Default.Infow("x-box-y", "api", api, "ext", ext)

		ctx.JSON(200, "x-box-y 留下买路钱")
		ctx.Abort()
		ctx.Next()
	}
}
