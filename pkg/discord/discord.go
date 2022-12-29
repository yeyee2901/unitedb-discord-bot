package discord

import (
	"github.com/rs/zerolog/log"
	"github.com/yeyee2901/unitedb-discord-bot/pkg/datasource"

	"github.com/bwmarrin/discordgo"
)

type DiscordBotService struct {
	ClientId           string
	ClientSecret       string
	Token              string
	RegisteredCommands []*discordgo.ApplicationCommand
	*datasource.DataSource
}

func NewDiscordBotService(clientId, clientSecret, token string, ds *datasource.DataSource) *DiscordBotService {
	return &DiscordBotService{
		ClientId:     clientId,
		ClientSecret: clientSecret,
		Token:        token,
		DataSource:   ds,
	}
}

// initialize the bot, this also run the bot by listening to discord updates
// sent through WebSocket
func (dc *DiscordBotService) Init(s *discordgo.Session) {
	log.Info().Msg("discord.Init")

	// set log level
	if dc.Config.Discord.Mode == "production" {
		s.LogLevel = discordgo.LogError
	} else {
		s.LogLevel = discordgo.LogInformational
	}

	// only retrieve guild chats
	s.Identify.Intents = discordgo.IntentGuildMessages

	// for login notification
	s.AddHandlerOnce(func(s2 *discordgo.Session, r *discordgo.Ready) {
		log.Info().Interface("servers", s2.State.Guilds).Msg("handler.login")
		return
	})

	// login to discord
	if err := s.Open(); err != nil {
		panic(err)
	}

	// register all commands
	dc.registerCommands(s)
}

func (dc *DiscordBotService) DeInit(s *discordgo.Session) {
	log.Info().Msg("discord.DeInit")
	dc.unregisterCommands(s)
}
