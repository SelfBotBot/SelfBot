package command

import (
	"selfbot/discord/command/handlers"
	"selfbot/discord/command/handlers/info"
	voiceH "selfbot/discord/command/handlers/voice"
	"selfbot/discord/voice"

	"github.com/rs/zerolog"

	"github.com/bwmarrin/discordgo"
)

const Prefix = "/"

type Manager struct {
	l                zerolog.Logger
	commands         map[string]handlers.Handler
	shutdownHandlers []func()
}

func NewCommandManager(l zerolog.Logger, s *discordgo.Session, vm voice.Manager) (Manager, error) {
	var ret = Manager{
		l:        l.With().Str("owner", "CommandManager").Logger(),
		commands: make(map[string]handlers.Handler),
	}
	ret.shutdownHandlers = append(ret.shutdownHandlers, s.AddHandler(ret.discordHandleMessageCreate))

	// Register commands.
	ret.commands["info"] = new(info.Handler)
	ret.commands["join"] = &voiceH.JoinHandler{VoiceManager: vm}
	ret.commands["leave"] = &voiceH.LeaveHandler{VoiceManager: vm}
	ret.commands["play"] = &voiceH.PlayHandler{VoiceManager: vm}
	ret.commands["sound_files"] = &voiceH.SoundsHandler{VoiceManager: vm}

	return ret, nil
}

func (cm *Manager) discordHandleMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID || m.Author.Bot {
		return
	}
	isCmd, cmdName, args := processCommand(m.Content)
	if !isCmd {
		return
	}

	cmd, ok := cm.commands[cmdName]
	if !ok {
		return // Command doesn't exist, ignore.
	}

	var shouldReact = cmd.ShouldReact()
	if shouldReact {
		_ = React(s, m, "⏳")
		defer func() {
			_ = Unreact(s, m, "⏳")
		}()
	}

	defer cm.recoveryFunc(s, m, cmdName)
	if err := cmd.Handle(s, m, args...); err != nil {
		cm.error(s, m, cmdName, err)
	} else if shouldReact {
		_ = React(s, m, "✅") // Great success!
	}
}
