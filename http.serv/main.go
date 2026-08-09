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

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
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

	dto := HelloDTO{
		Name: name,
		Age:  age,
	}

	json.NewEncoder(w).Encode(dto)
}

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("/", HelloHandler)

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
