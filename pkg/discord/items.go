package discord

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/yeyee2901/unitedb-discord-bot/pkg/models"
)

// formats the battle items into discordgo.InteractionResponseData according to the number of results from database
func formatBattleItems(items []models.BattleItem) (resp *discordgo.InteractionResponseData, err error) {
	if len(items) == 1 {
		m1 := fmt.Sprintf("🟢 **%s** (Tier **%s**)\n\n", strings.ToUpper(items[0].Name), items[0].Tier)
		m2 := fmt.Sprintf("🕒 %d seconds\n", items[0].Cooldown)
		m3 := fmt.Sprintf("🔓 Player Level %d \n\n", items[0].TrainerLevel)
		m4 := fmt.Sprintf("__Description__\n%s\n", items[0].Description)

		return &discordgo.InteractionResponseData{
			Content: m1 + m2 + m3 + m4,
		}, nil
	}

	return nil, fmt.Errorf("Unimplemented")
}
