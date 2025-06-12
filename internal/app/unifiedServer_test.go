package app

import (
	"testing"

	"github.com/VladSnap/shortener/internal/config"
	"github.com/VladSnap/shortener/internal/data/repos"
	"github.com/VladSnap/shortener/internal/handlers"
	"github.com/VladSnap/shortener/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewUnifiedShortenerServerWithOptions(t *testing.T) {
	cfg := &config.Options{
		BaseURL: "http://localhost:8080",
	}
	resMng := services.NewResourceManager()
	defer resMng.Cleanup()

	t.Run("Creates server with individual options", func(t *testing.T) {
		// Создаем необходимые зависимости
		repo := repos.NewShortLinkRepo()
		shorterService := services.NewNaiveShorterService(repo)
		deleteWorker := handlers.NewDeleteWorker(shorterService)
		deleteWorker.RunWork()
		defer deleteWorker.Close()

		// Создаем обработчики
		postHandler := handlers.NewPostHandler(shorterService, cfg.BaseURL)
		getHandler := handlers.NewGetHandler(shorterService)
		shortenHandler := handlers.NewShortenHandler(shorterService, cfg.BaseURL)
		pingHandler := handlers.NewGetPingHandler(cfg)
		batchHandler := handlers.NewBatchHandler(shorterService, cfg.BaseURL)
		urlsHandler := handlers.NewUrlsHandler(shorterService, cfg.BaseURL)
		deleteHandler := handlers.NewDeleteHandler(deleteWorker)
		getStatsHandler := handlers.NewGetStatsHandler(cfg, shorterService)

		// Создаем сервер с отдельными опциями
		server, err := NewUnifiedShortenerServer(
			cfg,
			WithUnifiedPostHandler(postHandler),
			WithUnifiedGetHandler(getHandler),
			WithUnifiedShortenHandler(shortenHandler),
			WithUnifiedPingHandler(pingHandler),
			WithUnifiedBatchHandler(batchHandler),
			WithUnifiedUrlsHandler(urlsHandler),
			WithUnifiedDeleteHandler(deleteHandler),
			WithUnifiedGetStatsHandler(getStatsHandler),
			WithGRPCHandler(shorterService, deleteWorker, cfg.BaseURL, cfg),
		)

		require.NoError(t, err)
		assert.NotNil(t, server)
		assert.Equal(t, cfg, server.opts)
		assert.Equal(t, postHandler, server.postHandler)
		assert.Equal(t, getHandler, server.getHandler)
		assert.NotNil(t, server.grpcHandler)
	})
	t.Run("Creates empty server without options", func(t *testing.T) {
		server, err := NewUnifiedShortenerServer(cfg)

		require.NoError(t, err)
		assert.NotNil(t, server)
		assert.Equal(t, cfg, server.opts)
		assert.Nil(t, server.postHandler)
		assert.Nil(t, server.grpcHandler)
	})
}

func TestUnifiedServerOptionsValidation(t *testing.T) {
	cfg := &config.Options{
		BaseURL: "http://localhost:8080",
	}

	t.Run("Individual options set fields correctly", func(t *testing.T) {
		repo := repos.NewShortLinkRepo()
		shorterService := services.NewNaiveShorterService(repo)
		deleteWorker := handlers.NewDeleteWorker(shorterService)
		deleteWorker.RunWork()
		defer deleteWorker.Close()

		postHandler := handlers.NewPostHandler(shorterService, cfg.BaseURL)
		getHandler := handlers.NewGetHandler(shorterService)

		server := &UnifiedShortenerServer{opts: cfg}

		// Применяем отдельные опции
		postOption := WithUnifiedPostHandler(postHandler)
		getOption := WithUnifiedGetHandler(getHandler)
		grpcOption := WithGRPCHandler(shorterService, deleteWorker, cfg.BaseURL, cfg)

		err := postOption(server)
		require.NoError(t, err)
		err = getOption(server)
		require.NoError(t, err)
		err = grpcOption(server)
		require.NoError(t, err)

		assert.Equal(t, postHandler, server.postHandler)
		assert.Equal(t, getHandler, server.getHandler)
		assert.NotNil(t, server.grpcHandler)
	})
}
