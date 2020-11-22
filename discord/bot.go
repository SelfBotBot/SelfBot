package discord

import (
	"fmt"
	"selfbot/config"
	"selfbot/discord/command"
	"selfbot/discord/voice"
	"selfbot/sound"

	"github.com/rs/zerolog"

	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	l              zerolog.Logger
	Session        *discordgo.Session
	VoiceManager   *voice.Manager
	CommandManager *command.Manager
}

func NewBot(l zerolog.Logger, dConf config.Discord, soundStore sound.Store) (Bot, error) {
	var ret = Bot{
		l: l,
	}

	sesh, err := discordgo.New("Bot " + dConf.Token)
	if err != nil {
		return Bot{}, fmt.Errorf("discord NewBot: discordgo init: %w", err)
	}
	sesh.AddHandler(ret.ready)
	ret.Session = sesh

	voiceManager, err := voice.NewManager(l, sesh, soundStore)
	if err != nil {
		return Bot{}, fmt.Errorf("discord NewBot: %w", err)
	}
	ret.VoiceManager = voiceManager

	commandManager, err := command.NewCommandManager(l, sesh, ret.VoiceManager)
	if err != nil {
		return Bot{}, fmt.Errorf("discord NewBot: %w", err)
	}
	ret.CommandManager = commandManager

	return ret, nil
}

func (b *Bot) ready(s *discordgo.Session, _ *discordgo.Ready) {
	_ = s.UpdateStatus(0, "/soundboard | /sb")
}

func (b *Bot) Close() error {
	if err := b.CommandManager.Close(); err != nil {
		b.l.Error().Err(err).Msg("unable to close command manager.")
	}

	if err := b.VoiceManager.Close(); err != nil {
		b.l.Error().Err(err).Msg("unable to close voice manager.")
	}

	return b.Session.Close()
}
