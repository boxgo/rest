package logger

import (
	"context"
	"fmt"
	"time"

	lg "github.com/boxgo/logger"
	"github.com/gin-gonic/gin"
	"github.com/teris-io/shortid"
	"go.uber.org/zap"
)

type (
	// Logger gin中间件
	Logger struct {
		RequestBodyLimit  int `config:"requestBodyLimit"`
		RequestQueryLimit int `config:"requestQueryLimit"`
		ResponseBodyLimit int `config:"responseBodyLimit"`
	}
)

var (
	// Default logger
	Default = &Logger{}
)

// Name 日志中间件配置名称
func (l *Logger) Name() string {
	return "middleware.logger"
}

// ConfigWillLoad 配置文件将要加载
func (l *Logger) ConfigWillLoad(context.Context) {

}

// ConfigDidLoad 配置文件已经加载。做一些默认值设置
func (l *Logger) ConfigDidLoad(context.Context) {
	if l.RequestBodyLimit == 0 {
		l.RequestBodyLimit = 2000
	}

	if l.RequestQueryLimit == 0 {
		l.RequestQueryLimit = 2000
	}

	if l.ResponseBodyLimit == 0 {
		l.ResponseBodyLimit = 2000
	}
}

// Logger zap
func (l *Logger) Logger(logs ...*lg.Logger) gin.HandlerFunc {
	var log *lg.Logger
	if len(logs) == 0 {
		log = lg.Default
	} else {
		log = logs[0]
	}

	return func(ctx *gin.Context) {
		start := time.Now()
		requestID, _ := shortid.Generate()
		method := ctx.Request.Method
		path := ctx.Request.URL.Path
		query := ctx.Request.URL.RawQuery
		userAgent := ctx.Request.UserAgent()
		clientIP := ctx.ClientIP()
		body := readBody(ctx)

		if l.RequestQueryLimit < len(query) {
			query = query[:l.RequestQueryLimit]
		}
		if l.RequestBodyLimit < len(body) {
			body = body[:l.RequestBodyLimit]
		}

		ctx.Set("requestId", requestID)
		ctx.Set("traceBizId", fmt.Sprintf("%s %s", ctx.Request.Method, ctx.Request.URL.Path))

		log.TraceRaw(ctx).Info(">>>", []zap.Field{
			zap.String("method", method),
			zap.String("path", path),
			zap.String("ip", clientIP),
			zap.String("user-agent", userAgent),
			zap.String("query", query),
			zap.String("body", string(body)),
		}...)

		bodyWriter := newBodyWriter(ctx)
		ctx.Writer = bodyWriter

		ctx.Next()

		end := time.Now()
		latency := end.Sub(start)
		resp := bodyWriter.body.String()

		if l.ResponseBodyLimit < len(resp) {
			resp = resp[:l.ResponseBodyLimit]
		}

		for _, err := range ctx.Errors {
			log.TraceRaw(ctx).Info("xxx",
				zap.String("method", method),
				zap.String("path", path),
				zap.String("ip", clientIP),
				zap.String("user-agent", userAgent),
				zap.String("query", query),
				zap.String("body", string(body)),
				zap.Int("status", ctx.Writer.Status()),
				zap.String("err", err.Error()),
			)
		}

		log.TraceRaw(ctx).Info("<<<",
			[]zap.Field{
				zap.String("method", method),
				zap.String("path", path),
				zap.String("ip", clientIP),
				zap.String("user-agent", userAgent),
				zap.String("query", query),
				zap.String("body", string(body)),
				zap.Int("status", ctx.Writer.Status()),
				zap.Duration("latency", latency),
				zap.String("resp", resp),
			}...)
	}
}
