package service

import (
	"context"
	"fmt"

	"github.com/qasimhafeez/taskapi/internal/domain"
	"github.com/qasimhafeez/taskapi/internal/repo"
	"github.com/qasimhafeez/taskapi/pkg/errs"
)

type TaskService struct{
	repo repo.TaskRepo
}

func NewTaskService(r repo.TaskRepo) *TaskService {
	return &TaskService{repo:r}
} 

func(s *TaskService) Create (ctx context.Context, t domain.Task) (domain.Task, error) {
	if t.Status == "" {
		t.Status = domain.DefaultStatus
	}

	if err:= t.Validate(); err != nil {
		return domain.Task{}, fmt.Errorf("%w: %v", errs.ErrBadRequest, err)
	}

	return s.repo.Create(ctx, t)
}

func(s *TaskService) Get (ctx context.Context, id domain.ID) (domain.Task, error){
	return s.repo.Get(ctx, id)
}

func(s *TaskService) List (ctx context.Context, limit, offset int) ([]domain.Task, error){
	return s.repo.List(ctx, limit, offset)
}

func(s *TaskService) Update (ctx context.Context, t domain.Task) (domain.Task, error){
	if err := t.Validate(); err != nil {
		return domain.Task{}, fmt.Errorf("%w: %v", errs.ErrBadRequest, err)
	}

	return s.repo.Update(ctx, t)
}

func(s *TaskService) Delete (ctx context.Context, id domain.ID) error {
	return s.repo.Delete(ctx, id)
}