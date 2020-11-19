package discord

import (
	"fmt"
	"os"
	"selfbot/config"
	"selfbot/discord/command"
	"selfbot/discord/voice"

	"github.com/rs/zerolog"

	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	Session        *discordgo.Session
	VoiceManager   voice.Manager
	CommandManager command.Manager

	//Sounds     map[string][][]byte
	//Sessions   map[string]*Session
	//stopping   bool
	//infoModule *info.InfoModule
}

func NewBot(dConf config.Discord) (Bot, error) {
	var ret = Bot{}

	sesh, err := discordgo.New("Bot " + dConf.Token)
	if err != nil {
		return Bot{}, fmt.Errorf("discord NewBot: discordgo init: %w", err)
	}
	sesh.AddHandler(ret.ready)
	ret.Session = sesh

	voiceManager, err := voice.NewManager(sesh)
	if err != nil {
		return Bot{}, fmt.Errorf("discord NewBot: %w", err)
	}
	ret.VoiceManager = voiceManager

	commandManager, err := command.NewCommandManager(zerolog.New(os.Stdout), sesh, ret.VoiceManager)
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
		// TODO nOPE
		fmt.Println(err)
	}

	if err := b.VoiceManager.Close(); err != nil {
		fmt.Println(err)
	}

	return b.Session.Close()
}
