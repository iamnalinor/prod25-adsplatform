package model

import "github.com/google/uuid"

type Client struct {
	Id       uuid.UUID `json:"client_id" binding:"required,uuid" db:"id"`
	Login    string    `json:"login" binding:"required" db:"login"`
	Age      *int      `json:"age" binding:"required,gte=0" db:"age"`
	Location string    `json:"location" binding:"required" db:"location"`
	Gender   string    `json:"gender" binding:"required,oneof=MALE FEMALE" db:"gender"`
}
