package data

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/VladSnap/shortener/internal/log"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type DatabaseShortener struct {
	*sql.DB
	Dsn string
}

func NewDatabaseShortener(dsn string) (*DatabaseShortener, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed open database: %w", err)
	}

	ds := &DatabaseShortener{db, dsn}
	return ds, nil
}

func (ds *DatabaseShortener) Close() error {
	err := ds.DB.Close()
	if err != nil {
		return fmt.Errorf("failed database connection close: %w", err)
	}
	log.Zap.Info("database connection closed")

	return nil
}

func (ds *DatabaseShortener) InitDatabase() error {
	driver, err := postgres.WithInstance(ds.DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to initialize postgres driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to initialize migrations: %w", err)
	}

	err = m.Up()
	noApply := errors.Is(err, migrate.ErrNoChange)

	if err != nil && !noApply {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	if !noApply {
		log.Zap.Info("Migrations applied successfully")
	}
	return nil
}
