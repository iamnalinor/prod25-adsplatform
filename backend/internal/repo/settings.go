package repo

import (
	"backend/internal/model"
	"github.com/jmoiron/sqlx"
	"log"
)

type SettingsRepo struct {
	db     *sqlx.DB
	cached model.Settings
}

func NewSettingsRepo(db *sqlx.DB) *SettingsRepo {
	r := &SettingsRepo{db: db}
	settings, err := r.Get()
	if err != nil {
		log.Fatalf("get settings: %s\n", err)
	}
	r.cached = settings
	return r
}

func (r *SettingsRepo) Get() (s model.Settings, err error) {
	err = r.db.Get(&s, `SELECT "current_date" FROM settings LIMIT 1`)
	return
}

func (r *SettingsRepo) GetCached() model.Settings {
	return r.cached
}

func (r *SettingsRepo) Update(settings model.Settings) error {
	_, err := r.db.Exec(`UPDATE settings SET "current_date" = $1 WHERE id > 0`, settings.CurrentDate)
	r.cached = settings
	return err
}
