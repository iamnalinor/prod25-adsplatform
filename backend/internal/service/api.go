package service

import (
	"backend/internal/repo"
	"log"
	"time"
)

type requestInfo struct {
	endpoint string
	duration time.Duration
}

type ApiService struct {
	apiRepo repo.Api
	queue   chan requestInfo
}

func NewApiService(apiRepo repo.Api) *ApiService {
	s := &ApiService{apiRepo, make(chan requestInfo, 5000)}
	go s.worker()
	return s
}

// worker reads requests from queue and adds them to database.
func (s *ApiService) worker() {
	for r := range s.queue {
		err := s.apiRepo.AddRequest(r.endpoint, float64(r.duration)/float64(time.Millisecond))
		if err != nil {
			log.Printf("requests logger middleware: add request: %s\n", err)
		}
	}
}

// LogRequest adds request to database. This function is non-blocking.
func (s *ApiService) LogRequest(method, path string, duration time.Duration) {
	if path == "/ping" {
		return
	}
	s.queue <- requestInfo{method + " " + path, duration}
}
