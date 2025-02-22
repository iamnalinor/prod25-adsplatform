package model

type CurrentDate struct {
	CurrentDate *int `json:"current_date" binding:"required,gte=0"`
}
