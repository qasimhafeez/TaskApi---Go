package domain

import (
	"errors"
	"strings"
	"time"
)

type ID string
type Status string

const (
	StatusTodo Status = "StatusTodo"
	StatusInProgress Status = "StatusInProgress"
	StatusDone Status = "StatusDone"
	DefaultStatus = StatusTodo
)

type Task struct {
	ID ID `json:"id"`
	Title string `json:"title"`
	Notes string `json:"notes"`
	Status Status `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (t *Task) Touch () {t.UpdatedAt = time.Now().UTC()}

var (
	ErrInvalidTitle = errors.New("Invalid Title")
	ErrInvalidStatus = errors.New("Invalid Status")
)

func (t *Task) Validate() error {
	if strings.TrimSpace(t.Title) == "" {
		return ErrInvalidTitle
	}
	switch t.Status {
	case StatusTodo,StatusInProgress, StatusDone:
	default:
		return ErrInvalidStatus
	}
	return nil
}