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

	"my-go-project/http.serv/internal/handler"
	"my-go-project/http.serv/internal/service"
	"my-go-project/http.serv/internal/storage"
)

type HelloDTO struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	fileStorage := storage.FileStorage{}
	taskService := service.NewTaskService(fileStorage)
	taskHandler := handler.NewTaskHandler(taskService)

	var err error

	mux := http.NewServeMux()

	mux.HandleFunc("/", HelloHandler)
	mux.HandleFunc("/tasks", taskHandler.Tasks)
	mux.HandleFunc("/tasks/{id}", taskHandler.TaskByID)

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

		dto = HelloDTO{Name: name, Age: age}

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
