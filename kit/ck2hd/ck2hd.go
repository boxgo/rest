package ck2hd

import (
	"bytes"
	"context"

	"github.com/boxgo/box/minibox"
	"github.com/gin-gonic/gin"
)

type (
	// CookieHacker 将cookie转换为http header
	CookieHacker struct {
		CookieName string `config:"cookieName"` // 原始的cookie名称
		HeaderName string `config:"headerName"` // 目标的header名称

		app  minibox.App
		name string
	}

	bodyWriter struct {
		gin.ResponseWriter
		body       *bytes.Buffer
		headerName string
	}
)

var (
	// Default 默认的cookie转换header功能
	Default = New("cookie2header")
)

// Name 配置文件名称
func (c *CookieHacker) Name() string {
	return "middleware." + c.name
}

// Hijacker support session id by header
func (c *CookieHacker) Hijacker() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 如果cookie不存在，那么直接返回了，不进行hack
		cookie, _ := ctx.Cookie(c.CookieName)
		if cookie != "" {
			ctx.Next()
			return
		}

		// 请求header中有token存在，那么设置到请求的cookie中方便将session解析出来
		token := ctx.GetHeader(c.HeaderName)
		if token != "" {
			ctx.Request.Header.Add("Cookie", token)
		}

		// hack ResponseWriter
		bodyWriter := &bodyWriter{body: bytes.NewBufferString(""), ResponseWriter: ctx.Writer, headerName: c.HeaderName}
		ctx.Writer = bodyWriter

		ctx.Next()
	}
}

// Exts 获取app信息
func (c *CookieHacker) Exts() []minibox.MiniBox {
	return []minibox.MiniBox{&c.app}
}

func (c *CookieHacker) ConfigWillLoad(context.Context) {

}

func (c *CookieHacker) ConfigDidLoad(context.Context) {
	if c.CookieName == "" {
		c.CookieName = c.app.AppName
	}
}

func (w bodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)

	setCookieValue := w.Header().Get("set-cookie")
	if setCookieValue != "" {
		w.Header().Add(w.headerName, setCookieValue)
	}

	return w.ResponseWriter.Write(b)
}

// New 返回一个CookieHacker
func New(name string) *CookieHacker {
	return &CookieHacker{
		name: name,
	}
}
