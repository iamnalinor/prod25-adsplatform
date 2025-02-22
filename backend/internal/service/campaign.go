package service

import (
	"backend/internal/model"
	"backend/internal/repo"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type CampaignService struct {
	campaignRepo repo.Campaign
	aiSvc        *AiService
	settingsSvc  *SettingsService
}

func (s *CampaignService) Create(advertiserId uuid.UUID, req model.CampaignCreateRequest) (model.Campaign, error) {
	var taskId *uuid.UUID
	if s.settingsSvc.ModerationEnabled() {
		taskIdRaw, err := s.aiSvc.SubmitModeration(req.AdTitle, req.AdText)
		if err != nil {
			return model.Campaign{}, fmt.Errorf("submit moderation task: %w", err)
		}
		taskId = &taskIdRaw
	}

	campaign := model.Campaign{
		Id:                    uuid.New(),
		CreatedAt:             time.Now(),
		AdvertiserId:          advertiserId,
		CampaignCreateRequest: req,
		ModerationTaskId:      taskId,
	}
	if err := s.campaignRepo.Add(campaign); err != nil {
		return model.Campaign{}, fmt.Errorf("create campaign: %w", err)
	}
	return campaign, nil
}

func (s *CampaignService) GetById(id uuid.UUID) (model.Campaign, error) {
	campaign, err := s.campaignRepo.GetById(id)
	if err != nil {
		return model.Campaign{}, fmt.Errorf("get campaign by id: %w", err)
	}
	return campaign, nil
}

func (s *CampaignService) GetPaginated(advertiserId uuid.UUID, size int, page int) ([]model.Campaign, error) {
	campaigns, err := s.campaignRepo.GetList(advertiserId, size, page)
	if err != nil {
		return nil, fmt.Errorf("get campaign list: %w", err)
	}
	return campaigns, nil
}

func (s *CampaignService) Update(campaign *model.Campaign, req model.CampaignCreateRequest) error {
	if (campaign.AdTitle != req.AdTitle || campaign.AdText != req.AdText) && s.settingsSvc.ModerationEnabled() {
		taskId, err := s.aiSvc.SubmitModeration(req.AdTitle, req.AdText)
		if err != nil {
			return fmt.Errorf("submit moderation task: %w", err)
		}
		campaign.ModerationTaskId = &taskId
		campaign.ModerationResult = nil
	}
	campaign.CampaignCreateRequest = req

	if err := s.campaignRepo.Update(*campaign); err != nil {
		return fmt.Errorf("update campaign: %w", err)
	}
	return nil
}

func (s *CampaignService) Delete(id uuid.UUID) error {
	if err := s.campaignRepo.Delete(id); err != nil {
		return fmt.Errorf("delete campaign: %w", err)
	}
	return nil
}

func (s *CampaignService) GetModerationFailed(size int, page int) ([]model.Campaign, error) {
	campaigns, err := s.campaignRepo.GetModerationFailed(size, page)
	if err != nil {
		return nil, fmt.Errorf("get campaign list with failed moderation: %w", err)
	}
	return campaigns, nil
}
