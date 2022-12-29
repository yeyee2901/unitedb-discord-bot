package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

// register all bot commands
func (dc *DiscordBotService) registerCommands(s *discordgo.Session) {
	log.Info().Msg("discord.RegisterCommands")

	// for development purposes
	var guildId string
	if dc.Config.Discord.Mode != "production" {
		guildId = dc.Config.Discord.Servers.Dev
	} else {
		guildId = dc.Config.Discord.Servers.Pokemon
	}

	// list of all commands
	commands := []*discordgo.ApplicationCommand{
		{
			Name:        "hello",
			Description: "Make the bot say hello",
		},
	}

	// register command to discord API
	registeredCommand := make([]*discordgo.ApplicationCommand, len(commands))
	for i, cmd := range commands {
		// appId = Bot ID
		if _, err := s.ApplicationCommandCreate(s.State.User.ID, guildId, cmd); err != nil {
			log.Error().Err(err).Msg("discord.RegisterCommands.Failure")
			panic(err)
		} else {
			registeredCommand[i] = cmd
		}
	}

	// register command handler
	commandHandlers := map[string]func(*discordgo.Session, *discordgo.InteractionCreate){
		"hello": dc.HelloCommand,
	}

	s.AddHandler(func(s2 *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s2, i)
		}
	})
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
// /battle-items {name}
//
// search for battle items
func (dc *DiscordBotService) BattleItemsCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	log.Info().Msg("discord.BattleItemsCommand")
}
