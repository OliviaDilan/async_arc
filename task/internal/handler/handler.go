package handler

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"

	"github.com/OliviaDilan/async_arc/pkg/auth"
	"github.com/OliviaDilan/async_arc/task/internal/task"
)

type Handler struct {
	taskRepo task.Repository
	authClient auth.Client
}

func NewHandler(taskRepo task.Repository, authClient auth.Client) *Handler {
	return &Handler{
		taskRepo: taskRepo,
		authClient: authClient,
	}
}

func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {

	var req createTaskRequest

	json.NewDecoder(r.Body).Decode(&req)
	defer r.Body.Close()

	if req.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	createdTask, err := h.taskRepo.Create(req.Title); 
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := createTaskResponse{
		Task: taskResponse{
			ID:          createdTask.ID,
			Title:       createdTask.Title,
			Assignee:    createdTask.Assignee,
			Status:      string(createdTask.Status),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Failed to encode users: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetTasks(w http.ResponseWriter, r *http.Request) {

	tasks, err := h.taskRepo.GetAll()
	if err != nil {
		log.Printf("Failed to get tasks: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var tasksRes []taskResponse
	for _, t := range tasks {
		tasksRes = append(tasksRes, taskResponse{
			ID:          t.ID,
			Title:       t.Title,
			Assignee:    t.Assignee,
			Status:      string(t.Status),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(getTasksResponse{Tasks: tasksRes}); err != nil {
		log.Printf("Failed to encode tasks: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) AssignTasks(w http.ResponseWriter, r *http.Request) {

	users, err := h.authClient.GetUsers()
	if err != nil {
		log.Printf("Failed to get users: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var filteredUsers []*auth.User
	for _, u := range users {
		if u.Role == auth.RoleDeveloper {
			filteredUsers = append(filteredUsers, u)
		}
	}

	tasks, err := h.taskRepo.GetByStatus(task.StatusOpen)
	if err != nil {
		log.Printf("Failed to get tasks: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, task := range tasks {
		assignee := filteredUsers[rand.Intn(len(filteredUsers))].Username
		if err := h.taskRepo.Assign(task.ID, assignee); err != nil {
			log.Printf("Failed to assign task: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetTasksByAssignee(w http.ResponseWriter, r *http.Request) {
	//Get assignee from query
	assignee := r.URL.Query().Get("assignee")

	if assignee == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tasks, err := h.taskRepo.GetByAssignee(assignee)
	if err != nil {
		log.Printf("Failed to get tasks: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var tasksRes []taskResponse
	for _, t := range tasks {
		tasksRes = append(tasksRes, taskResponse{
			ID:          t.ID,
			Title:       t.Title,
			Assignee:    t.Assignee,
			Status:      string(t.Status),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(getTasksByAssigneeResponse{Tasks: tasksRes}); err != nil {
		log.Printf("Failed to encode tasks: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) CloseTask(w http.ResponseWriter, r *http.Request) {

	var req closeTaskRequest

	json.NewDecoder(r.Body).Decode(&req)
	defer r.Body.Close()

	if req.TaskID <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.taskRepo.Close(req.TaskID); err != nil {
		log.Printf("Failed to close task: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}