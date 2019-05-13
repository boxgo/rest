package metrics

import (
	"context"
	"strconv"
	"time"

	"github.com/boxgo/box/minibox"
	"github.com/boxgo/metrics"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

type (
	// Metrics config
	Metrics struct {
		RequestTotalName           string `config:"requestTotalName"`
		RequestDurationSecondsName string `config:"requestDurationSeconds"`
		RequestSizeBytesName       string `config:"requestSizeBytes"`
		ResponseSizeBytesName      string `config:"responseSizeBytes"`

		name                string
		metrics             *metrics.Metrics
		requestURLMappingFn func(*gin.Context) string
		reqCounter          *prometheus.CounterVec
		reqDurationSummary  *prometheus.SummaryVec
		reqSizeSummary      *prometheus.SummaryVec
		resSizeSummary      *prometheus.SummaryVec
	}
)

var (
	// Default gin metrics
	Default = New("middleware.metrics")
	labels  = []string{"code", "retcode", "method", "url", "handler"}
)

// Name config
func (g *Metrics) Name() string {
	return g.name
}

// Exts app
func (g *Metrics) Exts() []minibox.MiniBox {
	return []minibox.MiniBox{g.metrics}
}

// ConfigWillLoad before load
func (g *Metrics) ConfigWillLoad(context.Context) {

}

// ConfigDidLoad after load
func (g *Metrics) ConfigDidLoad(context.Context) {
	g.requestURLMappingFn = urlMapping

	if g.RequestTotalName == "" {
		g.RequestTotalName = "http_request_total"
	}
	if g.RequestDurationSecondsName == "" {
		g.RequestDurationSecondsName = "http_request_duration_seconds"
	}
	if g.RequestSizeBytesName == "" {
		g.RequestSizeBytesName = "http_request_size_bytes"
	}
	if g.ResponseSizeBytesName == "" {
		g.ResponseSizeBytesName = "http_response_size_bytes"
	}

	g.reqCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: g.metrics.Namespace,
			Subsystem: g.metrics.Subsystem,
			Name:      g.RequestTotalName,
			Help:      "How many HTTP requests processed, partitioned by status code and HTTP method.",
		},
		labels,
	)

	g.reqDurationSummary = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: g.metrics.Namespace,
			Subsystem: g.metrics.Subsystem,
			Name:      g.RequestDurationSecondsName,
			Help:      "The HTTP request latencies in seconds.",
		},
		labels,
	)

	g.reqSizeSummary = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: g.metrics.Namespace,
			Subsystem: g.metrics.Subsystem,
			Name:      g.RequestSizeBytesName,
			Help:      "The HTTP request sizes in bytes.",
		},
		labels,
	)

	g.resSizeSummary = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: g.metrics.Namespace,
			Subsystem: g.metrics.Subsystem,
			Name:      g.ResponseSizeBytesName,
			Help:      "The HTTP response sizes in bytes.",
		},
		labels,
	)

	prometheus.MustRegister(g.reqCounter)
	prometheus.MustRegister(g.reqDurationSummary)
	prometheus.MustRegister(g.reqSizeSummary)
	prometheus.MustRegister(g.resSizeSummary)
}

// Handler is prometheus middleware of gin
func (g *Metrics) Handler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		reqSize := computeRequestSize(ctx.Request)

		ctx.Next()

		end := time.Now()
		elapsed := end.Sub(start)
		status := strconv.Itoa(ctx.Writer.Status())
		resSize := float64(ctx.Writer.Size())
		url := g.requestURLMappingFn(ctx)
		retcode := strconv.Itoa(ctx.GetInt("retcode"))
		labels := []string{status, retcode, ctx.Request.Method, url, ctx.HandlerName()}

		g.reqCounter.WithLabelValues(labels...).Inc()
		g.reqSizeSummary.WithLabelValues(labels...).Observe(reqSize)
		g.resSizeSummary.WithLabelValues(labels...).Observe(resSize)
		g.reqDurationSummary.WithLabelValues(labels...).Observe(elapsed.Seconds())
	}
}

// New a gin prom
func New(name string, ms ...*metrics.Metrics) *Metrics {

	if len(ms) == 0 {
		return &Metrics{
			name:    name,
			metrics: metrics.Default,
		}
	}

	return &Metrics{
		name:    name,
		metrics: ms[0],
	}
}
