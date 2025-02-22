package model

type Settings struct {
	CurrentDate       int  `db:"current_date"`
	ModerationEnabled bool `db:"moderation_enabled"`
}
