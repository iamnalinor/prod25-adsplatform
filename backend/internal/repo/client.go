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

type ClientRepo struct {
	db *sqlx.DB
}

func (r *ClientRepo) GetById(id uuid.UUID) (client model.Client, err error) {
	err = r.db.Get(&client, "SELECT * FROM clients WHERE id = $1", id)
	if errors.Is(err, sql.ErrNoRows) {
		err = ErrNotFound
	}
	return
}

func (r *ClientRepo) GetMany(ids []uuid.UUID) (map[uuid.UUID]model.Client, error) {
	var slice []model.Client
	if err := r.db.Select(&slice, "SELECT * FROM clients WHERE id = ANY($1::uuid[])", pq.Array(ids)); err != nil {
		return nil, err
	}

	clients := make(map[uuid.UUID]model.Client)
	for _, client := range slice {
		clients[client.Id] = client
	}
	return clients, nil
}

func (r *ClientRepo) UpsertMany(clients []model.Client) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	stmt, err := tx.Prepare(`INSERT INTO clients (id, login, age, location, gender) VALUES ($1, $2, $3, $4, $5) 
                                    ON CONFLICT (id) DO UPDATE SET (login, age, location, gender) = ($2, $3, $4, $5)`)
	if err != nil {
		return fmt.Errorf("prepare statement: %w", err)
	}
	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)
	for _, client := range clients {
		_, err = stmt.Exec(client.Id, client.Login, client.Age, client.Location, client.Gender)
		if err != nil {
			return fmt.Errorf("exec statement: %w", err)
		}
	}
	return tx.Commit()
}
