package model

import "github.com/google/uuid"

type MlScore struct {
	ClientId     uuid.UUID `json:"client_id" binding:"required,uuid" db:"client_id"`
	AdvertiserId uuid.UUID `json:"advertiser_id" binding:"required,uuid" db:"advertiser_id"`
	Score        *int      `json:"score" binding:"required" db:"score"`
}
