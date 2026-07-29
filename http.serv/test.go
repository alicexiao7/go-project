package main

import (
	//"context" //lkz shutdown
	"net/http"
)

func base(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello")) // []byte("Hello")  строку в массив байТов
}

func main() {
	mux := http.NewServeMux() //станд маршрутизатор
	mux.HandleFunc("/", base)
	serv := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	serv.ListenAndServe()

	/*ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)

	defer cancel() //чистка контекста*/
	//serv.Shutdown(ctx)
}
