package web

import (
	"fmt"
	"html/template"
	"net/http"
	"selfbot/config"
	"selfbot/discord/voice"
	"selfbot/sound"
	"selfbot/web/easter_egg"

	"github.com/gin-contrib/static"

	"github.com/gorilla/sessions"

	ginzerolog "github.com/dn365/gin-zerolog"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog"
)

type Server struct {
	l zerolog.Logger
	e *gin.Engine

	cfg      config.Config
	launcher ServerLauncher
	rdb      *redis.Client

	query        sound.Store
	sessionStore sessions.Store
	voiceManager *voice.Manager
}

// Handler an interface for sections of the site that handle.
type Handler interface {
	RegisterHandlers() error
}

func NewServer(l zerolog.Logger, rdb *redis.Client, cfg config.Config, voiceManager *voice.Manager) (Server, error) {
	var ret = Server{
		l:            l,
		cfg:          cfg,
		rdb:          rdb,
		voiceManager: voiceManager,
	}

	ret.launcher = ServerLauncher{
		s:       &ret,
		closers: make(map[string]func() error),
	}
	ret.e = gin.New()

	// Load the HTML templates
	// Templating
	ret.e.SetFuncMap(template.FuncMap{
		"comments": func(s string) template.HTML { return template.HTML(s) },
		"ASCII":    easter_egg.GetAscii,
	})
	ret.e.LoadHTMLGlob(cfg.Web.TemplateGlob)

	// Static files to load
	ret.e.Use(static.Serve("/", static.LocalFile(cfg.Web.StaticFilePath, false)))

	ret.e.Use(ginzerolog.Logger("gin"), gin.Recovery())
	_ = ret.AddPreMiddleware()
	ret.SetupDiscordOAuth()

	v1 := ret.e.Group("/api/v1/")
	{
		//v1.GET("/tree/:user", ret.v1TreeGETHandle)
		//v1User := v1.Group("/user/:user")
		//{
		//	v1User.GET("/", ret.v1UserGETHandle)
		//}
		v1.GET("/", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, map[string]string{"working": "true"})
		})
	}

	discordAuth := ret.e.Group("/auth/discord")
	{
		discordAuth.Use(discordOAuthContextSetter)
		discordAuth.GET("/", ret.oauthDiscordIndex)
		discordAuth.GET("/callback", ret.oauthDiscordCallback)
	}

	pages := &Pages{s: &ret}
	pages.RegisterHandlers()

	_ = ret.AddPostMiddleware()

	return ret, nil
}

func (s *Server) Start(listenAddr string, tls bool) error {
	var err error
	if !tls {
		err = s.e.Run(listenAddr)
	} else {
		err = s.launcher.RunAutoTLS()
	}

	if err != nil {
		return fmt.Errorf("web: server start: %w", err)
	}
	return nil
}

func (s *Server) Close() error {
	_ = s.launcher.Close() // launcher Close doesn't actually return an error.
	return nil
}
