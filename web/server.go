package web

import (
	"fmt"
	"html/template"
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

	closers []func() error
}

// Handler an interface for sections of the site that handle.
type Handler interface {
	RegisterHandlers() error
}

func NewServer(l zerolog.Logger, rdb *redis.Client, cfg config.Config, voiceManager *voice.Manager) (*Server, error) {
	var ret = &Server{
		l:            l,
		cfg:          cfg,
		rdb:          rdb,
		voiceManager: voiceManager,
	}

	ret.launcher = ServerLauncher{
		s:       ret,
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

	discordOAuthG := &DiscordOAuth{s: ret}
	discordOAuthG.RegisterHandlers()

	pagesG := &Pages{s: ret}
	pagesG.RegisterHandlers()

	boardG := &Board{s: ret}
	boardG.RegisterHandlers()

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
	for _, v := range s.closers {
		if err := v(); err != nil {
			s.l.Error().Err(err).Msg("Unable to close closer")
		}
	}

	_ = s.launcher.Close() // launcher Close doesn't actually return an error.
	return nil
}
