package command

import (
	"errors"
	"fmt"
	"selfbot/discord/feedback"

	"github.com/bwmarrin/discordgo"
)

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
