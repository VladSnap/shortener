package app

import (
	//"fmt"

	"net/http"

	"github.com/VladSnap/shortener/internal/config"
	"github.com/VladSnap/shortener/internal/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Handler interface {
	Handle(res http.ResponseWriter, req *http.Request)
}

type ShortenerServer interface {
	RunServer() error
}

type ChiShortenerServer struct {
	opts        *config.Options
	postHandler Handler
	getHandler  Handler
}

func NewChiShortenerServer(opts *config.Options,
	postHandler Handler,
	getHandler Handler) *ChiShortenerServer {

	server := new(ChiShortenerServer)
	server.opts = opts
	server.postHandler = postHandler
	server.getHandler = getHandler
	return server
}

func (server *ChiShortenerServer) RunServer() error {
	var httpListener = server.initServer()
	return http.ListenAndServe(server.opts.ListenAddress, httpListener)
}

func (server *ChiShortenerServer) initServer() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middlewares.LogMiddleware)
	//r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/", server.postHandler.Handle)
	r.Get("/{id}", server.getHandler.Handle)

	return r
}
