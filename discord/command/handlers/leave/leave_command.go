package leave

import (
	"selfbot/discord/command/handlers"
	"selfbot/discord/voice"

	"github.com/bwmarrin/discordgo"
)

var _ handlers.Handler = &Handler{}

type Handler struct {
	VoiceManager voice.Manager
}

func (h *Handler) Handle(s *discordgo.Session, m *discordgo.MessageCreate, args ...string) error {
	return h.VoiceManager.Leave(m.GuildID)
}

func (h Handler) ShouldReact() bool {
	return false
}
