package app

import (
	//"fmt"
	"net/http"

	"github.com/VladSnap/shortener/internal/handlers"
)

func RunServer() error {
	var httpListener = initServer()
	return http.ListenAndServe(`:8080`, httpListener)
}

func initServer() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, handlers.PostHandler)
	mux.HandleFunc(`/{id}`, handlers.GetHandler)

	return mux
}
