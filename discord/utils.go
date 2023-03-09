package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func (bot *DiscordBotService) handlePanic(err any, message string, s *discordgo.Session, i *discordgo.InteractionCreate) {
	bot.logger.Error().Err(fmt.Errorf("%+v", err)).Msg("Bot panicked when processing BattleItemsCommand")

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "I'm sorry, looks like something went wrong :(",
		},
	})
}
