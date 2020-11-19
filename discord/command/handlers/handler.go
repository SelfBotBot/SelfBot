package handlers

import "github.com/bwmarrin/discordgo"

type Handler interface {
	Handle(s *discordgo.Session, m *discordgo.MessageCreate, args ...string) error
	ShouldReact() bool
}
