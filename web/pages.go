package web

import (
	"errors"
	"net/http"
	"selfbot/discord/feedback"

	uuid "github.com/satori/go.uuid"

	//"github.com/gin-contrib/sessions"
	"selfbot/web/viewdata"

	"github.com/gin-gonic/gin"
)

const PrivacyPolicy = `
<html><head><title>Privacy</title></head>
<body>
	<h1>privacy policy</h1>
	<p>We aren't currently logging any information as of yet however in the near future, your google account, amazon account and discord account informations will be stored and used for the usage of this.</p>
</body>
</html>'`

type Pages struct {
	s *Server
	Handler
}

func (p *Pages) RegisterHandlers() error {
	p.s.e.GET("/", p.handleIndex)
	p.s.e.GET("/tos", p.handleTos)
	p.s.e.GET("/board", p.handleBoard)
	p.s.e.POST("/board/play/:guildID/:soundID", p.handleSoundPlay)
	return nil
}

func (p *Pages) handlePanel(ctx *gin.Context) {
	//sess := sessions.Default(ctx)
	v := viewdata.Default(ctx)
	v.Set("Title", "Panel")
	v.HTML(200, "pages/index.html")
}

func (p *Pages) handleIndex(ctx *gin.Context) {
	v := viewdata.Default(ctx)
	v.Set("Title", "Home")
	v.HTML(200, "pages/index.html")
}

func (p *Pages) handleBoard(ctx *gin.Context) {
	v := viewdata.Default(ctx)
	v.Set("Title", "Soundboard")
	var sounds = map[string]string{}
	for k, v := range p.s.voiceManager.Sounds {
		sounds[k.String()] = v.Name
	}
	v.Set("Sounds", sounds)
	v.HTML(200, "pages/board.html")
}

func (p *Pages) handleSoundPlay(ctx *gin.Context) {
	guildID := ctx.Param("guildID")
	soundID := ctx.Param("soundID")
	err := p.s.voiceManager.Play(
		guildID,
		uuid.Must(uuid.FromString(soundID)),
	)
	if err != nil {
		if errors.Is(err, feedback.ErrorBotNotInVoice) {
			ctx.String(http.StatusUnprocessableEntity, feedback.ErrorBotNotInVoice.Message)
			return
		}
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (p *Pages) handleTos(ctx *gin.Context) {
	v := viewdata.Default(ctx)
	v.Set("Title", "Terms Of Service")
	v.HTML(200, "pages/tos.html")
}

func (p *Pages) handlePrivacyPolicy(ctx *gin.Context) {
	ctx.String(200, PrivacyPolicy)
}
