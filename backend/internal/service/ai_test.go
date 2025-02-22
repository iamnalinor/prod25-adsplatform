package service

import (
	"backend/internal/model"
	"backend/internal/repo"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAiRepo struct {
	mock.Mock
}

func (r *MockAiRepo) AddTask(task model.AiTask) error {
	args := r.Called(task)
	return args.Error(0)
}

func (r *MockAiRepo) GetTask(id uuid.UUID) (model.AiTask, error) {
	args := r.Called(id)
	return args.Get(0).(model.AiTask), args.Error(1)
}

func (r *MockAiRepo) GetIncompleteTasks() ([]model.AiTask, error) {
	args := r.Called()
	return args.Get(0).([]model.AiTask), args.Error(1)
}

func (r *MockAiRepo) AddResult(result model.AiTaskResult) error {
	args := r.Called(result)
	return args.Error(0)
}

func (r *MockAiRepo) GetResult(taskId uuid.UUID) (model.AiTaskResult, error) {
	args := r.Called(taskId)
	return args.Get(0).(model.AiTaskResult), args.Error(1)
}

func TestGetTask(t *testing.T) {
	mockRepo := new(MockAiRepo)
	service := &AiService{aiRepo: mockRepo}

	taskID := uuid.New()
	task2ID := uuid.New()
	task3ID := uuid.New()

	mockRepo.On("GetTask", taskID).Return(model.AiTask{Id: taskID, CreatedAt: time.Now(), Type: model.AiTaskTypeSuggest}, nil)
	mockRepo.On("GetResult", taskID).Return(model.AiTaskResult{Answer: `["text 1", "text 2", "text 3"]`}, nil)

	mockRepo.On("GetTask", task2ID).Return(model.AiTask{Id: task2ID, CreatedAt: time.Now(), Type: model.AiTaskTypeModeration}, nil)
	mockRepo.On("GetResult", task2ID).Return(model.AiTaskResult{Answer: `{"acceptable": true, "reason": ""}`}, nil)

	mockRepo.On("GetTask", task3ID).Return(model.AiTask{Id: task3ID, CreatedAt: time.Now(), Type: model.AiTaskTypeSuggest}, nil)
	mockRepo.On("GetResult", task3ID).Return(model.AiTaskResult{}, repo.ErrNotFound)

	mockRepo.On("GetTask", uuid.Nil).Return(model.AiTask{}, repo.ErrNotFound)

	resp, err := service.GetTask(taskID)
	assert.NoError(t, err)
	assert.Equal(t, taskID, resp.Id)
	assert.True(t, resp.Completed)
	assert.Equal(t, resp.Suggestions, []string{"text 1", "text 2", "text 3"})

	resp, err = service.GetTask(task2ID)
	assert.NoError(t, err)
	assert.Equal(t, task2ID, resp.Id)
	assert.True(t, resp.Completed)
	assert.Equal(t, resp.Moderation, &model.AiModerationResult{Acceptable: true, Reason: ""})

	resp, err = service.GetTask(task3ID)
	assert.NoError(t, err)
	assert.Equal(t, task3ID, resp.Id)
	assert.False(t, resp.Completed)

	resp, err = service.GetTask(uuid.Nil)
	assert.Error(t, err)

	mockRepo.AssertExpectations(t)
}

func TestSubmitSuggestText(t *testing.T) {
	mockRepo := new(MockAiRepo)
	service := &AiService{aiRepo: mockRepo, ollamaSvc: &OllamaService{}}

	mockRepo.On("AddTask", mock.Anything).Return(nil).Once()
	_, err := service.SubmitSuggestText("name", "title", "")
	assert.NoError(t, err)

	mockRepo.On("AddTask", mock.Anything).Return(errors.New("error")).Once()
	_, err = service.SubmitSuggestText("name", "title", "")
	assert.Error(t, err)

	mockRepo.AssertExpectations(t)
}

func TestSubmitModeration(t *testing.T) {
	mockRepo := new(MockAiRepo)
	service := &AiService{aiRepo: mockRepo, ollamaSvc: &OllamaService{}}

	mockRepo.On("AddTask", mock.Anything).Return(nil).Once()
	_, err := service.SubmitModeration("title", "text")
	assert.NoError(t, err)

	mockRepo.On("AddTask", mock.Anything).Return(errors.New("error")).Once()
	_, err = service.SubmitModeration("title", "text")
	assert.Error(t, err)

	mockRepo.AssertExpectations(t)
}
