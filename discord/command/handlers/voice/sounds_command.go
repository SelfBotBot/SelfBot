package voice

import (
	"selfbot/discord/command/handlers"
	"selfbot/discord/discordio"
	"selfbot/discord/voice"

	"github.com/bwmarrin/discordgo"
)

var _ handlers.Handler = &SoundsHandler{}

type SoundsHandler struct {
	VoiceManager voice.Manager
}

func (h *SoundsHandler) Handle(s *discordgo.Session, m *discordgo.MessageCreate, args ...string) error {
	writer := discordio.NewMessageWriter(s, m)
	writer.CodeBlock = false

	writer.Write([]byte("Here's a list of available sounds!"))
	for _, v := range h.VoiceManager.ListSounds() {
		writer.Write([]byte("`/play " + v + "`"))
	}

	return writer.Close()
}

func (h SoundsHandler) ShouldReact() bool {
	return true
}
