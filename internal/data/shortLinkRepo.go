package data

import (
	"github.com/VladSnap/shortener/internal/helpers"
)

type ShortLinkRepo interface {
	CreateShortLink(url string) string
	GetURL(key string) string
}

type InMemoryShortLinkRepo struct {
	links map[string]string
}

const linkKeyLength = 8

func NewShortLinkRepo() *InMemoryShortLinkRepo {
	repo := new(InMemoryShortLinkRepo)
	repo.links = make(map[string]string)
	return repo
}

func (repo *InMemoryShortLinkRepo) CreateShortLink(url string) string {
	key := helpers.RandStringRunes(linkKeyLength)
	repo.links[key] = url
	return key
}

func (repo *InMemoryShortLinkRepo) GetURL(key string) string {
	return repo.links[key]
}
