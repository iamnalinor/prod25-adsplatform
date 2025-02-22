package repo

import (
	"backend/internal/model"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Advertiser interface {
	GetById(id uuid.UUID) (model.Advertiser, error)
	GetMany(ids []uuid.UUID) (map[uuid.UUID]model.Advertiser, error)
	UpsertMany(advertisers []model.Advertiser) error
}

type Ai interface {
	AddTask(task model.AiTask) error
	GetTask(id uuid.UUID) (model.AiTask, error)
	GetIncompleteTasks() ([]model.AiTask, error)
	AddResult(result model.AiTaskResult) error
	GetResult(taskId uuid.UUID) (model.AiTaskResult, error)
}

type Api interface {
	AddRequest(endpoint string, durationMs float64) error
}

type Client interface {
	GetById(id uuid.UUID) (model.Client, error)
	GetMany(ids []uuid.UUID) (map[uuid.UUID]model.Client, error)
	UpsertMany(clients []model.Client) error
}

type Campaign interface {
	Add(campaign model.Campaign) error
	GetList(advertiserId uuid.UUID, size int, page int) ([]model.Campaign, error)
	GetById(id uuid.UUID) (model.Campaign, error)
	Update(campaign model.Campaign) error
	Delete(id uuid.UUID) error
	GetStats(advertiserId uuid.UUID, campaignId uuid.UUID) (model.CampaignStats, error)
	GetStatsDaily(advertiserId uuid.UUID, campaignId uuid.UUID) ([]model.CampaignStats, error)
	GetModerationFailed(size int, page int) ([]model.Campaign, error)
	GetAdCandidates(clientId uuid.UUID, limitsThreshold float64) ([]model.AdCandidate, error)
	AddAdImpression(impression model.AdImpression) error
	GetAdImpression(clientId, campaignId uuid.UUID) (model.AdImpression, error)
	AddAdClick(click model.AdClick) error
}

type MlScore interface {
	Upsert(score model.MlScore) error
}

type Settings interface {
	Get() (model.Settings, error)
	GetCached() model.Settings
	Update(settings model.Settings) error
}

type Repositories struct {
	Advertiser Advertiser
	Ai         Ai
	Api        Api
	Client     Client
	Campaign   Campaign
	MlScore    MlScore
	Settings   Settings
}

func NewRepositories(db *sqlx.DB) *Repositories {
	return &Repositories{
		Advertiser: &AdvertiserRepo{db},
		Ai:         &AiRepo{db},
		Api:        &ApiRepo{db},
		Client:     &ClientRepo{db},
		Campaign:   &CampaignRepo{db},
		MlScore:    &MlScoreRepo{db},
		Settings:   NewSettingsRepo(db),
	}
}
