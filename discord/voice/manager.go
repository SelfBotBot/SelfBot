package voice

import (
	"fmt"
	"selfbot/discord/feedback"
	"time"

	"github.com/rs/zerolog"

	"github.com/bwmarrin/discordgo"
)

type Manager struct {
	l                zerolog.Logger
	Sounds           map[string][][]byte
	SessionByGuild   map[string]*Session
	SessionByChannel map[string]*Session

	stopping         bool
	shutdownHandlers []func()
}

func NewManager(l zerolog.Logger, s *discordgo.Session) (Manager, error) {
	var ret = Manager{
		l:                l.With().Str("owner", "VoiceManager").Logger(),
		SessionByGuild:   make(map[string]*Session),
		SessionByChannel: make(map[string]*Session),
		Sounds:           make(map[string][][]byte),
	}
	ret.shutdownHandlers = append(ret.shutdownHandlers, s.AddHandlerOnce(ret.handleVoiceStateUpdate))
	s.State.TrackVoice = true
	s.State.TrackMembers = true // TODO get rid of this.

	return ret, nil
}

func (m *Manager) Join(s *discordgo.Session, guildID string, userID string) error {
	channelID, err := findUserVoice(s, guildID, userID)
	if err != nil {
		return fmt.Errorf("join voice: %w", err)
	}

	vs, err := NewSession(s, m, guildID, channelID)
	if err != nil {
		_ = vs.Stop() // Call this incase we're still alive?
		return fmt.Errorf("join voice: %w", feedback.Wrap(feedback.ErrorFatalError, err))
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
		return fmt.Errorf("leave voice: %w", feedback.Wrap(feedback.ErrorFatalError, err))
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
