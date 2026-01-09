package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/MrGreenboom/go-task-manager/internal/model"
	"github.com/MrGreenboom/go-task-manager/internal/service"
)

type TaskHandler struct {
	svc *service.TaskService
}

func NewTaskHandler(svc *service.TaskService) *TaskHandler {
	return &TaskHandler{svc: svc}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func (h *TaskHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/tasks", h.handleTasks)
	mux.HandleFunc("/tasks/", h.handleTaskByID) // /tasks/{id}
}

func (h *TaskHandler) handleTasks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		tasks, err := h.svc.List(r.Context())
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, tasks)

	case http.MethodPost:
		var t model.Task
		if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
			writeError(w, http.StatusBadRequest, "invalid json")
			return
		}
		id, err := h.svc.Create(r.Context(), &t)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeJSON(w, http.StatusCreated, map[string]int64{"id": id})

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *TaskHandler) handleTaskByID(w http.ResponseWriter, r *http.Request) {
	// path: /tasks/{id}
	idStr := strings.TrimPrefix(r.URL.Path, "/tasks/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	switch r.Method {
	case http.MethodGet:
		t, err := h.svc.GetByID(r.Context(), id)
		if err != nil {
			writeError(w, http.StatusNotFound, "task not found")
			return
		}
		writeJSON(w, http.StatusOK, t)

	case http.MethodPut:
		var t model.Task
		if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
			writeError(w, http.StatusBadRequest, "invalid json")
			return
		}
		t.ID = id
		if err := h.svc.Update(r.Context(), &t); err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		w.WriteHeader(http.StatusNoContent)

	case http.MethodDelete:
		if err := h.svc.Delete(r.Context(), id); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		w.WriteHeader(http.StatusNoContent)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
