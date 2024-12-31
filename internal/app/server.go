package app

import (
	//"fmt"

	"net/http"

	"github.com/VladSnap/shortener/internal/config"
	"github.com/VladSnap/shortener/internal/data"
	"github.com/VladSnap/shortener/internal/handlers"
	"github.com/VladSnap/shortener/internal/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func RunServer(opts *config.Options) error {
	var httpListener = initServer(opts)
	return http.ListenAndServe(opts.ListenAddress, httpListener)
}

func initServer(opts *config.Options) *chi.Mux {
	shortLinkRepo := data.NewShortLinkRepo()
	getHandler := handlers.NewGetHandler(shortLinkRepo)
	postHandler := handlers.NewPostHandler(shortLinkRepo, opts.BaseURL)

	r := chi.NewRouter()
	r.Use(middlewares.TimerTrace)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/", postHandler.Handle)
	r.Get("/{id}", getHandler.Handle)

	return r
}
