package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

type botCommand struct {
	DiscordGoCommand *discordgo.ApplicationCommand
	Handler          func(*discordgo.Session, *discordgo.InteractionCreate)
}

// HELLO COMMAND
//
// command:
//
//	/hello
//
// make the bot say hello to check if it's alive
func (dc *DiscordBotService) HelloCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	log.Info().Msg("discord.HelloCommand")

	// respond to the interaction, it's like returning a JSON response
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Hello world!",
		},
	})
}

// BATTLE ITEMS COMMAND
//
// /battle-items {name=...} {tier=...}
//
// search for battle items
func (dc *DiscordBotService) BattleItemsCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	log.Info().Interface("interaction_data", i).Msg("discord.BattleItemsCommand")

	// parse command options passed in
	var (
		optFilterName string
		optFilterTier string
	)
	opt := i.ApplicationCommandData().Options
	for i := range opt {
		switch opt[i].Name {
		case "name":
			optFilterName = opt[i].StringValue()

		case "tier":
			optFilterTier = opt[i].StringValue()
		}
	}

	// respond data
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "You entered these filters -> name: " + optFilterName + ", tier: " + optFilterTier,
		},
	})
}
