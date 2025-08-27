package service

import (
	"context"
	"testing"

	"github.com/qasimhafeez/taskapi/internal/domain"
	"github.com/qasimhafeez/taskapi/internal/repo"
)

func TestCreateAndGet(t *testing.T){
	r := repo.NewInMem()
	svc := NewTaskService(r)
	
	created, err := svc.Create(context.Background(), domain.Task{
		ID: "1",
		Title: "task 1",
	})

	if err != nil {
		t.Fatalf("create err: %v", err)
	}

	got, err := svc.Get(context.Background(), domain.ID(created.ID))
	if err != nil {
		t.Fatalf("get err: %v", err)
	}
	if got.Title != "task 1" {
		t.Fatalf("wrong title err: %v", err)
	}
}