package discord

import (
	"fmt"

	"github.com/yeyee2901/unitedb-discord-bot/config"
	"github.com/yeyee2901/unitedb-discord-bot/datasource"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"
)

type BotCommand struct {
	Command *discordgo.ApplicationCommand
	Handler InteractionCreateEventHandler
}

type InteractionCreateEventHandler func(*discordgo.Session, *discordgo.InteractionCreate)

// DiscordBotService is the discord bot application logic
type DiscordBotService struct {
	clientID     string
	clientSecret string
	token        string
	redisStore   datasource.RedisStore
	logger       *zerolog.Logger
	config       *config.AppConfig

	registeredCommands map[string]BotCommand
}

// NewDiscordBotService constructs a new instance of the application. This also
// opens the discord websocket connection. On successful, it will return the
// master Service & the active discord session
func NewDiscordBotService(
	clientID string,
	clientSecret string,
	token string,
	redisStore datasource.RedisStore,
	logger *zerolog.Logger,
	cfg *config.AppConfig,
) (*DiscordBotService, *discordgo.Session, error) {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, nil, err
	}

	// set logging mode
	if cfg.Bot.Mode == "development" {
		session.LogLevel = discordgo.LogInformational | discordgo.LogWarning
	} else {
		session.LogLevel = discordgo.LogWarning
	}

	// open the websocket connection
	err = session.Open()
	if err != nil {
		return nil, nil, err
	}
	logger.Info().Msg("Successfully connected to discord")

	return &DiscordBotService{
		clientID:           clientID,
		clientSecret:       clientSecret,
		token:              token,
		redisStore:         redisStore,
		logger:             logger,
		config:             cfg,
		registeredCommands: make(map[string]BotCommand),
	}, session, nil
}

// RegisterCommands register the discord commands to the active discord
// session, so that it appears in the discord autocompletion when user is
// typing. On successful, it returns the array of registered command IDs
func (bot *DiscordBotService) RegisterCommands(session *discordgo.Session) ([]string, error) {
	var commandData = map[string]BotCommand{
		"battle-items": {
			Command: &discordgo.ApplicationCommand{
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
			Handler: bot.BattleItemsCommand,
		},
	}

	// for development purposes
	var guildID string
	if bot.config.Bot.Mode == "development" {
		id, ok := bot.config.Bot.ServerID["grand_lord_bidoof"]
		if !ok {
			return nil, fmt.Errorf("invalid guild ID for key `grand_lord_bidoof` in setting.yaml")
		}
		guildID = id
	} else {
		id, ok := bot.config.Bot.ServerID["poke_pokemon"]
		if !ok {
			return nil, fmt.Errorf("invalid guild ID for key `poke_pokemon` in setting.yaml")
		}
		guildID = id
	}

	var registeredCommandID []string
	for k := range commandData {
		// NOTE: appId = Bot ID
		bot.logger.Info().Str("command", k).Str("guild_id", guildID).Msg("registering command")
		successCmd, err := session.ApplicationCommandCreate(session.State.User.ID, guildID, commandData[k].Command)
		if err != nil {
			bot.logger.Error().Err(err).Msg("Failed to register command: " + k)
			return nil, err
		}

		// required for cleanup
		registeredCommandID = append(registeredCommandID, successCmd.ID)
	}

	// Register command handler
	bot.registeredCommands = commandData
	session.AddHandler(bot.interactionCreateEvent)

	bot.logger.Info().Msg("Successfully registered all commands!")

	return registeredCommandID, nil
}

// Close closes the `session` connection to the discord. User should call
// UnregisterCommands() first before closing the session
func (bot *DiscordBotService) Close(session *discordgo.Session) {
	err := session.Close()
	if err != nil {
		bot.logger.Warn().Err(err).Msg("Failed to close connection")
	}

	bot.logger.Info().Msg("Successfully closed the discord connection")
}

// UnregisterCommands unregisters all commands on the `session`
func (s *DiscordBotService) UnregisterCommands(session *discordgo.Session, commandIDs ...string) error {
	s.logger.Info().Interface("command_ids", commandIDs).Msg("Unregistering command")

	// for development purposes
	var guildId string
	if s.config.Bot.Mode == "development" {
		id, ok := s.config.Bot.ServerID["grand_lord_bidoof"]
		if !ok {
			return fmt.Errorf("invalid guild ID for key `grand_lord_bidoof` in setting.yaml")
		}
		guildId = id
	} else {
		id, ok := s.config.Bot.ServerID["poke_pokemon"]
		if !ok {
			return fmt.Errorf("invalid guild ID for key `poke_pokemon` in setting.yaml")
		}
		guildId = id
	}

	for _, cmdId := range commandIDs {
		if err := session.ApplicationCommandDelete(session.State.User.ID, guildId, cmdId); err != nil {
			return err
		}
	}

	return nil
}
