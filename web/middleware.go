package web

import (
	"context"
	"fmt"
	"regexp"
	"selfbot/web/key"
	"strings"
	"time"

	"github.com/gorilla/sessions"
	"github.com/rbcervilla/redisstore/v8"

	"github.com/SilverCory/gin-redisgo-cooldowns"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
	"github.com/utrack/gin-csrf"
)

var allowedInRegisterRegex = regexp.MustCompile(`(?i)^(/(logout|register|tos|privacy|((js|css|img|auth)/*.)))|/$`)

const CSP = `
default-src 'self';
img-src 'self' https://cdnjs.cloudflare.com/ https://placekitten.com/ https://cdn.discordapp.com/;
script-src 'self' https://cdnjs.cloudflare.com/ajax/libs/cookieconsent2/3.0.3/cookieconsent.min.js 'sha256-SplWdsqEBp8LjzZSKYaEfDXhXSi0/oXXxAnQSYREAuI=';
style-src 'self' https://cdnjs.cloudflare.com/ajax/libs/cookieconsent2/3.0.3/cookieconsent.min.css 'unsafe-inline';
`

type Middleware struct {
	s *Server
}

var m *Middleware

func (s *Server) AddPreMiddleware() (err error) {
	m = &Middleware{s}

	if err = m.setupSessions(); err != nil {
		return
	}
	return
}

func (s *Server) AddPostMiddleware() (err error) {
	m.setupCors()
	//m.setupCsrf()
	m.setupSecurity()
	m.setupRegisterRedirect()
	m.setupIPCooldowns()

	return
}

func (m *Middleware) setupCsrf() {
	m.s.e.Use(csrf.Middleware(csrf.Options{
		Secret: m.s.cfg.Web.CSRFSecret,
		ErrorFunc: func(c *gin.Context) {

			if c.Request.URL.Path == "/CSPReport" || strings.HasPrefix(c.Request.URL.Path, "/alexamemes/") {
				return
			}

			c.String(400, "CSRF token mismatch")
			c.Abort()
		},
	}))
}

func (m *Middleware) setupSessions() (err error) {
	conf := m.s.cfg.Redis
	if !(conf.Enabled) {
		m.s.sessionStore = sessions.NewCookieStore([]byte("dankest_selfbot_ever")) // not a good idea lol. (we use redis...)
		return
	}

	redisStore, err := redisstore.NewRedisStore(context.Background(), m.s.rdb)
	if err != nil {
		return fmt.Errorf("web: setup sessions: unable to create redis store: %w", err)
	}
	redisStore.KeyPrefix("selfbot.sessions.sesh:")
	redisStore.Options(sessions.Options{
		Secure:   true,
		Domain:   "sb.cory.red",
		MaxAge:   86400 * 60,
		HttpOnly: true,
	})
	m.s.sessionStore = redisStore

	m.s.e.Use(func(ctx *gin.Context) {
		sesh, _ := m.s.sessionStore.Get(ctx.Request, key.SessionCore)
		for k, v := range sesh.Values {
			ctx.Set(k.(string), v)
		}

		ctx.Next()

		// Reload just to ensure we have the latest version.
		sesh, _ = m.s.sessionStore.Get(ctx.Request, key.SessionCore)
		for k, v := range ctx.Keys {
			if key.IsContext(k) {
				sesh.Values[k] = v
			}
		}
		if err := sesh.Save(ctx.Request, ctx.Writer); err != nil {
			panic(err)
		}
	})

	return nil
}

func (m *Middleware) setupIPCooldowns() {
	m.s.e.Use(gin_redisgo_cooldowns.NewRateLimit(m.s.rdb, "selfbot.cooldown.general.ip:", 100, time.Second*5, nil))
}

func (m *Middleware) setupCors() {
	if gin.IsDebugging() {
		return
	}

	m.s.e.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://cdnjs.cloudflare.com", "https://placekitten.com", "https://sb.cory.red"},
		AllowMethods:     []string{"GET"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))
}

func (m *Middleware) setupSecurity() {
	sec := secure.New(secure.Options{
		AllowedHosts:            []string{"sb.cory.red"},
		SSLRedirect:             false,
		SSLTemporaryRedirect:    false,
		SSLHost:                 "sb.cory.red",
		STSSeconds:              86400,
		STSIncludeSubdomains:    true,
		STSPreload:              true,
		ForceSTSHeader:          true,
		FrameDeny:               true,
		CustomFrameOptionsValue: "SAMEORIGIN",
		ContentTypeNosniff:      true,
		BrowserXssFilter:        true,
		ContentSecurityPolicy:   CSP,
		HostsProxyHeaders:       []string{"X-Forwarded-For"},

		IsDevelopment: gin.IsDebugging(),
	})

	secureFunc := func() gin.HandlerFunc {
		return func(c *gin.Context) {
			err := sec.Process(c.Writer, c.Request)

			// If there was an error, do not continue.
			if err != nil {
				c.Abort()
				return
			}

			// Avoid header rewrite if response is a redirection.
			if status := c.Writer.Status(); status > 300 && status < 399 {
				c.Abort()
			}
		}
	}()

	m.s.e.Use(secureFunc)
}

func (m *Middleware) setupRegisterRedirect() {
	m.s.e.Use(func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		if allowedInRegisterRegex.MatchString(path) {
			ctx.Next()
			return
		}

		//userIf, ok1 := ctx.Get(key.ContextUser)
		//user, ok2 := userIf.(goth.User)
		//if !ok1 || !ok2 {
		//	ctx.Next()
		//	return
		//}

		//user.
		//fmt.Println(user)
		//if !user.Agreed {
		//	sess.Set(key.ContextRedirect, path)
		//	ctx.Redirect(302, "/register")
		//	return
		//}

		ctx.Next()
	})
}
