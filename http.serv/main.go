package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"my-go-project/http.serv/internal/handler"
	"my-go-project/http.serv/internal/service"
	"my-go-project/http.serv/internal/storage"
)

func main() {
	// 1. Создаём конкретное хранилище
	fileStorage := storage.FileStorage{}

	// 2. Передаём storage в service
	taskService, err := service.NewTaskService(fileStorage)
	if err != nil {
		log.Println("service init error:", err)
		return
	}

	// 3. Передаём service в handler
	taskHandler := handler.NewTaskHandler(taskService)

	// 4. Регистрируем маршруты
	mux := http.NewServeMux()

	mux.HandleFunc("/tasks", taskHandler.Tasks)
	mux.HandleFunc("/tasks/{id}", taskHandler.TaskByID)

	// 5. Создаём HTTP-сервер
	serv := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		log.Println("server started")

		err := serv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Println("server error:", err)
		}
	}()

	// graceful shutdown
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

	log.Println("server stopped")
}