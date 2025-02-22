package service

import (
	"backend/config"
	"backend/internal/model"
	"backend/internal/repo"
	"context"
	"encoding/json"
	"fmt"
	"github.com/ollama/ollama/api"
	"log"
	"net/http"
	"net/url"
	"time"
)

type OllamaService struct {
	model            string
	client           *api.Client
	aiRepo           repo.Ai
	enabled          bool
	suggestionsQueue chan model.AiTask
	otherQueue       chan model.AiTask
}

func NewOllamaService(env config.Environment, aiRepo repo.Ai) (*OllamaService, error) {
	urlParsed, err := url.Parse(env.OllamaHost)
	if err != nil {
		return nil, fmt.Errorf("parse ollama url: %w", err)
	}

	s := &OllamaService{
		model:            env.OllamaModel,
		client:           api.NewClient(urlParsed, http.DefaultClient),
		aiRepo:           aiRepo,
		enabled:          !env.RunningInCI,
		suggestionsQueue: make(chan model.AiTask, 1000),
		otherQueue:       make(chan model.AiTask, 5000),
	}
	if s.enabled {
		go s.init()
	}

	return s, nil
}

func (s *OllamaService) SubmitTask(task model.AiTask) {
	if !s.enabled {
		return
	}

	queue := s.otherQueue
	if task.Type == model.AiTaskTypeSuggest {
		queue = s.suggestionsQueue
	}

	select {
	case queue <- task:
	default:
		go func() {
			queue <- task
		}()
	}
}

func (s *OllamaService) init() {
	tasks, err := s.aiRepo.GetIncompleteTasks()
	if err != nil {
		log.Printf("ollama service: failed to get incomplete tasks: %s\n", err)
	} else {
		for _, task := range tasks {
			s.SubmitTask(task)
		}
	}

	req := &api.PullRequest{Model: s.model, Stream: new(bool)}
	callback := func(api.ProgressResponse) error {
		return nil
	}

	start := time.Now()
	for {
		// Sometimes, the download speed becomes too slow after several minutes. Restarting
		// the download resolves the issue. Ollama caches the downloaded files, so it does not
		// disrupt the progress.
		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(2*time.Minute))
		err := s.client.Pull(ctx, req, callback)
		cancel()
		if err == nil {
			break
		}

		log.Printf("ollama service: failed to pull model %s: %s, retrying\n", s.model, err)
		<-time.After(3)
	}

	log.Printf("ollama service: initialized after %s\n", time.Since(start))

	go s.worker(s.suggestionsQueue)
	go s.worker(s.otherQueue)
}

func (s *OllamaService) worker(queue <-chan model.AiTask) {
	for task := range queue {
		s.runTask(task)
	}
}

const systemPrompt = "You are a helpful assistant. Always respond in Russian. Current date: %s"

func (s *OllamaService) runTask(task model.AiTask) {
	answer := ""
	req := &api.GenerateRequest{
		Model:  s.model,
		Prompt: task.Prompt,
		System: fmt.Sprintf(systemPrompt, time.Now().Format("2006-01-02")),
		Format: json.RawMessage(task.Format),
		Stream: new(bool),
	}
	callback := func(resp api.GenerateResponse) error {
		answer = resp.Response
		return nil
	}

	log.Printf("ollama service: started working on task %s\n", task.Id)
	for i := range 5 {
		err := s.client.Generate(context.Background(), req, callback)
		if err == nil {
			break
		}

		if i == 4 {
			log.Printf("OllamaService.runTask: failed to generate: %s (task %s) after 5 attempts\n", err, task.Id)
			return
		}

		log.Printf("OllamaService.runTask: failed to generate: %s (task %s), retrying\n", err, task.Id)
		<-time.After(3)
	}

	log.Printf("ollama service: task %s done\n", task.Id)

	err := s.aiRepo.AddResult(model.AiTaskResult{
		TaskId:    task.Id,
		CreatedAt: time.Now(),
		Answer:    answer,
	})
	if err != nil {
		log.Printf("ollama service: failed to save result for task %s: %s\n", task.Id, err)
	}
}
