package discord

import (
	"github.com/bwmarrin/discordgo"
	unitepb "github.com/yeyee2901/unitedb-api-proto/gen/go/unitedb/v1"
)

// BattleItemsCommand adalah event handler interaction create untuk command /battle-item
//
// This function may panic because of the API design
func (dc *DiscordBotService) BattleItemsCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	defer func() {
		if err := recover(); err != nil {
			dc.handlePanic(err, "Bot panicked when handling /battle-items", s, i)
		}
	}()

	// parse arguments
	req := new(unitepb.GetBattleItemRequest)
	cmdOptions := i.ApplicationCommandData().Options
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

	// TODO: dial connection to the gRPC API

	// respond to the interaction, it's like returning a JSON response
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "name: " + *req.Name + ", tier: " + *req.Tier,
		},
	})
}
