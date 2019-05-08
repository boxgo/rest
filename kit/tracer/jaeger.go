package tracer

import (
	"context"
	"runtime"

	"github.com/boxgo/box/minibox"
	"github.com/boxgo/tracer"
	"github.com/gin-gonic/gin"
	opentracing "github.com/opentracing/opentracing-go"
)

type (
	// Tracer rest
	Tracer struct {
		name   string
		tracer *tracer.Tracer
	}
)

var (
	// Default tracer
	Default = New("middleware.tracer")
)

// Name config name
func (t *Tracer) Name() string {
	return t.name
}

// Exts app
func (t *Tracer) Exts() []minibox.MiniBox {
	return []minibox.MiniBox{t.tracer}
}

// ConfigWillLoad before load
func (t *Tracer) ConfigWillLoad(context.Context) {

}

// ConfigDidLoad after load
func (t *Tracer) ConfigDidLoad(context.Context) {

}

// Jaeger middleware
func (t *Tracer) Jaeger() gin.HandlerFunc {
	ctxSpanKey := "tracing-span"

	return func(ctx *gin.Context) {
		var span opentracing.Span
		method := ctx.Request.Method
		path := ctx.Request.URL.Path
		operation := method + " " + path

		if ctxSpan, ok := ctx.Get(ctxSpanKey); ok {
			span = startSpanWithParent(ctxSpan.(opentracing.Span).Context(), operation, method, path)
		} else {
			span = startSpanWithHeader(t.tracer.Tracer(), &ctx.Request.Header, operation, method, path)
		}

		ctx.Set(ctxSpanKey, span)
		defer span.Finish()

		ctx.Next()

		span.SetTag("http.url", ctx.Request.URL.Path)
		span.SetTag("http.method", ctx.Request.Method)
		span.SetTag("http.status_code", ctx.Writer.Status())
		span.SetTag("requestId", ctx.Value("requestId"))
		span.SetTag("uid", ctx.Value("uid"))
		span.SetTag("current-goroutines", runtime.NumGoroutine())
	}
}

// New a tracer
func New(name string, ts ...*tracer.Tracer) *Tracer {
	t := &Tracer{
		name: name,
	}

	if len(ts) == 0 {
		t.tracer = tracer.Default
	} else {
		t.tracer = ts[0]
	}

	return t
}
