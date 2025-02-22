package repo

import (
	"backend/internal/model"
	"github.com/jmoiron/sqlx"
)

type MlScoreRepo struct {
	db *sqlx.DB
}

func (r *MlScoreRepo) Upsert(s model.MlScore) (err error) {
	_, err = r.db.Exec(`INSERT INTO ml_scores (client_id, advertiser_id, score) VALUES ($1, $2, $3)
							  ON CONFLICT (client_id, advertiser_id) DO UPDATE SET score = $3`,
		s.ClientId, s.AdvertiserId, s.Score)
	return
}
