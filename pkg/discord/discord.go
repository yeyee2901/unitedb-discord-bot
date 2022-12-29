package discord

import "github.com/yeyee2901/unitedb-discord-bot/pkg/datasource"

type DiscordBotService struct {
	ClientId     string
	ClientSecret string
	Token        string

	*datasource.DataSource
}

func NewDiscordBotService(clientId, clientSecret, token string, ds *datasource.DataSource) *DiscordBotService {
	return &DiscordBotService{clientId, clientSecret, token, ds}
}
