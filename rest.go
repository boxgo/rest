package rest

import (
	"context"
	"net/http"
	"strings"

	"github.com/boxgo/swaggerfiles"
	"github.com/gin-gonic/gin"
)

type (
	// Server Powered by gin
	Server struct {
		Addr string `config:"addr" help:"Rest server listen addr. default is :8080"`
		Doc  bool   `config:"doc" help:"Whether to generate an api document"`
		Mode string `config:"mode" help:"Gin mode: debug,release,test. default is release"`

		engine *gin.Engine
		server *http.Server
		spec   *Spec
		plugin *PluginRegister
	}

	// API description information
	API struct {
		Method      string
		Path        string
		Tags        []string
		Summary     string
		Description string
		Deprecated  bool
		Parameters  Parameters
		Responses   Responses
		Extension   map[string]interface{}
		Handlers    gin.HandlersChain `json:"-"`
	}

	apiURL struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}
)

// Name box name
func (server *Server) Name() string {
	return "rest"
}

// ConfigWillLoad before config load
func (server *Server) ConfigWillLoad(ctx context.Context) {
	if server.engine == nil {
		panic("Rest server must new by NewServer function.")
	}
}

// ConfigDidLoad after config load
func (server *Server) ConfigDidLoad(ctx context.Context) {
	if server.Addr == "" {
		server.Addr = ":8080"
	}

	if server.Mode == "" {
		server.Mode = gin.ReleaseMode
	}

	gin.SetMode(server.Mode)
}

// Mount mount a plugin
func (server *Server) Mount(plugins ...Plugin) *Server {
	server.plugin.Mount(plugins...)

	return server
}

// Use attaches a global middleware to the router.
// ie. the middleware attached though Use() will be included in the handlers chain for every single request.
// Even 404, 405, static files...
// For example, this is the right place for a logger or error management middleware.
func (server *Server) Use(middleware ...gin.HandlerFunc) {
	server.engine.Use(middleware...)
}

// ServeHTTP conforms to the http.Handler interface.
func (server *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	server.engine.ServeHTTP(w, req)
}

// Describe a api
func (server *Server) Describe(apis ...API) {
	for _, api := range apis {
		pluginHandlers := server.plugin.Middlewarify(api)
		finalHandlers := append(pluginHandlers, api.Handlers...)

		server.spec.DescribeAPI(api.Path, api.Method, Operation{
			Tags:        api.Tags,
			Summary:     api.Summary,
			Description: api.Description,
			Deprecated:  api.Deprecated,
			Parameters:  api.Parameters,
			Responses:   api.Responses,
		})

		server.engine.Handle(strings.ToUpper(api.Method), api.Path, finalHandlers...)
	}
}

// Serve box serve handler
func (server *Server) Serve(ctx context.Context) error {
	server.server = &http.Server{
		Addr:    server.Addr,
		Handler: server.engine,
	}

	server.serveDoc()

	if err := server.server.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}

	return nil
}

// Shutdown box shutdown handler
func (server *Server) Shutdown(ctx context.Context) error {
	return server.server.Shutdown(ctx)
}

// serveDoc serve swagger api doc
func (server *Server) serveDoc() {
	if !server.Doc {
		return
	}

	apiJSONURL := "/swagger/api.json"
	apiCfgURL := "/swagger/config.json"

	server.engine.GET("/swagger/*any", func(ctx *gin.Context) {
		switch ctx.Request.URL.Path {
		case apiJSONURL:
			ctx.JSON(200, server.spec)
		case apiCfgURL:
			ctx.JSON(200, map[string]interface{}{
				"urls":                   []apiURL{apiURL{Name: "api-doc", URL: apiJSONURL}},
				"displayOperationId":     true,
				"displayRequestDuration": true,
			})
		default:
			gin.WrapH(swaggerfiles.Handler)(ctx)
		}
	})
}

// NewServer new a rest server
func NewServer() *Server {
	server := &Server{
		engine: gin.New(),
		plugin: NewPluginRegister(),
		spec:   NewSpec(),
	}

	return server
}
