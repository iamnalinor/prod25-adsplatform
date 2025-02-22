package service

import (
	"backend/internal/model"
	"backend/internal/repo"
	"backend/pkg/sliceutil"
	"fmt"
	"github.com/google/uuid"
)

type AdvertiserService struct {
	advertiserRepo repo.Advertiser
	clientRepo     repo.Client
	mlScoreRepo    repo.MlScore
}

func (s *AdvertiserService) GetById(id uuid.UUID) (model.Advertiser, error) {
	return s.advertiserRepo.GetById(id)
}

func (s *AdvertiserService) AddBulk(advertisers []model.Advertiser) ([]model.Advertiser, error) {
	err := s.advertiserRepo.UpsertMany(advertisers)
	if err != nil {
		return nil, fmt.Errorf("add bulk advertisers: %w", err)
	}

	// input can contain duplicates
	// saving the order but returning only last one when duplicated

	ids := make([]uuid.UUID, len(advertisers))
	for i, advertiser := range advertisers {
		ids[i] = advertiser.Id
	}
	ids = sliceutil.DeduplicateLast(ids)

	resultMap, err := s.advertiserRepo.GetMany(ids)
	if err != nil {
		return nil, fmt.Errorf("get added bulk advertisers: %w", err)
	}

	result := make([]model.Advertiser, len(ids))
	for i, id := range ids {
		result[i] = resultMap[id]
	}

	return result, nil
}

func (s *AdvertiserService) AddMlScore(score model.MlScore) error {
	// if some of the entities don't exist, ErrNotFound will be embedded in returned error
	_, err := s.advertiserRepo.GetById(score.AdvertiserId)
	if err != nil {
		return fmt.Errorf("get advertiser by id: %w", err)
	}
	_, err = s.clientRepo.GetById(score.ClientId)
	if err != nil {
		return fmt.Errorf("get client by id: %w", err)
	}

	err = s.mlScoreRepo.Upsert(score)
	if err != nil {
		return fmt.Errorf("upsert mlScore: %w", err)
	}
	return nil
}
