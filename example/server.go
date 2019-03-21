package main

import (
	"github.com/boxgo/rest"
	"github.com/gin-gonic/gin"
	"github.com/go-openapi/spec"
)

func newRestServer() *rest.Server {
	server := rest.NewServer()

	server.Mount(&PluginX{}, &PluginY{})

	op := spec.
		NewOperation("1").
		WithTags("test1").
		WithSummary("你好").
		WithDefaultResponse(spec.NewResponse().WithDescription("msg")).
		AddParam(spec.QueryParam("id").Typed("int", "csv").WithDefault("1").WithMaximum(100, false).WithMinimum(0, false).AsRequired())
	op.AddExtension("x-box-x", "abc")

	server.Describe(rest.API{
		Method: "get",
		Path:   "/test1",
		// Operation: op,
	}, func(ctx *gin.Context) {
		ctx.JSON(200, map[string]string{"test1": "test1"})
	})

	op1 := spec.
		NewOperation("1").
		WithTags("test2").
		WithSummary("你好1").
		WithDefaultResponse(spec.NewResponse().WithDescription("msg")).
		AddParam(spec.QueryParam("id").Typed("int", "csv").WithDefault("1").WithMaximum(100, false).WithMinimum(0, false).AsRequired())
	op1.AddExtension("x-box-y", map[string]string{"time": "123"})

	server.Describe(rest.API{
		Method:    "get",
		Path:      "/test2",
		Operation: op1,
	}, func(ctx *gin.Context) {
		ctx.JSON(200, map[string]string{"test2": "test2"})
	})

	op2 := spec.
		NewOperation("1").
		WithTags("test3").
		WithSummary("你好1").
		WithDefaultResponse(spec.NewResponse().WithDescription("msg")).
		AddParam(spec.QueryParam("id").Typed("int", "csv").WithDefault("1").WithMaximum(100, false).WithMinimum(0, false).AsRequired())
	op2.AddExtension("x-box-y", map[string]string{"time": "123"})

	server.Describe(rest.API{
		Method:    "post",
		Path:      "/test3",
		Operation: op2,
	}, func(ctx *gin.Context) {
		ctx.JSON(200, map[string]string{"test3": "test3"})
	})

	return server
}
