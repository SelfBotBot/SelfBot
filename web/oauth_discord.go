package web

import (
	"context"
	"net/http"

	"github.com/markbates/goth"
	discordOauth "github.com/markbates/goth/providers/discord"

	"github.com/bwmarrin/discordgo"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
)

// oauthContextSetter will set the oauth context for all requests for discord.
func discordOAuthContextSetter(ctx *gin.Context) {
	ctx.Request = ctx.Request.WithContext(
		context.WithValue(
			ctx.Request.Context(),
			"provider",
			"discord",
		),
	)
}

func (s *Server) SetupDiscordOAuth() {
	goth.UseProviders(
		discordOauth.New(
			s.cfg.DiscordOAuth.Key,
			s.cfg.DiscordOAuth.Secret,
			s.cfg.DiscordOAuth.Callback,
			discordOauth.ScopeIdentify,
			discordOauth.ScopeEmail,
			discordOauth.ScopeGuilds,
		),
	)
	gothic.Store = s.sessionStore
}

func (s *Server) oauthDiscordIndex(ctx *gin.Context) {
	if gothUser, err := gothic.CompleteUserAuth(ctx.Writer, ctx.Request); err == nil {
		ctx.JSON(http.StatusOK, gothUser)
		return
	} else {
		gothic.BeginAuthHandler(ctx.Writer, ctx.Request)
	}
}

func (s *Server) oauthDiscordCallback(ctx *gin.Context) {
	user, err := gothic.CompleteUserAuth(ctx.Writer, ctx.Request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, []string{"oauthComplete", err.Error()})
		return
	}
	_, err = discordgo.New("Bearer " + user.AccessToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, []string{"discordgo new", err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)
}
