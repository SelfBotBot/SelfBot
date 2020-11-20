package voice

import (
	"fmt"
	"selfbot/discord/feedback"

	"github.com/bwmarrin/discordgo"
)

func findUserVoice(s *discordgo.Session, guildID string, userID string) (channelID string, err error) {
	g, err := s.State.Guild(guildID)
	if err != nil {
		return "", fmt.Errorf(
			"findUserVoice: %w: ",
			feedback.Wrap(feedback.ErrorNoUserInVoice, err),
		)
	}

	for _, vs := range g.VoiceStates {
		if vs.UserID == userID {
			return vs.ChannelID, nil
		}
	}

	return "", feedback.ErrorNoUserInVoice
}
