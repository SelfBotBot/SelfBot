package command

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog"

	"github.com/bwmarrin/discordgo"
)

// processCommand will split the message content into the command and the arguments.
// the bool will also be whether it has the prefix, if it doesn't it won't return anything.
func processCommand(content string) (bool, string, []string) {
	content = strings.TrimSpace(content)
	if !strings.HasPrefix(content, Prefix) {
		return false, "", nil
	}
	content = strings.TrimPrefix(content, Prefix)

	cmdParts := strings.Split(content, " ")
	if len(cmdParts) > 1 {
		return true, strings.ToLower(cmdParts[0]), cmdParts[1:]
	}
	return true, strings.ToLower(cmdParts[0]), nil
}

func React(s *discordgo.Session, m *discordgo.MessageCreate, emoji string) error {
	err := s.MessageReactionAdd(m.ChannelID, m.ID, emoji)
	if err != nil {
		return fmt.Errorf("react: %w", err)
	}
	return nil
}

func Unreact(s *discordgo.Session, m *discordgo.MessageCreate, emoji string) error {
	err := s.MessageReactionRemove(m.ChannelID, m.ID, emoji, s.State.User.ID)
	if err != nil {
		return fmt.Errorf("unreact: %w", err)
	}
	return nil
}

func LogDataMessage(ev *zerolog.Event, m *discordgo.MessageCreate) *zerolog.Event {
	return ev.
		Str("for_type", "MessageCreate").
		Str("user_id", m.Author.ID).
		Str("channel_id", m.ChannelID).
		Str("message_id", m.ID).
		Str("message_content", fmt.Sprintf("%.20s", m.Content))
}
