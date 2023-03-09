package discord

import "github.com/bwmarrin/discordgo"

func (bot *DiscordBotService) interactionCreateEvent(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// do not respond if it's from private chat
	if len(i.GuildID) == 0 {
		return
	}

	// do not respond if it's not from command channel
	sourceChannel, err := s.Channel(i.ChannelID)
	if err != nil {
		bot.logger.Error().Err(err).Msg("discord.event.interactionCreateEvent.getChannel")
		return
	}
	if sourceChannel.Name != "bot-command" {
		return
	}

	// find the matching handler and call it
	if cmd, ok := bot.registeredCommands[i.ApplicationCommandData().Name]; ok {
		cmd.Handler(s, i)
	}
}
