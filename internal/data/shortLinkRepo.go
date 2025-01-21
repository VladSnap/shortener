package data

type ShortLinkRepo interface {
	CreateShortLink(shortID string, fullURL string) error
	GetURL(shortID string) string
}

type InMemoryShortLinkRepo struct {
	links map[string]string
}

func NewShortLinkRepo() *InMemoryShortLinkRepo {
	repo := new(InMemoryShortLinkRepo)
	repo.links = make(map[string]string)
	return repo
}

func (repo *InMemoryShortLinkRepo) CreateShortLink(shortID string, fullURL string) error {
	repo.links[shortID] = fullURL
	return nil
}

func (repo *InMemoryShortLinkRepo) GetURL(shortID string) string {
	return repo.links[shortID]
}
