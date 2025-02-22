package service

import (
	"backend/config"
	"backend/internal/repo"
	"fmt"
)

type Services struct {
	Ad         *AdService
	Advertiser *AdvertiserService
	Ai         *AiService
	Api        *ApiService
	Campaign   *CampaignService
	Client     *ClientService
	Image      *ImageService
	Ollama     *OllamaService
	Settings   *SettingsService
	Stats      *StatsService
}

func NewServices(repos *repo.Repositories, env config.Environment) (*Services, error) {
	settingsSvc := &SettingsService{repos.Settings}
	ollamaSvc, err := NewOllamaService(env, repos.Ai)
	if err != nil {
		return nil, fmt.Errorf("create ollama service: %w", err)
	}
	aiSvc := &AiService{repos.Ai, ollamaSvc}
	return &Services{
		Ad:         &AdService{repos.Campaign, repos.Settings},
		Advertiser: &AdvertiserService{repos.Advertiser, repos.Client, repos.MlScore},
		Ai:         aiSvc,
		Api:        NewApiService(repos.Api),
		Campaign:   &CampaignService{repos.Campaign, aiSvc, settingsSvc},
		Client:     &ClientService{repos.Client},
		Image:      &ImageService{campaignRepo: repos.Campaign, mediaFsPath: env.MediaFsPath, mediaBaseUrl: env.MediaBaseUrl},
		Ollama:     ollamaSvc,
		Settings:   settingsSvc,
		Stats:      &StatsService{repos.Campaign},
	}, nil
}
