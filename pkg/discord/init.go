package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

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
