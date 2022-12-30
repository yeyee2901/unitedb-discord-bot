package models

type BattleItem struct {
	// ID of the item, this is only known in bot scope, the unite-db.com API
	// does not provide this
	Id uint64 `json:"id" db:"id"`

	// name of the item
	Name string `json:"name" db:"name"`

	// description of the item
	Description string `json:"description" db:"description"`

	// Tier of the item in the current patch
	Tier string `json:"tier" db:"tier"`

	// cooldown in seconds
	Cooldown uint16 `json:"cooldown" db:"cooldown"`

	// player level condition to unlock this item
	TrainerLevel uint64 `json:"level" db:"trainer_level"`
}

type HeldItem struct {
	// Name of the item
	Name string `json:"name"`

	// Item bonus #1
	Bonus1 string `json:"bonus1"`

	// Item bonus #2 (may be missing in some items)
	Bonus2 string `json:"bonus2"`

	// Stats of the item, the important part is only the `label` field,
	// since it can be used for item tag
	Stats []heldItemStats `json:"stats"`

	// true description of the item (explains the effect of the item in
	// general)
	Description1 string `json:"description1"`

	// Don't know what this is used for, but maybe it can come in handy later
	// so I'm including it
	Description3 string `json:"description3"`

	// tier of the item in the current patch
	Tier string `json:"tier"`

	// item bonus at level 1, corresponds to `bonus1` & `bonus2` field
	Level1 string `json:"level1"`

	// item bonus at level 10, corresponds to `bonus1` & `bonus2` field
	Level10 string `json:"level10"`

	// item bonus at level 20, corresponds to `bonus1` & `bonus2` field
	Level20 string `json:"level20"`
}

type heldItemStats struct {
	// label of the stats, can be used for tagging the item
	Label string `json:"label"`

	// whether this stats is in percentage or not
	Percent bool `json:"percent"`

	// increments per level increase
	Increment uint16 `json:"increment"`

	// XXX: don't know
	Initial uint16 `json:"initial"`

	// XXX: don't know
	Start uint16 `json:"start"`

	// XXX: don't know
	Skip uint16 `json:"skip"`

	// XXX: don't know
	InitialDiff uint16 `json:"initial_diff"`
}
