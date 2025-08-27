package repo

import (
	"context"

	"github.com/qasimhafeez/taskapi/internal/domain"
)

type TaskRepo interface {
	Create(ctx context.Context, t domain.Task) (domain.Task, error)
	Get(ctx context.Context, t domain.ID) (domain.Task, error)
	List(ctx context.Context, limit, offset int) ([]domain.Task, error)
	Update(ctx context.Context, t domain.Task) (domain.Task, error)
	Delete(ctx context.Context, id domain.ID) error
}