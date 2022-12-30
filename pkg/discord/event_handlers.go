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

// event handler for when user inputs slash command
// WARN: had to do it like this, because interaction handler behaves like a middleware chain
// ex: after the first matching interaction was found, it continues to call the next handler in chain
// so in this case, the interaction create event only needs one handler, but we have to match
// inside that incoming event handler, what kind of interaction it is, after that we can call
// the corresponding interaction handler, passing in the interaction data & session data
func (dc *DiscordBotService) interactionCreateEvent(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// do not respond if it's from private chat
	if len(i.GuildID) == 0 {
		return
	}

	// do not respond if it's not from command channel
	sourceChannel, err := s.Channel(i.ChannelID)
	if err != nil {
		log.Error().Err(err).Msg("discord.event.interactionCreateEvent.getChannel")
		return
	}
	if sourceChannel.Name != "bot-command" {
		return
	}

	// find the matching handler and call it
	if cmd, ok := dc.BotCommands[i.ApplicationCommandData().Name]; ok {
		cmd.Handler(s, i)
	}
}
