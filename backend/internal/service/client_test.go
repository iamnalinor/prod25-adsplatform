package service

import (
	"backend/internal/model"
	"backend/internal/repo"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockClientRepo struct {
	mock.Mock
}

func (m *MockClientRepo) GetById(id uuid.UUID) (model.Client, error) {
	args := m.Called(id)
	return args.Get(0).(model.Client), args.Error(1)
}

func (m *MockClientRepo) GetMany(ids []uuid.UUID) (map[uuid.UUID]model.Client, error) {
	args := m.Called(ids)
	return args.Get(0).(map[uuid.UUID]model.Client), args.Error(1)
}

func (m *MockClientRepo) UpsertMany(clients []model.Client) error {
	args := m.Called(clients)
	return args.Error(0)
}

func TestClientService_GetById(t *testing.T) {
	cRepo := &MockClientRepo{}
	cService := &ClientService{cRepo}

	clientId := uuid.New()
	age := 42
	client := model.Client{Id: clientId, Login: "solution", Age: &age, Location: "Ducat Place II", Gender: "MALE"}
	cRepo.On("GetById", clientId).Return(client, nil)
	cRepo.On("GetById", uuid.Nil).Return(model.Client{}, repo.ErrNotFound)

	got, err := cService.GetById(clientId)
	if got != client || err != nil {
		t.Errorf("GetById() got = %v %v, want %v nil", got, err, client)
	}

	got, err = cService.GetById(uuid.Nil)
	if (got != model.Client{}) || !repo.IsNotFound(err) {
		t.Errorf("GetById() got = %v %v, want {} ErrNoRows", got, err)
	}

	cRepo.AssertExpectations(t)
}
