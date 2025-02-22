package service

import (
	"backend/internal/model"
	"backend/internal/repo"
	"fmt"
	"github.com/google/uuid"
)

type StatsService struct {
	campaignRepo repo.Campaign
}

func (s *StatsService) GetStatsCampaign(campaign model.Campaign) (model.CampaignStats, error) {
	stats, err := s.campaignRepo.GetStats(campaign.AdvertiserId, campaign.Id)
	if err != nil {
		return model.CampaignStats{}, fmt.Errorf("get stats (for campaign): %w", err)
	}
	return stats, nil
}

func (s *StatsService) GetStatsCampaignDaily(campaign model.Campaign) ([]model.CampaignStats, error) {
	stats, err := s.campaignRepo.GetStatsDaily(campaign.AdvertiserId, campaign.Id)
	if err != nil {
		return nil, fmt.Errorf("get stats daily (for campaign): %w", err)
	}
	return stats, nil
}

func (s *StatsService) GetStatsAdvertiser(advertiser model.Advertiser) (model.CampaignStats, error) {
	stats, err := s.campaignRepo.GetStats(advertiser.Id, uuid.Nil)
	if err != nil {
		return model.CampaignStats{}, fmt.Errorf("get stats (for advertiser): %w", err)
	}
	return stats, nil
}

func (s *StatsService) GetStatsAdvertiserDaily(advertiser model.Advertiser) ([]model.CampaignStats, error) {
	stats, err := s.campaignRepo.GetStatsDaily(advertiser.Id, uuid.Nil)
	if err != nil {
		return nil, fmt.Errorf("get stats daily (for advertiser): %w", err)
	}
	return stats, nil
}
