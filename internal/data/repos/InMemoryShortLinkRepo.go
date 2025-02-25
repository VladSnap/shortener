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

func (repo *InMemoryShortLinkRepo) Add(ctx context.Context, link *data.ShortLinkData) (*data.ShortLinkData, error) {
	repo.links[link.ShortURL] = link
	return link, nil
}

func (repo *InMemoryShortLinkRepo) AddBatch(ctx context.Context, links []*data.ShortLinkData) (
	[]*data.ShortLinkData, error) {
	for _, link := range links {
		repo.links[link.ShortURL] = link
	}
	return links, nil
}

func (repo *InMemoryShortLinkRepo) Get(ctx context.Context, shortID string) (*data.ShortLinkData, error) {
	return repo.links[shortID], nil
}

func (repo *InMemoryShortLinkRepo) GetAllByUserID(ctx context.Context, userID string) (
	[]*data.ShortLinkData, error) {
	return make([]*data.ShortLinkData, 0), nil
}

func (repo *InMemoryShortLinkRepo) DeleteBatch(ctx context.Context, shortIDs []string) error {
	for _, sid := range shortIDs {
		repo.links[sid].IsDeleted = true
	}
	return nil
}
