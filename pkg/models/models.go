package models

type BattleItem struct {
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	Tier        string `json:"tier" db:"tier"`
	Cooldown    uint16 `json:"cooldown" db:"cooldown"`
	UnlockLevel uint64 `json:"level" db:"trainer_level"`
}
