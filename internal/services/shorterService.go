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

// Генерирует мок для ShortLinkRepo
//go:generate mockgen -destination=mock_services_test.go -package services github.com/VladSnap/shortener/internal/services ShortLinkRepo

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
	id, shortID, err := createNewIds()
	if err != nil {
		return "", fmt.Errorf("failed create ids: %w", err)
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

func (service *NaiveShorterService) CreateShortLinkBatch(originalLinks []*OriginalLink) ([]*ShortedLink, error) {
	dataModels := make([]*data.ShortLinkData, 0, len(originalLinks))
	createdModels := make([]*ShortedLink, 0, len(originalLinks))

	if len(originalLinks) == 0 {
		return createdModels, nil
	}

	for _, ol := range originalLinks {
		id, shortID, err := createNewIds()
		if err != nil {
			return nil, fmt.Errorf("failed create ids: %w", err)
		}

		dm := &data.ShortLinkData{
			UUID:        id.String(),
			ShortURL:    shortID,
			OriginalURL: ol.URL,
		}
		dataModels = append(dataModels, dm)
		cm := &ShortedLink{
			UUID:         id.String(),
			CorelationID: ol.CorelationID,
			OriginalURL:  ol.URL,
			URL:          shortID,
		}
		createdModels = append(createdModels, cm)
	}

	_, err := service.shortLinkRepo.AddBatch(context.TODO(), dataModels)
	if err != nil {
		return nil, fmt.Errorf("failed add batch in repo: %w", err)
	}
	// todo: Тут по хорошему надо обновить ShortURL в моделях, если в репозитории будет логика проверки дублей

	return createdModels, nil
}

func createNewIds() (id uuid.UUID, shortID string, err error) {
	id, err = uuid.NewRandom()
	if err != nil {
		err = fmt.Errorf("failed create random: %w", err)
		return
	}
	shortID, err = helpers.RandStringRunes(shortIDLength)
	if err != nil {
		err = fmt.Errorf("failed create short url: %w", err)
		return
	}
	return
}
