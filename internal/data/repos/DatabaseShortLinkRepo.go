package repos

import (
	"context"
	"fmt"

	"github.com/VladSnap/shortener/internal/data"
	"github.com/VladSnap/shortener/internal/data/models"
	"github.com/google/uuid"
)

type DatabaseShortLinkRepo struct {
	database *data.DatabaseShortener
}

func NewDatabaseShortLinkRepo(database *data.DatabaseShortener) *DatabaseShortLinkRepo {
	repo := new(DatabaseShortLinkRepo)
	repo.database = database
	return repo
}

func (repo *DatabaseShortLinkRepo) CreateShortLink(shortID string, fullURL string) error {
	id, err := uuid.NewRandom()
	if err != nil {
		return fmt.Errorf("failed create random: %w", err)
	}

	link := models.ShortLinkData{UUID: id.String(), ShortURL: shortID, OriginalURL: fullURL}

	sql := "INSERT INTO public.short_links (uuid, short_url, orig_url) VALUES ($1, $2, $3)"
	_, err = repo.database.ExecContext(context.Background(), sql, link.UUID, link.ShortURL, link.OriginalURL)
	if err != nil {
		return fmt.Errorf("failed insert to public.short_links new row: %w", err)
	}

	return nil
}

func (repo *DatabaseShortLinkRepo) GetURL(shortID string) (string, error) {
	sql := `SELECT uuid, short_url, orig_url
            FROM public.short_links
			WHERE short_url = $1`
	row := repo.database.QueryRowContext(context.Background(), sql, shortID)

	link := models.ShortLinkData{}
	// порядок переменных должен соответствовать порядку колонок в запросе
	err := row.Scan(&link.UUID, &link.ShortURL, &link.OriginalURL)
	if err != nil {
		return "", fmt.Errorf("failed select from public.short_links: %w", err)
	}

	return link.OriginalURL, nil
}
