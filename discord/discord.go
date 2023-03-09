package discord

import (
	"github.com/yeyee2901/unitedb-discord-bot/datasource"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"
)

// DiscordBotService is the discord bot application logic
type DiscordBotService struct {
	ClientId     string
	ClientSecret string
	Token        string
	Redis        datasource.RedisStore
	Session      *discordgo.Session
	Logger       *zerolog.Logger
}

// NewDiscordBotService constructs a new instance of the application. This also
// opens the discord websocket connection
func NewDiscordBotService(
	clientID string,
	clientSecret string,
	token string,
	redisStore datasource.RedisStore,
	logger *zerolog.Logger,
) (*DiscordBotService, error) {
	s, err := discordgo.New(token)
	if err != nil {
		return nil, err
	}

	// open the websocket connection
	err = s.Open()
	if err != nil {
		return nil, err
	}
	logger.Info().Msg("Successfully connected to discord")

	return &DiscordBotService{
		ClientId:     clientID,
		ClientSecret: clientSecret,
		Token:        token,
		Redis:        redisStore,
		Session:      s,
		Logger:       logger,
	}, nil
}

// Close closes the connection to the discord
func (s *DiscordBotService) Close() {
	err := s.Session.Close()
	if err != nil {
		s.Logger.Warn().Err(err).Msg("Failed to close connection")
	}

	s.Logger.Info().Msg("Successfully closed the discord connection")
}
