package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

// event handler for new message
func (dc *DiscordBotService) messageCreateEvent(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// ignore
	log.Info().Interface("new_message", m.Message).Msg("event.messageCreate")
}
