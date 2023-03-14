package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (bot *DiscordBotService) handlePanic(err any, message string, s *discordgo.Session, i *discordgo.InteractionCreate) {
	bot.logger.Error().Err(fmt.Errorf("%+v", err)).Msg(message)

	sendError := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "I'm sorry, looks like something went wrong :(",
		},
	})

	if sendError != nil {
		bot.logger.Error().Err(sendError).Msg("Failed to send message to discord")
	}
}

func (bot *DiscordBotService) handleAPIConnectionFailure(err error, s *discordgo.Session, i *discordgo.InteractionCreate) {
	msg := "Failed to connect to the upstream API. :("
	bot.logger.Error().Err(err).Str("type", "upstream").Msg(msg)

	sendError := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "🔴 " + msg,
		},
	})

	if sendError != nil {
		bot.logger.Error().Err(sendError).Msg("Failed to send message to discord")
	}
}

func (bot *DiscordBotService) handleGrpcError(grpcErr *status.Status, s *discordgo.Session, i *discordgo.InteractionCreate) {
	var sendError error

	switch grpcErr.Code() {
	case codes.NotFound:
		sendError = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "🔴  Content not found :(",
			},
		})

	case codes.DeadlineExceeded:
		sendError = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "🔴  Server timeout :(",
			},
		})

	case codes.InvalidArgument:
		sendError = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "🔴  Invalid argument.",
			},
		})

	default:
		bot.logger.Warn().Str("grpc", grpcErr.String()).Str("type", "grpc").Msg("gRPC error")
		sendError = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "🔴  ERROR: Unknown. Please contact the adminstrator https://github.com/yeyee2901 :(",
			},
		})
	}

	if sendError != nil {
		bot.logger.Error().Err(sendError).Msg("Failed to send message to discord")
	}
}

func (bot *DiscordBotService) handleGenericError(err error, logMessage string, userMessage string, s *discordgo.Session, i *discordgo.InteractionCreate) {
	bot.logger.Error().Err(err).Str("type", "generic").Msg(logMessage)

	sendError := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "🔴 " + userMessage,
		},
	})

	if sendError != nil {
		bot.logger.Error().Err(sendError).Msg("Failed to send message to discord")
	}
}
