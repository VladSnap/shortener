package repos

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/VladSnap/shortener/internal/data"
	"github.com/VladSnap/shortener/internal/log"
	"go.uber.org/zap"
)

// DatabaseShortLinkRepo - Репозиторий для доступа к БД сокращателя ссылок.
type DatabaseShortLinkRepo struct {
	database *data.DatabaseShortener
}

// NewDatabaseShortLinkRepo - Создает новую структуру DatabaseShortLinkRepo с указателем.
func NewDatabaseShortLinkRepo(database *data.DatabaseShortener) *DatabaseShortLinkRepo {
	repo := new(DatabaseShortLinkRepo)
	repo.database = database
	return repo
}

// Add - Сохраняет структуру сокращенной ссылки в БД.
func (repo *DatabaseShortLinkRepo) Add(ctx context.Context, link *data.ShortLinkData) (
	*data.ShortLinkData, error) {
	sqlText := "INSERT INTO public.short_links (uuid, short_url, orig_url, user_id, is_deleted)" +
		"VALUES ($1, $2, $3, $4, $5) " +
		"ON CONFLICT (orig_url) DO UPDATE " +
		"SET orig_url = short_links.orig_url " +
		"RETURNING short_links.short_url"

	//nolint:execinquery // use ON CONFLICT and Return value
	row := repo.database.QueryRowContext(ctx, sqlText, link.UUID, link.ShortURL,
		link.OriginalURL, toNullString(link.UserID), link.IsDeleted)
	if row.Err() != nil {
		return nil, fmt.Errorf("failed insert to public.short_links new row: %w", row.Err())
	}
	var shortURL string
	err := row.Scan(&shortURL)
	if err != nil {
		return nil, fmt.Errorf("failed scan insert result from public.short_links new row: %w", err)
	}
	if shortURL == link.ShortURL {
		return link, nil
	} else {
		return nil, data.NewDuplicateError(shortURL) //nolint:wrapcheck // is new error
	}
}

// AddBatch - Сохраняет пачку структур сокращенных ссылок в БД.
func (repo *DatabaseShortLinkRepo) AddBatch(ctx context.Context, links []*data.ShortLinkData) (
	[]*data.ShortLinkData, error) {
	tx, err := repo.database.BeginTx(ctx, nil)
	isCommited := false
	if err != nil {
		return nil, fmt.Errorf("failed begin db transaction before insert batch operation: %w", err)
	}
	defer func() {
		if !isCommited {
			err := tx.Rollback()
			if err != nil {
				log.Zap.Error("unable to rollback transaction after failed insert batch operation", zap.Error(err))
			}
		}
	}()

	stmt, err := tx.PrepareContext(ctx,
		"INSERT INTO public.short_links (uuid, short_url, orig_url, user_id, is_deleted)"+
			" VALUES($1, $2, $3, $4, $5)")
	if err != nil {
		return nil, fmt.Errorf("failed prepare insert: %w", err)
	}
	defer func() {
		err := stmt.Close()
		if err != nil {
			log.Zap.Error("unable to stmt close after insert batch operation", zap.Error(err))
		}
	}()

	for _, link := range links {
		_, err := stmt.ExecContext(ctx, link.UUID, link.ShortURL, link.OriginalURL,
			toNullString(link.UserID), link.IsDeleted)
		if err != nil {
			return nil, fmt.Errorf("failed exec insert batch: %w", err)
		}
	}
	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("failed commit insert batch transaction: %w", err)
	}
	isCommited = true

	return links, nil
}

// Get - Читает полную ссылку по сокращенной ссылке.
func (repo *DatabaseShortLinkRepo) Get(ctx context.Context, shortID string) (*data.ShortLinkData, error) {
	sqlText := `SELECT * FROM public.short_links WHERE short_url = $1`
	row := repo.database.QueryRowContext(ctx, sqlText, shortID)

	link := data.ShortLinkData{}
	var userID sql.NullString
	err := row.Scan(&link.UUID, &link.ShortURL, &link.OriginalURL, &userID, &link.IsDeleted)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("failed select from public.short_links: %w", err)
	}
	if userID.Valid {
		link.UserID = userID.String
	}

	return &link, nil
}

// GetAllByUserID - Получить все сокращенные ссылки указанного пользователя.
func (repo *DatabaseShortLinkRepo) GetAllByUserID(ctx context.Context, userID string) (
	[]*data.ShortLinkData, error) {
	sqlText := `SELECT * FROM public.short_links WHERE user_id = $1`
	rows, err := repo.database.QueryContext(ctx, sqlText, toNullString(userID))
	if err != nil {
		return nil, fmt.Errorf("failed select from public.short_links: %w", err)
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Zap.Error("failed rows close for select request", zap.Error(err))
		}
	}()

	const startSizeLinks int = 10
	links := make([]*data.ShortLinkData, 0, startSizeLinks)
	for rows.Next() {
		link := data.ShortLinkData{}
		// порядок переменных должен соответствовать порядку колонок в запросе
		err := rows.Scan(&link.UUID, &link.ShortURL, &link.OriginalURL, &link.UserID, &link.IsDeleted)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("failed scan select from public.short_links: %w", err)
		}
		links = append(links, &link)
	}

	if err := rows.Err(); err != nil {
		log.Zap.Error("last error encountered by Rows.Scan", zap.Error(err))
	}
	return links, nil
}

// DeleteBatch - Удаляет пачку структур сокращенных ссылок из БД.
func (repo *DatabaseShortLinkRepo) DeleteBatch(ctx context.Context, shortIDs []data.DeleteShortData) error {
	tx, err := repo.database.BeginTx(ctx, nil)
	isCommited := false
	if err != nil {
		return fmt.Errorf("failed begin db transaction before update batch operation: %w", err)
	}
	defer func() {
		if !isCommited {
			err := tx.Rollback()
			if err != nil {
				log.Zap.Error("unable to rollback transaction after failed update batch operation", zap.Error(err))
			}
		}
	}()

	stmt, err := tx.PrepareContext(ctx,
		"UPDATE public.short_links SET is_deleted=true WHERE is_deleted != true and short_url = $1 and user_id = $2")
	if err != nil {
		return fmt.Errorf("failed prepare batch update: %w", err)
	}
	defer func() {
		err := stmt.Close()
		if err != nil {
			log.Zap.Error("unable to stmt close after failed update batch operation", zap.Error(err))
		}
	}()

	for _, shortID := range shortIDs {
		_, err := stmt.ExecContext(ctx, shortID.ShortURL, shortID.UserID)
		if err != nil {
			return fmt.Errorf("failed exec batch update: %w", err)
		}
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed commit batch update transaction: %w", err)
	}
	isCommited = true

	return nil
}

func toNullString(input string) sql.NullString {
	if input == "" {
		return sql.NullString{String: "", Valid: false}
	}
	return sql.NullString{String: input, Valid: true}
}
