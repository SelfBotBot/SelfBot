package command

import (
	"errors"
	"fmt"
	"selfbot/discord/command/handlers"
	"selfbot/discord/command/handlers/info"
	voiceH "selfbot/discord/command/handlers/voice"
	"selfbot/discord/feedback"
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
	ret.commands["sounds"] = &voiceH.SoundsHandler{VoiceManager: vm}

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
	case feedback.UserError, feedback.WrappedUserError:
		cm.handleFeedback(s, m, cmdName, rec)
		return
	case error:
		var (
			err        = rec.(error)
			wrappedErr feedback.WrappedUserError
			userErr    feedback.UserError
		)
		if errors.As(err, &wrappedErr) {
			cm.handleFeedback(s, m, cmdName, wrappedErr)
			return
		} else if errors.As(err, &userErr) {
			cm.handleFeedback(s, m, cmdName, userErr)
			return
		}
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

func (cm *Manager) handleFeedback(s *discordgo.Session, m *discordgo.MessageCreate, cmdName string, rec interface{}) {
	var userErr feedback.UserError
	var cause error

	if wrappedErr, ok := rec.(feedback.WrappedUserError); ok {
		userErr = wrappedErr.UserError
		cause = wrappedErr.Cause
	} else {
		userErr = rec.(feedback.UserError)
	}

	message := userErr.Message
	if cause != nil {
		message = fmt.Sprintf("%s \n\t(`%s`)", message, cause)
	}
	_, _ = s.ChannelMessageSend(m.ChannelID, message)
}
