package model

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type AiTaskType string

const (
	AiTaskTypeSuggest    AiTaskType = "suggestText"
	AiTaskTypeModeration AiTaskType = "moderation"
)

type AiTask struct {
	Id        uuid.UUID  `db:"id"`
	CreatedAt time.Time  `db:"created_at"`
	Type      AiTaskType `db:"type"`
	Prompt    string     `db:"prompt"`
	Format    string     `db:"format"`
}

type AiTaskResult struct {
	TaskId    uuid.UUID `db:"task_id"`
	CreatedAt time.Time `db:"created_at"`
	Answer    string    `db:"answer"`
}

type AiModerationResult struct {
	Acceptable bool   `json:"acceptable"`
	Reason     string `json:"reason"`
}

// Scan implements the sql.Scanner interface for AiModerationResult.
func (r *AiModerationResult) Scan(value any) error {
	if value == nil {
		*r = AiModerationResult{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("value %v must be of type []byte", value)
	}

	return json.Unmarshal(bytes, r)
}

type AiTaskResponse struct {
	Id          uuid.UUID           `json:"task_id"`
	CreatedAt   time.Time           `json:"created_at"`
	Completed   bool                `json:"completed"`
	Suggestions []string            `json:"suggestions,omitempty"`
	Moderation  *AiModerationResult `json:"moderation,omitempty"`
}
