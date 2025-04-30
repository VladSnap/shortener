// Package app хранит файлы для инициализации приложения shortener. Так же тут конфигурируется DI.
package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/VladSnap/shortener/internal/config"
	"github.com/VladSnap/shortener/internal/log"
	"github.com/VladSnap/shortener/internal/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

// Handler - Интерфейс обработчика http запросов.
type Handler interface {
	// Handle - Обработка запроса.
	Handle(res http.ResponseWriter, req *http.Request)
}

// ShortenerServer - Интерфейс сервера сокращателя ссылок.
type ShortenerServer interface {
	// RunServer - Запускает сервер.
	RunServer() error
}

// ChiShortenerServer - Импоементация сервера сокращателя ссылок ShortenerServer.
type ChiShortenerServer struct {
	opts           *config.Options
	postHandler    Handler
	getHandler     Handler
	shortenHandler Handler
	pingHandler    Handler
	batchHandler   Handler
	urlsHandler    Handler
	deleteHandler  Handler
}

// NewChiShortenerServer - Создает новую структуру ChiShortenerServer с указателем.
func NewChiShortenerServer(opts *config.Options,
	postHandler Handler,
	getHandler Handler,
	shortenHandler Handler,
	pingHandler Handler,
	batchHandler Handler,
	urlsHandler Handler,
	deleteHandler Handler) *ChiShortenerServer {
	server := new(ChiShortenerServer)
	server.opts = opts
	server.postHandler = postHandler
	server.getHandler = getHandler
	server.shortenHandler = shortenHandler
	server.pingHandler = pingHandler
	server.batchHandler = batchHandler
	server.urlsHandler = urlsHandler
	server.deleteHandler = deleteHandler
	return server
}

// RunServer - Запускает сервер.
func (server *ChiShortenerServer) RunServer() error {
	var httpListener = server.initServer()
	// Создаем сервер.
	serv := &http.Server{Addr: server.opts.ListenAddress, Handler: httpListener}
	// Горутина для прослушивания сигналов завершения.
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		log.Zap.Info("Termination signal received. Stopping server....")
		if err := serv.Shutdown(context.Background()); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Zap.Error("Error while stopping the server", zap.Error(err))
		}
	}()
	// Запускаем прослушивание запросов.
	err := serv.ListenAndServe()
	if err != nil {
		return fmt.Errorf("failed server listen: %w", err)
	}
	return nil
}

func (server *ChiShortenerServer) initServer() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middlewares.LogMiddleware)
	r.Use(middlewares.GzipMiddleware)
	r.Use(middleware.Recoverer)

	r.Get("/{id}", server.getHandler.Handle)
	r.Get("/ping", server.pingHandler.Handle)

	// Роутинги с аутентификацией.
	r.Group(func(r chi.Router) {
		r.Use(middlewares.AuthMiddleware(server.opts))
		r.Post("/", server.postHandler.Handle)
		r.Post("/api/shorten", server.shortenHandler.Handle)
		r.Post("/api/shorten/batch", server.batchHandler.Handle)
		r.Get("/api/user/urls", server.urlsHandler.Handle)
		r.Delete("/api/user/urls", server.deleteHandler.Handle)
	})

	if server.opts.Performance {
		r.Handle("/debug/pprof/*", http.DefaultServeMux)
	}
	return r
}
