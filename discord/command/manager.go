package command

import (
	"selfbot/discord/command/handlers"
	"selfbot/discord/command/handlers/info"
	"selfbot/discord/command/handlers/join"
	feedback2 "selfbot/discord/feedback"
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
		commands: make(map[string]handlers.Handler),
		l:        l.With().Str("owner", "Manager").Logger(),
	}
	ret.shutdownHandlers = append(ret.shutdownHandlers, s.AddHandlerOnce(ret.discordHandleMessageCreate))

	// Register command.
	ret.commands["info"] = new(info.Handler)
	ret.commands["join"] = &join.Handler{VoiceManager: vm}
	ret.commands["leave"] = &join.Handler{VoiceManager: vm}
	ret.commands["play"] = &join.Handler{VoiceManager: vm}
	ret.commands["sounds"] = &join.Handler{VoiceManager: vm}

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

func (cm *Manager) recoveryFunc(s *discordgo.Session, m *discordgo.MessageCreate, cmdName string) {
	var rec = recover()
	if rec != nil {
		cm.error(s, m, cmdName, rec)
	}
}

// errorFunc is a function to recover from a command issue.
func (cm *Manager) error(s *discordgo.Session, m *discordgo.MessageCreate, cmdName string, rec interface{}) {
	if rec == nil {
		return
	}

	if err := React(s, m, "❌"); err != nil {
		LogDataMessage(cm.l.Error(), m).
			Err(err).
			Msg("Unable to recovery react to message.")
	}

	logEvent := LogDataMessage(cm.l.Error(), m)
	switch rec.(type) {
	case feedback2.UserError:
		// TODO output error message to user.
	case error:
		logEvent = logEvent.Err(rec.(error))
	default:
		logEvent = logEvent.Interface("error", rec)
	}
	logEvent.Msg("Command processing error occurred.")
}

func (cm *Manager) Close() error {
	for _, v := range cm.shutdownHandlers {
		v() // Call shutdown handlers.
	}
	return nil
}
