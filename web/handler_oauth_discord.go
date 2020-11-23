package web

import (
	"context"
	"net/http"
	"selfbot/web/key"
	"strings"

	"github.com/markbates/goth"
	discordOauth "github.com/markbates/goth/providers/discord"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
)

const SessionKeyUser = "discord_user"

var _ Handler = &DiscordOAuth{}

type DiscordOAuth struct {
	s *Server
}

func (doa *DiscordOAuth) RegisterHandlers() error {
	doa.Setup()

	doa.s.e.GET("/logout", doa.logoutHandler)
	doa.s.e.GET("/login", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusFound, "/auth/discord")
	})

	discordAuth := doa.s.e.Group("/auth/discord")
	{
		discordAuth.Use(doa.contextSetter)
		discordAuth.GET("/", doa.indexHandler)
		discordAuth.GET("/callback", doa.callbackHandler)
	}
	return nil
}

func (doa *DiscordOAuth) Setup() {
	goth.UseProviders(
		discordOauth.New(
			doa.s.cfg.DiscordOAuth.Key,
			doa.s.cfg.DiscordOAuth.Secret,
			doa.s.cfg.DiscordOAuth.Callback,
			discordOauth.ScopeIdentify,
			discordOauth.ScopeEmail,
			discordOauth.ScopeGuilds,
		),
	)
	gothic.Store = doa.s.sessionStore
}

func (doa *DiscordOAuth) indexHandler(ctx *gin.Context) {
	if gothUser, err := gothic.CompleteUserAuth(ctx.Writer, ctx.Request); err == nil {
		ctx.JSON(http.StatusOK, gothUser)
		ctx.Set("User", gothUser)
		return
	} else {
		redirectTo := ctx.Param("redirectTo")
		if redirectTo == "" || !strings.HasPrefix(redirectTo, "/") {
			redirectTo = "/"
		}
		ctx.Set(key.ContextRedirect, redirectTo)
		gothic.BeginAuthHandler(ctx.Writer, ctx.Request)
	}
}

func (doa *DiscordOAuth) callbackHandler(ctx *gin.Context) {
	user, err := gothic.CompleteUserAuth(ctx.Writer, ctx.Request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, []string{"oauthComplete", err.Error()})
		return
	}

	ctx.Set(key.ContextUser, user)
	redirectTo := ctx.GetString(key.ContextRedirect)
	if redirectTo == "" || !strings.HasPrefix(redirectTo, "/") {
		redirectTo = "/"
	}
	ctx.Redirect(http.StatusFound, redirectTo)
}

func (doa *DiscordOAuth) logoutHandler(ctx *gin.Context) {
	gothic.Logout(ctx.Writer, ctx.Request)
	sesh, _ := doa.s.sessionStore.Get(ctx.Request, key.SessionCore)
	sesh.Values = map[interface{}]interface{}{}
	ctx.Keys = map[string]interface{}{}
	ctx.Redirect(http.StatusFound, "/")
}

// contextSetter will set the oauth context for all requests for discord.
func (doa *DiscordOAuth) contextSetter(ctx *gin.Context) {
	ctx.Request = ctx.Request.WithContext(
		context.WithValue(
			ctx.Request.Context(),
			"provider",
			"discord",
		),
	)
}
