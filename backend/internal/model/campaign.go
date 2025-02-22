package model

import (
	"github.com/google/uuid"
	"time"
)

type CampaignTargeting struct {
	Gender   *string `json:"gender" db:"targeting_gender" binding:"omitempty,oneof=MALE FEMALE ALL"`
	AgeFrom  *int    `json:"age_from" db:"targeting_age_from"`
	AgeTo    *int    `json:"age_to" db:"targeting_age_to"`
	Location *string `json:"location" db:"targeting_location"`
}

type CampaignCreateRequest struct {
	ImpressionsLimit  *int     `json:"impressions_limit" db:"impressions_limit" binding:"required,gte=0"`
	ClicksLimit       *int     `json:"clicks_limit" db:"clicks_limit" binding:"required,gte=0"`
	CostPerImpression *float64 `json:"cost_per_impression" db:"cost_per_impression" binding:"required,gte=0"`
	CostPerClick      *float64 `json:"cost_per_click" db:"cost_per_click" binding:"required,gte=0"`
	AdTitle           string   `json:"ad_title" db:"ad_title" binding:"required"`
	AdText            string   `json:"ad_text" db:"ad_text" binding:"required"`
	StartDate         *int     `json:"start_date" db:"start_date" binding:"required,gte=0"`
	EndDate           *int     `json:"end_date" db:"end_date" binding:"required,gte=0"`
	CampaignTargeting `json:"targeting"`
}

type Campaign struct {
	Id           uuid.UUID `json:"campaign_id" db:"id"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	AdvertiserId uuid.UUID `json:"advertiser_id" db:"advertiser_id"`
	CampaignCreateRequest
	ImagePath        string              `json:"image_path" db:"image_path"`
	ModerationTaskId *uuid.UUID          `json:"-" db:"moderation_task_id"`
	ModerationResult *AiModerationResult `json:"moderation_result" db:"moderation_result"`
}

type GetCampaignsRequest struct {
	Size int `form:"size" binding:"gte=0"`
	Page int `form:"page" binding:"gte=0"`
}

type CampaignStats struct {
	ImpressionsCount int     `json:"impressions_count" db:"impressions_count"`
	ClicksCount      int     `json:"clicks_count" db:"clicks_count"`
	Conversion       float64 `json:"conversion" db:"conversion"`
	SpentImpressions float64 `json:"spent_impressions" db:"spent_impressions"`
	SpentClicks      float64 `json:"spent_clicks" db:"spent_clicks"`
	SpentTotal       float64 `json:"spent_total" db:"spent_total"`
	Date             *int    `json:"date,omitempty" db:"date"`
}

type Ad struct {
	Id           uuid.UUID `json:"ad_id" db:"ad_id"`
	Title        string    `json:"ad_title" db:"ad_title"`
	Text         string    `json:"ad_text" db:"ad_text"`
	AdvertiserId uuid.UUID `json:"advertiser_id" db:"advertiser_id"`
	ImagePath    string    `json:"image_path" db:"image_path"`
}

type AdCandidate struct {
	Ad
	MlScore           int     `json:"ml_score" db:"ml_score"`
	CostPerImpression float64 `json:"cost_per_impression" db:"cost_per_impression"`
	ImpressionsCount  int     `json:"impressions_count" db:"impressions_count"`
	ImpressionsLimit  int     `json:"impressions_limit" db:"impressions_limit"`
	Viewed            bool    `json:"viewed" db:"viewed"`
	CostPerClick      float64 `json:"cost_per_click" db:"cost_per_click"`
	ClicksCount       int     `json:"clicks_count" db:"clicks_count"`
	ClicksLimit       int     `json:"clicks_limit" db:"clicks_limit"`
	Clicked           bool    `json:"clicked" db:"clicked"`
}

type AdImpression struct {
	ClientId   uuid.UUID `json:"client_id" db:"client_id"`
	CampaignId uuid.UUID `json:"campaign_id" db:"campaign_id"`
	Spent      float64   `json:"spent" db:"spent"`
	Date       int       `json:"date" db:"date"`
}

type AdClick struct {
	ClientId   uuid.UUID `json:"client_id" db:"client_id"`
	CampaignId uuid.UUID `json:"campaign_id" db:"campaign_id"`
	Spent      float64   `json:"spent" db:"spent"`
	Date       int       `json:"date" db:"date"`
}
