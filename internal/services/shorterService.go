package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/VladSnap/shortener/internal/constants"
	"github.com/VladSnap/shortener/internal/data"
	"github.com/VladSnap/shortener/internal/helpers"
	"github.com/google/uuid"
)

type ShortLinkRepo interface {
	Add(ctx context.Context, link *data.ShortLinkData) (*data.ShortLinkData, error)
	AddBatch(ctx context.Context, links []*data.ShortLinkData) ([]*data.ShortLinkData, error)
	Get(ctx context.Context, shortID string) (*data.ShortLinkData, error)
	GetAllByUserID(ctx context.Context, userID string) ([]*data.ShortLinkData, error)
	DeleteBatch(ctx context.Context, shortIDs []string) error
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

func (service *NaiveShorterService) CreateShortLink(ctx context.Context,
	originalURL string, userID string) (*ShortedLink, error) {
	id, shortID, err := createNewIds()
	if err != nil {
		return nil, fmt.Errorf("failed create ids: %w", err)
	}
	newLink := data.NewShortLinkData(id.String(), shortID, originalURL, userID)
	createdLink, err := service.shortLinkRepo.Add(ctx, newLink)
	if err != nil {
		var duplErr *data.DuplicateShortLinkError
		if errors.As(err, &duplErr) {
			res := NewShortedLink("", "", "", duplErr.ShortURL, true, false)
			return res, nil
		}
		return nil, fmt.Errorf("failed create short link object: %w", err)
	}
	// Если короткие ссылки разные, значит был найден дубль и возвращено его значение.
	isDuplicate := shortID != createdLink.ShortURL
	res := NewShortedLink(createdLink.UUID, "", createdLink.OriginalURL, createdLink.ShortURL, isDuplicate, false)
	return res, nil
}

func (service *NaiveShorterService) GetURL(ctx context.Context, shortID string) (*ShortedLink, error) {
	link, err := service.shortLinkRepo.Get(ctx, shortID)
	if err != nil {
		return nil, fmt.Errorf("failed get url from repo: %w", err)
	} else if link != nil {
		return NewShortedLink(link.UUID, "", link.OriginalURL, link.ShortURL, false, link.IsDeleted), nil
	}
	return nil, nil
}

func (service *NaiveShorterService) CreateShortLinkBatch(ctx context.Context,
	originalLinks []*OriginalLink, userID string) ([]*ShortedLink, error) {
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
		dm := data.NewShortLinkData(id.String(), shortID, ol.URL, userID)
		dataModels = append(dataModels, dm)
		cm := NewShortedLink(id.String(), ol.CorelationID, ol.URL, shortID, false, false)
		createdModels = append(createdModels, cm)
	}

	_, err := service.shortLinkRepo.AddBatch(ctx, dataModels)
	if err != nil {
		return nil, fmt.Errorf("failed add batch in repo: %w", err)
	}
	// todo: Тут по хорошему надо обновить ShortURL в моделях, если в репозитории будет логика проверки дублей

	return createdModels, nil
}

func (service *NaiveShorterService) GetAllByUserID(ctx context.Context, userID string) (
	[]*ShortedLink, error) {
	links, err := service.shortLinkRepo.GetAllByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed GetAllByUserId: %w", err)
	}

	shortedLinks := make([]*ShortedLink, 0, len(links))
	for _, sl := range links {
		shortedLink := NewShortedLink(sl.UUID, "", sl.OriginalURL, sl.ShortURL, false, sl.IsDeleted)
		shortedLinks = append(shortedLinks, shortedLink)
	}

	return shortedLinks, nil
}

func (service *NaiveShorterService) DeleteBatch(ctx context.Context, shortIDs []string) error {
	// simple implement
	err := service.shortLinkRepo.DeleteBatch(ctx, shortIDs)
	if err != nil {
		return fmt.Errorf("failed DeleteBatch in repo: %w", err)
	}
	return nil
}

func createNewIds() (id uuid.UUID, shortID string, err error) {
	id, err = uuid.NewRandom()
	if err != nil {
		err = fmt.Errorf("failed create random: %w", err)
		return
	}
	shortID, err = helpers.RandStringRunes(constants.ShortIDLength)
	if err != nil {
		err = fmt.Errorf("failed create short url: %w", err)
		return
	}
	return
}
