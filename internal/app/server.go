package app

import (
	//"fmt"

	"net/http"

	"github.com/go-chi/chi/v5"
"github.com/go-chi/chi/v5/middleware"
	"github.com/VladSnap/shortener/internal/data"
	"github.com/VladSnap/shortener/internal/handlers"
)

func RunServer() error {
	var httpListener = initServer()
	return http.ListenAndServe(`:8080`, httpListener)
}

func initServer() *chi.Mux {
	shortLinkRepo := data.NewShortLinkRepo()
	getHandler := handlers.NewGetHandler(shortLinkRepo)
	postHandler := handlers.NewPostHandler(shortLinkRepo)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)

	r.Post("/", postHandler.Handle)
	r.Get("/{id}", getHandler.Handle)

	return r
}
