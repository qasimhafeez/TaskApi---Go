package httpx

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/qasimhafeez/taskapi/internal/domain"
	"github.com/qasimhafeez/taskapi/internal/service"
	"github.com/qasimhafeez/taskapi/pkg/errs"
)

type Handler struct{
	svc *service.TaskService
	log *slog.Logger
}

func NewHandler(svc *service.TaskService, log *slog.Logger) *Handler{
	return &Handler{svc: svc, log: log}
}

func (h *Handler) Router() http.Handler{
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", h.health)
	mux.HandleFunc("GET /v1/tasks", h.listTask)
	mux.HandleFunc("POST /v1/tasks", h.createTask)
	mux.HandleFunc("GET /v1/tasks/{id}", h.getTask)
	mux.HandleFunc("PUT /v1/tasks/{id}", h.updateTask)
	mux.HandleFunc("DELETE /v1/tasks/{id}", h.deleteTask)
	return withLogging(h.log, mux)
}

func (h *Handler) health (w http.ResponseWriter, r *http.Request){
	writeJSON(w, http.StatusOK, map[string]any{"ok":true})
}

func (h *Handler) createTask (w http.ResponseWriter, r *http.Request){
	var in struct {
		ID string `json:"id"`
		Title string `json:"title"`
		Notes string `json:"notes"`
		Status domain.Status `json:"status"`
	 }

	 if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	 }

	 ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	 defer cancel()

	 task, err := h.svc.Create(ctx, domain.Task{
		ID: domain.ID(in.ID),
		Title: in.Title,
		Notes: in.Notes,
		Status: in.Status,
	 })

	 if err != nil {
		status := http.StatusInternalServerError
		switch {
		case errors.Is(err, errs.ErrBadRequest):
			status = http.StatusBadRequest
		case errors.Is(err, errs.ErrConflict):
			status = http.StatusConflict
		}
		writeError(w, status, err.Error())
		return
	 }
	 writeJSON(w, http.StatusCreated, task)
}

func (h *Handler) getTask(w http.ResponseWriter, r *http.Request){
	id := r.PathValue("id")
	ctx, cancel  := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	task, err := h.svc.Get(ctx, domain.ID(id))
	if err != nil {
		if errors.Is(err, errs.ErrNotFound){
			writeError(w, http.StatusNotFound, "task not found")
			return
		}
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, task)
}

func (h *Handler) listTask (w http.ResponseWriter, r *http.Request){
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	taskList, err := h.svc.List(ctx, 50, 0)
	if err != nil{
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, taskList)
}

func (h *Handler) updateTask(w http.ResponseWriter, r *http.Request){
	id := r.PathValue("id")
	var in domain.Task
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	in.ID = domain.ID(id)

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	task, err := h.svc.Update(ctx, in)
	switch{
	case err == nil:
		writeJSON(w, http.StatusOK, task)
	case errors.Is(err, errs.ErrNotFound):
		writeError(w, http.StatusNotFound, "not found")
	case errors.Is(err, errs.ErrBadRequest):
		writeError(w, http.StatusBadRequest, err.Error())
	default:
		writeError(w, http.StatusInternalServerError, err.Error())
	}
}

func (h *Handler) deleteTask(w http.ResponseWriter, r *http.Request){
	id := r.PathValue("id")
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if err := h.svc.Delete(ctx, domain.ID(id)); err != nil {
		if (errors.Is(err, errs.ErrNotFound)){
			writeError(w, http.StatusNotFound, "not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
	}
	writeJSON(w, http.StatusOK, nil)
}

