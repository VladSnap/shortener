package handlers

import (
	"context"

	"github.com/VladSnap/shortener/internal/services"
)

// Генерирует мок для ShorterService
//go:generate mockgen -destination=mocks/shorterService_mock.go -package mocks github.com/VladSnap/shortener/internal/handlers ShorterService

// ShorterService - интерфейс сервиса сокращателя ссылок, который реализует основную бизнес логику данного приложения.
type ShorterService interface {
	// CreateShortLink - Создает объект сокращенной ссылки для конкретного пользователя.
	CreateShortLink(ctx context.Context, originalURL string, userID string) (*services.ShortedLink, error)
	// CreateShortLinkBatch - Создает пачку объектов сокращенной ссылки для конкретного пользователя.
	CreateShortLinkBatch(ctx context.Context, originalLinks []*services.OriginalLink, userID string) (
		[]*services.ShortedLink, error)
	// GetURL - Читает полный URL по идентификатору сокращенной ссылки.
	GetURL(ctx context.Context, shortID string) (*services.ShortedLink, error)
	// GetAllByUserID - Читает все сокращенные ссылки конкретного пользователя.
	GetAllByUserID(ctx context.Context, userID string) ([]*services.ShortedLink, error)
	// DeleteBatch - Удаляет одной пачкой сокращенные ссылки.
	DeleteBatch(ctx context.Context, shortIDs []services.DeleteShortID) error
	// GetStats - Получает статистику о пользователях и всех ссылках.
	GetStats(ctx context.Context) (*services.Stats, error)
}
