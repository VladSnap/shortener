package repos

import (
	"context"

	"github.com/VladSnap/shortener/internal/data"
)

type InMemoryShortLinkRepo struct {
	links map[string]*data.ShortLinkData
}

func NewShortLinkRepo() *InMemoryShortLinkRepo {
	repo := new(InMemoryShortLinkRepo)
	repo.links = make(map[string]*data.ShortLinkData)
	return repo
}

func (repo *InMemoryShortLinkRepo) CreateShortLink(link *data.ShortLinkData) (*data.ShortLinkData, error) {
	repo.links[link.ShortURL] = link
	return link, nil
}

func (repo *InMemoryShortLinkRepo) AddBatch(ctx context.Context, links []*data.ShortLinkData) (
	[]*data.ShortLinkData, error) {
	return nil, nil
}

func (repo *InMemoryShortLinkRepo) GetURL(shortID string) (*data.ShortLinkData, error) {
	return repo.links[shortID], nil
}
