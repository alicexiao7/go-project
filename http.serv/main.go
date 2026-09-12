package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	
	"syscall"
	"time"

	"my-go-project/http.serv/internal/model"
	"my-go-project/http.serv/internal/storage"
)

type HelloDTO struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var tasks = make([]model.TaskDTO, 0)
var nextID = 1

func findTaskIndexByID(id int) int {
	for i, task := range tasks {
		if task.ID == id {
			return i
		}
	}
	return -1
}

func TaskByIDHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:

		i := findTaskIndexByID(id)

		if i == -1 {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tasks[i])

	case http.MethodDelete:

		i := findTaskIndexByID(id)

		if i == -1 {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}

		tasks = append(tasks[:i], tasks[i+1:]...)

		err := storage.SaveTasks(tasks)
		if err != nil {
			log.Println("save error", err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"message:" : "delete ok",
		})


	case http.MethodPut:

		var updatedTask model.TaskDTO

		err = json.NewDecoder(r.Body).Decode(&updatedTask)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		if updatedTask.Title == "" {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		i := findTaskIndexByID(id)

		if i == -1 {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}

		updatedTask.ID = id
		tasks[i] = updatedTask

		err := storage.SaveTasks(tasks)
		if err != nil {
			log.Println("save error", err)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(updatedTask)

	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
}

func TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
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

		if newTask.Title == "" {
			log.Println("empty title")
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		task := model.TaskDTO{
			ID:          nextID,
			Title:       newTask.Title,
			IsCompleted: false,
		}

		tasks = append(tasks, task)
		nextID += 1

		err = storage.SaveTasks(tasks)
		if err != nil {
			log.Println("save error", err)
		}

		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(http.StatusCreated)

		json.NewEncoder(w).Encode(task)

	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	var dto HelloDTO
	switch r.Method {
	case http.MethodGet:
		name := r.URL.Query().Get("name")
		ageS := r.URL.Query().Get("age")
		age := 0

		if ageS != "" {
			var err error
			age, err = strconv.Atoi(ageS)
			if err != nil {
				http.Error(w, "Bad Request", http.StatusBadRequest)
				return
			}
		}

		dto = HelloDTO{
			Name: name,
			Age:  age,
		}
	case http.MethodPost:
		err := json.NewDecoder(r.Body).Decode(&dto)

		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return

		}
	default:

		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(dto)
}

func updateNextID() {
	nextID = 1

	for _, task := range tasks {
		if task.ID >= nextID {
			nextID = task.ID + 1
		}
	}
}
func main() {
	var err error
	tasks, err = storage.LoadTasks()
	if err != nil {
		log.Println("Load error:", err)
		return
	}
	updateNextID()

	mux := http.NewServeMux()

	mux.HandleFunc("/", HelloHandler)
	mux.HandleFunc("/tasks", TaskHandler)
	mux.HandleFunc("/tasks/{id}", TaskByIDHandler)

	serv := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		log.Println("start")

		err := serv.ListenAndServe()

		if err != nil {
			log.Println("server error:", err)
		}
	}()

	stop := make(chan os.Signal, 1)

	signal.Notify(
		stop,
		os.Interrupt,
		syscall.SIGTERM,
	)

	<-stop

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)

	defer cancel()

	err = serv.Shutdown(ctx)

	if err != nil {
		log.Println("shutdown error:", err)
	}

	log.Println("server stop")
}
