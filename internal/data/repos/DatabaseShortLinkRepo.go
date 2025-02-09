package repos

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/VladSnap/shortener/internal/data"
)

type DatabaseShortLinkRepo struct {
	database *data.DatabaseShortener
}

func NewDatabaseShortLinkRepo(database *data.DatabaseShortener) *DatabaseShortLinkRepo {
	repo := new(DatabaseShortLinkRepo)
	repo.database = database
	return repo
}

func (repo *DatabaseShortLinkRepo) CreateShortLink(link *data.ShortLinkData) (*data.ShortLinkData, error) {
	// Пробуем найти по оригинальной ссылке сокращенную, чтобы не делать попытку записи,
	// т.к. в таблице есть ограничение на уникальность поля orig_url.
	existLink, ok, err := repo.getShortLinkByOriginalURL(link.OriginalURL)

	if err != nil {
		return nil, fmt.Errorf("failed getShortLinkByOriginalURL: %w", err)
	} else if ok {
		return existLink, nil // Вернем найденный результат, чтобы возвратить сокращенную ссылку в ответ на запрос.
	}

	sqlText := "INSERT INTO public.short_links (uuid, short_url, orig_url) VALUES ($1, $2, $3)"
	_, err = repo.database.ExecContext(context.Background(), sqlText, link.UUID, link.ShortURL, link.OriginalURL)
	if err != nil {
		return nil, fmt.Errorf("failed insert to public.short_links new row: %w", err)
	}
	return link, nil
}

func (repo *DatabaseShortLinkRepo) GetURL(shortID string) (*data.ShortLinkData, error) {
	return repo.GetShortLink(shortID)
}

func (repo *DatabaseShortLinkRepo) GetShortLink(shortID string) (*data.ShortLinkData, error) {
	sqlText := `SELECT uuid, short_url, orig_url
            FROM public.short_links
			WHERE short_url = $1`
	row := repo.database.QueryRowContext(context.Background(), sqlText, shortID)

	link := data.ShortLinkData{}
	// порядок переменных должен соответствовать порядку колонок в запросе
	err := row.Scan(&link.UUID, &link.ShortURL, &link.OriginalURL)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed select from public.short_links: %w", err)
	}

	return &link, nil
}

func (repo *DatabaseShortLinkRepo) getShortLinkByOriginalURL(originalURL string) (*data.ShortLinkData, bool, error) {
	sqlText := `SELECT uuid, short_url, orig_url
            FROM public.short_links
			WHERE orig_url = $1`
	row := repo.database.QueryRowContext(context.Background(), sqlText, originalURL)

	link := data.ShortLinkData{}
	// порядок переменных должен соответствовать порядку колонок в запросе
	err := row.Scan(&link.UUID, &link.ShortURL, &link.OriginalURL)
	if err != nil && err != sql.ErrNoRows {
		return nil, false, fmt.Errorf("failed select ByOriginalURL from public.short_links: %w", err)
	}
	if err == sql.ErrNoRows {
		return &link, false, nil
	}
	return &link, true, nil
}
