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
	"time"

	"github.com/VladSnap/shortener/internal/config"
	"github.com/VladSnap/shortener/internal/log"
	"github.com/VladSnap/shortener/internal/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"golang.org/x/crypto/acme/autocert"
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
	opts            *config.Options
	postHandler     Handler
	getHandler      Handler
	shortenHandler  Handler
	pingHandler     Handler
	batchHandler    Handler
	urlsHandler     Handler
	deleteHandler   Handler
	getStatsHandler Handler
}

// NewChiShortenerServer - Создает новую структуру ChiShortenerServer с указателем.
func NewChiShortenerServer(opts *config.Options,
	postHandler Handler,
	getHandler Handler,
	shortenHandler Handler,
	pingHandler Handler,
	batchHandler Handler,
	urlsHandler Handler,
	deleteHandler Handler,
	getStatsHandler Handler) *ChiShortenerServer {
	server := new(ChiShortenerServer)
	server.opts = opts
	server.postHandler = postHandler
	server.getHandler = getHandler
	server.shortenHandler = shortenHandler
	server.pingHandler = pingHandler
	server.batchHandler = batchHandler
	server.urlsHandler = urlsHandler
	server.deleteHandler = deleteHandler
	server.getStatsHandler = getStatsHandler
	return server
}

// RunServer - Запускает сервер.
func (server *ChiShortenerServer) RunServer() error {
	var httpListener = server.initRouter()
	// Создаем сервер.
	serv := &http.Server{Addr: server.opts.ListenAddress, Handler: httpListener}
	idleConnsClosed := runHandleGracefulShutdown(serv)
	err := server.runListen(serv, idleConnsClosed)
	if err != nil {
		return fmt.Errorf("failed run listen: %w", err)
	}
	return nil
}

func (server *ChiShortenerServer) initRouter() *chi.Mux {
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

	r.Group(func(r chi.Router) {
		r.Get("/api/internal/stats", server.getStatsHandler.Handle)
	})

	if server.opts.Performance != nil && *server.opts.Performance {
		r.Handle("/debug/pprof/*", http.DefaultServeMux)
	}
	return r
}

func (server *ChiShortenerServer) runListen(serv *http.Server,
	idleConnsClosed chan struct{}) error {
	var err error
	// Запускаем прослушивание запросов.
	if server.opts.EnableHTTPS != nil && *server.opts.EnableHTTPS {
		err = listenTLS(serv)
	} else {
		err = serv.ListenAndServe()
	}
	if err != nil {
		return fmt.Errorf("failed server listen: %w", err)
	}
	// Ждём завершения процедуры graceful shutdown чтобы закрыть все соединения.
	<-idleConnsClosed
	return nil
}

func listenTLS(serv *http.Server) error {
	// конструируем менеджер TLS-сертификатов
	manager := &autocert.Manager{
		// директория для хранения сертификатов
		Cache: autocert.DirCache("cache-dir"),
		// функция, принимающая Terms of Service издателя сертификатов
		Prompt: autocert.AcceptTOS,
		// перечень доменов, для которых будут поддерживаться сертификаты
		HostPolicy: autocert.HostWhitelist("shortener.local"),
	}
	serv.TLSConfig = manager.TLSConfig()
	err := serv.ListenAndServeTLS("", "")
	return fmt.Errorf("failed listenTLS: %w", err)
}

// runHandleGracefulShutdown - Обрабатывает завершение работы сервера.
// Принимает сигнал завершения и завершает работу сервера.
// После завершения работы сервера, закрывает канал idleConnsClosed,
// чтобы основной поток мог продолжить выполнение.
func runHandleGracefulShutdown(serv *http.Server) chan struct{} {
	idleConnsClosed := make(chan struct{})
	// Горутина для прослушивания сигналов завершения.
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		<-sigChan

		const shtdwnTimeout = 30 * time.Second
		shtdwnctx, cancel := context.WithTimeout(context.Background(), shtdwnTimeout)
		defer cancel()
		log.Zap.Info("Termination signal received. Stopping server....")
		if err := serv.Shutdown(shtdwnctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Zap.Error("Error while stopping the server", zap.Error(err))
		}
		// Сообщаем основному потоку,
		// что все сетевые соединения обработаны и закрыты
		close(idleConnsClosed)
	}()
	return idleConnsClosed
}
