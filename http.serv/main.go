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
)

type HelloDTO struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type TaskDTO struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

var tasks = make([]TaskDTO, 0)
var nextID = 1

func TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tasks)

	case http.MethodPost:
		var newTask TaskDTO
		err := json.NewDecoder(r.Body).Decode(&newTask)

		if err != nil {
			log.Println("decode error:", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return

		}

		if newTask.Title == "" {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		task := TaskDTO{
			ID:        nextID,
			Title:     newTask.Title,
			Completed: false,
		}

		tasks = append(tasks, task)
		nextID += 1

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

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("/", HelloHandler)
	mux.HandleFunc("/tasks", TaskHandler)

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

	err := serv.Shutdown(ctx)

	if err != nil {
		log.Println("shutdown error:", err)
	}

	log.Println("server stop")
}
