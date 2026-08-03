package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type HelloDTO struct {
	Name string
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	name := r.URL.Query().Get("name")

	dto := HelloDTO{
		Name: name, //зачем запятая??
	}
	if dto.Name != "" {
		w.Write([]byte("Hello " + dto.Name))
	} else {
		w.Write([]byte("Hello"))
	}

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
