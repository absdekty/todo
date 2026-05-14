package rest

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"net/http"
	"todo/internal/model"
	"todo/internal/service"
)

type RestHandler struct {
	serv service.ServiceI
}

func NewHandler(serv service.ServiceI) *RestHandler {
	return &RestHandler{serv: serv}
}

func (h *RestHandler) mainHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Default Page"))
}

func (h *RestHandler) getAllTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.serv.GetTasks(r.Context())
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":  tasks,
		"count": len(tasks),
	})
}

func (h *RestHandler) getTaskByID(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "id")
	task, err := h.serv.GetTaskByID(r.Context(), taskID)
	if err != nil {
		if errors.Is(err, model.ErrTaskInvalidID) {
			http.Error(w, "Invalid task ID format", http.StatusBadRequest)
			return
		}

		if errors.Is(err, model.ErrTaskNotExist) {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}

		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(TaskResponse{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Completed:   task.Completed,
		CreatedAt:   task.CreatedAt})
}

func (h *RestHandler) createTask(w http.ResponseWriter, r *http.Request) {
	var req TaskCreate
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON format", http.StatusBadRequest)
		return
	}

	task, err := h.serv.CreateTask(r.Context(), req.Title, req.Description)
	if err != nil {
		if errors.Is(err, model.ErrTaskTitleTooShort) || errors.Is(err, model.ErrTaskTitleTooLong) {
			http.Error(w, "Incorrect title: symbols 3-15", http.StatusBadRequest)
			return
		}

		if errors.Is(err, model.ErrTaskDescriptionTooLong) {
			http.Error(w, "Incorrect description: symbols >500", http.StatusBadRequest)
			return
		}

		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(TaskResponse{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Completed:   task.Completed,
		CreatedAt:   task.CreatedAt})
}

func (h *RestHandler) putTask(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "id")
	task, err := h.serv.GetTaskByID(r.Context(), taskID)
	if err != nil {
		if errors.Is(err, model.ErrTaskInvalidID) {
			http.Error(w, "Invalid task ID format", http.StatusBadRequest)
			return
		}

		if errors.Is(err, model.ErrTaskNotExist) {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}

		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var req TaskUpdateFull
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON format", http.StatusBadRequest)
		return
	}

	if err := task.SetTitle(req.Title); err != nil {
		http.Error(w, "Incorrect title: symbols 3-15", http.StatusBadRequest)
		return
	}

	if err := task.SetDescription(req.Description); err != nil {
		http.Error(w, "Incorrect description: symbols >500", http.StatusBadRequest)
		return
	}

	task.SetCompleted(req.Completed)

	err = h.serv.UpdateTask(r.Context(), task)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(TaskResponse{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Completed:   task.Completed,
		CreatedAt:   task.CreatedAt})
}

func (h *RestHandler) patchTask(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "id")
	task, err := h.serv.GetTaskByID(r.Context(), taskID)
	if err != nil {
		if errors.Is(err, model.ErrTaskInvalidID) {
			http.Error(w, "Invalid task ID format", http.StatusBadRequest)
			return
		}

		if errors.Is(err, model.ErrTaskNotExist) {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}

		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "invalid JSON format", http.StatusBadRequest)
		return
	}

	if title, ok := updates["title"].(string); ok && title != "" {
		task.SetTitle(title)
	}
	if desc, ok := updates["desc"].(string); ok {
		task.SetDescription(desc)
	}
	if completed, ok := updates["completed"].(bool); ok {
		task.SetCompleted(completed)
	}

	if err := h.serv.UpdateTask(r.Context(), task); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(TaskResponse{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Completed:   task.Completed,
		CreatedAt:   task.CreatedAt})
}

func (h *RestHandler) deleteTask(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "id")
	if err := h.serv.DeleteTask(r.Context(), taskID); err != nil {
		if errors.Is(err, model.ErrTaskInvalidID) {
			http.Error(w, "Invalid task ID format", http.StatusBadRequest)
			return
		}

		if errors.Is(err, model.ErrTaskNotExist) {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}

		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
