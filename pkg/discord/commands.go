package discord

import (
	"time"

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

	// fetch from database
	battleItems, err := dc.GetBattleItems(optFilterName, optFilterTier)
	if err != nil {
		internalBotError(err, i, s, "discord.BattleItemsCommand.GetBattleItems")
		return
	}

	// early return if no result found
	if len(battleItems) == 0 {
		userInteractionError(i, s, "No result found for that search term.", "discord.BattleItemsCommand.not-found", 10)
		return
	}

	// format the battle items
	formattedData, err := formatBattleItems(battleItems)
	if err != nil {
		internalBotError(err, i, s, "discord.BattleItemsCommand.formatBattleItems")
		return
	}

	// send to discord
	resp := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: formattedData,
	}
	if err := s.InteractionRespond(i.Interaction, resp); err != nil {
		internalBotError(err, i, s, "discord.BattleItemsCommand.InteractionRespond")
		return
	}
}

func internalBotError(err error, i *discordgo.InteractionCreate, s *discordgo.Session, logSubject string) {
	log.Error().Err(err).Str("interaction_id", i.ID).Msg(logSubject)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "🔴 Something went wrong inside. :(",
		},
	})
}

func userInteractionError(i *discordgo.InteractionCreate, s *discordgo.Session, message string, logSubject string, delAfterSeconds int) {
	log.Info().Str("interaction_id", i.ID).Msg(logSubject)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "🟡 " + message + "\n\n*note: this message will be automatically deleted.",
		},
	})

	// delete the message after that, else do not delete
	if delAfterSeconds > 0 {
		time.AfterFunc(time.Duration(delAfterSeconds)*time.Second, func() {
			s.InteractionResponseDelete(i.Interaction)
		})
	}
}
