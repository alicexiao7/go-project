
package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"my-go-project/http.serv/internal/model"
	"my-go-project/http.serv/internal/service"
)

type TaskHandler struct {
	service *service.TaskService
}

func NewTaskHandler(s *service.TaskService) *TaskHandler {
	return &TaskHandler{
		service: s,
	}
}

func (h *TaskHandler) Tasks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		tasks, err := h.service.GetTasks()
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tasks)

	case http.MethodPost:
		var newTask model.TaskDTO

		err := json.NewDecoder(r.Body).Decode(&newTask)
		if err != nil {
			log.Println("decode error:", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		task, err := h.service.CreateTask(newTask.Title)
		if err != nil {
			log.Println("create task error:", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(task)

	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func (h *TaskHandler) TaskByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		task, err := h.service.GetTaskByID(id)
		if err != nil {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(task)

	case http.MethodDelete:
		err := h.service.DeleteTask(id)
		if err != nil {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "delete ok",
		})

	case http.MethodPut:
		var updatedTask model.TaskDTO

		err := json.NewDecoder(r.Body).Decode(&updatedTask)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		if updatedTask.Title == "" {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		task, err := h.service.UpdateTask(id, updatedTask)
		if err != nil {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(task)

	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}
