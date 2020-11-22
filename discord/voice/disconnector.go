package voice

import (
	"github.com/bwmarrin/discordgo"
)

func (m *Manager) handleVoiceStateUpdate(s *discordgo.Session, vsu *discordgo.VoiceStateUpdate) {
	sesh, ok := m.SessionByGuild[vsu.GuildID]
	if !ok {
		return
	}

	sGuild, err := s.State.Guild(vsu.GuildID)
	if err != nil {
		m.l.Error().Err(err).Msg("state guild failed")
		return
	}

	var found = false
	for _, v := range sGuild.VoiceStates {
		if v.UserID == vsu.UserID || v.UserID == s.State.User.ID {
			continue
		}

		member, err := s.State.Member(vsu.GuildID, vsu.UserID)
		if err == discordgo.ErrStateNotFound || err == discordgo.ErrNilState {
			m.l.Warn().
				Str("guild_id", vsu.GuildID).
				Str("user_id", vsu.UserID).
				Msg("getting user via api.")

			if st, err := s.User(vsu.UserID); err != nil {
				m.l.Error().Err(err).Msg("session user failed")
				continue
			} else if st.Bot {
				continue
			}
		} else if err != nil {
			m.l.Error().Err(err).Msg("state member failed")
			continue
		} else if member.User.Bot {
			continue
		}

		found = true // A user that's not a bot is in the channel.
	}

	if found {
		return
	}
	if err := sesh.Stop(); err != nil {
		m.l.Error().Err(err).Msg("voice session stop")
		return
	}
}
