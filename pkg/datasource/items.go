package datasource

import (
	"strings"

	"github.com/yeyee2901/unitedb-discord-bot/pkg/models"
)

// TEST: ok
// Get battle items, you can also use filter
// - filterName -> will query for items that match the pattern (using LIKE comparison)
// - filterTier -> will filter out items for that tier (using equal sign '=' comparison)
func (ds *DataSource) GetBattleItemsByName(filterName, filterTier string) (res []models.BattleItem, err error) {
	var queryReplacer []any
	var filterString []string

	baseQuery := `
        SELECT
            id, name, description, tier, cooldown, trainer_level
        FROM
            pokemon_battle_items
    `

	if len(filterName) != 0 {
		filterString = append(filterString, " name LIKE ?")
		queryReplacer = append(queryReplacer, "%"+filterName+"%")
	}

	if len(filterTier) != 0 {
		filterString = append(filterString, " tier = ?")
		queryReplacer = append(queryReplacer, filterTier)
	}

	// check whether the filter is filled or not,
	// otherwise use the base query
	if len(filterString) != 0 {
		q := baseQuery + " WHERE " + strings.Join(filterString, " AND ")
		err = ds.DB.Select(&res, q, queryReplacer...)
	} else {
		err = ds.DB.Select(&res, baseQuery)
	}

	return res, err
}
