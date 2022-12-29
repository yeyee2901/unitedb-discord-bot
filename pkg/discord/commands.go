package discord

import (
	"fmt"

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
	var botCommands = []botCommand{
		{
			DiscordGoCommand: &discordgo.ApplicationCommand{
				Name:        "hello",
				Description: "make the bot say hello",
			},
			Handler: dc.HelloCommand,
		},
	}

	// for development purposes
	var guildId string
	if dc.Config.Discord.Mode != "production" {
		guildId = dc.Config.Discord.Servers.Dev
	} else {
		guildId = dc.Config.Discord.Servers.Pokemon
	}

	registeredCommand := make([]*discordgo.ApplicationCommand, len(botCommands))

	for i := range botCommands {
		// register command to discord API
		// NOTE: appId = Bot ID
		if _, err := s.ApplicationCommandCreate(s.State.User.ID, guildId, botCommands[i].DiscordGoCommand); err != nil {
			log.Error().Err(err).Msg("discord.RegisterCommands.Failure")
			panic(err)
		} else {
			// required for cleanup
			registeredCommand[i] = botCommands[i].DiscordGoCommand

			// register command handler
			s.AddHandler(botCommands[i].Handler)
		}
	}

	dc.RegisteredCommands = registeredCommand
}

// should be called on cleanup function
func (dc *DiscordBotService) unregisterCommands(s *discordgo.Session) {
	// for development purposes
	var guildId string
	if dc.Config.Discord.Mode != "production" {
		guildId = dc.Config.Discord.Servers.Dev
	} else {
		guildId = dc.Config.Discord.Servers.Pokemon
	}

	log.Info().Msg("discord.unregisterCommands")
	for _, cmd := range dc.RegisteredCommands {
		if err := s.ApplicationCommandDelete(s.State.User.ID, guildId, cmd.ID); err != nil {
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
	log.Info().Msg("discord.BattleItemsCommand")
	fmt.Printf("i.ApplicationCommandData().Options: %v\n", i.ApplicationCommandData().Options)
}
