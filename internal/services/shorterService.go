package services

import (
	"context"
	"fmt"

	"github.com/VladSnap/shortener/internal/data"
	"github.com/VladSnap/shortener/internal/helpers"
	"github.com/google/uuid"
)

type ShortLinkRepo interface {
	CreateShortLink(link *data.ShortLinkData) (*data.ShortLinkData, error)
	AddBatch(ctx context.Context, links []*data.ShortLinkData) ([]*data.ShortLinkData, error)
	GetURL(shortID string) (*data.ShortLinkData, error)
}

type NaiveShorterService struct {
	shortLinkRepo ShortLinkRepo
}

func NewNaiveShorterService(repo ShortLinkRepo) *NaiveShorterService {
	service := new(NaiveShorterService)
	service.shortLinkRepo = repo
	return service
}

const shortIDLength = 8

func (service *NaiveShorterService) CreateShortLink(originalURL string) (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("failed create random: %w", err)
	}
	shortID, err := helpers.RandStringRunes(shortIDLength)
	if err != nil {
		return "", fmt.Errorf("failed create short url: %w", err)
	}
	newLink := &data.ShortLinkData{UUID: id.String(), ShortURL: shortID, OriginalURL: originalURL}
	createdLink, err := service.shortLinkRepo.CreateShortLink(newLink)
	if err != nil {
		return "", fmt.Errorf("failed create short link object: %w", err)
	}
	// Важно вернуть сокращенную ссылку из created объекта, т.к. мы могли не создавать его повторно, если он существует
	return createdLink.ShortURL, nil
}

func (service *NaiveShorterService) GetURL(shortID string) (string, error) {
	link, err := service.shortLinkRepo.GetURL(shortID)
	if err != nil {
		return "", fmt.Errorf("failed get url from repo: %w", err)
	} else if link != nil {
		return link.OriginalURL, nil
	}
	return "", nil
}
