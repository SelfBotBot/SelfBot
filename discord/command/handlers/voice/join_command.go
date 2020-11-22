package voice

import (
	"selfbot/discord/command/handlers"
	"selfbot/discord/voice"

	"github.com/bwmarrin/discordgo"
)

var _ handlers.Handler = &JoinHandler{}

type JoinHandler struct {
	VoiceManager *voice.Manager
}

func (h *JoinHandler) Handle(s *discordgo.Session, m *discordgo.MessageCreate, args ...string) error {
	return h.VoiceManager.Join(s, m.GuildID, m.Author.ID)
}

func (h *JoinHandler) ShouldReact() bool {
	return false
}
