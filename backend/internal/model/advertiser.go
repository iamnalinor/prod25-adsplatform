package model

import "github.com/google/uuid"

type Advertiser struct {
	Id   uuid.UUID `json:"advertiser_id" binding:"required,uuid" db:"id"`
	Name string    `json:"name" binding:"required" db:"name"`
}
