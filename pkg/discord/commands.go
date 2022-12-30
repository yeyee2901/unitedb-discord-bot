package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

type botCommand struct {
	DiscordGoCommand *discordgo.ApplicationCommand
	Handler          func(*discordgo.Session, *discordgo.InteractionCreate)
}

// register all bot commands
func (dc *DiscordBotService) registerCommands(s *discordgo.Session) {
	log.Info().Msg("discord.RegisterCommands")

	// list all bot commands
	var botCommands = map[string]botCommand{
		"hello": {
			DiscordGoCommand: &discordgo.ApplicationCommand{
				Name:        "hello",
				Description: "Make the bot say hello",
			},
			Handler: dc.HelloCommand,
		},
		"battle-items": {
			DiscordGoCommand: &discordgo.ApplicationCommand{
				Name:        "battle-items",
				Description: "Search for battle items",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "name",
						Description: "Filter battle items by name",
						Required:    false,
						MaxLength:   15,
					},
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "tier",
						Description: "Filter battle items by tier",
						Required:    false,
						MaxLength:   1,
					},
				},
			},
			Handler: dc.BattleItemsCommand,
		},
	}

	// for development purposes
	var guildId string
	if dc.Config.Discord.Mode != "production" {
		guildId = dc.Config.Discord.Servers.Dev
	} else {
		guildId = dc.Config.Discord.Servers.Pokemon
	}

	for k := range botCommands {
		// register command to discord API
		// NOTE: appId = Bot ID
		if successCmd, err := s.ApplicationCommandCreate(s.State.User.ID, guildId, botCommands[k].DiscordGoCommand); err != nil {
			log.Error().Err(err).Msg("discord.RegisterCommands.Failure")
			panic(err)
		} else {
			// required for cleanup
			dc.CommandIdArray = append(dc.CommandIdArray, successCmd.ID)
		}
	}

	// attach to object, required for handlers register
	dc.BotCommands = botCommands

	// Register command handler
	s.AddHandler(dc.interactionCreateEvent)
}

// should be called on cleanup function
func (dc *DiscordBotService) unregisterCommands(s *discordgo.Session) {
	log.Info().Interface("command_ids", dc.CommandIdArray).Msg("discord.unregisterCommand")

	// for development purposes
	var guildId string
	if dc.Config.Discord.Mode != "production" {
		guildId = dc.Config.Discord.Servers.Dev
	} else {
		guildId = dc.Config.Discord.Servers.Pokemon
	}

	for _, cmdId := range dc.CommandIdArray {
		if err := s.ApplicationCommandDelete(s.State.User.ID, guildId, cmdId); err != nil {
			panic(err)
		}
	}
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
