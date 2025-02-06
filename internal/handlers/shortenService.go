package handlers

type ShorterService interface {
	CreateShortLink(fullURL string) (string, error)
	GetURL(shortID string) (string, error)
}
