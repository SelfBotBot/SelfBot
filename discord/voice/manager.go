package voice

import (
	"fmt"
	"selfbot/discord/feedback"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Manager struct {
	Sounds           map[string][][]byte
	SessionByGuild   map[string]*Session
	SessionByChannel map[string]*Session

	stopping         bool
	shutdownHandlers []func()
}

func NewManager(s *discordgo.Session) (Manager, error) {
	var ret = Manager{}
	ret.shutdownHandlers = append(ret.shutdownHandlers, s.AddHandlerOnce(ret.handleVoiceStateUpdate))
	s.State.TrackVoice = true
	s.State.TrackMembers = true // TODO get rid of this.

	return Manager{}, nil
}

func (m *Manager) Join(s *discordgo.Session, guildID string, userID string) error {
	channelID, err := findUserVoice(s, userID, guildID)
	if err != nil {
		return fmt.Errorf("join voice: %w", err)
	}

	vs, err := NewSession(s, m, guildID, channelID)
	if err != nil {
		//s.ChannelMessageSend(c.ID, "Unable to join VC.\n```"+err.Error()+"```")
		_ = vs.Stop() // Call this incase we're still alive?
		return nil
	}

	go vs.StartLoop()
	m.SessionByChannel[channelID] = vs
	m.SessionByGuild[guildID] = vs
	go vs.SetBuffer(welcome)
	return nil
}

func (m *Manager) Leave(guildID string) error {
	sesh, ok := m.SessionByGuild[guildID]
	if !ok {
		return feedback.ErrorNotInVoice
	}

	if err := sesh.Stop(); err != nil {
		return fmt.Errorf("leave voice: %w", err)
	}
	return nil
}

func (m *Manager) Play(guildID, soundName string) error {
	sound, ok := m.Sounds[soundName]
	if !ok {
		return feedback.ErrorSoundNotFound
	}

	ses, ok := m.SessionByGuild[guildID]
	if !ok {
		return feedback.ErrorNotInVoice
	}
	ses.SetBuffer(sound)
	return nil
}

func (m *Manager) Close() error {
	m.stopping = true
	for _, v := range m.SessionByGuild {
		go v.Stop()
		time.Sleep(75 * time.Millisecond)
	}
	time.Sleep(2 * time.Second)
	return nil
}
