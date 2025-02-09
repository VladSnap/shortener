package handlers

import "github.com/VladSnap/shortener/internal/services"

type ShorterService interface {
	CreateShortLink(originalURL string) (string, error)
	CreateShortLinkBatch(originalLinks []*services.OriginalLink) ([]*services.ShortedLink, error)
	GetURL(shortID string) (string, error)
}
