package service

import (
	"backend/internal/model"
	"backend/internal/repo"
	"backend/pkg/sliceutil"
	"fmt"
	"github.com/google/uuid"
)

type ClientService struct {
	clientRepo repo.Client
}

func (s *ClientService) GetById(id uuid.UUID) (model.Client, error) {
	return s.clientRepo.GetById(id)
}

func (s *ClientService) AddBulk(clients []model.Client) ([]model.Client, error) {
	err := s.clientRepo.UpsertMany(clients)
	if err != nil {
		return nil, fmt.Errorf("add bulk clients: %w", err)
	}

	ids := make([]uuid.UUID, len(clients))
	for i, client := range clients {
		ids[i] = client.Id
	}
	ids = sliceutil.DeduplicateLast(ids)

	resultMap, err := s.clientRepo.GetMany(ids)
	if err != nil {
		return nil, fmt.Errorf("get added bulk clients: %w", err)
	}

	result := make([]model.Client, len(ids))
	for i, id := range ids {
		result[i] = resultMap[id]
	}

	return result, nil
}
