package datasource

import "time"

type HeldItem struct {
	// ID of the item
	Id uint64 `db:"id"`

	// Stat bonus #1
	Bonus1 string `db:"bonus1"`

	// Stat bonus #2 (may not be present in some items)
	Bonus2 string `db:"bonus2"`

	// Name of the item
	Name string `db:"name"`

	// description of the item effects, long text
	// corresponds to `description1` field
	Description string `db:"description"`

	// tier of the item in the current patch
	Tier string `db:"tier"`

	// Stats at level 1
	Level1 string `db:"level1"`

	// Stats at level 10
	Level10 string `db:"level10"`

	// Stats at level 1
	Level20 string `db:"level20"`

	// item tags
	Tags string `db:"tags"`

	// last updated time
	LastUpdate time.Time `db:"last_update"`
}
