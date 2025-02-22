package repo

import (
	"github.com/jmoiron/sqlx"
)

type ApiRepo struct {
	db *sqlx.DB
}

func (r *ApiRepo) AddRequest(endpoint string, durationMs float64) error {
	_, err := r.db.Exec(`INSERT INTO api_requests (endpoint, duration_ms) VALUES ($1, $2)`, endpoint, durationMs)
	return err
}
