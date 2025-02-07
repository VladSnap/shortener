package repos

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/VladSnap/shortener/internal/data"
	"github.com/VladSnap/shortener/internal/data/models"
)

type DatabaseShortLinkRepo struct {
	database *data.DatabaseShortener
}

func NewDatabaseShortLinkRepo(database *data.DatabaseShortener) *DatabaseShortLinkRepo {
	repo := new(DatabaseShortLinkRepo)
	repo.database = database
	return repo
}

func (repo *DatabaseShortLinkRepo) CreateShortLink(link *models.ShortLinkData) (*models.ShortLinkData, error) {
	// Пробуем найти по оригинальной ссылке сокращенную, чтобы не делать попытку записи,
	// т.к. в таблице есть ограничение на уникальность поля orig_url.
	existLink, err := repo.getShortLinkByOriginalURL(link.OriginalURL)

	if err != nil {
		return nil, fmt.Errorf("failed getShortLinkByOriginalURL: %w", err)
	} else if existLink != nil {
		return existLink, nil // Вернем найденный результат, чтобы возвратить сокращенную ссылку в ответ на запрос.
	}

	sql := "INSERT INTO public.short_links (uuid, short_url, orig_url) VALUES ($1, $2, $3)"
	_, err = repo.database.ExecContext(context.Background(), sql, link.UUID, link.ShortURL, link.OriginalURL)
	if err != nil {
		return nil, fmt.Errorf("failed insert to public.short_links new row: %w", err)
	}
	return link, nil
}

func (repo *DatabaseShortLinkRepo) GetURL(shortID string) (*models.ShortLinkData, error) {
	return repo.GetShortLink(shortID)
}

func (repo *DatabaseShortLinkRepo) GetShortLink(shortID string) (*models.ShortLinkData, error) {
	sqlText := `SELECT uuid, short_url, orig_url
            FROM public.short_links
			WHERE short_url = $1`
	row := repo.database.QueryRowContext(context.Background(), sqlText, shortID)

	link := models.ShortLinkData{}
	// порядок переменных должен соответствовать порядку колонок в запросе
	err := row.Scan(&link.UUID, &link.ShortURL, &link.OriginalURL)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed select from public.short_links: %w", err)
	}

	return &link, nil
}

func (repo *DatabaseShortLinkRepo) getShortLinkByOriginalURL(originalURL string) (*models.ShortLinkData, error) {
	sqlText := `SELECT uuid, short_url, orig_url
            FROM public.short_links
			WHERE orig_url = $1`
	row := repo.database.QueryRowContext(context.Background(), sqlText, originalURL)

	link := models.ShortLinkData{}
	// порядок переменных должен соответствовать порядку колонок в запросе
	err := row.Scan(&link.UUID, &link.ShortURL, &link.OriginalURL)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed select ByOriginalURL from public.short_links: %w", err)
	}

	return &link, nil
}
