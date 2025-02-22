package repo

import (
	"backend/internal/model"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type AdvertiserRepo struct {
	db *sqlx.DB
}

func (r *AdvertiserRepo) GetById(id uuid.UUID) (adv model.Advertiser, err error) {
	err = r.db.Get(&adv, "SELECT * FROM advertisers WHERE id = $1", id)
	if errors.Is(err, sql.ErrNoRows) {
		err = ErrNotFound
	}
	return
}

func (r *AdvertiserRepo) GetMany(ids []uuid.UUID) (map[uuid.UUID]model.Advertiser, error) {
	var slice []model.Advertiser
	if err := r.db.Select(&slice, "SELECT * FROM advertisers WHERE id = ANY($1::uuid[])", pq.Array(ids)); err != nil {
		return nil, err
	}

	advs := make(map[uuid.UUID]model.Advertiser)
	for _, adv := range slice {
		advs[adv.Id] = adv
	}
	return advs, nil
}

func (r *AdvertiserRepo) UpsertMany(advertisers []model.Advertiser) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	stmt, err := tx.Prepare(`INSERT INTO advertisers (id, name) VALUES ($1, $2) ON CONFLICT (id) DO UPDATE SET name = $2`)
	if err != nil {
		return fmt.Errorf("prepare statement: %w", err)
	}
	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)
	for _, adv := range advertisers {
		_, err = stmt.Exec(adv.Id, adv.Name)
		if err != nil {
			return fmt.Errorf("exec statement: %w", err)
		}
	}
	return tx.Commit()
}
