package app

import (
	"testing"

	"github.com/VladSnap/shortener/internal/config"
	"github.com/VladSnap/shortener/internal/data/repos"
	"github.com/VladSnap/shortener/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestServerOptions тестирует создание и применение опций сервера.
func TestServerOptions(t *testing.T) {
	cfg := &config.Options{
		BaseURL: "http://localhost:8080",
	}
	resMng := services.NewResourceManager()
	defer func() {
		if err := resMng.Cleanup(); err != nil {
			t.Errorf("failed to cleanup ResourceManager: %v", err)
		}
	}()

	t.Run("NewServerOptions creates valid options", func(t *testing.T) {
		options := NewServerOptions(cfg, resMng)

		assert.NotNil(t, options)
		assert.Equal(t, cfg, options.GetConfig())
		assert.Equal(t, resMng, options.GetResourceManager())
	})

	t.Run("WithShortLinkRepo sets repository", func(t *testing.T) {
		options := NewServerOptions(cfg, resMng)
		repo := repos.NewShortLinkRepo()

		err := options.Apply(WithShortLinkRepo(repo))
		require.NoError(t, err)
		assert.Equal(t, repo, options.GetShortLinkRepo())
	})

	t.Run("Multiple options can be applied", func(t *testing.T) {
		options := NewServerOptions(cfg, resMng)
		repo := repos.NewShortLinkRepo()
		service := services.NewNaiveShorterService(repo)

		err := options.Apply(
			WithShortLinkRepo(repo),
			WithShorterService(service),
		)
		require.NoError(t, err)

		assert.Equal(t, repo, options.GetShortLinkRepo())
		assert.Equal(t, service, options.GetShorterService())
	})
}

// TestServerBuilder тестирует создание сервера с использованием паттерна Builder.
func TestServerBuilder(t *testing.T) {
	cfg := &config.Options{
		BaseURL:            "http://localhost:8080",
		DataBaseConnString: "",
		FileStoragePath:    "",
	}
	resMng := services.NewResourceManager()
	defer func() {
		if err := resMng.Cleanup(); err != nil {
			t.Errorf("failed to cleanup ResourceManager: %v", err)
		}
	}()

	t.Run("Builder creates valid server", func(t *testing.T) {
		server, err := NewServerBuilder(cfg, resMng).
			WithRepository().
			WithServices().
			WithHandlers().
			Build()

		require.NoError(t, err)
		assert.NotNil(t, server)

		// Проверяем, что это UnifiedShortenerServer
		unifiedServer, ok := server.(*UnifiedShortenerServer)
		assert.True(t, ok)
		assert.NotNil(t, unifiedServer)
	})

	t.Run("Builder fails without required dependencies", func(t *testing.T) {
		builder := NewServerBuilder(cfg, resMng)

		// Пытаемся собрать сервер без настройки зависимостей
		server, err := builder.Build()

		assert.Error(t, err)
		assert.Nil(t, server)
		assert.Contains(t, err.Error(), "is not configured")
	})

	t.Run("Builder works with partial configuration", func(t *testing.T) {
		builder := NewServerBuilder(cfg, resMng)

		// Настраиваем только репозиторий
		builder.WithRepository()
		options := builder.GetOptions()
		assert.NotNil(t, options.GetShortLinkRepo())
		assert.Nil(t, options.GetShorterService())
	})
}

// TestServerBuilder тестирует создание сервера с использованием паттерна Builder.
func TestCreateServerFunction(t *testing.T) {
	cfg := &config.Options{
		BaseURL:            "http://localhost:8080",
		DataBaseConnString: "",
		FileStoragePath:    "",
	}
	resMng := services.NewResourceManager()
	defer func() {
		if err := resMng.Cleanup(); err != nil {
			t.Errorf("failed to cleanup ResourceManager: %v", err)
		}
	}()

	t.Run("CreateServer creates valid server", func(t *testing.T) {
		server, err := CreateServer(cfg, resMng)

		require.NoError(t, err)
		assert.NotNil(t, server)
	})
}

// BenchmarkServerCreation измеряет производительность создания сервера.
func BenchmarkServerCreation(b *testing.B) {
	cfg := &config.Options{
		BaseURL:            "http://localhost:8080",
		DataBaseConnString: "",
		FileStoragePath:    "",
	}

	for b.Loop() {
		resMng := services.NewResourceManager()
		server, err := CreateServer(cfg, resMng)
		if err != nil {
			b.Fatal(err)
		}
		if server == nil {
			b.Fatal("server is nil")
		}
		err = resMng.Cleanup()
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkBuilderPattern измеряет производительность builder pattern.
func BenchmarkBuilderPattern(b *testing.B) {
	cfg := &config.Options{
		BaseURL:            "http://localhost:8080",
		DataBaseConnString: "",
		FileStoragePath:    "",
	}

	for b.Loop() {
		resMng := services.NewResourceManager()
		server, err := NewServerBuilder(cfg, resMng).
			WithRepository().
			WithServices().
			WithHandlers().
			Build()

		if err != nil {
			b.Fatal(err)
		}
		if server == nil {
			b.Fatal("server is nil")
		}
		err = resMng.Cleanup()
		if err != nil {
			b.Fatal(err)
		}
	}
}
