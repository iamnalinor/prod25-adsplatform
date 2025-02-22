package service

import (
	"backend/internal/model"
	"backend/internal/repo"
	"backend/pkg/floatutil"
	"fmt"
	"github.com/google/uuid"
	"slices"
)

type AdService struct {
	campaignRepo repo.Campaign
	settingsRepo repo.Settings
}

const (
	revenueViewWeight  = 0.5 / (0.5 + 0.25) * 0.5
	revenueClickWeight = 0.5 / (0.5 + 0.25) * 0.5
	mlScoreWeight      = 0.25 / (0.5 + 0.25)
)

const limitsThreshold = 1.04

func compareAdCandidates(a, b model.AdCandidate) int {
	A := 0.0
	B := 0.0

	if a.Viewed && !b.Viewed {
		B += revenueViewWeight
	} else if !a.Viewed && b.Viewed {
		A += revenueViewWeight
	} else if !a.Viewed && !b.Viewed {
		A += floatutil.Norm(a.CostPerImpression, b.CostPerImpression) * revenueViewWeight
		B += floatutil.Norm(b.CostPerImpression, a.CostPerImpression) * revenueViewWeight
	}

	if a.Clicked && !b.Clicked {
		B += revenueViewWeight
	} else if !a.Clicked && b.Clicked {
		A += revenueViewWeight
	} else if !a.Clicked && !b.Clicked {
		A += floatutil.Norm(a.CostPerClick, b.CostPerClick) * revenueClickWeight
		B += floatutil.Norm(b.CostPerClick, a.CostPerClick) * revenueClickWeight
	}

	A += floatutil.Norm(float64(a.MlScore), float64(b.MlScore)) * mlScoreWeight
	B += floatutil.Norm(float64(b.MlScore), float64(a.MlScore)) * mlScoreWeight

	if A < B {
		return -1
	} else if A > B {
		return 1
	}
	return 0
}

// chooseAdCandidate chooses the most suitable ad for user to display. If no ad is suitable,
// the last return value is false. It modifies the input slice.
func chooseAdCandidate(candidates []model.AdCandidate) (model.AdCandidate, bool) {
	if len(candidates) == 0 {
		return model.AdCandidate{}, false
	}

	slices.SortFunc(candidates, compareAdCandidates)
	return candidates[len(candidates)-1], true
}

func (s *AdService) GetAd(client model.Client) (model.Ad, error) {
	candidates, err := s.campaignRepo.GetAdCandidates(client.Id, limitsThreshold)
	if err != nil {
		return model.Ad{}, fmt.Errorf("get ad candidates: %w", err)
	}

	candidate, ok := chooseAdCandidate(candidates)
	if !ok {
		return model.Ad{}, repo.ErrNotFound
	}

	currentDate := s.settingsRepo.GetCached().CurrentDate

	if err := s.campaignRepo.AddAdImpression(model.AdImpression{
		ClientId:   client.Id,
		CampaignId: candidate.Id,
		Spent:      candidate.CostPerImpression,
		Date:       currentDate,
	}); err != nil {
		return model.Ad{}, fmt.Errorf("add impression: %w", err)
	}

	return candidate.Ad, nil
}

func (s *AdService) GetAdCandidates(client model.Client) ([]model.AdCandidate, error) {
	candidates, err := s.campaignRepo.GetAdCandidates(client.Id, limitsThreshold)
	if err != nil {
		return nil, fmt.Errorf("get ad candidates: %w", err)
	}
	return candidates, nil
}

func (s *AdService) IsAdViewed(clientId, campaignId uuid.UUID) (bool, error) {
	_, err := s.campaignRepo.GetAdImpression(clientId, campaignId)
	if err != nil {
		if repo.IsNotFound(err) {
			return false, nil
		}
		return false, fmt.Errorf("get ad impression: %w", err)
	}
	return true, nil
}

func (s *AdService) ClickAd(client model.Client, campaign model.Campaign) error {
	err := s.campaignRepo.AddAdClick(model.AdClick{
		ClientId:   client.Id,
		CampaignId: campaign.Id,
		Spent:      *campaign.CostPerClick,
		Date:       s.settingsRepo.GetCached().CurrentDate,
	})
	if err != nil {
		return fmt.Errorf("add click: %w", err)
	}
	return nil
}
