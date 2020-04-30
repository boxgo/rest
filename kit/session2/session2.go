package session2

import (
	"context"

	"github.com/boxgo/box/minibox"
	"github.com/boxgo/redisstore"
	"github.com/boxgo/redisstore/serializer"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
)

type (
	Session struct {
		CookieName string `config:"cookieName" help:"cookie name"`
		KeyPair    string `config:"keyPair" help:"cookie value encrypt key pair"`
		KeyPrefix  string `config:"keyPrefix" help:"redis save key prefix"`
		MaxLen     int    `config:"maxLen" help:"redis save value max len. -1 not limit"`
		app        minibox.App
	}
)

var (
	Default = &Session{}
)

func (s *Session) Name() string {
	return "middleware.session"
}

// Exts 获取app信息
func (s *Session) Exts() []minibox.MiniBox {
	return []minibox.MiniBox{&s.app}
}

// ConfigWillLoad 配置文件将要加载
func (s *Session) ConfigWillLoad(context.Context) {

}

// ConfigDidLoad 配置文件已经加载。做一些默认值设置
func (s *Session) ConfigDidLoad(context.Context) {
	if s.CookieName == "" {
		s.CookieName = s.app.AppName
	}

	if s.CookieName == "" {
		s.CookieName = "session"
	}

	if s.KeyPrefix == "" {
		s.KeyPrefix = s.app.AppName + "_session_"
	}

	if s.MaxLen == 0 {
		s.MaxLen = 10240
	}
}

func (s *Session) Cookie() gin.HandlerFunc {
	return sessions.Sessions(s.CookieName, cookie.NewStore([]byte(s.KeyPair)))
}

func (s *Session) Redis(client redis.UniversalClient, sessionSerializer ...serializer.SessionSerializer) gin.HandlerFunc {
	var ss serializer.SessionSerializer

	if len(sessionSerializer) == 0 {
		ss = &serializer.GobSerializer{}
	} else {
		ss = sessionSerializer[0]
	}

	st, _ := redisstore.NewStoreWithUniversalClient(
		client,
		redisstore.WithMaxLength(s.MaxLen),
		redisstore.WithKeyPrefix(s.KeyPrefix),
		redisstore.WithKeyPairs([]byte(s.KeyPair)),
		redisstore.WithSerializer(ss),
	)

	return sessions.Sessions(s.CookieName, &redisStore{
		RedisStore: st,
	})
}
