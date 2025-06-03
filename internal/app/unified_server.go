package app

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/VladSnap/shortener/internal/config"
	grpchandlers "github.com/VladSnap/shortener/internal/grpc/handlers"
	"github.com/VladSnap/shortener/internal/grpc/interceptors"
	"github.com/VladSnap/shortener/internal/handlers"
	"github.com/VladSnap/shortener/internal/log"
	"github.com/VladSnap/shortener/internal/middlewares"
	pb "github.com/VladSnap/shortener/proto"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"golang.org/x/crypto/acme/autocert"
	"google.golang.org/grpc"
)

const (
	// Server configuration constants.
	shutdownTimeout   = 30 * time.Second
	signalChannelSize = 1
	errorChannelSize  = 2
	httpServerCount   = 1
	grpcServerCount   = 1
)

// UnifiedShortenerServer implements ShortenerServer interface and supports both HTTP and gRPC.
type UnifiedShortenerServer struct {
	opts            *config.Options
	postHandler     Handler
	getHandler      Handler
	shortenHandler  Handler
	pingHandler     Handler
	batchHandler    Handler
	urlsHandler     Handler
	deleteHandler   Handler
	getStatsHandler Handler
	grpcHandler     *grpchandlers.ShortenerGRPCHandler
}

// NewUnifiedShortenerServer creates a new unified server that supports both HTTP and gRPC.
func NewUnifiedShortenerServer(
	opts *config.Options,
	postHandler Handler,
	getHandler Handler,
	shortenHandler Handler,
	pingHandler Handler,
	batchHandler Handler,
	urlsHandler Handler,
	deleteHandler Handler,
	getStatsHandler Handler,
	service handlers.ShorterService,
	deleteWorker handlers.DeleterWorker,
) *UnifiedShortenerServer {
	grpcHandler := grpchandlers.NewShortenerGRPCHandler(service, deleteWorker, opts.BaseURL, opts)

	return &UnifiedShortenerServer{
		opts:            opts,
		postHandler:     postHandler,
		getHandler:      getHandler,
		shortenHandler:  shortenHandler,
		pingHandler:     pingHandler,
		batchHandler:    batchHandler,
		urlsHandler:     urlsHandler,
		deleteHandler:   deleteHandler,
		getStatsHandler: getStatsHandler,
		grpcHandler:     grpcHandler,
	}
}

// RunServer starts both HTTP and gRPC servers.
func (server *UnifiedShortenerServer) RunServer() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	errChan := make(chan error, errorChannelSize)

	// Start HTTP server
	wg.Add(httpServerCount)
	go func() {
		defer wg.Done()
		if err := server.runHTTPServer(ctx); err != nil {
			errChan <- fmt.Errorf("HTTP server error: %w", err)
		}
	}()

	// Start gRPC server
	wg.Add(grpcServerCount)
	go func() {
		defer wg.Done()
		if err := server.runGRPCServer(ctx); err != nil {
			errChan <- fmt.Errorf("gRPC server error: %w", err)
		}
	}()

	// Handle graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, signalChannelSize)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		<-sigChan
		log.Zap.Info("Received shutdown signal")
		cancel()
	}()

	// Wait for either an error or context cancellation
	select {
	case err := <-errChan:
		cancel()
		return err
	case <-ctx.Done():
		log.Zap.Info("Shutting down servers...")
	}

	// Wait for all servers to shut down
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	// Wait for shutdown to complete or timeout
	select {
	case <-done:
		log.Zap.Info("All servers shut down successfully")
		return nil
	case <-time.After(shutdownTimeout):
		return errors.New("timeout waiting for servers to shut down")
	}
}

// runHTTPServer starts the HTTP server.
func (server *UnifiedShortenerServer) runHTTPServer(ctx context.Context) error {
	httpRouter := server.initRouter()
	httpServer := &http.Server{
		Addr:    server.opts.ListenAddress,
		Handler: httpRouter,
	}

	// Start HTTP server in a goroutine
	go func() {
		log.Zap.Info("Starting HTTP server", zap.String("address", server.opts.ListenAddress))
		var err error
		if server.opts.EnableHTTPS != nil && *server.opts.EnableHTTPS {
			err = server.listenTLS(httpServer)
		} else {
			err = httpServer.ListenAndServe()
		}
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Zap.Error("HTTP server failed", zap.Error(err))
		}
	}()

	// Wait for context cancellation, then shut down
	<-ctx.Done()
	log.Zap.Info("Shutting down HTTP server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("failed to shutdown HTTP server: %w", err)
	}

	log.Zap.Info("HTTP server shut down successfully")
	return nil
}

// runGRPCServer starts the gRPC server
func (server *UnifiedShortenerServer) runGRPCServer(ctx context.Context) error {
	lis, err := net.Listen("tcp", server.opts.GRPCAddress)
	if err != nil {
		return fmt.Errorf("failed to listen on gRPC address %s: %w", server.opts.GRPCAddress, err)
	}

	// Create gRPC server with interceptors
	var grpcServer *grpc.Server
	if server.opts.TrustedSubnet != "" {
		// Add trusted subnet interceptor for stats endpoint
		grpcServer = grpc.NewServer(
			grpc.ChainUnaryInterceptor(
				interceptors.LoggingInterceptor(),
				interceptors.AuthInterceptor(server.opts),
				interceptors.TrustedSubnetInterceptor(server.opts.TrustedSubnet),
			),
		)
	} else {
		// No trusted subnet configured
		grpcServer = grpc.NewServer(
			grpc.ChainUnaryInterceptor(
				interceptors.LoggingInterceptor(),
				interceptors.AuthInterceptor(server.opts),
			),
		)
	}

	pb.RegisterShortenerServiceServer(grpcServer, server.grpcHandler)

	// Start gRPC server in a goroutine
	go func() {
		log.Zap.Info("Starting gRPC server", zap.String("address", server.opts.GRPCAddress))
		if err := grpcServer.Serve(lis); err != nil {
			log.Zap.Error("gRPC server failed", zap.Error(err))
		}
	}()

	// Wait for context cancellation, then shut down
	<-ctx.Done()
	log.Zap.Info("Shutting down gRPC server...")

	// Graceful stop with timeout
	stopped := make(chan struct{})
	go func() {
		grpcServer.GracefulStop()
		close(stopped)
	}()

	select {
	case <-stopped:
		log.Zap.Info("gRPC server shut down gracefully")
	case <-time.After(30 * time.Second):
		log.Zap.Warn("gRPC server graceful shutdown timeout, forcing stop")
		grpcServer.Stop()
	}

	return nil
}

// initRouter initializes the HTTP router (same as ChiShortenerServer)
func (server *UnifiedShortenerServer) initRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middlewares.LogMiddleware)
	r.Use(middlewares.GzipMiddleware)
	r.Use(middleware.Recoverer)

	r.Get("/{id}", server.getHandler.Handle)
	r.Get("/ping", server.pingHandler.Handle)

	// Routes with authentication
	r.Group(func(r chi.Router) {
		r.Use(middlewares.AuthMiddleware(server.opts))
		r.Post("/", server.postHandler.Handle)
		r.Post("/api/shorten", server.shortenHandler.Handle)
		r.Post("/api/shorten/batch", server.batchHandler.Handle)
		r.Get("/api/user/urls", server.urlsHandler.Handle)
		r.Delete("/api/user/urls", server.deleteHandler.Handle)
	})

	r.Group(func(r chi.Router) {
		r.Use(middlewares.TrustedSubnetMiddleware(server.opts.TrustedSubnet))
		r.Get("/api/internal/stats", server.getStatsHandler.Handle)
	})

	if server.opts.Performance != nil && *server.opts.Performance {
		r.Handle("/debug/pprof/*", http.DefaultServeMux)
	}
	return r
}

// listenTLS starts HTTPS server with automatic certificate management.
func (server *UnifiedShortenerServer) listenTLS(serv *http.Server) error {
	// Configure TLS certificate manager
	manager := &autocert.Manager{
		// Directory for certificate storage
		Cache: autocert.DirCache("cache-dir"),
		// Function that accepts Terms of Service from certificate issuer
		Prompt: autocert.AcceptTOS,
		// List of domains for which certificates will be supported
		HostPolicy: autocert.HostWhitelist("shortener.local"),
	}
	serv.TLSConfig = manager.TLSConfig()
	err := serv.ListenAndServeTLS("", "")
	return fmt.Errorf("failed listenTLS: %w", err)
}
