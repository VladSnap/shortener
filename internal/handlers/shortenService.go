package handlers

import (
	"context"

	"github.com/VladSnap/shortener/internal/services"
)

// Генерирует мок для ShorterService
//go:generate mockgen -destination=mock_handlers_test.go -package handlers github.com/VladSnap/shortener/internal/handlers ShorterService

type ShorterService interface {
	CreateShortLink(ctx context.Context, originalURL string, userID string) (*services.ShortedLink, error)
	CreateShortLinkBatch(ctx context.Context, originalLinks []*services.OriginalLink, userID string) (
		[]*services.ShortedLink, error)
	GetURL(ctx context.Context, shortID string) (*services.ShortedLink, error)
	GetAllByUserID(ctx context.Context, userID string) ([]*services.ShortedLink, error)
	DeleteBatch(ctx context.Context, shortIDs []string, userID string) error
}
