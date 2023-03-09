package discord

import (
	"context"
	"time"

	"github.com/bwmarrin/discordgo"
	unitepb "github.com/yeyee2901/unitedb-api-proto/gen/go/unitedb/v1"
	"github.com/yeyee2901/unitedb-discord-bot/unitedbapi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

// BattleItemsCommand adalah event handler interaction create untuk command /battle-item
//
// This function may panic because of the API design
func (dc *DiscordBotService) BattleItemsCommand(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	commandCtx, cancel := context.WithTimeout(context.Background(), time.Duration(dc.config.API.TimeoutSeconds)*time.Second)
	defer cancel()

	// this function may panic because of how discordgo API design is, so I'm
	// putting this here
	defer func() {
		if err := recover(); err != nil {
			dc.handlePanic(err, "Bot panicked when handling /battle-items", session, interaction)
		}
	}()

	// parse arguments
	req := new(unitepb.GetBattleItemRequest)
	cmdOptions := interaction.ApplicationCommandData().Options
	for idx := range cmdOptions {
		switch cmdOptions[idx].Name {
		case "name":
			itemName := cmdOptions[idx].StringValue()
			req.Name = &itemName

		case "tier":
			itemTier := cmdOptions[idx].StringValue()
			req.Tier = &itemTier
		}
	}

	// dial connection to the gRPC API
	api, err := unitedbapi.NewGrpcUniteDBAPI(
		&dc.config.API,
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
	)

	if err != nil {
		dc.handleAPIConnectionFailure(err, session, interaction)
		return
	}

	resp, err := api.GetBattleItem(commandCtx, req)
	if err != nil {
		grpcError, isGrpcErr := status.FromError(err)
		if isGrpcErr {
			dc.handleGrpcError(grpcError, session, interaction)
			return
		}

		// not a grpc error
		dc.handleGenericError(
			err,
			"Unknown / unhandled error type",
			"Something went wrong inside",
			session,
			interaction,
		)
		return
	}

	formattedData, err := formatBattleItem(resp.Data)
	if err != nil {
		dc.handleGenericError(
			err,
			"Failed while sorting",
			"Something went wrong inside :(",
			session,
			interaction,
		)
		return
	}

	// send to discord
	discordResp := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: formattedData,
	}

	// respond to the interaction, it's like returning a JSON response
	err = session.InteractionRespond(interaction.Interaction, discordResp)
	if err != nil {
		// not a grpc error
		dc.handleGenericError(
			err,
			"Failed to send data to discord",
			"ERROR: while sending the item to discord, something went wrong :(",
			session,
			interaction,
		)
	}
}
