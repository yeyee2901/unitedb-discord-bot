package datasource

import (
	"strings"
)

// TEST: ok
// Get battle items, you can also use filter
// - filterName -> will query for items that match the pattern (using LIKE comparison)
// - filterTier -> will filter out items for that tier (using equal sign '=' comparison)
func (ds *DataSource) GetBattleItems(filterName, filterTier string) (res []BattleItem, err error) {
	var queryReplacer []any
	var filterString []string

	baseQuery := `
        SELECT
            id, name, description, tier, cooldown, trainer_level, last_updated
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

// Bulk insert battle items into the database
func (ds *DataSource) InsertBattleItems(battleItems []BattleItem) error {
	// insert the data
	q := `
        INSERT INTO pokemon_battle_items
            (name, description, tier, cooldown, trainer_level)
        VALUES
            (:name, :description, :tier, :cooldown, :trainer_level)
    `

	tx, err := ds.DB.Beginx()
	if err != nil {
		return err
	}

	if _, err := tx.NamedExec(q, battleItems); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
