package repo

import (
	"context"
	"sync"
	"time"

	"github.com/qasimhafeez/taskapi/internal/domain"
	"github.com/qasimhafeez/taskapi/pkg/errs"
)

// in mem struct
type InMem struct {
	mu sync.RWMutex
	store map[domain.ID]domain.Task
}

func NewInMem() *InMem {
	return &InMem {store: make(map[domain.ID]domain.Task)}
}

func (r *InMem) Create (ctx context.Context, t domain.Task) (domain.Task, error) {
	// For Context Done
	select {
	case <-ctx.Done():
		return domain.Task{}, ctx.Err()
	default:
	}

	// locking and release
	r.mu.Lock()
	defer r.mu.Unlock()

	// Duplicate Check 

	if _, exists := r.store[t.ID]; exists {
		return domain.Task{}, errs.ErrConflict
	}

	now := time.Now().UTC()
	t.CreatedAt, t.UpdatedAt = now, now
	r.store[t.ID] = t
	return t, nil
}

func (r *InMem) Get(ctx context.Context, id domain.ID) (domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if t, ok := r.store[id]; ok {
		return t, nil
	}
	return domain.Task{}, errs.ErrNotFound
}

func (r *InMem) List (ctx context.Context, limit, offset int) ([]domain.Task, error){
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]domain.Task, 0, len(r.store))
	for _,t := range r.store{
		out = append(out, t)
	}

	if offset > len(out) {
		return []domain.Task{}, nil
	}

	end := offset + limit
	if limit <=0 || end > len(out){
		end = len(out)
	}
	return out[offset:end], nil
}

func (r *InMem) Update (ctx context.Context, t domain.Task) (domain.Task, error){
	r.mu.Lock()
	defer r.mu.Unlock()

	if _,ok := r.store[t.ID]; !ok {
		return domain.Task{}, errs.ErrNotFound
	}

	t.UpdatedAt = time.Now().UTC()
	r.store[t.ID] = t
	return t, nil
}

func (r *InMem) Delete (ctx context.Context, id domain.ID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.store[id]; !exists {
		return errs.ErrNotFound
	}

	delete(r.store, id)
	return nil
}