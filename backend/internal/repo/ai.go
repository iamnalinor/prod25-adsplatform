package repo

import (
	"backend/internal/model"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type AiRepo struct {
	db *sqlx.DB
}

func (r *AiRepo) AddTask(task model.AiTask) (err error) {
	_, err = r.db.Exec(`INSERT INTO ai_tasks (id, created_at, type, prompt, "format") 
		VALUES ($1, $2, $3, $4, $5)`,
		task.Id, task.CreatedAt, task.Type, task.Prompt, task.Format)
	return
}

func (r *AiRepo) GetTask(id uuid.UUID) (task model.AiTask, err error) {
	err = r.db.Get(&task, `SELECT * FROM ai_tasks WHERE id = $1`, id)
	if errors.Is(err, sql.ErrNoRows) {
		err = ErrNotFound
	}
	return
}

func (r *AiRepo) GetIncompleteTasks() ([]model.AiTask, error) {
	tasks := make([]model.AiTask, 0)
	err := r.db.Select(&tasks, `SELECT * FROM ai_tasks WHERE id NOT IN (SELECT task_id FROM ai_task_results)`)
	return tasks, err
}

func (r *AiRepo) AddResult(result model.AiTaskResult) (err error) {
	_, err = r.db.Exec(`INSERT INTO ai_task_results (task_id, created_at, answer) VALUES ($1, $2, $3)`,
		result.TaskId, result.CreatedAt, result.Answer)
	return
}

func (r *AiRepo) GetResult(taskId uuid.UUID) (res model.AiTaskResult, err error) {
	err = r.db.Get(&res, `SELECT * FROM ai_task_results WHERE task_id = $1`, taskId)
	if errors.Is(err, sql.ErrNoRows) {
		err = ErrNotFound
	}
	return
}
