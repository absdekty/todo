package rest

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"net/http"
	"time"
	"todo/internal/model"
	"todo/pkg/logger"
)

type RestHandler struct {
	task    ServiceTask
	user    ServiceUser
	token   ServiceToken
	metrics *Metrics
}

func NewHandler(task ServiceTask, user ServiceUser, token ServiceToken, metrics *Metrics) *RestHandler {
	return &RestHandler{task: task, user: user, token: token, metrics: metrics}
}

func (h *RestHandler) mainHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Default Page"))
}

func (h *RestHandler) GetMetrics(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"total_requests":   h.metrics.GetTotalRequests(),
		"active_requests":  h.metrics.GetActiveRequests(),
		"errors_total":     h.metrics.GetErrorsTotal(),
		"errors_by_status": h.metrics.GetErrorsByStatus()})
}

func (h *RestHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req UserCreate
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Debug.Printf("registeruser[json]: %v", err)
		http.Error(w, "invalid JSON format", http.StatusBadRequest)
		return
	}

	if err := h.user.CreateUser(r.Context(), req.Name, req.Password); err != nil {
		if errors.Is(err, model.ErrUserNameTooShort) || errors.Is(err, model.ErrUserNameTooLong) {
			logger.Debug.Printf("registeruser[name]: %v", err)
			http.Error(w, "Name must be 3-10 characters", http.StatusBadRequest)
			return
		}

		if errors.Is(err, model.ErrUserPasswordTooShort) || errors.Is(err, model.ErrUserPasswordTooLong) {
			logger.Debug.Printf("registeruser[password]: %v", err)
			http.Error(w, "Password must be 8-16 characters", http.StatusBadRequest)
			return
		}

		if errors.Is(err, model.ErrUserAlreadyExist) {
			logger.Debug.Printf("registeruser[exist]: %v", err)
			http.Error(w, "User already exists", http.StatusConflict)
			return
		}

		logger.Error.Printf("registeruser[other]: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *RestHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var req UserCreate
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Debug.Printf("loginuser[json]: %v", err)
		http.Error(w, "invalid JSON format", http.StatusBadRequest)
		return
	}

	userID, err := h.user.Login(r.Context(), req.Name, req.Password)
	if err != nil {
		if errors.Is(err, model.ErrUserNotExist) {
			logger.Debug.Printf("loginuser[not exist]: %v", err)
			http.Error(w, "User not exist", http.StatusNotFound)
			return
		}

		if errors.Is(err, model.ErrUserInvalidPW) {
			logger.Debug.Printf("loginuser[wrong pw]: %v", err)
			http.Error(w, "Invalid password", http.StatusUnauthorized)
			return
		}

		logger.Error.Printf("loginuser[login other]: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	accessToken, refreshToken, err := h.token.GenerateTokens(r.Context(), userID)
	if err != nil {
		logger.Error.Printf("loginuser[token]: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		MaxAge:   int(7 * 24 * time.Hour / time.Second),
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id":      userID,
		"access_token": accessToken,
		"token_type":   "Bearer",
	})
}

func (h *RestHandler) RefreshTokens(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		logger.Debug.Printf("refreshtoken[cookie]: %v", err)
		http.Error(w, "refresh token required", http.StatusUnauthorized)
		return
	}

	refreshToken := cookie.Value

	newAccessToken, newRefreshToken, err := h.token.RefreshTokens(r.Context(), refreshToken)
	if err != nil {
		logger.Debug.Printf("refreshtoken[token]: %v", err)
		http.Error(w, "invalid refresh token", http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    newRefreshToken,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		MaxAge:   int(7 * 24 * time.Hour / time.Second),
	})

	json.NewEncoder(w).Encode(map[string]interface{}{
		"access_token": newAccessToken,
		"token_type":   "Bearer",
	})
}

func (h *RestHandler) LogoutUser(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err == nil {
		if err := h.token.RevokeRefreshToken(r.Context(), cookie.Value); err != nil {
			logger.Error.Printf("logoutuser: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		MaxAge:   -1,
	})

	w.WriteHeader(http.StatusOK)
}

func (h *RestHandler) getAllTasks(w http.ResponseWriter, r *http.Request) {
	userID, err := GetUserID(r.Context())
	if err != nil {
		logger.Debug.Printf("getalltasks[userid]: %v", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	tasks, err := h.task.GetUserTasks(r.Context(), userID)
	if err != nil {
		if errors.Is(err, model.ErrUserNotExist) {
			logger.Debug.Printf("getalltasks[not exist]: %v", err)
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		logger.Error.Printf("getalltasks: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
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
	userID, err := GetUserID(r.Context())
	if err != nil {
		logger.Debug.Printf("getbyid[userid]: %v", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	taskID := chi.URLParam(r, "id")
	task, err := h.task.GetTaskByID(r.Context(), userID, taskID)
	if err != nil {
		if errors.Is(err, model.ErrUserNotExist) {
			logger.Debug.Printf("getbyid[user not exist]: %v", err)
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		if errors.Is(err, model.ErrTaskInvalidID) {
			logger.Debug.Printf("getbyid[invalid taskid]: %v", err)
			http.Error(w, "Invalid task ID format", http.StatusBadRequest)
			return
		}

		if errors.Is(err, model.ErrTaskNotExist) {
			logger.Debug.Printf("getbyid[task not exist]: %v", err)
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}

		logger.Error.Printf("gettaskbyid: %v", err)
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
		CreatedAt:   task.CreatedAt,
	})
}

func (h *RestHandler) createTask(w http.ResponseWriter, r *http.Request) {
	userID, err := GetUserID(r.Context())
	if err != nil {
		logger.Debug.Printf("createtask[userid]: %v", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req TaskCreate
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Debug.Printf("createtask[json]: %v", err)
		http.Error(w, "invalid JSON format", http.StatusBadRequest)
		return
	}

	task, err := h.task.CreateTask(r.Context(), userID, req.Title, req.Description)
	if err != nil {
		if errors.Is(err, model.ErrUserNotExist) {
			logger.Debug.Printf("createtask[user not exist]: %v", err)
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		if errors.Is(err, model.ErrTaskTitleTooShort) || errors.Is(err, model.ErrTaskTitleTooLong) {
			logger.Debug.Printf("createtask[title]: %v", err)
			http.Error(w, "Incorrect title: symbols 3-15", http.StatusBadRequest)
			return
		}

		if errors.Is(err, model.ErrTaskDescriptionTooLong) {
			logger.Debug.Printf("createtask[desc too long]: %v", err)
			http.Error(w, "Incorrect description: symbols >500", http.StatusBadRequest)
			return
		}

		logger.Error.Printf("createtask: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(TaskResponse{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Completed:   task.Completed,
		CreatedAt:   task.CreatedAt,
	})
}

func (h *RestHandler) putTask(w http.ResponseWriter, r *http.Request) {
	userID, err := GetUserID(r.Context())
	if err != nil {
		logger.Debug.Printf("puttask[userid]: %v", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	taskID := chi.URLParam(r, "id")
	task, err := h.task.GetTaskByID(r.Context(), userID, taskID)
	if err != nil {
		if errors.Is(err, model.ErrUserNotExist) {
			logger.Debug.Printf("puttask[user not exist]: %v", err)
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		if errors.Is(err, model.ErrTaskInvalidID) {
			logger.Debug.Printf("puttask[invalid taskid]: %v", err)
			http.Error(w, "Invalid task ID format", http.StatusBadRequest)
			return
		}

		if errors.Is(err, model.ErrTaskNotExist) {
			logger.Debug.Printf("puttask[task not exist]: %v", err)
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}

		logger.Error.Printf("puttask: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var req TaskUpdateFull
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Debug.Printf("puttask[json]: %v", err)
		http.Error(w, "invalid JSON format", http.StatusBadRequest)
		return
	}

	if err := task.SetTitle(req.Title); err != nil {
		logger.Debug.Printf("puttask[title]: %v", err)
		http.Error(w, "Incorrect title: symbols 3-15", http.StatusBadRequest)
		return
	}

	if err := task.SetDescription(req.Description); err != nil {
		logger.Debug.Printf("puttask[desc]: %v", err)
		http.Error(w, "Incorrect description: symbols >500", http.StatusBadRequest)
		return
	}

	task.SetCompleted(req.Completed)

	err = h.task.UpdateTask(r.Context(), task)
	if err != nil {
		logger.Error.Printf("puttask: %v", err)
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
		CreatedAt:   task.CreatedAt,
	})
}

func (h *RestHandler) patchTask(w http.ResponseWriter, r *http.Request) {
	userID, err := GetUserID(r.Context())
	if err != nil {
		logger.Debug.Printf("patchtask[userid]: %v", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	taskID := chi.URLParam(r, "id")
	task, err := h.task.GetTaskByID(r.Context(), userID, taskID)
	if err != nil {
		if errors.Is(err, model.ErrUserNotExist) {
			logger.Debug.Printf("patchtask[user not exist]: %v", err)
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		if errors.Is(err, model.ErrTaskInvalidID) {
			logger.Debug.Printf("patchtask[invalid taskid]: %v", err)
			http.Error(w, "Invalid task ID format", http.StatusBadRequest)
			return
		}

		if errors.Is(err, model.ErrTaskNotExist) {
			logger.Debug.Printf("patchtask[task not exist]: %v", err)
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}

		logger.Error.Printf("patchtask: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		logger.Debug.Printf("patchtask[json]: %v", err)
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

	if err := h.task.UpdateTask(r.Context(), task); err != nil {
		logger.Error.Printf("patchtask: %v", err)
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
		CreatedAt:   task.CreatedAt,
	})
}

func (h *RestHandler) deleteTask(w http.ResponseWriter, r *http.Request) {
	userID, err := GetUserID(r.Context())
	if err != nil {
		logger.Debug.Printf("deletetask[userid]: %v", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	taskID := chi.URLParam(r, "id")
	if err := h.task.DeleteTask(r.Context(), userID, taskID); err != nil {
		if errors.Is(err, model.ErrUserNotExist) {
			logger.Debug.Printf("deletetask[user not exist]: %v", err)
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		if errors.Is(err, model.ErrTaskInvalidID) {
			logger.Debug.Printf("deletetask[invalid taskid]: %v", err)
			http.Error(w, "Invalid task ID format", http.StatusBadRequest)
			return
		}

		if errors.Is(err, model.ErrTaskNotExist) {
			logger.Debug.Printf("deletetask[task not exist]: %v", err)
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}

		logger.Error.Printf("deletetask: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
