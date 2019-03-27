package main

import (
	"github.com/boxgo/rest"
	"github.com/gin-gonic/gin"
)

func newRestServer() *rest.Server {
	server := rest.NewServer()

	server.Mount(&PluginX{}, &PluginY{}, &PluginZ{})

	server.Describe(rest.API{
		Method:      "get",
		Path:        "/test1",
		Tags:        []string{"NB"},
		Summary:     "Summary",
		Description: "Describe",
		Parameters: rest.Parameters{
			rest.Parameter{
				Name:        "id",
				In:          "query",
				Type:        "string",
				Description: "id",
			},
		},
		Extension: map[string]interface{}{
			"x-box-y": 1,
		},
		Handlers: gin.HandlersChain{func(ctx *gin.Context) {
			ctx.JSON(200, map[string]string{"test1": "test1"})
		}},
	})

	return server
}
