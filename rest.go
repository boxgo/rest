package rest

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-openapi/spec"
)

type (
	// Server Powered by gin
	Server struct {
		Addr string `config:"addr" desc:"Rest server listen addr. default is :8080"`
		Doc  bool   `config:"doc" desc:"Whether to generate an api document"`
		Mode string `config:"mode" desc:"Gin mode: debug,release,test. default is release"`

		engine *gin.Engine
		server *http.Server
		spec   OpenAPI
		plugin PluginRegister
	}

	// API description information
	API struct {
		Method    string
		Path      string
		Operation *spec.Operation
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
func (server *Server) Describe(api API, handlers ...gin.HandlerFunc) {
	pluginHandlers := server.plugin.Middlewarify(api)
	finalHandlers := append(pluginHandlers, handlers...)

	server.engine.Handle(strings.ToUpper(api.Method), api.Path, finalHandlers...)
	server.spec.RegisterOperation(api)
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

	server.engine.GET("/swagger/index.html", func(ctx *gin.Context) {
		ctx.Data(200, "text/html", []byte(swaggerHTML))
	})
	server.engine.GET("/swagger/api.json", func(ctx *gin.Context) {
		ctx.JSON(200, server.spec)
	})
	server.engine.GET("/swagger/config.json", func(ctx *gin.Context) {
		ctx.JSON(200, map[string]interface{}{
			"urls":                   []apiURL{apiURL{Name: "api-doc", URL: "/swagger/api.json"}},
			"displayOperationId":     true,
			"displayRequestDuration": true,
		})
	})
}

// NewServer new a rest server
func NewServer() *Server {
	server := &Server{
		engine: gin.New(),
		spec:   OpenAPI{},
		plugin: &pluginRegister{},
	}

	return server
}