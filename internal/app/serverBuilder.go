package app

import (
	"errors"
	"fmt"

	"github.com/VladSnap/shortener/internal/config"
	"github.com/VladSnap/shortener/internal/data"
	"github.com/VladSnap/shortener/internal/data/repos"
	"github.com/VladSnap/shortener/internal/handlers"
	"github.com/VladSnap/shortener/internal/services"
)

// ServerBuilder предоставляет fluent interface для создания сервера.
type ServerBuilder struct {
	options *ServerOptions
}

// NewServerBuilder создает новый билдер сервера.
func NewServerBuilder(cfg *config.Options, resMng *services.ResourceManager) *ServerBuilder {
	return &ServerBuilder{
		options: NewServerOptions(cfg, resMng),
	}
}

// CreateServer - Создает структуру интерфейса ShortenerServer используя Options pattern.
func CreateServer(opts *config.Options, resMng *services.ResourceManager) (ShortenerServer, error) {
	return NewServerBuilder(opts, resMng).
		WithRepository().
		WithServices().
		WithHandlers().
		Build()
}

// Apply применяет опции к builder'у.
func (sb *ServerBuilder) Apply(options ...ServerOption) *ServerBuilder {
	for _, option := range options {
		if err := option(sb.options); err != nil {
			continue
		}
	}
	return sb
}

// WithRepository настраивает репозиторий на основе конфигурации.
func (sb *ServerBuilder) WithRepository() *ServerBuilder {
	cfg := sb.options.GetConfig()
	resMng := sb.options.GetResourceManager()

	var shortLinkRepo services.ShortLinkRepo

	switch {
	case cfg.DataBaseConnString != "":
		database, err := data.NewDatabaseShortener(cfg.DataBaseConnString)
		if err != nil {
			// В реальном приложении лучше возвращать ошибку,
			// но для совместимости с существующим кодом используем panic
			panic(fmt.Errorf("failed create DatabaseShortener: %w", err))
		}
		resMng.Register(database.Close)
		err = database.InitDatabase()
		if err != nil {
			panic(fmt.Errorf("failed init Database: %w", err))
		}
		shortLinkRepo = repos.NewDatabaseShortLinkRepo(database)
	case cfg.FileStoragePath != "":
		fileRepo, err := repos.NewFileShortLinkRepo(cfg.FileStoragePath)
		if err != nil {
			panic(fmt.Errorf("failed create FileShortLinkRepo: %w", err))
		}
		resMng.Register(fileRepo.Close)
		shortLinkRepo = fileRepo
	default:
		shortLinkRepo = repos.NewShortLinkRepo()
	}

	err := sb.options.Apply(WithShortLinkRepo(shortLinkRepo))
	if err != nil {
		panic(fmt.Errorf("failed Apply ShortLinkRepo: %w", err))
	}

	return sb
}

// WithServices создает и настраивает все необходимые сервисы.
func (sb *ServerBuilder) WithServices() *ServerBuilder {
	shorterService := services.NewNaiveShorterService(sb.options.GetShortLinkRepo())
	deleteWorker := handlers.NewDeleteWorker(shorterService)

	sb.options.GetResourceManager().Register(deleteWorker.Close)
	deleteWorker.RunWork()

	err := sb.options.Apply(
		WithShorterService(shorterService),
		WithDeleteWorker(deleteWorker),
	)
	if err != nil {
		panic(fmt.Errorf("failed Apply Services: %w", err))
	}
	return sb
}

// WithHandlers создает и настраивает все обработчики HTTP запросов.
func (sb *ServerBuilder) WithHandlers() *ServerBuilder {
	cfg := sb.options.GetConfig()
	shorterService := sb.options.GetShorterService()
	deleteWorker := sb.options.GetDeleteWorker()

	postHandler := handlers.NewPostHandler(shorterService, cfg.BaseURL)
	getHandler := handlers.NewGetHandler(shorterService)
	shortenHandler := handlers.NewShortenHandler(shorterService, cfg.BaseURL)
	pingHandler := handlers.NewGetPingHandler(cfg)
	batchHandler := handlers.NewBatchHandler(shorterService, cfg.BaseURL)
	urlsHandler := handlers.NewUrlsHandler(shorterService, cfg.BaseURL)
	deleteHandler := handlers.NewDeleteHandler(deleteWorker)
	getStatsHandler := handlers.NewGetStatsHandler(cfg, shorterService)

	err := sb.options.Apply(
		WithPostHandler(postHandler),
		WithGetHandler(getHandler),
		WithShortenHandler(shortenHandler),
		WithPingHandler(pingHandler),
		WithBatchHandler(batchHandler),
		WithUrlsHandler(urlsHandler),
		WithDeleteHandler(deleteHandler),
		WithGetStatsHandler(getStatsHandler),
	)
	if err != nil {
		panic(fmt.Errorf("failed Apply Handlers: %w", err))
	}

	return sb
}

// Build создает окончательный экземпляр сервера.
func (sb *ServerBuilder) Build() (ShortenerServer, error) {
	// Проверяем, что все необходимые зависимости установлены
	if sb.options.shortLinkRepo == nil {
		return nil, errors.New("shortLinkRepo is not configured")
	}
	if sb.options.shorterService == nil {
		return nil, errors.New("shorterService is not configured")
	}
	if sb.options.deleteWorker == nil {
		return nil, errors.New("deleteWorker is not configured")
	}

	// Проверяем, что все обработчики установлены
	if sb.options.postHandler == nil || sb.options.getHandler == nil || sb.options.shortenHandler == nil ||
		sb.options.pingHandler == nil || sb.options.batchHandler == nil || sb.options.urlsHandler == nil ||
		sb.options.deleteHandler == nil || sb.options.getStatsHandler == nil {
		return nil, errors.New("not all handlers are configured")
	}

	// Создаем unified server используя отдельные опции
	server, err := NewUnifiedShortenerServer(
		sb.options.GetConfig(),
		WithUnifiedPostHandler(sb.options.postHandler),
		WithUnifiedGetHandler(sb.options.getHandler),
		WithUnifiedShortenHandler(sb.options.shortenHandler),
		WithUnifiedPingHandler(sb.options.pingHandler),
		WithUnifiedBatchHandler(sb.options.batchHandler),
		WithUnifiedUrlsHandler(sb.options.urlsHandler),
		WithUnifiedDeleteHandler(sb.options.deleteHandler),
		WithUnifiedGetStatsHandler(sb.options.getStatsHandler),
		WithGRPCHandler(
			sb.options.GetShorterService(),
			sb.options.GetDeleteWorker(),
			sb.options.GetConfig().BaseURL,
			sb.options.GetConfig(),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create unified server: %w", err)
	}

	return server, nil
}

// GetOptions возвращает текущие опции (для тестирования или отладки).
func (sb *ServerBuilder) GetOptions() *ServerOptions {
	return sb.options
}
