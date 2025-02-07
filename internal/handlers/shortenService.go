package handlers

type ShorterService interface {
	CreateShortLink(originalURL string) (string, error)
	GetURL(shortID string) (string, error)
}
