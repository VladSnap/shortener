package app

import (
	//"fmt"
	"net/http"

	"github.com/VladSnap/shortener/internal/data"
	"github.com/VladSnap/shortener/internal/handlers"
)

func RunServer() error {
	var httpListener = initServer()
	return http.ListenAndServe(`:8080`, httpListener)
}

func initServer() *http.ServeMux {
	shortLinkRepo := data.NewShortLinkRepo()
	getHandler := handlers.NewGetHandler(shortLinkRepo)
	postHandler := handlers.NewPostHandler(shortLinkRepo)

	mux := http.NewServeMux()
	mux.HandleFunc(`/`, postHandler.Handle)
	mux.HandleFunc(`/{id}`, getHandler.Handle)

	return mux
}
