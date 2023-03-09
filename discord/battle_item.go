package discord

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	unitepb "github.com/yeyee2901/unitedb-api-proto/gen/go/unitedb/v1"
	"github.com/yeyee2901/unitedb-discord-bot/helper"
)

func formatBattleItem(items []*unitepb.BattleItem) (*discordgo.InteractionResponseData, error) {
	// for single result, show the details
	if len(items) == 1 {
		m1 := fmt.Sprintf("🟢 **%s** (Tier **%s**)\n\n", strings.ToUpper(items[0].Name), items[0].Tier)
		m2 := fmt.Sprintf("🕒 %d seconds\n", items[0].Cooldown)
		m3 := fmt.Sprintf("🔓 Player Level %d \n\n", items[0].TrainerLevel)
		m4 := fmt.Sprintf("__Description__\n%s\n", items[0].Description)

		return &discordgo.InteractionResponseData{
			Content: m1 + m2 + m3 + m4,
		}, nil
	}

	// get all tier
	tierList := []string{}
	for i := range items {
		tierList = append(tierList, items[i].Tier)
	}
	tierList = helper.RemoveDuplicate(tierList)

	// sort by order
	order := []string{
		1: "S",
		2: "A",
		3: "B",
		4: "C",
		5: "D",
	}
	tierList, err := helper.SortUsingOrder(tierList, order)
	if err != nil {
		return nil, err
	}

	// group the items according to tier
	groupedItems := []string{}
	for _, tier := range tierList {
		itemInThisTier := []string{}
		for i := range items {
			if items[i].Tier == tier {
				itemInThisTier = append(itemInThisTier, items[i].Name)
			}
		}

		groupedItems = append(groupedItems, fmt.Sprintf("__Tier **%s**__\n%s", tier, strings.Join(itemInThisTier, ", ")))
	}

	return &discordgo.InteractionResponseData{
		Content: "🟢 **RESULT**\n\n" + strings.Join(groupedItems, "\n\n"),
	}, nil
}
