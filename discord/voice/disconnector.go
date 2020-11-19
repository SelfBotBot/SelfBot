package voice

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func (m *Manager) handleVoiceStateUpdate(s *discordgo.Session, vsu *discordgo.VoiceStateUpdate) {
	sesh, ok := m.SessionByChannel[vsu.ChannelID]
	if !ok {
		return
	}

	g, err := s.State.Guild(vsu.GuildID)
	if err != nil {
		fmt.Printf("%#v\n", s.State.Guilds)
		return
	}

	for _, vs := range g.VoiceStates {
		// Ignore not this channel
		// Ignore myself
		// Ignore deaf
		if vs.ChannelID != vsu.ChannelID || vs.UserID == s.State.User.ID || vs.Deaf || vs.SelfDeaf {
			continue
		}

		// Ignore bots.
		if u, err := s.State.Member(vs.GuildID, vs.UserID); err != nil {
			if err == discordgo.ErrStateNotFound {
				// TODO do I really want to handle this?
			}
			continue
		} else if u.User.Bot {
			continue
		}

		return // This is a normal user, we can return.
	}

	sesh.Stop()
}
