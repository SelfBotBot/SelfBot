package voice

import (
	"selfbot/discord/command/handlers"
	"selfbot/discord/feedback"
	"selfbot/discord/voice"

	"github.com/bwmarrin/discordgo"
)

var _ handlers.Handler = &PlayHandler{}

type PlayHandler struct {
	VoiceManager voice.Manager
}

func (h *PlayHandler) Handle(s *discordgo.Session, m *discordgo.MessageCreate, args ...string) error {
	if len(args) != 1 {
		// TODO USAGE
		return feedback.ErrorSoundNotFound
	}

	if err := h.VoiceManager.Play(m.GuildID, args[0]); err != nil {
		return err
	}
	go s.ChannelMessageDelete(m.ChannelID, m.ID)
	return nil
}

func (h PlayHandler) ShouldReact() bool {
	return false
}
