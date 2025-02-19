package repos

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/VladSnap/shortener/internal/data"
	"github.com/VladSnap/shortener/internal/log"
)

type DatabaseShortLinkRepo struct {
	database *data.DatabaseShortener
}

func NewDatabaseShortLinkRepo(database *data.DatabaseShortener) *DatabaseShortLinkRepo {
	repo := new(DatabaseShortLinkRepo)
	repo.database = database
	return repo
}

func (repo *DatabaseShortLinkRepo) Add(ctx context.Context, link *data.ShortLinkData) (
	*data.ShortLinkData, error) {
	sqlText := "INSERT INTO public.short_links (uuid, short_url, orig_url, user_id) VALUES ($1, $2, $3, $4) " +
		"ON CONFLICT (orig_url) DO UPDATE " +
		"SET orig_url = short_links.orig_url " +
		"RETURNING short_links.short_url"

	//nolint:execinquery // use ON CONFLICT and Return value
	row := repo.database.QueryRowContext(ctx, sqlText, link.UUID, link.ShortURL, link.OriginalURL, toNullString(link.UserID))
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

func (repo *DatabaseShortLinkRepo) AddBatch(ctx context.Context, links []*data.ShortLinkData) (
	[]*data.ShortLinkData, error) {
	tx, err := repo.database.BeginTx(ctx, nil)
	isCommited := false
	if err != nil {
		return nil, fmt.Errorf("failed begin db transaction: %w", err)
	}
	defer func() {
		if !isCommited {
			err := tx.Rollback()
			if err != nil {
				log.Zap.Errorf("failed Rollback: %w", err)
			}
		}
	}()

	stmt, err := tx.PrepareContext(ctx,
		"INSERT INTO public.short_links (uuid, short_url, orig_url, user_id)"+
			" VALUES($1, $2, $3, $4)")
	if err != nil {
		return nil, fmt.Errorf("failed prepare insert: %w", err)
	}
	defer func() {
		err := stmt.Close()
		if err != nil {
			log.Zap.Errorf("failed stmt Close: %w", err)
		}
	}()

	for _, link := range links {
		_, err := stmt.ExecContext(ctx, link.UUID, link.ShortURL, link.OriginalURL, toNullString(link.UserID))
		if err != nil {
			return nil, fmt.Errorf("failed exec insert: %w", err)
		}
	}
	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("failed commit transaction: %w", err)
	}
	isCommited = true

	return links, nil
}

func (repo *DatabaseShortLinkRepo) Get(ctx context.Context, shortID string) (*data.ShortLinkData, error) {
	sqlText := `SELECT * FROM public.short_links WHERE short_url = $1`
	row := repo.database.QueryRowContext(ctx, sqlText, shortID)

	link := data.ShortLinkData{}
	var userID sql.NullString
	err := row.Scan(&link.UUID, &link.ShortURL, &link.OriginalURL, &userID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("failed select from public.short_links: %w", err)
	}
	if userID.Valid {
		link.UserID = userID.String
	}

	return &link, nil
}

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
			log.Zap.Errorf("failed rows Close: %w", err)
		}
	}()

	const startSizeLinks int = 10
	links := make([]*data.ShortLinkData, 0, startSizeLinks)
	for rows.Next() {
		link := data.ShortLinkData{}
		// порядок переменных должен соответствовать порядку колонок в запросе
		err := rows.Scan(&link.UUID, &link.ShortURL, &link.OriginalURL, &link.UserID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("failed scan select from public.short_links: %w", err)
		}
		links = append(links, &link)
	}

	if err := rows.Err(); err != nil {
		log.Zap.Errorf("last error encountered by Rows.Scan: %w", err)
	}
	return links, nil
}

func toNullString(input string) sql.NullString {
	if input == "" {
		return sql.NullString{String: "", Valid: false}
	}
	return sql.NullString{String: input, Valid: true}
}
