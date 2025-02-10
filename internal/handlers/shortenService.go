package handlers

import "github.com/VladSnap/shortener/internal/services"

// Генерирует мок для ShorterService
//go:generate mockgen -destination=mock_handlers_test.go -package handlers github.com/VladSnap/shortener/internal/handlers ShorterService

type ShorterService interface {
	CreateShortLink(originalURL string) (string, error)
	CreateShortLinkBatch(originalLinks []*services.OriginalLink) ([]*services.ShortedLink, error)
	GetURL(shortID string) (string, error)
}
