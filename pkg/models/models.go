package models

type BattleItem struct {
	Id           uint64 `json:"id" db:"id"`
	Name         string `json:"name" db:"name"`
	Description  string `json:"description" db:"description"`
	Tier         string `json:"tier" db:"tier"`
	Cooldown     uint16 `json:"cooldown" db:"cooldown"`
	TrainerLevel uint64 `json:"level" db:"trainer_level"`
}
