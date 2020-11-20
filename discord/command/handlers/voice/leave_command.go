package voice

import (
	"selfbot/discord/command/handlers"
	"selfbot/discord/voice"

	"github.com/bwmarrin/discordgo"
)

var _ handlers.Handler = &LeaveHandler{}

type LeaveHandler struct {
	VoiceManager voice.Manager
}

func (h *LeaveHandler) Handle(s *discordgo.Session, m *discordgo.MessageCreate, args ...string) error {
	return h.VoiceManager.Leave(m.GuildID)
}

func (h LeaveHandler) ShouldReact() bool {
	return false
}
