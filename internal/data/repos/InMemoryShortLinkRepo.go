package repos

import (
	"github.com/VladSnap/shortener/internal/data/models"
)

type InMemoryShortLinkRepo struct {
	links map[string]*models.ShortLinkData
}

func NewShortLinkRepo() *InMemoryShortLinkRepo {
	repo := new(InMemoryShortLinkRepo)
	repo.links = make(map[string]*models.ShortLinkData)
	return repo
}

func (repo *InMemoryShortLinkRepo) CreateShortLink(link *models.ShortLinkData) (*models.ShortLinkData, error) {
	repo.links[link.ShortURL] = link
	return link, nil
}

func (repo *InMemoryShortLinkRepo) GetURL(shortID string) (*models.ShortLinkData, error) {
	return repo.links[shortID], nil
}
