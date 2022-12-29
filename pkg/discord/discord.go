package discord

import (
	"github.com/yeyee2901/unitedb-discord-bot/pkg/datasource"

	"github.com/bwmarrin/discordgo"
)

type DiscordBotService struct {
	ClientId     string
	ClientSecret string
	Token        string

	*datasource.DataSource
	Bot *discordgo.Session
}

func NewDiscordBotService(clientId, clientSecret, token string, ds *datasource.DataSource, dc *discordgo.Session) *DiscordBotService {
	return &DiscordBotService{
		clientId,
		clientSecret,
		token,
		ds,
		dc,
	}
}

func (dc *DiscordBotService) InitBot() {
	dc.Bot.Identify.Intents = discordgo.IntentGuildMessages
	dc.Bot.AddHandler(dc.messageCreateEvent)

	if dc.Config.Discord.Mode == "production" {
		dc.Bot.LogLevel = discordgo.LogError
	} else {
		dc.Bot.LogLevel = discordgo.LogInformational
	}
}
