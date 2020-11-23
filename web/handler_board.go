package web

import (
	"errors"
	"net/http"
	"selfbot/discord/feedback"
	"selfbot/web/key"
	"selfbot/web/ws_player"

	"github.com/gorilla/websocket"

	"github.com/bwmarrin/discordgo"
	"github.com/markbates/goth"

	"selfbot/web/viewdata"

	"github.com/gin-gonic/gin"
)

var _ Handler = &Board{}

var upper = websocket.Upgrader{
	ReadBufferSize: 512,
}

type Board struct {
	s *Server

	Handler
}

func (b *Board) RegisterHandlers() error {
	boardRoute := b.s.e.Group("/board")
	{
		boardRoute.GET("/", b.handleBoard)
		boardRoute.GET("/ws/start", b.handleWSAuth)
		boardRoute.GET("/ws", b.handleWebsocket)
	}
	return nil
}

func (b *Board) handleBoard(ctx *gin.Context) {
	v := viewdata.Default(ctx)
	v.Set("Title", "Soundboard")
	var sounds = map[string]string{}
	for k, v := range b.s.voiceManager.Sounds {
		sounds[k.String()] = v.Name
	}
	v.Set("Sounds", sounds)
	v.HTML(200, "pages/board.html")
}

func (b *Board) handleWSAuth(ctx *gin.Context) {
	oauthUser, ok := ctx.Get(key.ContextUser)
	if !ok {
		ctx.Status(http.StatusForbidden)
		return
	}

	user := oauthUser.(goth.User)
	s, err := discordgo.New("Bearer " + user.AccessToken)
	if err != nil {
		ctx.String(500, err.Error())
		return
	}
	userGuilds, err := s.UserGuilds(100, "", "")
	if err != nil {
		ctx.String(500, err.Error())
		return
	}
	s.Close()
	s = nil // Hopefully gc gets this lol.

	var guildIDs []string
	for _, v := range userGuilds {
		guildIDs = append(guildIDs, v.ID)
	}

	gid, err := b.s.voiceManager.FindUserGuild(user.UserID, guildIDs)
	if err != nil {
		ctx.String(500, err.Error())
		return
	}

	err = b.s.voiceManager.Join(gid, user.UserID)
	if err != nil {
		if !errors.Is(err, feedback.ErrorAlreadyInVoice) {
			ctx.String(500, err.Error())
			return
		}
	}

	authKey := ws_player.StartAuth(ws_player.NewVoicePlayer(gid, user, b.s.voiceManager))
	ctx.JSON(200, ws_player.PlayRequest{ID: authKey})
}

func (b *Board) handleWebsocket(ctx *gin.Context) {
	wsConn, err := upper.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.String(500, err.Error())
		return
	}

	player, ok := ws_player.IsAuthed(wsConn)
	if !ok {
		ctx.String(http.StatusUnauthorized, "Not authenticated")
		return
	}

	player.Start(wsConn)
}
