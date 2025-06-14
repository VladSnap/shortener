package app

import (
	"github.com/VladSnap/shortener/internal/config"
	"github.com/VladSnap/shortener/internal/handlers"
	"github.com/VladSnap/shortener/internal/services"
)

// ServerOptions содержит все зависимости для создания сервера.
type ServerOptions struct {
	config          *config.Options
	resourceManager *services.ResourceManager

	// Repositories
	shortLinkRepo services.ShortLinkRepo

	// Services
	shorterService handlers.ShorterService
	deleteWorker   handlers.DeleterWorker

	// Handlers
	postHandler     Handler
	getHandler      Handler
	shortenHandler  Handler
	pingHandler     Handler
	batchHandler    Handler
	urlsHandler     Handler
	deleteHandler   Handler
	getStatsHandler Handler
}

// ServerOption представляет функцию для настройки ServerOptions.
type ServerOption func(*ServerOptions) error

// NewServerOptions создает новый экземпляр ServerOptions с базовой конфигурацией.
func NewServerOptions(cfg *config.Options, resMng *services.ResourceManager) *ServerOptions {
	return &ServerOptions{
		config:          cfg,
		resourceManager: resMng,
	}
}

// WithShortLinkRepo устанавливает репозиторий для коротких ссылок.
func WithShortLinkRepo(repo services.ShortLinkRepo) ServerOption {
	return func(opts *ServerOptions) error {
		opts.shortLinkRepo = repo
		return nil
	}
}

// WithShorterService устанавливает сервис для сокращения ссылок.
func WithShorterService(service handlers.ShorterService) ServerOption {
	return func(opts *ServerOptions) error {
		opts.shorterService = service
		return nil
	}
}

// WithDeleteWorker устанавливает воркер для удаления ссылок.
func WithDeleteWorker(worker handlers.DeleterWorker) ServerOption {
	return func(opts *ServerOptions) error {
		opts.deleteWorker = worker
		return nil
	}
}

// WithPostHandler устанавливает обработчик POST запросов.
func WithPostHandler(handler Handler) ServerOption {
	return func(opts *ServerOptions) error {
		opts.postHandler = handler
		return nil
	}
}

// WithGetHandler устанавливает обработчик GET запросов.
func WithGetHandler(handler Handler) ServerOption {
	return func(opts *ServerOptions) error {
		opts.getHandler = handler
		return nil
	}
}

// WithShortenHandler устанавливает обработчик сокращения ссылок.
func WithShortenHandler(handler Handler) ServerOption {
	return func(opts *ServerOptions) error {
		opts.shortenHandler = handler
		return nil
	}
}

// WithPingHandler устанавливает обработчик ping запросов.
func WithPingHandler(handler Handler) ServerOption {
	return func(opts *ServerOptions) error {
		opts.pingHandler = handler
		return nil
	}
}

// WithBatchHandler устанавливает обработчик batch запросов.
func WithBatchHandler(handler Handler) ServerOption {
	return func(opts *ServerOptions) error {
		opts.batchHandler = handler
		return nil
	}
}

// WithUrlsHandler устанавливает обработчик URLs запросов.
func WithUrlsHandler(handler Handler) ServerOption {
	return func(opts *ServerOptions) error {
		opts.urlsHandler = handler
		return nil
	}
}

// WithDeleteHandler устанавливает обработчик delete запросов.
func WithDeleteHandler(handler Handler) ServerOption {
	return func(opts *ServerOptions) error {
		opts.deleteHandler = handler
		return nil
	}
}

// WithGetStatsHandler устанавливает обработчик статистики.
func WithGetStatsHandler(handler Handler) ServerOption {
	return func(opts *ServerOptions) error {
		opts.getStatsHandler = handler
		return nil
	}
}

// Apply применяет все переданные опции к ServerOptions.
func (so *ServerOptions) Apply(options ...ServerOption) error {
	for _, option := range options {
		if err := option(so); err != nil {
			return err
		}
	}
	return nil
}

// GetConfig возвращает конфигурацию.
func (so *ServerOptions) GetConfig() *config.Options {
	return so.config
}

// GetResourceManager возвращает менеджер ресурсов.
func (so *ServerOptions) GetResourceManager() *services.ResourceManager {
	return so.resourceManager
}

// GetShortLinkRepo возвращает репозиторий коротких ссылок.
func (so *ServerOptions) GetShortLinkRepo() services.ShortLinkRepo {
	return so.shortLinkRepo
}

// GetShorterService возвращает сервис сокращения ссылок.
func (so *ServerOptions) GetShorterService() handlers.ShorterService {
	return so.shorterService
}

// GetDeleteWorker возвращает воркер удаления.
func (so *ServerOptions) GetDeleteWorker() handlers.DeleterWorker {
	return so.deleteWorker
}
