package rest

import (
	"strings"

	"github.com/go-openapi/analysis"
	"github.com/go-openapi/spec"
)

type (
	OpenAPI struct {
		spec.Swagger
	}

	apiURL struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}
)

func (oai *OpenAPI) RegisterOperation(api API) {
	oai.Swagger.Swagger = "2.0"

	analysis.Mixin(&oai.Swagger, &spec.Swagger{
		SwaggerProps: spec.SwaggerProps{
			Paths: &spec.Paths{
				Paths: makePathItem(api),
			},
		},
	})
}

func makePathItem(api API) map[string]spec.PathItem {
	if api.Operation == nil {
		api.Operation = &spec.Operation{}
	}

	pathItem := spec.PathItemProps{}
	method := strings.ToUpper(api.Method)
	switch method {
	case "GET":
		pathItem = spec.PathItemProps{
			Get: api.Operation,
		}
	case "POST":
		pathItem = spec.PathItemProps{
			Post: api.Operation,
		}
	case "PUT":
		pathItem = spec.PathItemProps{
			Put: api.Operation,
		}
	case "DELETE":
		pathItem = spec.PathItemProps{
			Delete: api.Operation,
		}
	case "Options":
		pathItem = spec.PathItemProps{
			Options: api.Operation,
		}
	case "Head":
		pathItem = spec.PathItemProps{
			Head: api.Operation,
		}
	case "Patch":
		pathItem = spec.PathItemProps{
			Patch: api.Operation,
		}
	}

	return map[string]spec.PathItem{
		api.Path: spec.PathItem{
			PathItemProps: pathItem,
		},
	}
}
