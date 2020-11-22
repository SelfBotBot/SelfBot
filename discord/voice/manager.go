package voice

import (
	"fmt"
	"selfbot/discord/feedback"
	"selfbot/sound"
	"time"

	uuid "github.com/satori/go.uuid"

	"github.com/rs/zerolog"

	"github.com/bwmarrin/discordgo"
)

type Manager struct {
	l              zerolog.Logger
	Sounds         map[uuid.UUID]sound.Sound
	SessionByGuild map[string]*Session

	stopping         bool
	shutdownHandlers []func()
}

func NewManager(l zerolog.Logger, s *discordgo.Session, soundStore sound.Store) (*Manager, error) {
	var ret = &Manager{
		l:              l.With().Str("owner", "VoiceManager").Logger(),
		SessionByGuild: make(map[string]*Session),
		Sounds:         make(map[uuid.UUID]sound.Sound),
	}
	ret.shutdownHandlers = append(ret.shutdownHandlers, s.AddHandler(ret.handleVoiceStateUpdate))
	s.State.TrackVoice = true
	s.State.TrackMembers = true // TODO get rid of this.

	resp, err := soundStore.ListSounds(sound.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("voice manager: new manager %w", err)
	}

	for _, v := range resp.SoundIDs {
		loaded, err := soundStore.LoadSound(v)
		if err != nil {
			ret.l.Err(err).Str("uuid", v.String()).Msg("unable to load sound.")
			continue
		}

		ret.l.Info().Str("name", loaded.Name).Msg("loaded sound.")
		ret.Sounds[loaded.ID] = loaded
	}

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
	m.SessionByGuild[guildID] = vs
	go vs.SetBuffer(welcome)
	return nil
}

func (m *Manager) Leave(guildID string) error {
	sesh, ok := m.SessionByGuild[guildID]
	if !ok {
		return feedback.ErrorBotNotInVoice
	}

	if err := sesh.Stop(); err != nil {
		return fmt.Errorf("leave voice: %w", feedback.Wrap(feedback.ErrorFatalError, err))
	}
	return nil
}

func (m *Manager) Play(guildID string, soundName uuid.UUID) error {
	s, ok := m.Sounds[soundName]
	if !ok {
		return feedback.ErrorSoundNotFound
	}

	ses, ok := m.SessionByGuild[guildID]
	if !ok {
		return feedback.ErrorBotNotInVoice
	}
	ses.SetBuffer(s.Data)
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
